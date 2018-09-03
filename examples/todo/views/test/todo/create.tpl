{{define "title"}}New Todo{{end}}
{{define "body"}}
<div class="container">
	<div class="card">
		<div class="card-body">
			<form name="new_todo" method="post" action="/test/todo/new">
				<fieldset>
					<legend>New Todo</legend>
					
					<div class="form-group">
						<label for="Title">Title</label>
						<input type="text" name="Title" class="form-control"/>
					</div>
					<div class="form-group">
						<label for="Text">Text</label>
						<input type="text" name="Text" class="form-control"/>
					</div>
					<input type="submit" value="Create" class="btn btn-primary"/>
				</fieldset>
			</form>
		</div>
	</div>
</div>
{{end}}
