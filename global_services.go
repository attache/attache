package attache

import (
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

var (
	gsCache       = cache{}
	gsFormDecoder = schema.NewDecoder()
	gsSessions    = sessions.NewCookieStore()
)

// ViewCacheFor will return the cached ViewCache or a
// newly initialized ViewCache for the given ViewConfig.
// Any error encountered while initializing the new ViewCahce
// is returned.
func ViewCacheFor(conf ViewConfig) (ViewCache, error) {
	return gsCache.viewsFor(conf)
}

func ViewCacheRefresh(conf ViewConfig) error {
	newCache := viewCache{}
	if err := newCache.load(conf.Root, "", nil); err != nil {
		return err
	}

	gsCache.vcCache.put(conf, newCache)
	return nil
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
