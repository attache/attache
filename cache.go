package attache

import (
	"database/sql"
	"path/filepath"
	"sync"

	"github.com/gocraft/dbr"
)

type vcCache struct {
	sync.RWMutex
	have map[ViewConfig]viewCache
}

func (v *vcCache) lazy() {
	if v.have == nil {
		v.have = make(map[ViewConfig]viewCache, 1)
	}
}

func (v *vcCache) lookupOk(key ViewConfig) (viewCache, bool) {
	v.RLock()
	defer v.RUnlock()

	if v.have == nil {
		return nil, false
	}

	got, ok := v.have[key]
	return got, ok
}

func (v *vcCache) put(key ViewConfig, vc viewCache) {
	v.Lock()
	defer v.Unlock()
	v.lazy()
	v.have[key] = vc
}

type dbCache struct {
	sync.RWMutex
	have map[DBConfig]db
}

func (c *dbCache) lazy() {
	if c.have == nil {
		c.have = make(map[DBConfig]db, 1)
	}
}

func (c *dbCache) lookupOk(key DBConfig) (db, bool) {
	c.RLock()
	defer c.RUnlock()

	if c.have == nil {
		return db{}, false
	}

	got, ok := c.have[key]
	return got, ok
}

func (c *dbCache) put(key DBConfig, toPut db) {
	c.Lock()
	defer c.Unlock()
	c.lazy()
	c.have[key] = toPut
}

type cache struct {
	vcCache vcCache
	dbCache dbCache
}

func (c *cache) viewsFor(conf ViewConfig) (ViewCache, error) {
	conf.Root = filepath.Clean(conf.Root)

	if cached, ok := c.vcCache.lookupOk(conf); ok {
		return cached, nil
	}

	v := viewCache{}
	if err := v.load(conf.Root, "", nil); err != nil {
		return v, err
	}

	c.vcCache.put(conf, v)
	return v, nil
}

func (c *cache) dbFor(conf DBConfig) (DB, error) {
	if db, ok := c.dbCache.lookupOk(conf); ok {
		if err := db.conn.Ping(); err != nil {
			return nil, err
		}

		return db, nil
	}

	dbr.Open(conf.Driver, conf.DSN)
	conn, err := sql.Open(conf.Driver, conf.DSN)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(); err != nil {
		return nil, err
	}

	got := db{conn: conn}
	c.dbCache.put(conf, got)
	return got, nil
}
