{{define "title"}}Edit Todo{{end}}
{{define "body"}}
	<form name="edit_todo" method="post" action="/todo?id={{.Title}}" class="pure-form pure-form-stacked">
		<fieldset>
			<legend>Edit Todo</legend>
			
			<label>Title</label>
			<input type="text" value="{{.Title}}" readonly="true"/>
			<label for="Details">Details</label>
			<input type="text" name="Details" value="{{.Details}}" />
			<input type="submit" value="Update" class="pure-button pure-button-primary"/>
		</fieldset>
	</form>
{{end}}