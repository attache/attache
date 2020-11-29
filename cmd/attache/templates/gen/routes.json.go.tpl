package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"database/sql"
	"encoding/json"

	"github.com/attache/attache"
)

func (c *{{.ContextType}}) GET_{{.ScopeCamel}}{{.Name}}List() {
	w := c.ResponseWriter()
	all, err := c.DB().All(new(models.{{.Name}}))
	if err != nil {
		attache.ErrorFatal(err)
	}
	attache.RenderJSON(w, all)
}

func (c *{{.ContextType}}) GET_{{.ScopeCamel}}{{.Name}}() {
	w := c.ResponseWriter()
	r := c.Request()
	id := r.FormValue("id")
	var target models.{{.Name}}
	if err := c.DB().Get(&target, id); err != nil {
		if err == attache.ErrRecordNotFound {
			attache.Error(404)
		}
		attache.ErrorFatal(err)
	}
	attache.RenderJSON(w, target)
}

func (c *{{.ContextType}}) POST_{{.ScopeCamel}}{{.Name}}New() {
	w := c.ResponseWriter()
	r := c.Request()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		attache.ErrorFatal(err)
	}
	var target models.{{.Name}}
	if err := json.Unmarshal(body, &target); err != nil {
		attache.ErrorFatal(err)
	}
	if err := c.DB().Insert(&target); err != nil {
		attache.ErrorFatal(err)
	}
	w.WriteHeader(200)
}

func (c *{{.ContextType}}) POST_{{.ScopeCamel}}{{.Name}}() {
	w := c.ResponseWriter()
	r := c.Request()
	id := r.FormValue("id")
	var target models.{{.Name}}
	if err := c.DB().Get(&target, id); err != nil {
		if err == attache.ErrRecordNotFound {
			attache.Error(404)
		}
		attache.ErrorFatal(err)
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		attache.ErrorFatal(err)
	}
	if err := json.Unmarshal(body, &target); err != nil {
		attache.ErrorFatal(err)
	}	
	if err := c.DB().Update(&target); err != nil {
		attache.ErrorFatal(err)
	}
	w.WriteHeader(200)
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