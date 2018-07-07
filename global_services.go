package attache

import (
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

var (
	gsCache       = cache{}
	gsFormDecoder = schema.NewDecoder()
	gsSessions    = sessions.NewCookieStore()
)

func init() {
	gsFormDecoder.IgnoreUnknownKeys(true)
}
