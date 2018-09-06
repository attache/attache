[[- $scopePath := .ScopePath -]]
[[- with .Model -]]
{{define "title"}}Edit [[.Name]]{{end}}
{{define "body"}}
{{with .ViewData}}
<div class="container">
	<div class="card">
		<div class="card-body">
			<form name="edit_[[.Table]]" method="post" action="[[$scopePath]]/[[.Table]]?id={{.[[.KeyStructField]]}}">
				<fieldset>
					<legend>Edit [[.Name]]</legend>
					[[range .Fields]]
					<div class="form-group">
						<label for="[[.StructField]]">[[.StructField]]</label>
						<input type="text" name="[[.StructField]]"
								value="{{.[[.StructField]]}}" 
								[[if or .NoUpdate .Key]]readonly="true"[[end]]
								class="form-control"/>
					</div>
					[[end]]
					<a href="[[$scopePath]]/[[.Table]]/list" class="btn btn-default">Back</a>
					<input type="submit" value="Update" class="btn btn-primary"/>
				</fieldset>
			</form>
		</div>
	</div>
</div>
{{end}}
{{end}}
[[end]]