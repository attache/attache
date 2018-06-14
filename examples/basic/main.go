package main

import (
	"log"
	"net/http"

	"github.com/mccolljr/attache"

	_ "github.com/mattn/go-sqlite3"
)

type Ctx struct {
	Views attache.ViewCache
	DB    *attache.DB
}

func (c *Ctx) Init(w http.ResponseWriter, r *http.Request) { /* nothing yet */ }

// HasViews
func (c *Ctx) ViewRoot() string             { return "views" }
func (c *Ctx) SetViews(v attache.ViewCache) { c.Views = v }

// HasDB
func (c *Ctx) DBDriver() string     { return "sqlite3" }
func (c *Ctx) DBString() string     { return "test.db" }
func (c *Ctx) SetDB(db *attache.DB) { c.DB = db }

func (c *Ctx) GET_(w http.ResponseWriter, r *http.Request) { c.GET_Index(w, r) }

func (c *Ctx) GET_Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	c.Views.Get("index").Execute(w, nil)
}

func main() {
	app, err := attache.Bootstrap(&Ctx{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(app.Run())
}
