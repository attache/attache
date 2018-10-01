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
						<th>Description</th>
					</tr>
				</thead>
				<tbody>
				{{range .}}
					<tr>
						<td><a href="/todo?id={{.ID}}">{{.ID}}</a></td>
						<td>{{.Title}}</td>
						<td>{{.Description}}</td>
					</tr>
				{{end}}
				</tbody>
			</table>
		</div>
	</div>
</div>
{{end}}
{{end}}
