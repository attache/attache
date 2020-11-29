[[- $scopePath := .ScopePath -]]
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
					[[- range .Model.Fields]]
						[[- if not .NoSelect]]
						<th>[[.StructField]]</th>
						[[- end -]]
					[[- end]]
					</tr>
				</thead>
				<tbody>
				{{range .}}
					<tr>
					[[- $nameSnake := .NameSnake]]
					[[- range .Model.Fields]]
						[[- if not .NoSelect]]
						<td>
							[[- if .Key -]]<a href="[[$scopePath]]/[[$nameSnake]]?id={{.[[.StructField]]}}">[[- end -]]
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