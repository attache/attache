package attache

import (
	"fmt"
	"net/http"
)

type httpResult struct {
	code   int
	status string
}

func (x httpResult) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if x.code == 0 {
		x.code = 500
	}

	if x.status == "" {
		x.status = http.StatusText(x.code)
	}

	w.WriteHeader(x.code)
	w.Write([]byte(x.status))
}

func (x httpResult) String() string { return fmt.Sprintf("%d %s", x.code, x.status) }
