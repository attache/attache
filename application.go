package attache

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type (
	HandlerFunc = http.HandlerFunc
	RenderFunc  func(*http.Request) ([]byte, error)
)

func (fn RenderFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := fn(r)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

var (
	tHandlerFunc = reflect.TypeOf((*HandlerFunc)(nil)).Elem()
	tRenderFunc  = reflect.TypeOf((*RenderFunc)(nil)).Elem()
)

var ctxContextKey = struct{ x int }{0xfeef}

func reflectCall(method int, convertTo reflect.Type) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := r.Context().Value(ctxContextKey).(reflect.Value).Method(method)
		fn.Convert(convertTo).Interface().(http.Handler).ServeHTTP(w, r)
	}
}

type Application struct {
	router chi.Mux

	contextType reflect.Type
}

func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer a.recover(w, r)

	ctx := reflect.New(a.contextType)

	ictx := ctx.Interface()

	// initialize views when context has view capability
	if impl, ok := ictx.(HasViews); ok {
		views, err := viewsForRoot(impl.ViewRoot())
		if err != nil {
			panic(err)
		}
		impl.SetViews(views)
	}

	// initialize db when context has db capability
	if impl, ok := ictx.(HasDB); ok {
		db, err := openDB(impl.DBDriver(), impl.DBString())
		if err != nil {
			panic(err)
		}

		impl.SetDB(db)
	}

	// initialize context
	ictx.(Context).Init(w, r)

	// store context and execute handlers
	a.router.ServeHTTP(w, r.WithContext(
		context.WithValue(r.Context(), ctxContextKey, ctx),
	))
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

	if fn, ok := val.(http.HandlerFunc); ok {
		fn(w, r)
		return
	}

	log.Println("recovered: panic:", val)
	httpResult{500, ""}.ServeHTTP(w, r)
}

func (a *Application) Run() error { return http.ListenAndServe(":8080", a) }

func Bootstrap(ctxType Context) (*Application, error) {
	var (
		v = reflect.ValueOf(ctxType)
		t = v.Type()
		a = Application{
			router: *chi.NewMux(),
		}
	)

	a.router.Use(middleware.DefaultLogger)

	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("expecting pointer to a struct, got %T", ctxType)
	}

	a.contextType = t.Elem()

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)

		match := methodRx.FindStringSubmatch(m.Name)
		if match == nil {
			continue
		}

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
			if meth == "ALL" {
				a.router.Handle(path, reflectCall(i, convertTo))
			} else {
				a.router.Method(meth, path, reflectCall(i, convertTo))
			}
		}
	}

	return &a, nil
}

var methodRx = regexp.MustCompile(`^(GET|PUT|POST|PATCH|DELETE|HEAD|OPTIONS|TRACE|ALL)_(.*)$`)

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
