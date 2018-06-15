package attache

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
)

func Error(code int) { ErrorMessage(code, "") }

func ErrorFatal(err error) {
	_, file, line, _ := runtime.Caller(1)
	log.Println("fatal:", file, line, err)
	Error(500)
}

func ErrorMessage(code int, msg string, args ...interface{}) {
	panic(httpResult{
		code:   code,
		status: fmt.Sprintf(msg, args...),
	})
}

func Success() { Error(200) }

func RedirectPage(path string) {
	panic(http.RedirectHandler(path, http.StatusSeeOther))
}

func RedirectPermanent(path string) {
	panic(http.RedirectHandler(path, http.StatusPermanentRedirect))
}

func RedirectTemporary(path string) {
	panic(http.RedirectHandler(path, http.StatusTemporaryRedirect))
}
