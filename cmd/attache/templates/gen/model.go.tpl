package models

import (
	"database/sql"
	"github.com/attache/attache"
)

type {{.Name}} struct { {{range .Model.Fields}}{{.StructField}} {{.Type}} `db:"{{.Column}}"`;{{end}} }

func New{{.Name}}() attache.Record { return new({{.Name}}) }

func (m *{{.Name}}) Table() string { return {{printf "%q" .Table}} }

func (m *{{.Name}}) Key() (columns []string, values []interface{}) {
	columns = []string{
		{{- range .Model.Fields -}}
			{{if .Key}}{{ printf "%q" .Column }},{{end}}
		{{- end -}}
	}
	values = []interface{}{
		{{- range .Model.Fields -}}
			{{if .Key}}m.{{ .StructField }},{{end}}
		{{- end -}}
	}
	return
}

func (m *{{.Name}}) Insert() (columns []string, values []interface{}) {
	columns = []string{
		{{- range .Model.Fields -}}
		{{ if not .NoInsert }}{{ printf "%q" .Column }},{{end}}
		{{- end -}}
	}
	values = []interface{}{
		{{- range .Model.Fields -}}
		{{ if not .NoInsert }}m.{{.StructField}},{{end}}
		{{- end -}}
	}
	return
}

{{if .Model.DefaultKey -}}
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
		{{- range .Model.Fields -}}
		{{ if not .NoUpdate }}{{ printf "%q" .Column }},{{end}}
		{{- end -}}
	}
	values = []interface{}{
		{{- range .Model.Fields -}}
		{{ if not .NoUpdate }}m.{{.StructField}},{{end}}
		{{- end -}}
	}
	return
}