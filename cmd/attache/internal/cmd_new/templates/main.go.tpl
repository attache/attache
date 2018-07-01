package main

import (
    "net/http"
    "log"

    "github.com/mccolljr/attache"
)

type {{.Name}} struct {
    attache.BaseContext // required
    attache.DefaultFileServer
    attache.DefaultViews
    // attache.DefaultDB
    // attache.DefaultToken 
}

func (c *{{.Name}}) Init(w http.ResponseWriter, r *http.Request) {
    /* TODO: initialize context */
}

func (c *{{.Name}}) GET_Index(r *http.Request) ([]byte, error) {
    return c.Views().Render("index", nil)
}

func main() {
    app, err := attache.Bootstrap(&{{.Name}}{})
    if err != nil {
        log.Fatalln(err)
    }

    log.Fatalln(app.Run())
}
