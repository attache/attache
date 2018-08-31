package models

import (
	"database/sql"
	"github.com/mccolljr/attache"
)

type {{.Name}} struct { {{range .Fields}}{{.StructField}} {{.Type}};{{end}} }

func New{{.Name}}() attache.Storeable { return new({{.Name}}) }

func (m *{{.Name}}) Table() string { return {{printf "%q" .Table}} }

func (m *{{.Name}}) Insert() (columns []string, values []interface{}) {
	columns = []string{
		{{- range .Fields -}}
		{{ if not .NoInsert }}{{ printf "%q" .Column }},{{end}}
		{{- end -}}
	}
	values = []interface{}{
		{{- range .Fields -}}
		{{ if not .NoInsert }}m.{{.StructField}},{{end}}
		{{- end -}}
	}
	return
}

{{if .DefaultKey -}}
func (m *{{.Name}}) AfterInsert(result sql.Result) {
	id, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	m.ID = id
}
{{- end}}

func (m *{{.Name}}) Update() (columns []string, values []interface{}) {
	columns = []string{
		{{- range .Fields -}}
		{{ if not .NoUpdate }}{{ printf "%q" .Column }},{{end}}
		{{- end -}}
	}
	values = []interface{}{
		{{- range .Fields -}}
		{{ if not .NoUpdate }}m.{{.StructField}},{{end}}
		{{- end -}}
	}
	return
}

func (m *{{.Name}}) Select() (columns []string, into []interface{}) {
	columns = []string{
		{{- range .Fields -}}
		{{ if not .NoSelect }}{{ printf "%q" .Column }},{{end}}
		{{- end -}}
	}
	into = []interface{}{
		{{- range .Fields -}}
		{{ if not .NoSelect }}&m.{{.StructField}},{{end}}
		{{- end -}}
	}
	return
}

func (m *{{.Name}}) KeyColumns() []string { return []string{
	{{- range .Fields -}}
		{{if .Key}}{{ printf "%q" .Column }},{{end}}
	{{- end -}}
} }
func (m *{{.Name}}) KeyValues() []interface{} { return []interface{}{
	{{- range .Fields -}}
		{{if .Key}}m.{{ .StructField }},{{end}}
	{{- end -}}
} }