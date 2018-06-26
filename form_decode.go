package attache

import (
	"reflect"

	"github.com/gorilla/schema"
)

type FormConverter func(string) reflect.Value

func FormDecode(dst interface{}, src map[string][]string) error {
	return gsFormDecoder.Decode(dst, src)
}

func RegisterFormConverter(val interface{}, converterFunc FormConverter) {
	gsFormDecoder.RegisterConverter(val, schema.Converter(converterFunc))
}
