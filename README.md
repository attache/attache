![Attache](https://user-images.githubusercontent.com/8538024/46573302-529f9180-c961-11e8-97a0-051bcc399e13.png)

[![godoc](https://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/attache/attache)
[![Go Report Card](https://goreportcard.com/badge/github.com/attache/attache)](https://goreportcard.com/report/github.com/attache/attache)
[![CircleCI](https://circleci.com/gh/attache/attache.svg?style=svg)](https://circleci.com/gh/attache/attache)

_Since this project is still a work in progress, please make sure you test it out before deciding to use it._

# Preface

A couple of years ago, I got an idea. At the time, it was a fuzzy concept for "some kind of web framework" for Go. I ignored this idea for the longest time, because I felt Go was not conducive to that kind of project. That's not a slight at Go, I just felt that the language didn't really _need_ it. As time went on I started building code that I re-used between several personal (read: half-finished, never-published) projects. I started to coalesce those items into a package. Over time, that package became unruly. I started refactoring. I started standardizing. I started toying with weird ideas and making bad decisions, removing them, and trying something new.

Eventually, I ended up with Attache, a really lame play on words:

1. __Attachè__, _noun_: person on the staff of an ambassador, _typically with a specialized area of responsibility_
1. __Attache__, mispronounced like "Apache" ('cuz it's a web server)).

I've been using this "framework" (which I think of as more of a collection of tools) in some freelance work recently and I really enjoy it. It makes some things much easier by allowing me to avoid some boilerplate each time I start a web app that uses Go, but without forcing me into a particular application structure.

This is my longest running personal project. For better or for worse, I think it could use some new ideas and new eyes.

Documentation is a little sparse (I'm working on it, but it takes time and I have a full time job outside of this :( ).

I'm a little nervous. All of this is my own code right now. It feels weird to put so much out there for the community to see and judge. I like this, though, and I hope maybe you all will as well.

# Installation

```bash
$> go get -u github.com/attache/attache/...
```

# Usage

### Getting Started

The CLI is the canonical way to start a new Attache application. An example of a working application can be found in `examples/todo`

##### Create a new Attache application
```bash
$> attache new -n MyApp
$> cd ./my_app
```

Inside the `my_app` folder, you will find the basic structure of the application:

```
/my_app/            - root
|  /models/         - directory where Attache generates models
|  /views/          - directory where Attache generates views
|  |  index.tpl
|  |  layout.tpl
|  /web/            - directory where frontend files live
|  |  /dist/        - default public file server root 
|  |  |  /css/
|  |  |  /img/
|  |  |  /js/
|  |  /src/         - where TypeScript, Less, Sass, etc. should be written before compiling
|  |  |  /script/
|  |  |  /styles/
|  /secret/         - files that should be ignored by git
|  |  run.sh
|  |  schema.sql
|  main.go          - basic application structure & startup
|  .at-conf.json    - Attache CLI configuration - don't touch this
|  .gitignore
```

Run
```bash
# from within the application's root directory
$> ./secret/run.sh
```
Then, visit `http://localhost:8080/`, you should see a page with
> Welcome to your Attache application!

### Database Connection

To get Attache to manage a database connection for you, embed the `attache.DefaultDB` type into your content type. You can do so by uncommenting the auto-generated line. This will inspect 2 environment variables, `DB_DRIVER` and `DB_DSN`, and attempt to use these to establisha a database connection. This connection will then be available to your Context's methods via the `.DB()` method of the embedded `attache.DefaultDB` type _when a context is initialized by Attache to handle a request_

```go
type MyApp struct {
        // required
        attache.BaseContext

        // capabilities
        attache.DefaultFileServer
        attache.DefaultViews
        attache.DefaultDB // <--- UNCOMMENT THIS LINE
        // attache.DefaultSession // enable session storage
}
```

Once you have enabled database connectivity, you can begin to generate (or manually build)models to represent your schema.

Let's say you have `todos` table in your database.

```sql
CREATE TABLE `todos` (
	`title` TEXT PRIMARY KEY NOT NULL,
	`desc`  TEXT NOT NULL DEFAULT ""
);
```

You can use the CLI to generate a model, views, and routes.

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

##### Generate JSON routes (and no views)
If you want to generate the routes to serve JSON data rather than views, you can include the
`-json` flag in your command. This will generate routes for a JSON-based API, and prevent generation of views.
```bash
# from within the application's root directory
$> attache gen -json [...]
```

##### Use CLI plugin
The CLI is extensible. It can take advantage of dynamically linked plugins, which can be used like so:
```bash
# attempts to load plugin from $HOME/.attache/plugins/PLUGNAME.so
$> attache PLUGNAME [...]
```

# Contributing

#### If you'd like to help, here are the high-level TODOs:
1. Documentation Wiki
	- document the method naming conventions
	- document some example applications
	- etc...
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