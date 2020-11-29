package main

import (
	"fmt"
	"net/http"
	"database/sql"

	"github.com/attache/attache"
)

func (c *{{.ContextType}}) GET_{{.ScopeCamel}}{{.Name}}New() {
	{{ if .ScopeSnake -}}
	attache.RenderHTML(c, "{{.ScopeSnake}}.{{.NameSnake}}.create")
	{{ else -}}
	attache.RenderHTML(c, "{{.NameSnake}}.create")
	{{ end }}
}

func (c *{{.ContextType}}) GET_{{.ScopeCamel}}{{.Name}}List() {
	all, err := c.DB().All(new(models.{{.Name}}))
	if err != nil {
		attache.ErrorFatal(err)
	}
	c.SetViewData(all)
	{{ if .ScopeSnake -}}
	attache.RenderHTML(c, "{{.ScopeSnake}}.{{.NameSnake}}.list")
	{{ else -}}
	attache.RenderHTML(c, "{{.NameSnake}}.list")
	{{ end }}
}

func (c *{{.ContextType}}) GET_{{.ScopeCamel}}{{.Name}}() {
	r := c.Request()
	id := r.FormValue("id")
	var target models.{{.Name}}
	if err := c.DB().Get(&target, id); err != nil {
		if err == attache.ErrRecordNotFound {
			attache.Error(404)
		}
		attache.ErrorFatal(err)
	}
	c.SetViewData(target)
	{{ if .ScopeSnake -}}
	attache.RenderHTML(c, "{{.ScopeSnake}}.{{.NameSnake}}.update")
	{{ else -}}
	attache.RenderHTML(c, "{{.NameSnake}}.update")
	{{ end }}
}

func (c *{{.ContextType}}) POST_{{.ScopeCamel}}{{.Name}}New() {
	r := c.Request()
	if err := r.ParseForm(); err != nil {
		attache.ErrorFatal(err)
	}
	var target models.{{.Name}}
	if err := attache.FormDecode(&target, r.Form); err != nil {
		attache.ErrorFatal(err)
	}
	if err := c.DB().Insert(&target); err != nil {
		attache.ErrorFatal(err)
	}
	attache.RedirectPage(fmt.Sprintf("{{.ScopePath}}/{{.NameSnake}}?id=%v", target.{{.Model.KeyStructField}}))
}

func (c *{{.ContextType}}) POST_{{.ScopeCamel}}{{.Name}}() {
	r := c.Request()
	id := r.FormValue("id")
	var target models.{{.Name}}
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
	attache.RedirectPage(fmt.Sprintf("{{.ScopePath}}/{{.NameSnake}}?id=%v", target.{{.Model.KeyStructField}}))
}

func (c *{{.ContextType}}) DELETE_{{.ScopeCamel}}{{.Name}}() {
	w := c.ResponseWriter()
	r := c.Request()
	id := r.FormValue("id")
	var target models.{{.Name}}
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