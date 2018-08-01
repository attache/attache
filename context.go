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
	Request() *http.Request
	ResponseWriter() http.ResponseWriter

	// require embedding of BaseContext
	baseContext() *BaseContext
}

// BaseContext must be embedded into any Context implementation,
// thus enforcing that all Context implementations are struct types
type BaseContext struct {
	baseReq *http.Request
	baseRw  http.ResponseWriter
}

// Request returns the *http.Request for the current request scope
func (b *BaseContext) Request() *http.Request { return b.baseReq }

// ResponseWriter returns the http.ResponseWriter for the current request scope
func (b *BaseContext) ResponseWriter() http.ResponseWriter { return b.baseRw }

// used internally to get a reference to the base context from the Context interface
func (b *BaseContext) baseContext() *BaseContext { return b }
