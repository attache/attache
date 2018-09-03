[[- $scopePath := .ScopePath -]]
[[- with .Model -]]
{{define "title"}}[[.Name]] List{{end}}
{{define "body"}}
{{with .ViewData}}
<div class="container">
	<div class="card">
		<div class="card-body">
			<h3>[[.Name]] List</h3>
			<table class="table table-bordered">
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
							[[- if .Key -]]<a href="[[$scopePath]]/[[$table]]?id={{.[[.StructField]]}}">[[- end -]]
							{{.[[.StructField]]}}
							[[- if .Key -]]</a>[[- end -]]
						</td>
						[[- end -]]
					[[- end]]
					</tr>
				{{end}}
				</tbody>
			</table>
		</div>
	</div>
</div>
{{end}}
{{end}}
[[end]]