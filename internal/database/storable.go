package database

type Storable interface {
	Table() string

	Insert() (columns []string, values []interface{})
	Update() (columns []string, values []interface{})

	Select() (columns []string, into []interface{})

	KeyColumns() []string
	KeyValues() []interface{}
}
