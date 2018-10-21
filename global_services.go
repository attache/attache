package attache

import (
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

var (
	gsCache       = cache{}
	gsFormDecoder = schema.NewDecoder()
	gsSessions    = sessions.NewCookieStore()
)

// LoadEnvironment will attempt to load environment variables
// based on the given EnvironmentConfig
func LoadEnvironment(conf EnvironmentConfig) error {
	return godotenv.Load(conf.EnvPath)
}

// DBFor will returh the cached DB or a newly initialized DB
// for the given DBConfig. Any error encountered while
// initializing a new DB is returned.
func DBFor(conf DBConfig) (DB, error) {
	return gsCache.dbFor(conf)
}

func init() {
	gsFormDecoder.IgnoreUnknownKeys(true)
}
