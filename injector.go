package attache

import (
	"log"
	"net/http"
	"reflect"
)

type injector struct {
	app      *Application
	ctx      Context
	req      *http.Request
	res      http.ResponseWriter
	provided []reflect.Value
}

var (
	tRequest        = reflect.TypeOf((*http.Request)(nil))
	tResponseWriter = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	tContext        = reflect.TypeOf((*Context)(nil)).Elem()
)

func (i injector) getFor(typ reflect.Type) reflect.Value {
	// special cases
	switch true {
	// provide *http.Request
	case typ == tRequest:
		return reflect.ValueOf(i.req)

	// provide http.ResponseWriter
	case typ == tResponseWriter:
		return reflect.ValueOf(i.res)

	// provide concrete Context where applicable
	case i.app.contextType.AssignableTo(typ):
		return reflect.ValueOf(i.ctx)

	// try lookup or default to the zero value
	default:
		for _, x := range i.provided {
			xtyp := x.Type()

			if xtyp.AssignableTo(typ) {
				return x
			}

			if x.CanAddr() && reflect.PtrTo(xtyp).AssignableTo(typ) {
				return x.Addr()
			}
		}

		log.Printf("inject: unrecognized type %s, using zero value", typ)
		return reflect.Zero(typ)
	}
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
