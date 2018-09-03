{{define "title"}}Edit Todo{{end}}
{{define "body"}}
{{with .ViewData}}
<div class="container">
	<div class="card">
		<div class="card-body">
			<form name="edit_todo" method="post" action="/test/todo?id={{.ID}}">
				<fieldset>
					<legend>Edit Todo</legend>
					
					<div class="form-group">
						<label for="ID">ID</label>
						<input type="text" name="ID"
								value="{{.ID}}" 
								readonly="true"
								class="form-control"/>
					</div>
					
					<div class="form-group">
						<label for="Title">Title</label>
						<input type="text" name="Title"
								value="{{.Title}}" 
								
								class="form-control"/>
					</div>
					
					<div class="form-group">
						<label for="Text">Text</label>
						<input type="text" name="Text"
								value="{{.Text}}" 
								
								class="form-control"/>
					</div>
					
					<input type="submit" value="Update" class="btn btn-primary"/>
				</fieldset>
			</form>
		</div>
	</div>
</div>
{{end}}
{{end}}
