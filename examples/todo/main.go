package main

import (
	"log"
	"net/http"

	"github.com/mccolljr/attache"
)

type Todo struct {
	// required
	attache.BaseContext

	// capabilities
	attache.DefaultFileServer
	attache.DefaultViews
	// attache.DefaultDB
	// attache.DefaultSession
}

func (c *Todo) Init(w http.ResponseWriter, r *http.Request) {
	/* TODO: initialize context */
}

// GET /
func (c *Todo) GET_(w http.ResponseWriter, r *http.Request) {
	c.GET_Index(w, r)
}

// GET /index
func (c *Todo) GET_Index(w http.ResponseWriter, r *http.Request) {
	attache.RenderHTML(c, "index", w, nil)
}

func main() {
	// bootstrap application for context type Todo
	app, err := attache.Bootstrap(&Todo{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(app.Run())
}
