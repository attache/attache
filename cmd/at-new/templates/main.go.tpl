package main

import (
    "net/http"
    "log"

    "github.com/mccolljr/attache"
)

type Ctx struct {
    attache.DefaultDB
    attache.DefaultViews
    attache.DefaultToken
    attache.DefaultFileServer
}

func (c *Ctx) Init(w http.ResponseWriter, r *http.Request) {
    /* TODO: initialize context */
}

func (c *Ctx) GET_Index(r *http.Request) ([]byte, error) {
    return c.Views().Render("index", nil)
}

func main() {
    {{.Name}}, err := attache.Bootstrap(&Ctx{})
    if err != nil {
        log.Fatalln(err)
    }

    log.Fatalln({{.Name}}.Run())
}
