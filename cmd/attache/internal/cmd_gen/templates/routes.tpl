package main

import (
	"fmt"
	"net/http"
	"database/sql"

	"github.com/mccolljr/attache"
)

func (c *{{.ContextType}}) GET_{{.Model.Name}}New(w http.ResponseWriter, r *http.Request) {
	attache.RenderHTML(c, "{{.Model.Table}}.create", w, nil)
}

func (c *{{.ContextType}}) GET_{{.Model.Name}}List(w http.ResponseWriter, r *http.Request) {
	all, err := c.DB().All(func() attache.Storable{ return new(models.{{.Model.Name}}) })
	if err != nil && err != sql.ErrNoRows {
		attache.ErrorFatal(err)
	}

	attache.RenderHTML(c, "{{.Model.Table}}.list", w, all)
}

func (c *{{.ContextType}}) GET_{{.Model.Name}}(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.{{.Model.Name}})
	if err := c.DB().Find(target, id); err != nil {
		if err == sql.ErrNoRows {
			attache.Error(404)
		}
		
		attache.ErrorFatal(err)
	}


	attache.RenderHTML(c, "{{.Model.Table}}.update", w, &target)
}

func (c *{{.ContextType}}) POST_{{.Model.Name}}New(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		attache.ErrorFatal(err)
	}

	target := new(models.{{.Model.Name}})
	
	if err := attache.FormDecode(target, r.Form); err != nil {
		attache.ErrorFatal(err)
	}

	if err := c.DB().Insert(target); err != nil {
		attache.ErrorFatal(err)
	}

	attache.RedirectPage(fmt.Sprintf("/{{.Model.Table}}?id=%v", target.{{.Model.KeyStructField}}))
}

func (c *{{.ContextType}}) POST_{{.Model.Name}}(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.{{.Model.Name}})
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

	attache.RedirectPage(fmt.Sprintf("/{{.Model.Table}}?id=%v", target.{{.Model.KeyStructField}}))
}

func (c *{{.ContextType}}) DELETE_{{.Model.Name}}(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.{{.Model.Name}})
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