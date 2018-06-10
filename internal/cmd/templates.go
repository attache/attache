package main

const modelTemplate = `
package {{.Package}}

type {{.TypeName}} struct {
{{range .Columns -}}
	{{.Name}} {{.GoType}}
{{- end}}
}

func ()
`
