package main

import (
    "net/http"
    "log"
    
    // database drivers
    // (remove lines you don't need)
    _ "github.com/mattn/go-sqlite3" // Sqlite3
    _ "github.com/go-sql-driver/mysql" // MySQL
    _ "github.com/lib/pq" // PostgreSQL

    // attache
    "github.com/attache/attache"
)

type {{.Name}} struct {
    // required
    attache.BaseContext 

    // capabilities
    attache.DefaultEnvironment
    attache.DefaultFileServer
    attache.DefaultViews
    // attache.DefaultDB // enable database connectivity
    // attache.DefaultSession // enable session storage
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
    attache.RenderHTML(c, "index")
}

func main() {
    // bootstrap application for context type {{.Name}}
    app, err := attache.Bootstrap(&{{.Name}}{})
    if err != nil {
        log.Fatalln(err)
    }

    log.Fatalln(app.Run())
}
