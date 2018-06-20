package main

import (
	"log"
	"net/http"

	"github.com/mccolljr/attache"

	_ "github.com/mattn/go-sqlite3"
)

type Ctx struct {
	attache.BaseContext
}

func (c *Ctx) DBDriver() string { return "sqlite3" }
func (c *Ctx) DBString() string { return "test.db" }

func (c *Ctx) Init(w http.ResponseWriter, r *http.Request) { /* nothing yet */ }

func main() {
	app, err := attache.Bootstrap(&Ctx{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(app.Run())
}
