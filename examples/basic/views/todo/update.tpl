
{{define "title"}}Edit Todo{{end}}
{{define "body"}}
	<h1>Edit Todo</h1>
	<form name="edit_todo" method="post" action="/todo?id={{.ID}}">
	
		<div>
			<label for="Title">Title</label>
			<input type="text" name="Title" value="{{.Title}}" />
		</div>
		
		<div>
			<label for="Body">Body</label>
			<input type="text" name="Body" value="{{.Body}}" />
		</div>
		
		<input type="submit" value="Update"/>
	</form>
{{end}}
