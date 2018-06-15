package main

import (
	"html/template"
	"path/filepath"
	"strings"
)

type View struct {
	File string
	Body string
}

func viewsFor(m *Model) []View {
	createForm := template.Must(template.New("").Delims("[[", "]]").Parse(createFormTemplate))
	updateForm := template.Must(template.New("").Delims("[[", "]]").Parse(updateFormTemplate))
	listView := template.Must(template.New("").Delims("[[", "]]").Parse(listViewTemplate))

	var (
		create = &strings.Builder{}
		update = &strings.Builder{}
		list   = &strings.Builder{}
	)

	if err := createForm.Execute(create, m); err != nil {
		panic(err)
	}

	if err := updateForm.Execute(update, m); err != nil {
		panic(err)
	}

	if err := listView.Execute(list, m); err != nil {
		panic(err)
	}

	return []View{
		View{
			File: filepath.Join("views", m.Table, "create.tpl"),
			Body: create.String(),
		},

		View{
			File: filepath.Join("views", m.Table, "update.tpl"),
			Body: update.String(),
		},

		View{
			File: filepath.Join("views", m.Table, "list.tpl"),
			Body: list.String(),
		},
	}
}
