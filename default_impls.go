package attache

import (
	"os"
)

func envOrDefault(s, dflt string) string {
	if got := os.Getenv(s); got != "" {
		return got
	}
	return dflt
}

// DefaultViews is a type that can be embedded into a Context type
// to enable views with default configuration options
type DefaultViews struct {
	views ViewCache
}

func (d *DefaultViews) Views() ViewCache     { return d.views }
func (d *DefaultViews) SetViews(v ViewCache) { d.views = v }
func (d *DefaultViews) CONFIG_Views() ViewConfig {
	return ViewConfig{
		Root: "views",
	}
}

// DefaultDB is a type that can be embedded into a Context type
// to enable a database connection with default configuration
// options
type DefaultDB struct {
	db DB
}

func (d *DefaultDB) DB() DB      { return d.db }
func (d *DefaultDB) SetDB(db DB) { d.db = db }
func (d *DefaultDB) CONFIG_DB() DBConfig {
	return DBConfig{
		Driver: envOrDefault("DB_DRIVER", ""),
		DSN:    envOrDefault("DB_DSN", ""),
	}
}

// DefaultToken is a type that can be embedded into a Context type
// to enable managed JWTs with default configuration options
type DefaultToken struct {
	tok Token
}

func (d *DefaultToken) Token() Token     { return d.tok }
func (d *DefaultToken) SetToken(t Token) { d.tok = t }
func (d *DefaultToken) CONFIG_Token() TokenConfig {
	return TokenConfig{
		Secret: []byte(envOrDefault("TOKEN_SECRET", "")),
	}
}

// DefaultSession is a type that can be embedded into a Context type
// to enable user sessions
type DefaultSession struct {
	sess Session
}

func (d *DefaultSession) Session() Session      { return d.sess }
func (d *DefaultSession) SetSesstion(s Session) { d.sess = s }
func (d *DefaultSession) CONFIG_Session() SessionConfig {
	return SessionConfig{
		Name:   envOrDefault("SESSION_NAME", "AttacheSession"),
		Secret: []byte(envOrDefault("SESSION_SECRET", "")),
	}
}

// DefaultFileServer is a type that can be embedded into a Context type
// to enable a static file server with default configuration options
type DefaultFileServer struct{}

func (d DefaultFileServer) CONFIG_FileServer() FileServerConfig {
	return FileServerConfig{
		Root:     "web/dist",
		BasePath: "web",
	}
}
