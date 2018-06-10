
{{define "title"}}New User{{end}}
{{define "body"}}
	<h1>New User</h1>
	<form name="new_user" method="post" action="/user">
	
		<div>
			<label for="ID">ID</label>
			<input type="text" name="ID" />
		</div>
	
		<div>
			<label for="Username">Username</label>
			<input type="text" name="Username" />
		</div>
	
		<div>
			<label for="Firstname">Firstname</label>
			<input type="text" name="Firstname" />
		</div>
	
		<div>
			<label for="Lastname">Lastname</label>
			<input type="text" name="Lastname" />
		</div>
	
		<div>
			<label for="Password">Password</label>
			<input type="text" name="Password" />
		</div>
	
		<div>
			<label for="Created">Created</label>
			<input type="text" name="Created" />
		</div>
	
		<input type="submit" value="Create"/>
	</form>
{{end}}
