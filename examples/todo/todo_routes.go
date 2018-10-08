package main

import (
	"fmt"

	"github.com/attache/attache"
	"github.com/attache/attache/examples/todo/models"
)

func (c *Todo) GET_TodoNew() {
	attache.RenderHTML(c, "todo.create")

}

func (c *Todo) GET_TodoList() {
	all, err := c.DB().All(new(models.Todo))
	if err != nil {
		attache.ErrorFatal(err)
	}
	c.SetViewData(all)
	attache.RenderHTML(c, "todo.list")

}

func (c *Todo) GET_Todo() {
	r := c.Request()
	id := r.FormValue("id")
	var target models.Todo
	if err := c.DB().Get(&target, id); err != nil {
		if err == attache.ErrRecordNotFound {
			attache.Error(404)
		}
		attache.ErrorFatal(err)
	}
	c.SetViewData(target)
	attache.RenderHTML(c, "todo.update")

}

func (c *Todo) POST_TodoNew() {
	r := c.Request()
	if err := r.ParseForm(); err != nil {
		attache.ErrorFatal(err)
	}
	var target models.Todo
	if err := attache.FormDecode(&target, r.Form); err != nil {
		attache.ErrorFatal(err)
	}
	if err := c.DB().Insert(&target); err != nil {
		attache.ErrorFatal(err)
	}
	attache.RedirectPage(fmt.Sprintf("/todo?id=%v", target.ID))
}

func (c *Todo) POST_Todo() {
	r := c.Request()
	id := r.FormValue("id")
	var target models.Todo
	if err := c.DB().Get(&target, id); err != nil {
		if err == attache.ErrRecordNotFound {
			attache.Error(404)
		}
		attache.ErrorFatal(err)
	}
	if err := attache.FormDecode(&target, r.Form); err != nil {
		attache.ErrorFatal(err)
	}
	if err := c.DB().Update(&target); err != nil {
		attache.ErrorFatal(err)
	}
	attache.RedirectPage(fmt.Sprintf("/todo?id=%v", target.ID))
}

func (c *Todo) DELETE_Todo() {
	w := c.ResponseWriter()
	r := c.Request()
	id := r.FormValue("id")
	var target models.Todo
	if err := c.DB().Get(&target, id); err != nil {
		if err == attache.ErrRecordNotFound {
			w.WriteHeader(200)
			return
		}
		attache.ErrorFatal(err)
	}
	if err := c.DB().Delete(&target); err != nil {
		attache.ErrorFatal(err)
	}
	w.WriteHeader(200)
}
