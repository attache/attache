package models

type Todo struct {
	Title   string
	Details string
}

func (m *Todo) Table() string { return "todo" }

func (m *Todo) Insert() (columns []string, values []interface{}) {
	columns = []string{"title", "details"}
	values = []interface{}{m.Title, m.Details}
	return
}

func (m *Todo) Update() (columns []string, values []interface{}) {
	columns = []string{"details"}
	values = []interface{}{m.Details}
	return
}

func (m *Todo) Select() (columns []string, into []interface{}) {
	columns = []string{"title", "details"}
	into = []interface{}{&m.Title, &m.Details}
	return
}

func (m *Todo) KeyColumns() []string     { return []string{"title"} }
func (m *Todo) KeyValues() []interface{} { return []interface{}{m.Title} }
