package attache

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type httpResult struct {
	code int
	msg  string
	json bool
}

func (x httpResult) MarshalJSON() ([]byte, error) {
	buf := getbuf()
	defer putbuf(buf)

	fmt.Fprintf(buf, `{"error": %q}`, x.msg)
	return buf.Bytes(), nil
}

func (x httpResult) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if x.code == 0 {
		x.code = 500
	}

	w.WriteHeader(x.code)

	if x.msg != "" {
		if !x.json {
			w.Write([]byte(x.msg))
			return
		}

		if err := json.NewEncoder(w).Encode(x); err != nil {
			panic(err)
		}
	}
}

func (x httpResult) String() string { return fmt.Sprintf("%d %s", x.code, x.msg) }
