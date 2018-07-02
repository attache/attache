package attache

import (
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
		code:   code,
		status: fmt.Sprintf(msg, args...),
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

// RenderHTML immediately terminates the executing handler chain by
// rendering the view with the given name from the ViewCache attached
// to ctx using data as the View's data argument. If an error
// is encountered, a 500 Internal Server Error is written to w,
// otherwise the rendered template is written to w with a Content-Type
// of "text/html"
func RenderHTML(ctx HasViews, name string, w http.ResponseWriter, data interface{}) {
	buf, err := ctx.Views().Render(name, data)
	if err != nil {
		ErrorFatal(err)
	}

	w.Header().Set("content-type", "text/html")
	w.Write(buf)
}
