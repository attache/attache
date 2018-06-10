package models

import (
	"database/sql"
)

type User struct {
	ID        int64
	Username  string
	Firstname string
	Lastname  string
	Password  string
	Created   t
}

func (m *User) Table() string { return "user" }

func (m *User) Insert() (columns []string, values []interface{}) {
	columns = []string{"id", "username", "firstname", "lastname", "password", "created"}
	values = []interface{}{m.ID, m.Username, m.Firstname, m.Lastname, m.Password, m.Created}
	return
}

func (m *User) AfterInsert(result sql.Result) {
	id, err := result.LastInsertID()
	if err != nil {
		panic(err)
	}
	m.ID = id
}

func (m *User) Update() (columns []string, values []interface{}) {
	columns = []string{"username", "firstname", "lastname", "password", "created"}
	values = []interface{}{m.Username, m.Firstname, m.Lastname, m.Password, m.Created}
	return
}

func (m *User) Select() (columns []string, into []interface{}) {
	columns = []string{"id", "username", "firstname", "lastname", "password", "created"}
	into = []interface{}{&m.ID, &m.Username, &m.Firstname, &m.Lastname, &m.Password, &m.Created}
	return
}

func (m *User) KeyColumns() []string     { return []string{"id"} }
func (m *User) KeyValues() []interface{} { return []interface{}{m.ID} }
