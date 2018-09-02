package main

import (
	"fmt"
	"net/http"
	"database/sql"

	"github.com/mccolljr/attache"
)

func (c *{{.ContextType}}) GET_{{.Model.Name}}New() {
	attache.RenderHTML(c, "{{.Model.Table}}.create")
}

func (c *{{.ContextType}}) GET_{{.Model.Name}}List() {
	all, err := c.DB().All(func() attache.Storable{ return new(models.{{.Model.Name}}) })
	if err != nil && err != sql.ErrNoRows {
		attache.ErrorFatal(err)
	}
	c.SetViewData(all)
	attache.RenderHTML(c, "{{.Model.Table}}.list")
}

func (c *{{.ContextType}}) GET_{{.Model.Name}}() {
	r := c.Request()
	id := r.FormValue("id")
	var target models.{{.Model.Name}}
	if err := c.DB().Find(&target, id); err != nil {
		if err == sql.ErrNoRows {
			attache.Error(404)
		}
		attache.ErrorFatal(err)
	}
	c.SetViewData(target)
	attache.RenderHTML(c, "{{.Model.Table}}.update")
}

func (c *{{.ContextType}}) POST_{{.Model.Name}}New() {
	r := c.Request()
	if err := r.ParseForm(); err != nil {
		attache.ErrorFatal(err)
	}
	var target models.{{.Model.Name}}
	if err := attache.FormDecode(&target, r.Form); err != nil {
		attache.ErrorFatal(err)
	}
	if err := c.DB().Insert(&target); err != nil {
		attache.ErrorFatal(err)
	}
	attache.RedirectPage(fmt.Sprintf("/{{.Model.Table}}?id=%v", target.{{.Model.KeyStructField}}))
}

func (c *{{.ContextType}}) POST_{{.Model.Name}}() {
	r := c.Request()
	id := r.FormValue("id")
	var target models.{{.Model.Name}}
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
	attache.RedirectPage(fmt.Sprintf("/{{.Model.Table}}?id=%v", target.{{.Model.KeyStructField}}))
}

func (c *{{.ContextType}}) DELETE_{{.Model.Name}}() {
	w := c.ResponseWriter()
	r := c.Request()
	id := r.FormValue("id")
	var target models.{{.Model.Name}}
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