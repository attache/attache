package main

import (
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/mccolljr/attache"
)

type TodoApp struct {
	attache.BaseContext // required

	// default capability implementations
	attache.DefaultFileServer
	attache.DefaultViews
	attache.DefaultDB
	// attache.DefaultToken
}

func (c *TodoApp) Init(w http.ResponseWriter, r *http.Request) {
	/* TODO: initialize context */
}

// GET /
func (c *TodoApp) GET_(w http.ResponseWriter, r *http.Request) {
	c.GET_Index(w, r)
}

// GET /index
func (c *TodoApp) GET_Index(w http.ResponseWriter, r *http.Request) {
	attache.RenderHTML(c, "index", w, nil)
}

func main() {
	// bootstrap application for context type TodoApp
	app, err := attache.Bootstrap(&TodoApp{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(app.Run())
}
