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
	<h1>New [[.Name]]</h1>
	<form name="new_[[.Table]]" method="post" action="/[[.Table]]/new">
	[[range .Fields]]
		[[- if not .NoInsert]]
		<div>
			<label for="[[.StructField]]">[[.StructField]]</label>
			<input type="text" name="[[.StructField]]" />
		</div>
		[[end -]]
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
		[[- if not .NoUpdate]]
		<div>
			<label for="[[.StructField]]">[[.StructField]]</label>
			<input type="text" name="[[.StructField]]" value="{{.[[.StructField]]}}" [[if .Key]]readonly="true"[[end]]/>
		</div>
		[[end -]]
	[[end]]
		<input type="submit" value="Update"/>
	</form>
{{end}}
`

var routeTemplate = `
package main

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/mccolljr/attache/utils"
)

func init() { Register{{.Name}}(App) }

func Register{{.Name}}(r chi.Router) {
	// new form
	r.Get("/{{.Table}}/new", utils.RenderHTML("{{.Table}}.create"))

	// edit form
	r.Get("/{{.Table}}/{id}", utils.RenderFunc(
		"text/html",
		"{{.Table}}.update",
		func (r *http.Request) interface{} {
			id := chi.URLParam(r, "id")
			var target models.{{.Name}}
			if err := DB.Find(&target, id); err != nil {
				if err == sql.ErrNoRows {
					utils.Error(404)
				} else {
					log.Println(err)
					utils.Error(500)
				}
			}
			
			return &target
		},
	))

	// create handler
	r.Post("/{{.Table}}/new", func(w http.ResponseWriter, r *http.Request){
		if err := r.ParseForm(); err != nil {
			log.Println(err)
			utils.Error(500)
		}

		var target models.{{.Name}}
		
		if err := FormDecoder.Decode(&target, r.Form); err != nil {
			log.Println(err)
			utils.Error(500)
		}

		if err := DB.Insert(&target); err != nil {
			log.Println(err)
			utils.Error(500)
		}

		w.WriteHeader(200)
		fmt.Fprintf(w, "%v", target.{{.KeyStructField}})
	})

	// update handler
	r.Post("/{{.Table}}/{id}", func(w http.ResponseWriter, r *http.Request){
		if err := r.ParseForm(); err != nil {
			log.Println(err)
			utils.Error(500)
		}

		var (
			id = chi.URLParam(r, "id")
			target models.{{.Name}}
		)

		if err := DB.Find(&target, id); err != nil {
			if err == sql.ErrNoRows {
				utils.Error(404)
			} else {
				log.Println(err)
				utils.Error(500)
			}
		}
		
		if err := FormDecoder.Decode(&target, r.Form); err != nil {
			log.Println(err)
			utils.Error(500)
		}

		if err := DB.Update(&target); err != nil {
			log.Println(err)
			utils.Error(500)
		}

		w.WriteHeader(200)
	})

	// delete handler
	r.Delete("/{{.Table}}/{id}", func(w http.ResponseWriter, r *http.Request){
		var (
			id = chi.URLParam(r, "id")
			target models.{{.Name}}
		)

		if err := DB.Find(&target, id); err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(200)
				return
			} else {
				log.Println(err)
				utils.Error(500)
			}
		}

		if err := DB.Delete(&target); err != nil {
			log.Println(err)
			utils.Error(500)
		}

		w.WriteHeader(200)
	})
}
`
