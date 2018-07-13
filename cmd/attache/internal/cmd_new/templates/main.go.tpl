package main

import (
    "net/http"
    "log"

    "github.com/mccolljr/attache"
)

type {{.Name}} struct {
    /* required */
    attache.BaseContext 

    /* default capability implementations */
    attache.DefaultFileServer
    attache.DefaultViews
    // attache.DefaultRequestResponse
    // attache.DefaultDB
    // attache.DefaultToken 
    // attache.DefaultSession
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
