module {{.Dir}}

go 1.15

require (
	github.com/attache/attache {{.Version}}
)

{{ if .LocalAttache -}}
replace (
	github.com/attache/attache {{.Version}} => {{ .LocalAttache }}
)
{{- end }}