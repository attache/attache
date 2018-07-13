package attache

import (
	"net/http"
)

var ctxContextKey = struct{ x int }{0xfeef}

// Context is the base set of methods that a concrete Context type
// must provide. It contains private members, thus requiring
// that the BaseContext type be embedded in all types implementing
// Context
type Context interface {
	Init(http.ResponseWriter, *http.Request)

	embedRequired()
}

// BaseContext must be embedded into any Context implementation,
// thus enforcing that all Context implementations are struct types
type BaseContext struct{}

func (BaseContext) embedRequired() {}
