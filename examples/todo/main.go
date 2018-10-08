package main

import (
	"log"
	"net/http"

	// database drivers
	// (remove lines you don't need)
	// _ "github.com/go-sql-driver/mysql" // MySQL
	// _ "github.com/lib/pq"              // PostgreSQL
	_ "github.com/mattn/go-sqlite3" // Sqlite3

	// attache
	"github.com/attache/attache"
)

type Todo struct {
	// required
	attache.BaseContext

	// capabilities
	attache.DefaultFileServer
	attache.DefaultViews
	attache.DefaultDB // enable database connectivity
	// attache.DefaultSession // enable session storage
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
	attache.RenderHTML(c, "index")
}

func main() {
	// bootstrap application for context type Todo
	app, err := attache.Bootstrap(&Todo{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(app.Run())
}
