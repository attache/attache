{{define "title"}}Edit Todo{{end}}
{{define "body"}}
	<form name="edit_todo" method="post" action="/todo?id={{.ID}}" class="pure-form pure-form-stacked">
		<fieldset>
			<legend>Edit Todo</legend>
			
			<label>ID</label>
			<input type="text" value="{{.ID}}" readonly="true"/>
			<label for="Title">Title</label>
			<input type="text" name="Title" value="{{.Title}}" />
			<label for="Info">Info</label>
			<input type="text" name="Info" value="{{.Info}}" />
			<input type="submit" value="Update" class="pure-button pure-button-primary"/>
		</fieldset>
	</form>
{{end}}