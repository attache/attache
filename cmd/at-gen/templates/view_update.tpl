{{define "title"}}Edit [[.Name]]{{end}}
{{define "body"}}
	<form name="edit_[[.Table]]" method="post" action="/[[.Table]]?id={{.[[.KeyStructField]]}}" class="pure-form pure-form-stacked">
		<fieldset>
			<legend>Edit [[.Name]]</legend>
			[[range .Fields]]
			[[- if not .NoUpdate]]
			<label for="[[.StructField]]">[[.StructField]]</label>
			<input type="text" name="[[.StructField]]" value="{{.[[.StructField]]}}" [[if .Key]]readonly="true"[[end]]/>
			[[- else]]
			<label>[[.StructField]]</label>
			<input type="text" value="{{.[[.StructField]]}}" readonly="true"/>
			[[- end -]]
			[[end]]
			<input type="submit" value="Update" class="pure-button pure-button-primary"/>
		</fieldset>
	</form>
{{end}}