
{{define "title"}}Edit User{{end}}
{{define "body"}}
	<h1>Edit User</h1>
	<form name="edit_user" method="post" action="/user/{{.ID}}">
	
		<div>
			<label for="ID">ID</label>
			<input type="text" name="ID" value="{{.ID}}" readonly="true"/>
		</div>
	
		<div>
			<label for="Username">Username</label>
			<input type="text" name="Username" value="{{.Username}}" />
		</div>
	
		<div>
			<label for="Firstname">Firstname</label>
			<input type="text" name="Firstname" value="{{.Firstname}}" />
		</div>
	
		<div>
			<label for="Lastname">Lastname</label>
			<input type="text" name="Lastname" value="{{.Lastname}}" />
		</div>
	
		<div>
			<label for="Password">Password</label>
			<input type="text" name="Password" value="{{.Password}}" />
		</div>
	
		<div>
			<label for="Created">Created</label>
			<input type="text" name="Created" value="{{.Created}}" />
		</div>
	
		<input type="submit" value="Update"/>
	</form>
{{end}}
