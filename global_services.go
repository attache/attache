package attache

import "github.com/gorilla/schema"

var (
	gsCache       = cache{}
	gsFormDecoder = schema.NewDecoder()
)

func init() {
	gsFormDecoder.IgnoreUnknownKeys(true)
}
