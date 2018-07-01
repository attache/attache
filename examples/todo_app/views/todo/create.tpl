{{define "title"}}New Todo{{end}}
{{define "body"}}
	<form name="new_todo" method="post" action="/todo/new" class="pure-form pure-form-stacked">
		<fieldset>
			<legend>New Todo</legend>
			
			<label for="Title">Title</label>
			<input type="text" name="Title" />
			<label for="Info">Info</label>
			<input type="text" name="Info" />
			<input type="submit" value="Create" class="pure-button pure-button-primary"/>
		</fieldset>
	</form>
{{end}}