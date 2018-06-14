package attache

import "net/http"

type Context interface {
	Init(http.ResponseWriter, *http.Request)
}
