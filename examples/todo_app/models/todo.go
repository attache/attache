package models

import (
	"database/sql"
)

type Todo struct {
	ID    int64
	Title string
	Info  string
}

func (m *Todo) Table() string { return "todo" }

func (m *Todo) Insert() (columns []string, values []interface{}) {
	columns = []string{"title", "info"}
	values = []interface{}{m.Title, m.Info}
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
	columns = []string{"title", "info"}
	values = []interface{}{m.Title, m.Info}
	return
}

func (m *Todo) Select() (columns []string, into []interface{}) {
	columns = []string{"id", "title", "info"}
	into = []interface{}{&m.ID, &m.Title, &m.Info}
	return
}

func (m *Todo) KeyColumns() []string     { return []string{"id"} }
func (m *Todo) KeyValues() []interface{} { return []interface{}{m.ID} }
