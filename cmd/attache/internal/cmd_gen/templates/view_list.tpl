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