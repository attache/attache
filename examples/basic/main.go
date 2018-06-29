package main

import (
	"log"
	"net/http"

	"github.com/mccolljr/attache"

	_ "github.com/mattn/go-sqlite3"
)

type Ctx struct {
	// base context (required)
	attache.BaseContext

	// embedded capabilities
	attache.DefaultDB
	attache.DefaultViews
	attache.DefaultFileServer
}

func (c *Ctx) CONFIG_DB() attache.DBConfig {
	return attache.DBConfig{
		Driver: "sqlite3",
		DSN:    "test.db",
	}
}

func (c *Ctx) Init(w http.ResponseWriter, r *http.Request) { /* nothing yet */ }

func main() {
	app, err := attache.Bootstrap(&Ctx{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(app.Run())
}
