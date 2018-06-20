package main

import (
	"fmt"
	"net/http"
	"database/sql"

	"github.com/go-chi/chi"
	"github.com/mccolljr/attache"
)

func (c *Ctx) GET_{{.Name}}New(r *http.Request) ([]byte, error) {
	return c.Views().Render("{{.Table}}.create", nil)
}

func (c *Ctx) GET_{{.Name}}List(r *http.Request) ([]byte, error) {
	all, err := c.DB().All(func() attache.Storable{ return new(models.{{.Name}}) })
	if err != nil && err != sql.ErrNoRows {
		attache.ErrorFatal(err)
	}

	return c.Views().Render("{{.Table}}.list", all)
}

func (c *Ctx) GET_{{.Name}}(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.{{.Name}})
	if err := c.DB().Find(target, id); err != nil {
		if err == sql.ErrNoRows {
			attache.Error(404)
		}
		
		attache.ErrorFatal(err)
	}

	data, err := c.Views().Render("{{.Table}}.update", &target)
	if err != nil {
		attache.ErrorFatal(err)
	}

	w.Header().Set("content-type", "text/html")
	w.Write(data)
}

func (c *Ctx) POST_{{.Name}}New(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		attache.ErrorFatal(err)
	}

	target := new(models.{{.Name}})
	
	if err := attache.FormDecode(target, r.Form); err != nil {
		attache.ErrorFatal(err)
	}

	if err := c.DB().Insert(target); err != nil {
		attache.ErrorFatal(err)
	}

	attache.RedirectPage(fmt.Sprintf("/{{.Table}}?id=%v", target.{{.KeyStructField}}))
}

func (c *Ctx) POST_{{.Name}}(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.{{.Name}})
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

	attache.RedirectPage(fmt.Sprintf("/{{.Table}}?id=%v", target.{{.KeyStructField}}))
}

func (c *Ctx) DELETE_{{.Name}}(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.{{.Name}})
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