package attache

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"runtime"
)

// Error is the equivalent of calling ErrorMessage with an empty message
func Error(code int) { ErrorMessage(code, "") }

// ErrorFatal logs err and then calls Error(500)
func ErrorFatal(err error) {
	_, file, line, ok := runtime.Caller(1)

	if ok {
		log.Printf("fatal: %s:%d %s", filepath.Base(file), line, err)
	} else {
		log.Println("fatal: (unknown loc)", err)
	}

	Error(500)
}

// ErrorMessage immediately terminates the executing handler chain with
// the given status code and status text
func ErrorMessage(code int, msg string, args ...interface{}) {
	panic(httpResult{
		code: code,
		msg:  fmt.Sprintf(msg, args...),
	})
}

// ErrorMessageJSON immediately terminates the executing handler chain with
// the given status code and a json body containing the status text
func ErrorMessageJSON(code int, msg string, args ...interface{}) {
	panic(httpResult{
		code: code,
		msg:  fmt.Sprintf(msg, args...),
		json: true,
	})
}

// Success immediately terminates the executing handler chain with
// a 200 OK
func Success() { Error(200) }

// RedirectPage immediately terminates the executing handler chain with a
// 303 (See Other)
func RedirectPage(path string) {
	panic(http.RedirectHandler(path, http.StatusSeeOther))
}

// RedirectPermanent immediately terminates the executing handler chain with a
// 308 (Permanent Redirect)
func RedirectPermanent(path string) {
	panic(http.RedirectHandler(path, http.StatusPermanentRedirect))
}

// RedirectTemporary immediately terminates the executing handler chain with a
// 307 (Temporary Redirect)
func RedirectTemporary(path string) {
	panic(http.RedirectHandler(path, http.StatusTemporaryRedirect))
}

// RenderHTML renders the view with the given name from the ViewCache attached
// to ctx using data as the Execute's data argument.
// The rendered template is written to w with a Content-Type of "text/html".
// ErrorFatal is called for any error encountered
func RenderHTML(ctx interface {
	Context
	HasViews
}, name string) {
	buf := getbuf()
	defer putbuf(buf)

	if err := ctx.Views().Get(name).Execute(buf, ctx); err != nil {
		panic(err)
	}

	ctx.ResponseWriter().Header().Set("content-type", "text/html")
	buf.WriteTo(ctx.ResponseWriter())
}

// RenderJSON marshals data to JSON, then writes the data to w with a
// Content-Type of "application/json". ErrorFatal is called for any
// error encountered
func RenderJSON(w http.ResponseWriter, data interface{}) {
	buf, err := json.Marshal(data)
	if err != nil {
		ErrorFatal(err)
	}

	w.Header().Set("content-type", "application/json")
	w.Write(buf)
}
