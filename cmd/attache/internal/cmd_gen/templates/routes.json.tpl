package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"database/sql"
	"encoding/json"

	"github.com/mccolljr/attache"
)

func (c *{{.ContextType}}) GET_{{.Model.Name}}List(w http.ResponseWriter, r *http.Request) {
	all, err := c.DB().All(func() attache.Storable{ return new(models.{{.Model.Name}}) })
	if err != nil && err != sql.ErrNoRows {
		attache.ErrorFatal(err)
	}

	attache.RenderJSON(w, all)
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

	attache.RenderJSON(w, target)
}

func (c *{{.ContextType}}) POST_{{.Model.Name}}New(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		attache.ErrorFatal(err)
	}

	target := new(models.{{.Model.Name}})
	
	if err := json.Unmarshal(body, target); err != nil {
		attache.ErrorFatal(err)
	}

	if err := c.DB().Insert(target); err != nil {
		attache.ErrorFatal(err)
	}

	w.WriteHeader(200)
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

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		attache.ErrorFatal(err)
	}
		
	if err := json.Unmarshal(body, target); err != nil {
		attache.ErrorFatal(err)
	}	

	if err := c.DB().Update(target); err != nil {
		attache.ErrorFatal(err)
	}

	w.WriteHeader(200)
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