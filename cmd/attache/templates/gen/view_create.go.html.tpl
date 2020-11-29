[[- $scopePath := .ScopePath -]]
{{define "title"}}New [[.Name]]{{end}}
{{define "body"}}
<div class="container">
	<div class="card">
		<div class="card-body">
			<form name="new_[[.NameSnake]]" method="post" action="[[$scopePath]]/[[.NameSnake]]/new">
				<fieldset>
					<legend>New [[.Name]]</legend>
					[[range .Model.Fields]]
					[[- if not .NoInsert]]
					<div class="form-group">
						<label for="[[.StructField]]">[[.StructField]]</label>
						<input type="text" name="[[.StructField]]" class="form-control"/>
					</div>
					[[- end -]]
					[[end]]
					<input type="submit" value="Create" class="btn btn-primary"/>
				</fieldset>
			</form>
		</div>
	</div>
</div>
{{end}}