
{{define "title"}}New Todo{{end}}
{{define "body"}}
	<h1>New Todo</h1>
	<form name="new_todo" method="post" action="/todo/new">
	
		<div>
			<label for="Title">Title</label>
			<input type="text" name="Title" />
		</div>
		
		<div>
			<label for="Body">Body</label>
			<input type="text" name="Body" />
		</div>
		
		<input type="submit" value="Create"/>
	</form>
{{end}}
