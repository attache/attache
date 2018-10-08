package models

import (
	"database/sql"

	"github.com/attache/attache"
)

type Todo struct {
	ID          int64  `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
}

func NewTodo() attache.Record { return new(Todo) }

func (m *Todo) Table() string { return "todo" }

func (m *Todo) Key() (columns []string, values []interface{}) {
	columns = []string{"id"}
	values = []interface{}{m.ID}
	return
}

func (m *Todo) Insert() (columns []string, values []interface{}) {
	columns = []string{"title", "description"}
	values = []interface{}{m.Title, m.Description}
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
	columns = []string{"title", "description"}
	values = []interface{}{m.Title, m.Description}
	return
}
