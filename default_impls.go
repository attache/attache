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

// CONFIG_Environment implements HasEnvironment for DefaultEnvironment
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

// CONFIG_Views implements HasViews for DefaultViews
func (d *DefaultViews) CONFIG_Views() ViewConfig {
	return ViewConfig{
		Driver: envOrDefault("VIEW_DRIVER", "attache"),
		Root:   envOrDefault("VIEW_ROOT", "views"),
	}
}

// Views implements HasViews for DefaultViews
func (d *DefaultViews) Views() ViewCache { return d.views }

// SetViews implements HasViews for DefaultViews
func (d *DefaultViews) SetViews(v ViewCache) { d.views = v }

// ViewData implements HasViews for DefaultViews
func (d *DefaultViews) ViewData() interface{} { return d.viewData }

// SetViewData implements HasViews for DefaultViews
func (d *DefaultViews) SetViewData(data interface{}) { d.viewData = data }

// DefaultDB is a type that can be embedded into a Context type
// to enable a database connection with default configuration
// options
type DefaultDB struct {
	db DB
}

// CONFIG_DB implements HasDB for DefaultDB
func (d *DefaultDB) CONFIG_DB() DBConfig {
	return DBConfig{
		Driver: envOrDefault("DB_DRIVER", ""),
		DSN:    envOrDefault("DB_DSN", ""),
	}
}

// DB implements HasDB for DefaultDB
func (d *DefaultDB) DB() DB { return d.db }

// SetDB implements HasDB for DefaultDB
func (d *DefaultDB) SetDB(db DB) { d.db = db }

// DefaultSession is a type that can be embedded into a Context type
// to enable user sessions
type DefaultSession struct {
	sess Session
}

// CONFIG_Session implements HasSession for DefaultSession
func (d *DefaultSession) CONFIG_Session() SessionConfig {
	return SessionConfig{
		Name:   envOrDefault("SESSION_NAME", "AttacheSession"),
		Secret: []byte(envOrDefault("SESSION_SECRET", "")),
	}
}

// Session implements HasSession for DefaultSession
func (d *DefaultSession) Session() Session { return d.sess }

// SetSession implements HasSession for DefaultSession
func (d *DefaultSession) SetSession(s Session) { d.sess = s }

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
