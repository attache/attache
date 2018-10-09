package attache

import (
	"net/http"
	"os"
)

var (
	_ HasViews       = (*DefaultViews)(nil)
	_ HasSession     = (*DefaultSession)(nil)
	_ HasDB          = (*DefaultDB)(nil)
	_ HasEnvironment = (*DefaultEnvironment)(nil)
)

// DefaultEnvironment is a type that can be embedded into a Context type
// to enable auto-load of environment variables with default
// configuration options
type DefaultEnvironment struct{}

func (*DefaultEnvironment) CONFIG_Environment() EnvironmentConfig {
	return EnvironmentConfig{
		EnvPath: envOrDefault("ENV_FILE", "secret/dev.env"),
	}
}

// DefaultViews is a type that can be embedded into a Context type
// to enable views with default configuration options
type DefaultViews struct {
	views    ViewCache
	viewData interface{}
}

func (d *DefaultViews) CONFIG_Views() ViewConfig {
	return ViewConfig{
		Root: "views",
	}
}

func (d *DefaultViews) Views() ViewCache     { return d.views }
func (d *DefaultViews) setViews(v ViewCache) { d.views = v }

func (d *DefaultViews) ViewData() interface{}        { return d.viewData }
func (d *DefaultViews) SetViewData(data interface{}) { d.viewData = data }

// DefaultDB is a type that can be embedded into a Context type
// to enable a database connection with default configuration
// options
type DefaultDB struct {
	db DB
}

func (d *DefaultDB) CONFIG_DB() DBConfig {
	return DBConfig{
		Driver: envOrDefault("DB_DRIVER", ""),
		DSN:    envOrDefault("DB_DSN", ""),
	}
}

func (d *DefaultDB) DB() DB      { return d.db }
func (d *DefaultDB) setDB(db DB) { d.db = db }

// DefaultSession is a type that can be embedded into a Context type
// to enable user sessions
type DefaultSession struct {
	sess Session
}

func (d *DefaultSession) Session() Session     { return d.sess }
func (d *DefaultSession) setSession(s Session) { d.sess = s }
func (d *DefaultSession) CONFIG_Session() SessionConfig {
	return SessionConfig{
		Name:   envOrDefault("SESSION_NAME", "AttacheSession"),
		Secret: []byte(envOrDefault("SESSION_SECRET", "")),
	}
}

// DefaultFileServer is a type that can be embedded into a Context type
// to enable a static file server with default configuration options
type DefaultFileServer struct{}

// MOUNT_Web provides a static file server under the path /web/*
func (DefaultFileServer) MOUNT_Web() (http.Handler, error) {
	return http.FileServer(http.Dir("web/dist")), nil
}

func envOrDefault(s, dflt string) string {
	if got := os.Getenv(s); got != "" {
		return got
	}

	return dflt
}
