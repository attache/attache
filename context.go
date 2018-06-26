package attache

import (
	"net/http"
)

var ctxContextKey = struct{ x int }{0xfeef}

type Context interface {
	Init(http.ResponseWriter, *http.Request)

	embeddedBase()
}

type BaseContext struct{}

func (b *BaseContext) embeddedBase() {}
