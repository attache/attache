package attache

import (
	"fmt"
	"net/http"
	"reflect"
)

type injector struct {
	app *Application
	ctx Context
	req *http.Request
	res http.ResponseWriter
}

var (
	tRequest        = reflect.TypeOf((*http.Request)(nil))
	tResponseWriter = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	tContext        = reflect.TypeOf((*Context)(nil)).Elem()
)

func (i injector) getFor(typ reflect.Type) reflect.Value {
	fmt.Println(typ)
	// special cases
	switch true {
	case typ == tRequest:
		return reflect.ValueOf(i.req)
	case typ == tResponseWriter:
		return reflect.ValueOf(i.res)
	case typ == tContext:
		fallthrough
	case typ.Kind() == reflect.Ptr && typ.Elem() == i.app.contextType:
		fmt.Println("context")
		return reflect.ValueOf(i.ctx)
	}

	// default
	return reflect.Zero(typ)
}

func (i *injector) apply(fn reflect.Value) {
	typ := fn.Type()
	args := make([]reflect.Value, typ.NumIn())
	for x := 0; x < typ.NumIn(); x++ {
		args[x] = i.getFor(typ.In(x))
	}
	// for now, ignore results
	_ = fn.Call(args)
}
