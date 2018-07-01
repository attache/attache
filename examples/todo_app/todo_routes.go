package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/mccolljr/attache"
	"github.com/mccolljr/attache/examples/todo_app/models"
)

func (c *TodoApp) GET_TodoNew(w http.ResponseWriter, r *http.Request) {
	attache.RenderHTML(c, "todo.create", w, nil)
}

func (c *TodoApp) GET_TodoList(w http.ResponseWriter, r *http.Request) {
	all, err := c.DB().All(func() attache.Storable { return new(models.Todo) })
	if err != nil && err != sql.ErrNoRows {
		attache.ErrorFatal(err)
	}

	attache.RenderHTML(c, "todo.list", w, all)
}

func (c *TodoApp) GET_Todo(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.Todo)
	if err := c.DB().Find(target, id); err != nil {
		if err == sql.ErrNoRows {
			attache.Error(404)
		}

		attache.ErrorFatal(err)
	}

	attache.RenderHTML(c, "todo.update", w, &target)
}

func (c *TodoApp) POST_TodoNew(w http.ResponseWriter, r *http.Request) {
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

	attache.RedirectPage(fmt.Sprintf("/todo?id=%v", target.ID))
}

func (c *TodoApp) POST_Todo(w http.ResponseWriter, r *http.Request) {
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

	attache.RedirectPage(fmt.Sprintf("/todo?id=%v", target.ID))
}

func (c *TodoApp) DELETE_Todo(w http.ResponseWriter, r *http.Request) {
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
