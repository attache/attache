{{define "title"}}Todo List{{end}}
{{define "body"}}
<h1>Todo List</h1>
<table class="pure-table pure-table-bordered">
	<thead>
		<tr>
			<th>ID</th>
			<th>Title</th>
			<th>Info</th>
		</tr>
	</thead>
	<tbody>
	{{range .}}
		<tr>
			<td><a href="/todo?id={{.ID}}">{{.ID}}</a></td>
			<td>{{.Title}}</td>
			<td>{{.Info}}</td>
		</tr>
	{{end}}
	</tbody>
</table>
{{end}}