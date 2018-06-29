package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/mccolljr/attache"
	"github.com/mccolljr/attache/examples/basic/models"
)

func (c *Ctx) GET_TodoNew(w http.ResponseWriter, r *http.Request) {
	data, err := c.Views().Render("todo.create", nil)
	if err != nil {
		attache.ErrorFatal(err)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func (c *Ctx) GET_TodoList(w http.ResponseWriter, r *http.Request) {
	log.Println(c)
	all, err := c.DB().All(func() attache.Storable { return new(models.Todo) })
	log.Println("successful query")
	if err != nil && err != sql.ErrNoRows {
		attache.ErrorFatal(err)
	}

	data, err := c.Views().Render("todo.list", all)
	if err != nil {
		attache.ErrorFatal(err)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

func (c *Ctx) GET_Todo(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.Todo)
	if err := c.DB().Find(target, id); err != nil {
		if err == sql.ErrNoRows {
			attache.Error(404)
		}

		attache.ErrorFatal(err)
	}

	data, err := c.Views().Render("todo.update", &target)
	if err != nil {
		attache.ErrorFatal(err)
	}

	w.Header().Set("content-type", "text/html")
	w.Write(data)
}

func (c *Ctx) POST_TodoNew(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		attache.ErrorFatal(err)
	}

	target := new(models.Todo)

	if err := attache.FormDecode(target, r.Form); err != nil {
		attache.ErrorFatal(err)
	}

	if err := c.DB().Insert(target); err != nil {
		attache.ErrorFatal(err)
	}

	attache.RedirectPage(fmt.Sprintf("/todo?id=%v", target.Title))
}

func (c *Ctx) POST_Todo(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.Todo)
	if err := c.DB().Find(target, id); err != nil {
		if err == sql.ErrNoRows {
			attache.Error(404)
		}

		attache.ErrorFatal(err)
	}

	if err := attache.FormDecode(target, r.Form); err != nil {
		attache.ErrorFatal(err)
	}

	if err := c.DB().Update(target); err != nil {
		attache.ErrorFatal(err)
	}

	attache.RedirectPage(fmt.Sprintf("/todo?id=%v", target.Title))
}

func (c *Ctx) DELETE_Todo(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.Todo)
	if err := c.DB().Find(target, id); err != nil {
		if err == sql.ErrNoRows {
			attache.Success()
		}

		attache.ErrorFatal(err)
	}

	if err := c.DB().Delete(target); err != nil {
		attache.ErrorFatal(err)
	}

	w.WriteHeader(200)
}
