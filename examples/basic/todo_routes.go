package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/mccolljr/attache"
	"github.com/mccolljr/attache/examples/basic/models"
)

func (c *Ctx) GET_TodoNew(r *http.Request) ([]byte, error) {
	return c.Views.Render("todo.create", nil)
}

func (c *Ctx) GET_Todo(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.Todo)
	if err := c.DB.Find(target, id); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
		} else {
			log.Println(err)
			w.WriteHeader(500)
		}
		return
	}

	data, err := c.Views.Render("todo.update", &target)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("content-type", "text/html")
	w.Write(data)
}

func (c *Ctx) POST_TodoNew(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	target := new(models.Todo)

	if err := attache.FormDecode(target, r.Form); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if err := c.DB.Insert(target); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	fmt.Fprintf(w, "%v", target.ID)
}

func (c *Ctx) POST_Todo(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.Todo)
	if err := c.DB.Find(target, id); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
		} else {
			log.Println(err)
			w.WriteHeader(500)
		}
		return
	}

	if err := attache.FormDecode(target, r.Form); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if err := c.DB.Update(target); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}

func (c *Ctx) DELETE_Todo(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.Todo)
	if err := c.DB.Find(target, id); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(200) // treat as success
		} else {
			log.Println(err)
			w.WriteHeader(500)
		}
		return
	}

	if err := c.DB.Delete(target); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
}
