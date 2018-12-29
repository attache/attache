package attache

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"unicode"

	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/securecookie"
)

// An Application routes HTTP traffic to an instance of its associated
// concrete Context type
type Application struct {
	r           router
	providers   stack
	contextType reflect.Type

	// NoLogging can be set to true to disable logging
	NoLogging  bool
	logHandler http.Handler
}

// ServeHTTP implements http.Handler for *Application.
func (a *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// call main handler func
	if a.NoLogging {
		a.baseHandler(w, r)
	} else {
		a.logHandler.ServeHTTP(w, r)
	}
}

func (a *Application) baseHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure a's recovery method runs.
	defer a.recover(w, r)

	// Try to locate the handler stack for the request's path.
	n := a.r.root.lookup(r.URL.Path)
	if n == nil || (!n.hasMount() && len(n.methods) == 0) {
		// No endpoint for the path; fail with a 404
		Error(404)
	}

	// Create a local copy of n's stack.
	var s stack
	s = append(s, n.guard...)

	// Short-circuit for mounted routes with no guards.
	if n.hasMount() && len(s) == 0 {
		n.mount.ServeHTTP(w, r)
		return
	}

	if n.hasMount() {
		// Add the mounted handler's ServeHTTP method to the stack
		s = append(s, reflect.ValueOf(n.mount.ServeHTTP))
	} else {
		if mainStack := n.methods[strings.ToUpper(r.Method)]; mainStack != nil {
			// Add the handlers from the main stack (optional BEFORE_..., [METHOD]_..., and optional AFTER_...)
			// to the request handler stack
			s = append(s, mainStack...)
		} else {
			// This particular HTTP method isn't allowed; fail with a 405
			Error(405)
		}
	}

	// Initialize an instance of the bootstrapped type.
	ctx := reflect.New(a.contextType.Elem()).Interface().(Context)

	// Initialize the context, or die with 500
	if err := initContextInstance(ctx, w, r); err != nil {
		log.Println(err)
		httpResult{code: 500}.ServeHTTP(w, r)
		return
	}

	// Create a new injector for this request.
	// TODO: should we cache these?
	injector := injector{
		app: a,
		ctx: ctx,
		req: r,
		res: w,
	}

	// Run all the registered providers and add their provided values to the injector.
	for _, x := range a.providers {
		result := x.Call(
			[]reflect.Value{
				reflect.ValueOf(ctx),
				reflect.ValueOf(r),
			},
		)

		injector.provided = append(injector.provided, reflect.ValueOf(result[0].Interface()))
	}

	// Execute the completed stack.
	for _, x := range s {
		injector.apply(x)
	}
}

// Initialize a Context instance for use
func initContextInstance(ictx Context, w http.ResponseWriter, r *http.Request) error {
	// set Request and ResponseWriter for this context
	bctx := ictx.baseContext()
	bctx.baseRw = w
	bctx.baseReq = r

	// Initialize views when context has view capability
	if impl, ok := ictx.(HasViews); ok {
		views, err := ViewCacheFor(impl.CONFIG_Views())
		if err != nil {
			return err
		}
		impl.SetViews(views)
	}

	// Initialize db when context has db capability
	if impl, ok := ictx.(HasDB); ok {
		db, err := DBFor(impl.CONFIG_DB())
		if err != nil {
			return err
		}

		impl.SetDB(db)
	}

	// Initialize session when context has session capability
	if impl, ok := ictx.(HasSession); ok {
		conf := impl.CONFIG_Session()

		s, err := gsSessions.Get(r, conf.Name)
		if err != nil {
			log.Println(err)
		}

		s.Options.HttpOnly = true
		impl.SetSession(Session{s})
	}

	// Initialize context
	ictx.Init(w, r)

	return nil
}

// recover is the deferred recovery handler run for each request
func (*Application) recover(w http.ResponseWriter, r *http.Request) {
	val := recover()
	if val == nil {
		return
	}

	if impl, ok := val.(http.Handler); ok {
		impl.ServeHTTP(w, r)
		return
	}

	pc := make([]uintptr, 4)
	pc = pc[:runtime.Callers(4, pc)]
	buf := strings.Builder{}
	for _, f := range pc {
		fn := runtime.FuncForPC(f)
		if fn != nil {
			fmt.Fprint(&buf, "\ntrace: ", fn.Name())
		}
	}

	log.Println("recovered: panic:", val, buf.String())

	httpResult{code: 500}.ServeHTTP(w, r)
}

// Run runs an HTTP server to handle requests to `a` on the
// default port, 8080
func (a *Application) Run() error { return a.RunAt(":8080") }

// RunAt runs an HTTP server to handle requests to `a`
func (a *Application) RunAt(addr string) error {
	return http.ListenAndServe(addr, a)
}

// RunWithServer mounts `a` to `s` and starts listening
func (a *Application) RunWithServer(s *http.Server) error {
	s.Handler = a
	return s.ListenAndServe()
}

// RunTLS runs an HTTP serverto handle requsts to `a` via TLS on the
// default port, 8443
func (a *Application) RunTLS(certFile, keyFile string) error {
	return a.RunAtTLS(":8443", certFile, keyFile)
}

// RunAtTLS runs an HTTP server to handle requests to `a` via TLS
func (a *Application) RunAtTLS(addr, certFile, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, a)
}

// RunWithServerTLS mounts `a` to `s` and starts listening via TLS
func (a *Application) RunWithServerTLS(s *http.Server, certFile, keyFile string) error {
	s.Handler = a
	return s.ListenAndServeTLS(certFile, keyFile)
}

var (
	methodRx = regexp.MustCompile(`^(GET|PUT|POST|PATCH|DELETE|HEAD|OPTIONS|TRACE|ALL)_(.*)$`)
)

// Bootstrap attempts to create an Application to serve requests for
// the provided concrete Context type. If an error is encountered
// during the bootstrapping process, it is returned.
// If a nil *Application is returned, the returned error will be non-nil.
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

	if err := bootstrapContextInit(ctxType); err != nil {
		return nil, err
	}

	if err := bootstrapRouter(&a, ctxType); err != nil {
		return nil, err
	}

	// set up logHandler
	a.logHandler = middleware.DefaultLogger(http.HandlerFunc(a.baseHandler))
	return &a, nil
}

func bootstrapContextInit(impl Context) error {
	// Attempt to load environment first, if supported by context
	if impl, ok := impl.(HasEnvironment); ok {
		conf := impl.CONFIG_Environment()
		if err := LoadEnvironment(conf); err != nil {
			return BootstrapError{Cause: err, Phase: "load environment"}
		}
	}

	// Attempt to load views, if supported by context
	if impl, ok := impl.(HasViews); ok {
		_, err := ViewCacheFor(impl.CONFIG_Views())
		if err != nil {
			return BootstrapError{Cause: err, Phase: "init views"}
		}
	}

	// Attempt db connection, if supported by context
	if impl, ok := impl.(HasDB); ok {
		_, err := DBFor(impl.CONFIG_DB())
		if err != nil {
			return BootstrapError{Cause: err, Phase: "init database"}
		}
	}

	// Examine session config for validity
	if impl, ok := impl.(HasSession); ok {
		conf := impl.CONFIG_Session()

		if len(conf.Secret) == 0 {
			return BootstrapError{Cause: errors.New("empty secret"), Phase: "check session config"}
		}

		gsSessions.Codecs = append(gsSessions.Codecs, securecookie.CodecsFromPairs(conf.Secret)...)
	}

	return nil
}

func bootstrapRouter(a *Application, impl Context) error {
	v := reflect.ValueOf(impl)
	t := v.Type()

	// Types only used for bootstrapping.
	// Defined in scope to unclutter the global namespace.
	type (
		// A guard represents a guard definition
		guard struct {
			path  string
			stack stack
		}

		// A route represents an endpoint definition
		route struct {
			method string
			path   string
			stack  stack
		}

		// A mount represents a mount definition
		mount struct {
			path    string
			handler http.Handler
		}
	)

	// Pre-allocate slices for guard, route, and mount definitions
	guards := make([]guard, 0, 32)
	routes := make([]route, 0, 32)
	mounts := make([]mount, 0, 32)

	// The function signature expected of MOUNT_ methods, as a reflect.Type
	mountFnTyp := reflect.TypeOf((func() (http.Handler, error))(nil))

	// The function signature expected of PROVIDE_ methods, as a reflect.Type
	provideFnTyp := reflect.TypeOf((func(*http.Request) interface{})(nil))

	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)

		// Provider methods
		if strings.HasPrefix(m.Name, "PROVIDE_") {
			fnTyp := v.Method(m.Index).Type()

			if !fnTyp.ConvertibleTo(provideFnTyp) {
				return BootstrapError{
					Cause: fmt.Errorf("%s does not have signature %s", m.Name, provideFnTyp),
					Phase: fmt.Sprint("check provider ", m.Name),
				}
			}

			a.providers = append(a.providers, m.Func)
		}

		// Route methods
		if match := methodRx.FindStringSubmatch(m.Name); match != nil {
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
			continue
		}

		// Guard methods
		if strings.HasPrefix(m.Name, "GUARD_") {
			path := pathForName(m.Name[6:] /* strip GUARD_ prefix */)

			guards = append(guards, guard{
				path:  path,
				stack: stack{m.Func},
			})
			continue
		}

		// Mount methods
		if strings.HasPrefix(m.Name, "MOUNT_") {
			path := pathForName(m.Name[6:])

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
					Cause: fmt.Errorf("error: %s", err),
					Phase: fmt.Sprintf("mount %s", path),
				}
			}

			mt.handler = h

			mounts = append(mounts, mt)
			continue
		}
	}

	// Bootstrap was called for a type that didn't provide any final request handlers (routes or mounts)ß.
	// This is most likely developer error.
	// Rather than silently continue, we'll warn the developer and fail.
	if len(routes)+len(mounts) == 0 {
		return BootstrapError{
			Phase: "register routes",
			Cause: errors.New("no routes defined"),
		}
	}

	// The order of insertion is important for guards, mounts, and routes.
	// In order to ensure correctness, we need to sort all 3 lists with
	// the same set of rules: by path length (short to long), then alphabetically.

	sort.SliceStable(guards, func(i, j int) bool {
		var (
			pathI, pathJ = guards[i].path, guards[j].path
			lenI, lenJ   = len(pathI), len(pathJ)
		)

		if lenI == lenJ {
			return pathI < pathJ
		}

		return lenI < lenJ
	})

	sort.SliceStable(mounts, func(i, j int) bool {
		var (
			pathI, pathJ = mounts[i].path, mounts[j].path
			lenI, lenJ   = len(pathI), len(pathJ)
		)

		if lenI == lenJ {
			return pathI < pathJ
		}

		return lenI < lenJ
	})

	sort.SliceStable(routes, func(i, j int) bool {
		var (
			pathI, pathJ = routes[i].path, routes[j].path
			lenI, lenJ   = len(pathI), len(pathJ)
		)

		if lenI == lenJ {
			return pathI < pathJ
		}

		return lenI < lenJ
	})

	// Once we've sorted the guards, mounts, and routes, we can
	// actually register them to the Application's router

	for _, g := range guards {
		if err := a.r.guard(g.path, g.stack); err != nil {
			return BootstrapError{
				Phase: fmt.Sprintf("guard %s", g.path),
				Cause: err,
			}
		}
	}

	for _, mt := range mounts {
		if err := a.r.mount(mt.path, mt.handler); err != nil {
			return BootstrapError{
				Phase: fmt.Sprintf("mount %s", mt.path),
				Cause: err,
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

	// Development: log the list of registered routes, etc.
	var b bytes.Buffer
	b.WriteString("\n======= ROUTES =======\n")
	dump(a.r.root, "", 0, &b)
	b.WriteString("======================")
	log.Println(b.String())

	return nil
}

// Calculates an HTTP request path based on a go method name.
func pathForName(name string) string {
	// Ignore (i.e remove) any underscores
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

	// Make sure we encode the last segment
	result.WriteByte('/')
	result.WriteString(strings.ToLower(name[start:]))

	return result.String()
}
