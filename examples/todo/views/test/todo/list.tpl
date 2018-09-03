{{define "title"}}Todo List{{end}}
{{define "body"}}
{{with .ViewData}}
<div class="container">
	<div class="card">
		<div class="card-body">
			<h3>Todo List</h3>
			<table class="table table-bordered">
				<thead>
					<tr>
						<th>ID</th>
						<th>Title</th>
						<th>Text</th>
					</tr>
				</thead>
				<tbody>
				{{range .}}
					<tr>
						<td><a href="/test/todo?id={{.ID}}">{{.ID}}</a></td>
						<td>{{.Title}}</td>
						<td>{{.Text}}</td>
					</tr>
				{{end}}
				</tbody>
			</table>
		</div>
	</div>
</div>
{{end}}
{{end}}
