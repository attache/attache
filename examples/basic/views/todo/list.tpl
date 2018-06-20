{{define "title"}}Todo List{{end}}
{{define "body"}}
<h1>Todo List</h1>
<table class="pure-table pure-table-bordered">
	<thead>
		<tr>
			<th>Title</th>
			<th>Details</th>
		</tr>
	</thead>
	<tbody>
	{{range .}}
		<tr>
			<td><a href="/todo?id={{.Title}}">{{.Title}}</a></td>
			<td>{{.Details}}</td>
		</tr>
	{{end}}
	</tbody>
</table>
{{end}}