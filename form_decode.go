package attache

import (
	"reflect"

	"github.com/gorilla/schema"
)

// A FormConverter is a function that can convert a string value from an
// HTTP form to a concrete type
type FormConverter func(string) reflect.Value

// FormDecode decodes the form represented by src into dst.
// dst must be a pointer
func FormDecode(dst interface{}, src map[string][]string) error {
	return gsFormDecoder.Decode(dst, src)
}

// RegisterFormConverter registers a FormConverter function that will
// return values of type reflect.TypeOf(val)
func RegisterFormConverter(val interface{}, converterFunc FormConverter) {
	gsFormDecoder.RegisterConverter(val, schema.Converter(converterFunc))
}
