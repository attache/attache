package main

import (
    "net/http"
    "log"

    "github.com/mccolljr/attache"
)

type {{.Name}} struct {
    attache.BaseContext // required

    // default capability implementations
    attache.DefaultFileServer
    attache.DefaultViews
    attache.DefaultDB
    attache.DefaultToken 
}

func (c *{{.Name}}) Init(w http.ResponseWriter, r *http.Request) {
    /* TODO: initialize context */
}

// GET /
func (c *{{.Name}}) GET_(w http.ResponseWriter, r *http.Request) {
    c.GET_Index(w, r)
}

// GET /index
func (c *{{.Name}}) GET_Index(w http.ResponseWriter, r *http.Request) {
    attache.RenderHTML(c, "index", w, nil)
}

func main() {
    // bootstrap application for context type {{.Name}}
    app, err := attache.Bootstrap(&{{.Name}}{})
    if err != nil {
        log.Fatalln(err)
    }

    log.Fatalln(app.Run())
}
