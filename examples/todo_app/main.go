package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
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

func (c *TodoApp) MOUNT_Testmt() http.Handler {
	mux := chi.NewMux()
	mux.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Println(r.URL.Path)
				h.ServeHTTP(w, r)
			},
		)
	})

	mux.Get("/test/{parm}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("matched 1")
		fmt.Fprintf(w, "%v", chi.URLParam(r, "parm"))
	})

	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("matched 2")
		fmt.Fprintf(w, "welcome")
	})
	return mux
}

func main() {
	// bootstrap application for context type TodoApp
	app, err := attache.Bootstrap(&TodoApp{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(app.Run())
}
