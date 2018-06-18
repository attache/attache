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
`

var createFormTemplate = `
{{define "title"}}New [[.Name]]{{end}}
{{define "body"}}
	<form name="new_[[.Table]]" method="post" action="/[[.Table]]/new" class="pure-form pure-form-stacked">
		<fieldset>
			<legend>New [[.Name]]</legend>
			[[range .Fields]]
			[[- if not .NoInsert]]
			<label for="[[.StructField]]">[[.StructField]]</label>
			<input type="text" name="[[.StructField]]" />
			[[- end -]]
			[[end]]
			<input type="submit" value="Create" class="pure-button pure-button-primary"/>
		</fieldset>
	</form>
{{end}}
`

var updateFormTemplate = `
{{define "title"}}Edit [[.Name]]{{end}}
{{define "body"}}
	<form name="edit_[[.Table]]" method="post" action="/[[.Table]]?id={{.[[.KeyStructField]]}}" class="pure-form pure-form-stacked">
		<fieldset>
			<legend>Edit [[.Name]]</legend>
			[[range .Fields]]
			[[- if not .NoUpdate]]
			<label for="[[.StructField]]">[[.StructField]]</label>
			<input type="text" name="[[.StructField]]" value="{{.[[.StructField]]}}" [[if .Key]]readonly="true"[[end]]/>
			[[- else]]
			<label>[[.StructField]]</label>
			<input type="text" value="{{.[[.StructField]]}}" readonly="true"/>
			[[- end -]]
			[[end]]
			<input type="submit" value="Update" class="pure-button pure-button-primary"/>
		</fieldset>
	</form>
{{end}}
`

var listViewTemplate = `
{{define "title"}}[[.Name]] List{{end}}
{{define "body"}}
<h1>[[.Name]] List</h1>
<table class="pure-table pure-table-bordered">
	<thead>
		<tr>
		[[- range .Fields]]
			[[- if not .NoSelect]]
			<th>[[.StructField]]</th>
			[[- end -]]
		[[- end]]
		</tr>
	</thead>
	<tbody>
	{{range .}}
		<tr>
		[[- $table := .Table]]
		[[- range .Fields]]
			[[- if not .NoSelect]]
			<td>
				[[- if .Key -]]<a href="/[[$table]]?id={{.[[.StructField]]}}">[[- end -]]
				{{.[[.StructField]]}}
				[[- if .Key -]]</a>[[- end -]]
			</td>
			[[- end -]]
		[[- end]]
		</tr>
	{{end}}
	</tbody>
</table>
{{end}}
`

var routeTemplate = `
package main

import (
	"fmt"
	"net/http"
	"database/sql"

	"github.com/go-chi/chi"
	"github.com/mccolljr/attache"
)

func (c *Ctx) GET_{{.Name}}New(r *http.Request) ([]byte, error) {
	return c.Views.Render("{{.Table}}.create", nil)
}

func (c *Ctx) GET_{{.Name}}List(r *http.Request) ([]byte, error) {
	all, err := c.DB.All(func() attache.Storable{ return new(models.{{.Name}}) })
	if err != nil && err != sql.ErrNoRows {
		attache.ErrorFatal(err)
	}

	return c.Views.Render("{{.Table}}.list", all)
}

func (c *Ctx) GET_{{.Name}}(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.{{.Name}})
	if err := c.DB.Find(target, id); err != nil {
		if err == sql.ErrNoRows {
			attache.Error(404)
		}
		
		attache.ErrorFatal(err)
	}

	data, err := c.Views.Render("{{.Table}}.update", &target)
	if err != nil {
		attache.ErrorFatal(err)
	}

	w.Header().Set("content-type", "text/html")
	w.Write(data)
}

func (c *Ctx) POST_{{.Name}}New(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		attache.ErrorFatal(err)
	}

	target := new(models.{{.Name}})
	
	if err := attache.FormDecode(target, r.Form); err != nil {
		attache.ErrorFatal(err)
	}

	if err := c.DB.Insert(target); err != nil {
		attache.ErrorFatal(err)
	}

	attache.RedirectPage(fmt.Sprintf("/{{.Table}}?id=%v", target.{{.KeyStructField}}))
}

func (c *Ctx) POST_{{.Name}}(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.{{.Name}})
	if err := c.DB.Find(target, id); err != nil {
		if err == sql.ErrNoRows {
			attache.Error(404)
		}
		
		attache.ErrorFatal(err)
	}

	if err := attache.FormDecode(target, r.Form); err != nil {
		attache.ErrorFatal(err)
	}

	if err := c.DB.Update(target); err != nil {
		attache.ErrorFatal(err)
	}

	attache.RedirectPage(fmt.Sprintf("/{{.Table}}?id=%v", target.{{.KeyStructField}}))
}

func (c *Ctx) DELETE_{{.Name}}(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	target := new(models.{{.Name}})
	if err := c.DB.Find(target, id); err != nil {
		if err == sql.ErrNoRows {
			attache.Success()
		}
		
		attache.ErrorFatal(err)
	}

	if err := c.DB.Delete(target); err != nil {
		attache.ErrorFatal(err)
	}

	w.WriteHeader(200)
}
`
