package attache

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"github.com/go-chi/chi/middleware"
)

// An Application routes HTTP traffic to an instance of its associated
// concrete Context type
type Application struct {
	r           router
	contextType reflect.Type
}

// ServeHTTP implements http.Handler for *Application
func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer a.recover(w, r)

	matched := a.r.root.lookup(r.URL.Path)
	if matched == nil || (!matched.isLeaf() && len(matched.methods) == 0) {
		Error(404)
	}

	stack := matched.stackFor(r.Method)
	if stack == nil {
		Error(405)
		return
	}

	// short-circuit context creation for
	// unguarded mounted routes, since we won't use it
	if matched.isLeaf() && len(stack) == 1 {
		stack[0].Call(
			[]reflect.Value{
				reflect.ValueOf(w),
				reflect.ValueOf(r),
			},
		)
		return
	}

	ctx := reflect.New(a.contextType.Elem()).Interface().(Context)

	// initialize the context, or die with 500
	if err := initContextInstance(ctx, w, r); err != nil {
		log.Println(err)
		httpResult{code: 500}.ServeHTTP(w, r)
		return
	}

	injector := injector{
		app: a,
		ctx: ctx,
		req: r,
		res: w,
	}

	for _, x := range stack {
		injector.apply(x)
	}
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

func (*Application) recover(w http.ResponseWriter, r *http.Request) {
	val := recover()
	if val == nil {
		return
	}

	if impl, ok := val.(http.Handler); ok {
		impl.ServeHTTP(w, r)
		return
	}

	pc := make([]uintptr, 12)
	pc = pc[:runtime.Callers(2, pc)]
	buf := strings.Builder{}
	for _, f := range pc {
		fn := runtime.FuncForPC(f)
		if fn != nil {
			fmt.Fprint(&buf, "\n", fn.Name())
		}
	}

	log.Println("recovered: panic:", val, buf.String())

	httpResult{code: 500}.ServeHTTP(w, r)
}

// Run runs an HTTP server serving requests for a on
// 0.0.0.0:8080
func (a *Application) Run() error {
	return http.ListenAndServe(":8080", middleware.DefaultLogger(a))
}

var (
	methodRx = regexp.MustCompile(`^(GET|PUT|POST|PATCH|DELETE|HEAD|OPTIONS|TRACE|ALL)_(.*)$`)
	mountRx  = regexp.MustCompile(`^MOUNT_(.*)$`)
)

// Bootstrap attempts to create an Application to serve requests for
// the provided concrete Context type. If an error is encountered
// during the bootstrapping process, it is returned.
// If a nil *Application is returned, the returned error will be non-nil
func Bootstrap(ctxType Context) (*Application, error) {
	var (
		v = reflect.ValueOf(ctxType)
		t = v.Type()
		a = Application{
			r: newrouter(),
		}
	)

	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("expecting pointer to a struct, got %T", ctxType)
	}

	a.contextType = t

	if err := bootstrapTryContextInit(ctxType); err != nil {
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
				Cause: errors.New(conf.Root + " is not a directory"),
				Phase: "init file server",
			}
		}

		basePath := path.Join("/", conf.BasePath)
		a.r.mount(basePath, http.FileServer(http.Dir(conf.Root)))
	}

	return nil
}

func bootstrapRoutes(a *Application, impl Context) error {
	v := reflect.ValueOf(impl)
	t := v.Type()
	found := 0

	type (
		route struct {
			method string
			path   string
			stack  stack
		}

		mount struct {
			path    string
			guard   reflect.Value
			handler http.Handler
		}
	)

	routes := make([]route, 0, 32)
	mounts := make([]mount, 0, 32)
	mountFnTyp := reflect.TypeOf((func() (http.Handler, error))(nil))

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)

		match := methodRx.FindStringSubmatch(m.Name)
		if match != nil {
			meth, path := match[1], pathForName(match[2])

			rt := route{
				method: meth,
				path:   path,
				stack:  make(stack, 0, 3),
			}

			if bm, ok := t.MethodByName("BEFORE_" + match[2]); ok {
				rt.stack = append(rt.stack, bm.Func)
			}

			rt.stack = append(rt.stack, m.Func)

			if am, ok := t.MethodByName("AFTER_" + match[2]); ok {
				rt.stack = append(rt.stack, am.Func)
			}

			routes = append(routes, rt)
			found++
			continue
		}

		match = mountRx.FindStringSubmatch(m.Name)
		if match != nil {
			path := pathForName(match[1])

			mt := mount{
				path: path,
			}

			fnVal := v.Method(m.Index)
			fnValTyp := fnVal.Type()

			if !fnValTyp.ConvertibleTo(mountFnTyp) {
				return BootstrapError{
					Cause: fmt.Errorf("%s does not have signature %s", m.Name, mountFnTyp),
					Phase: fmt.Sprintf("mount %s", path),
				}
			}

			h, err := fnVal.
				Convert(mountFnTyp).
				Interface().(func() (http.Handler, error))()

			if err != nil {
				return BootstrapError{
					Cause: fmt.Errorf("%s does not have signature %s", m.Name, mountFnTyp),
					Phase: fmt.Sprintf("mount %s", path),
				}
			}

			mt.handler = h

			if gm, ok := t.MethodByName("GUARD_" + match[1]); ok {
				mt.guard = gm.Func
			}

			mounts = append(mounts, mt)
			found++
			continue
		}
	}

	if found == 0 {
		return BootstrapError{
			Phase: "register routes",
			Cause: errors.New("no routes defined"),
		}
	}

	for _, mt := range mounts {
		if mt.guard.IsValid() {
			if err := a.r.mountGuarded(mt.path, mt.guard, mt.handler); err != nil {
				return BootstrapError{
					Phase: fmt.Sprintf("mount %s", mt.path),
					Cause: err,
				}
			}
		} else {
			if err := a.r.mount(mt.path, mt.handler); err != nil {
				return BootstrapError{
					Phase: fmt.Sprintf("mount %s", mt.path),
					Cause: err,
				}
			}
		}
	}

	for _, rt := range routes {
		if rt.method == "ALL" {
			if err := a.r.all(rt.path, rt.stack); err != nil {
				return BootstrapError{
					Phase: fmt.Sprintf("route %s %s", rt.method, rt.path),
					Cause: err,
				}
			}
		} else {
			if err := a.r.handle(rt.method, rt.path, rt.stack); err != nil {
				return BootstrapError{
					Phase: fmt.Sprintf("route %s %s", rt.method, rt.path),
					Cause: err,
				}
			}
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
