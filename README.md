# Attache

_Attache is in alpha release, and as such should not yet be considered for production usage_

#### If you'd like to help, here are the high-level TODOs:
1. Documentation of types and API
	- Record API
	- Capabilities API
	- Context Bootstrapping API
	- Naming conventions
		- `func (c *Ctx) GET_TestTest()` serves content for `GET /test/test`
2. Add CLI plugins for common use cases
	- Generate content for features such as Authorization, REST API for a database schema, etc
	- Embed other code generation tools
		- github.com/xo/xo
		- github.com/cheekybits/genny
		- _etc..._
3. Add CLI plugin subcommand for managing installation of plugins
	- maybe have a central repo with checksum verification?
4. Add context capabilities for 
	- Config _(maybe github.com/spf13/viper)_
5. Code cleanup and package re-org

### Installation
```bash
$> go get -u github.com/mccolljr/attache/...
```

### Usage
The CLI is the canonical way to start a new Attache application. An example of a working application can be found in `examples/todo`

### CLI Usage

##### Create a new Attache application
```bash
$> attache new -n MyApp
$> cd ./my_app
```

##### Create a model, view, and routes
```bash
# from within the application's root directory
$> attache gen -n Todo -t todos -f Title:string:key -f Desc:string
# creates   ./models/todos.go
#           ./views/todos/create.tpl
#           ./views/todos/list.tpl
#           ./views/todos/update.tpl
#           ./todos_routes.go
```

##### Create just models, just views, or just routes
```bash
# from within the application's root directory

# just generate the model
$> attache gen -model [...]

# just generate the views
$> attache gen -views [...]

# just generate the routes
$> attache gen -routes [...]

# just generate some combination
$> attache gen -routes -models [...]
```

##### Replace existing files
```bash
# from within the application's root directory
$> attache gen -replace [...]
```

##### Use CLI plugin
```bash
# attempts to load plugin from $HOME/.attache/plugins/PLUGNAME.so
$> attache PLUGNAME [...]
```
