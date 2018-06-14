package attache

import (
	"github.com/gorilla/schema"
)

var global_formDecoder = schema.NewDecoder()

func FormDecode(dst interface{}, src map[string][]string) error {
	return global_formDecoder.Decode(dst, src)
}

func FormConverter(val interface{}, converterFunc schema.Converter) {
	global_formDecoder.RegisterConverter(val, converterFunc)
}
