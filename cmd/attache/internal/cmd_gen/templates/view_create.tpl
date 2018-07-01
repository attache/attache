{{define "title"}}New [[.Name]]{{end}}
{{define "body"}}
	<form name="new_[[.Table]]" method="post" action="/[[.Table]]/new" class="pure-form pure-form-stacked">
		<fieldset>
			<legend>New [[.Name]]</legend>
			[[range .Fields]]
			[[- if not .NoInsert]]
			<label for="[[.StructField]]">[[.StructField]]</label>
			<input type="text" name="[[.StructField]]" />
			[[- end -]]
			[[end]]
			<input type="submit" value="Create" class="pure-button pure-button-primary"/>
		</fieldset>
	</form>
{{end}}