package main

import (
    "net/http"
    "log"
    
    // database drivers
    // (uncomment the lines you need)
    // _ "github.com/mattn/go-sqlite3"    // Sqlite3
    // _ "github.com/go-sql-driver/mysql" // MySQL
    // _ "github.com/jackc/pgx"           // PostgreSQL

    // attache
    "github.com/attache/attache"
)

type {{.Name}} struct {
    // required
    attache.BaseContext 

    // capabilities
    attache.DefaultEnvironment // loads environment variables from  the file pointed to by $ENV_FILE, defaults to secret/dev.env
    attache.DefaultDB          // connects to a database using $DB_DRIVER and $DB_DSN
    {{ if not .API -}}
    attache.DefaultFileServer  // serves static files from the web/dist directory
    attache.DefaultViews       // provides rendering capabilities for views in the views directory
    {{- end }}
    // attache.DefaultSession  // provides access to a secure cookie session
}

func (c *{{.Name}}) Init(w http.ResponseWriter, r *http.Request) {
    /* TODO: initialize context */
}

{{ if not .API -}}
// GET /
func (c *{{.Name}}) GET_() {
    c.GET_Index()
}

// GET /index
func (c *{{.Name}}) GET_Index() {
    attache.RenderHTML(c, "index")
}
{{- else -}}
// GET /health
func (c *{{.Name}}) GET_Health() {
    attache.RenderJSON(
        c.ResponseWriter(),
        map[string]interface{}{
            "time": time.Now().Unix(),
            "ok": true,
        },
    )
}
{{- end }}

func main() {
    // bootstrap application for context type {{.Name}}
    app, err := attache.Bootstrap(&{{.Name}}{})
    if err != nil {
        log.Fatalln(err)
    }

    log.Fatalln(app.Run(":5000"))
}
