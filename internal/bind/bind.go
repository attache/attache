package bind

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func init() {
	decoder.ZeroEmpty(false)
	decoder.SetAliasTag("json")
}

func Form(target interface{}, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	return decoder.Decode(target, r.Form)
}

func JSON(target interface{}, r *http.Request) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, target)
}
