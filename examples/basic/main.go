package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mccolljr/attache"
	"github.com/mccolljr/attache/examples/basic/models"

	_ "github.com/mattn/go-sqlite3"
)

type Ctx struct {
	Views attache.ViewCache
	DB    *attache.DB
}

func (c *Ctx) Init(w http.ResponseWriter, r *http.Request) { /* nothing yet */ }

// HasViews
func (c *Ctx) ViewRoot() string             { return "views" }
func (c *Ctx) SetViews(v attache.ViewCache) { c.Views = v }

// HasDB
func (c *Ctx) DBDriver() string     { return "sqlite3" }
func (c *Ctx) DBString() string     { return "test.db" }
func (c *Ctx) SetDB(db *attache.DB) { c.DB = db }

func (c *Ctx) GET_TodoList(w http.ResponseWriter, r *http.Request) {
	all := make([]*models.Todo, 0, 32)

	rows, err := c.DB.Query("SELECT id, title, body FROM todo")
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
		} else {
			log.Println(err)
			w.WriteHeader(500)
		}
		return
	}

	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Body); err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		all = append(all, &todo)
	}

	data, _ := json.Marshal(all)
	w.Header().Set("content-type", "application/json")
	w.Write(data)
}

func (c *Ctx) GET_Todo(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(404)
		return
	}

	var todo models.Todo
	if err := c.DB.Find(&todo, id); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
		} else {
			log.Println(err)
			w.WriteHeader(500)
		}
		return
	}

	data, _ := json.Marshal(todo)
	w.Header().Set("content-type", "application/json")
	w.Write(data)
}

func (c *Ctx) POST_Todo(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		c.createTodo(w, r)
	} else {
		c.updateTodo(id, w, r)
	}
}

func (c *Ctx) createTodo(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	var todo models.Todo
	if err := json.Unmarshal(payload, &todo); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if err := c.DB.Insert(&todo); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	data, _ := json.Marshal(todo)
	w.Header().Set("content-type", "application/json")
	w.Write(data)
}

func (c *Ctx) updateTodo(id string, w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	if err := c.DB.Find(&todo, id); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(404)
		} else {
			log.Println(err)
			w.WriteHeader(500)
		}
		return
	}

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if err := json.Unmarshal(payload, &todo); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	if err := c.DB.Update(&todo); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	data, _ := json.Marshal(todo)
	w.Header().Set("content-type", "application/json")
	w.Write(data)
}

func (c *Ctx) DELETE_Todo(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(404)
		return
	}

	var todo models.Todo
	if err := c.DB.Find(&todo, id); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(200) // consider it a success
		} else {
			log.Println(err)
			w.WriteHeader(500)
		}
		return
	}

	if err := c.DB.Delete(&todo); err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	data, _ := json.Marshal(todo)
	w.Header().Set("content-type", "application/json")
	w.Write(data)
}

func (c *Ctx) GET_(w http.ResponseWriter, r *http.Request) { c.GET_Index(w, r) }

func (c *Ctx) GET_Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	c.Views.Get("index").Execute(w, nil)
}

func main() {
	app, err := attache.Bootstrap(&Ctx{})
	if err != nil {
		log.Fatalln(err)
	}

	log.Fatalln(app.Run())
}
