package attache

import (
	"net/http"
	"reflect"
)

type injector struct {
	ctx   Context
	req   *http.Request
	res   http.ResponseWriter
	cache map[reflect.Type]reflect.Value
}

var (
	tRequest        = reflect.TypeOf((*http.Request)(nil))
	tResponseWriter = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	tContext        = reflect.TypeOf((*Context)(nil)).Elem()
	tString         = reflect.TypeOf("")
)

func (i injector) getFor(typ reflect.Type, name string) reflect.Value {
	if v := i.cache[typ]; v.IsValid() {
		return v
	}

	// special cases
	switch true {
	case typ == tRequest:
		return i.putCache(tRequest, reflect.ValueOf(i.req))
	case typ == tResponseWriter:
		return i.putCache(tResponseWriter, reflect.ValueOf(i.res))
	case typ == tContext:
		return i.putCache(tContext, reflect.ValueOf(i.ctx))
	case typ.Implements(tContext):
		return i.putCache(typ, reflect.ValueOf(i.ctx))
	case typ == tString:
		// TODO continue
	}

	// default
	return i.putCache(typ, reflect.Zero(typ))
}

func (i injector) putCache(
	typ reflect.Type,
	val reflect.Value,
) reflect.Value {
	i.cache[typ] = val
	return val
}
