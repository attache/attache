
{{define "title"}}Todo List{{end}}
{{define "body"}}
<h1>Todo List</h1>
<table>
	<thead>
		<tr>
		
			<th>ID</th>
			
			<th>Title</th>
			
			<th>Body</th>
			
		</tr>
	</thead>
	<tbody>
	{{range .}}
		<tr>
		
		
			<td><a href="/todo?id={{.ID}}">{{.ID}}</a></td>
			
			<td>{{.Title}}</td>
			
			<td>{{.Body}}</td>
			
		<tr>
	{{end}}
	</tbody>
<table>
{{end}}
