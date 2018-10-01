package main

import (
	"fmt"
	"net/http"
	"database/sql"

	"github.com/mccolljr/attache"
)

func (c *{{.ContextType}}) GET_{{.ScopeCamel}}{{.Model.Name}}New() {
	{{ if .ScopeSnake -}}
	attache.RenderHTML(c, "{{.ScopeSnake}}.{{.Model.Table}}.create")
	{{ else -}}
	attache.RenderHTML(c, "{{.Model.Table}}.create")
	{{ end }}
}

func (c *{{.ContextType}}) GET_{{.ScopeCamel}}{{.Model.Name}}List() {
	all, err := c.DB().All(new(models.{{.Model.Name}}))
	if err != nil {
		attache.ErrorFatal(err)
	}
	c.SetViewData(all)
	{{ if .ScopeSnake -}}
	attache.RenderHTML(c, "{{.ScopeSnake}}.{{.Model.Table}}.list")
	{{ else -}}
	attache.RenderHTML(c, "{{.Model.Table}}.list")
	{{ end }}
}

func (c *{{.ContextType}}) GET_{{.ScopeCamel}}{{.Model.Name}}() {
	r := c.Request()
	id := r.FormValue("id")
	var target models.{{.Model.Name}}
	if err := c.DB().Get(&target, id); err != nil {
		if err == attache.ErrRecordNotFound {
			attache.Error(404)
		}
		attache.ErrorFatal(err)
	}
	c.SetViewData(target)
	{{ if .ScopeSnake -}}
	attache.RenderHTML(c, "{{.ScopeSnake}}.{{.Model.Table}}.update")
	{{ else -}}
	attache.RenderHTML(c, "{{.Model.Table}}.update")
	{{ end }}
}

func (c *{{.ContextType}}) POST_{{.ScopeCamel}}{{.Model.Name}}New() {
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
	attache.RedirectPage(fmt.Sprintf("{{.ScopePath}}/{{.Model.Table}}?id=%v", target.{{.Model.KeyStructField}}))
}

func (c *{{.ContextType}}) POST_{{.ScopeCamel}}{{.Model.Name}}() {
	r := c.Request()
	id := r.FormValue("id")
	var target models.{{.Model.Name}}
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
	attache.RedirectPage(fmt.Sprintf("{{.ScopePath}}/{{.Model.Table}}?id=%v", target.{{.Model.KeyStructField}}))
}

func (c *{{.ContextType}}) DELETE_{{.ScopeCamel}}{{.Model.Name}}() {
	w := c.ResponseWriter()
	r := c.Request()
	id := r.FormValue("id")
	var target models.{{.Model.Name}}
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