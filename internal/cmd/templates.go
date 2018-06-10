package main

const modelTemplate = `
package models

import (
	"database/sql"
)

type {{.Name}} struct { {{range .Fields}}{{.StructField}} {{.Type}};{{end}} }

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
	id, err := result.LastInsertID()
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
`

var createFormTemplate = `
{{define "title"}}New [[.Name]]{{end}}
{{define "body"}}
	<h1>New [[.Name]]</h1>
	<form name="new_[[.Table]]" method="post" action="/[[.Table]]">
	[[range .Fields]]
		<div>
			<label for="[[.StructField]]">[[.StructField]]</label>
			<input type="text" name="[[.StructField]]" />
		</div>
	[[end]]
		<input type="submit" value="Create"/>
	</form>
{{end}}
`

var updateFormTemplate = `
{{define "title"}}Edit [[.Name]]{{end}}
{{define "body"}}
	<h1>Edit [[.Name]]</h1>
	<form name="edit_[[.Table]]" method="post" action="/[[.Table]]/{{.[[.KeyStructField]]}}">
	[[range .Fields]]
		<div>
			<label for="[[.StructField]]">[[.StructField]]</label>
			<input type="text" name="[[.StructField]]" value="{{.[[.StructField]]}}" [[if .Key]]readonly="true"[[end]]/>
		</div>
	[[end]]
		<input type="submit" value="Update"/>
	</form>
{{end}}
`
