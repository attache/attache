package main

import (
	"database/sql"
	"fmt"

	"github.com/mccolljr/attache"
	"github.com/mccolljr/attache/examples/todo/models"
)

func (c *Todo) GET_TestTodoNew() {
	attache.RenderHTML(c, "test.todo.create")

}

func (c *Todo) GET_TestTodoList() {
	all, err := c.DB().All(func() attache.Storable { return new(models.Todo) })
	if err != nil && err != sql.ErrNoRows {
		attache.ErrorFatal(err)
	}
	c.SetViewData(all)
	attache.RenderHTML(c, "test.todo.list")

}

func (c *Todo) GET_TestTodo() {
	r := c.Request()
	id := r.FormValue("id")
	var target models.Todo
	if err := c.DB().Find(&target, id); err != nil {
		if err == sql.ErrNoRows {
			attache.Error(404)
		}
		attache.ErrorFatal(err)
	}
	c.SetViewData(target)
	attache.RenderHTML(c, "test.todo.list")

}

func (c *Todo) POST_testTodoNew() {
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
	attache.RedirectPage(fmt.Sprintf("/test/todo?id=%v", target.ID))
}

func (c *Todo) POST_TestTodo() {
	r := c.Request()
	id := r.FormValue("id")
	var target models.Todo
	if err := c.DB().Find(&target, id); err != nil {
		if err == sql.ErrNoRows {
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
	attache.RedirectPage(fmt.Sprintf("/test/todo?id=%v", target.ID))
}

func (c *Todo) DELETE_TestTodo() {
	w := c.ResponseWriter()
	r := c.Request()
	id := r.FormValue("id")
	var target models.Todo
	if err := c.DB().Find(&target, id); err != nil {
		if err == sql.ErrNoRows {
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
