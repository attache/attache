package attache

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type (
	Middlewares = chi.Middlewares
	HandlerFunc = http.HandlerFunc

	RenderFunc         func(*http.Request) ([]byte, error)
	MiddlewareProvider func() Middlewares
)

func (fn RenderFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := fn(r)
	if err != nil {
		ErrorFatal(err)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

var (
	tHandlerFunc        = reflect.TypeOf((*HandlerFunc)(nil)).Elem()
	tRenderFunc         = reflect.TypeOf((*RenderFunc)(nil)).Elem()
	tMiddlewareProvider = reflect.TypeOf((*MiddlewareProvider)(nil)).Elem()
)

type reflectFn struct {
	index int
	typ   reflect.Type
}

func (f reflectFn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Context().
		Value(ctxContextKey).(reflect.Value).
		Method(f.index).
		Convert(f.typ).
		Interface().(http.Handler).
		ServeHTTP(w, r)
}

type mwlist []http.Handler

func (x mwlist) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, h := range x {
		h.ServeHTTP(w, r)
	}
}

func reffn(method int, convertTo reflect.Type) http.Handler {
	return reflectFn{index: method, typ: convertTo}
}

type Application struct {
	router chi.Mux

	contextType reflect.Type
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := reflect.New(a.contextType)

	ictx := ctx.Interface().(Context)

	// initialize the context, or die with 500
	if err := initContextInstance(ictx, w, r); err != nil {
		log.Println(err)
		httpResult{code: 500}.ServeHTTP(w, r)
		return
	}

	// store context and execute handlers
	a.router.ServeHTTP(w, r.WithContext(
		context.WithValue(r.Context(), ctxContextKey, ctx),
	))
}

func initContextInstance(ictx Context, w http.ResponseWriter, r *http.Request) error {
	// initialize views when context has view capability
	if impl, ok := ictx.(HasViews); ok {
		views, err := gsCache.viewsFor(impl.CONFIG_Views())
		if err != nil {
			return err
		}
		impl.SetViews(views)
	}

	// initialize db when context has db capability
	if impl, ok := ictx.(HasDB); ok {
		db, err := gsCache.dbFor(impl.CONFIG_DB())
		if err != nil {
			return err
		}

		impl.SetDB(db)
	}

	if impl, ok := ictx.(HasToken); ok {
		t := Token{
			conf: impl.CONFIG_Token(),
			Header: TokenHeader{
				Alg: "HS256",
				Typ: "JWT",
			},
			Claims: TokenClaims{},
		}

		cookie, _ := r.Cookie(t.conf.Cookie)

		if cookie != nil {
			if err := t.Decode([]byte(cookie.Value)); err != nil {
				t.ClearCookie(w)
				log.Println(err)
			}
		}

		impl.SetToken(t)
	}

	// initialize context
	ictx.Init(w, r)

	return nil
}

func (a *Application) recoveryMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer a.recover(w, r)
		h.ServeHTTP(w, r)
	})
}

func (*Application) recover(w http.ResponseWriter, r *http.Request) {
	val := recover()
	if val == nil {
		return
	}

	if impl, ok := val.(http.Handler); ok {
		impl.ServeHTTP(w, r)
		return
	}

	log.Println("recovered: panic:", val)
	httpResult{code: 500}.ServeHTTP(w, r)
}

func (a *Application) Run() error { return http.ListenAndServe(":8080", a) }

var (
	methodRx     = regexp.MustCompile(`^(GET|PUT|POST|PATCH|DELETE|HEAD|OPTIONS|TRACE|ALL)_(.*)$`)
	middlewareRx = regexp.MustCompile(`^USE_(.*)$`)
)

func Bootstrap(ctxType Context) (*Application, error) {
	var (
		v = reflect.ValueOf(ctxType)
		t = v.Type()
		a = Application{
			router: *chi.NewMux(),
		}
	)

	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("expecting pointer to a struct, got %T", ctxType)
	}

	a.contextType = t.Elem()

	if err := bootstrapTryContextInit(ctxType); err != nil {
		return nil, err
	}

	if err := bootstrapMiddleware(&a, ctxType); err != nil {
		return nil, err
	}

	if err := bootstrapFileServer(&a, ctxType); err != nil {
		return nil, err
	}

	if err := bootstrapRoutes(&a, ctxType); err != nil {
		return nil, err
	}

	return &a, nil
}

func bootstrapTryContextInit(impl Context) error {
	// attempt to load views, if supported by context
	if impl, ok := impl.(HasViews); ok {
		_, err := gsCache.viewsFor(impl.CONFIG_Views())
		if err != nil {
			return BootstrapError{Cause: err, Phase: "init views"}
		}
	}

	// attempt db connection, if supported by context
	if impl, ok := impl.(HasDB); ok {
		_, err := gsCache.dbFor(impl.CONFIG_DB())
		if err != nil {
			return BootstrapError{Cause: err, Phase: "init database"}
		}
	}

	// examine token config for validity
	if impl, ok := impl.(HasToken); ok {
		conf := impl.CONFIG_Token()

		if len(conf.Secret) == 0 {
			return BootstrapError{Cause: errors.New("empty secret"), Phase: "check token config"}
		}

		if len(conf.Cookie) == 0 {
			return BootstrapError{Cause: errors.New("empty cookie"), Phase: "check token config"}
		}
	}

	return nil
}

func bootstrapMiddleware(a *Application, impl Context) error {
	stack := Middlewares{
		middleware.DefaultLogger,
		a.recoveryMiddleware,
	}

	if impl, ok := impl.(HasMiddleware); ok {
		stack = append(stack, impl.Middleware()...)
	}

	a.router.Use(stack...)

	return nil
}

func bootstrapFileServer(a *Application, impl Context) error {
	if impl, ok := impl.(HasFileServer); ok {
		conf := impl.CONFIG_FileServer()
		info, err := os.Stat(conf.Root)
		if err != nil {
			return BootstrapError{
				Cause: err,
				Phase: "init file server",
			}
		}

		if !info.IsDir() {
			return BootstrapError{
				Cause: fmt.Errorf("bootstrap: static files: %s is not a directory", conf.Root),
				Phase: "init file server",
			}
		}

		basePath := path.Join("/", conf.BasePath)

		a.router.Handle(
			path.Join(basePath, "*"),
			http.StripPrefix(
				basePath,
				http.FileServer(http.Dir(conf.Root)),
			),
		)
	}

	return nil
}

func bootstrapRoutes(a *Application, impl Context) error {
	var (
		v     = reflect.ValueOf(impl)
		t     = v.Type()
		found = 0
	)

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)

		match := methodRx.FindStringSubmatch(m.Name)
		if match == nil {
			continue
		}

		found++

		meth, path := match[1], pathForName(match[2])

		mtyp := v.Method(i).Type()
		var convertTo reflect.Type
		switch true {
		case mtyp.ConvertibleTo(tHandlerFunc):
			convertTo = tHandlerFunc
		case mtyp.ConvertibleTo(tRenderFunc):
			convertTo = tRenderFunc
		}

		if convertTo != nil {
			list := make(mwlist, 0, 3)

			if m, ok := t.MethodByName("BEFORE_" + match[2]); ok {
				mtyp := v.Method(m.Index).Type()
				if mtyp.ConvertibleTo(tHandlerFunc) {
					list = append(list, reffn(m.Index, tHandlerFunc))
				}
			}

			list = append(list, reffn(i, convertTo))

			if m, ok := t.MethodByName("AFTER_" + match[2]); ok {
				mtyp := v.Method(m.Index).Type()
				if mtyp.ConvertibleTo(tHandlerFunc) {
					list = append(list, reffn(m.Index, tHandlerFunc))
				}
			}

			if meth == "ALL" {
				a.router.Handle(path, list)
			} else {
				a.router.Method(meth, path, list)
			}
		}
	}

	if found == 0 {
		return BootstrapError{
			Phase: "register routes",
			Cause: errors.New("no routes found"),
		}
	}

	return nil
}

func pathForName(name string) string {
	// sanitize name
	name = strings.Replace(name, "_", "", -1)
	result := strings.Builder{}
	size := len(name)
	start := 0
	lastUpper := false

	for i := 0; i < size; i++ {
		r := rune(name[i])
		if unicode.IsUpper(r) {
			if !lastUpper && i != start {
				result.WriteByte('/')
				result.WriteString(strings.ToLower(name[start:i]))
				start = i
			}
			lastUpper = true
		} else {
			if lastUpper && i-1 != start {
				result.WriteByte('/')
				result.WriteString(strings.ToLower(name[start : i-1]))
				start = i - 1
			}
			lastUpper = false
		}
	}

	result.WriteByte('/')
	result.WriteString(strings.ToLower(name[start:]))

	return result.String()
}
