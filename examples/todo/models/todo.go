package models

import (
	"database/sql"

	"github.com/mccolljr/attache"
)

type Todo struct {
	ID    int64
	Title string
	Text  string
}

func NewTodo() attache.Storable { return new(Todo) }

func (m *Todo) Table() string { return "todo" }

func (m *Todo) Insert() (columns []string, values []interface{}) {
	columns = []string{"title", "text"}
	values = []interface{}{m.Title, m.Text}
	return
}

func (m *Todo) AfterInsert(result sql.Result) {
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	m.ID = id
}

func (m *Todo) Update() (columns []string, values []interface{}) {
	columns = []string{"title", "text"}
	values = []interface{}{m.Title, m.Text}
	return
}

func (m *Todo) Select() (columns []string, into []interface{}) {
	columns = []string{"id", "title", "text"}
	into = []interface{}{&m.ID, &m.Title, &m.Text}
	return
}

func (m *Todo) KeyColumns() []string     { return []string{"id"} }
func (m *Todo) KeyValues() []interface{} { return []interface{}{m.ID} }
