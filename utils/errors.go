package utils

import (
	"fmt"
	"net/http"
)

type HTTPError struct {
	Code    int
	Message string
}

func Error(code int) { ErrorMessage(code, http.StatusText(code)) }

func ErrorMessage(code int, msg string, args ...interface{}) {
	panic(HTTPError{
		Code:    code,
		Message: fmt.Sprintf(msg, args...),
	})
}
