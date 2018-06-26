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

type DefaultToken struct {
	tok Token
}

func (d *DefaultToken) Token() Token     { return d.tok }
func (d *DefaultToken) SetToken(t Token) { d.tok = t }
func (d *DefaultToken) CONFIG_Token() TokenConfig {
	return TokenConfig{
		Secret: []byte(envOrDefault("TOKEN_SECRET", "")),
		Cookie: envOrDefault("TOKEN_COOKIE", "ATTACHE_TOKEN"),
		MaxAge: 10 * 60 * 1000,
	}
}

type DefaultFileServer struct{}

func (d DefaultFileServer) CONFIG_FileServer() FileServerConfig {
	return FileServerConfig{
		Root:     "web/dist",
		BasePath: "web",
	}
}
