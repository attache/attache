package attache

import (
	"sync"

	"github.com/gocraft/dbr"
)

type dbCache struct {
	sync.RWMutex
	have map[DBConfig]*dbr.Connection
}

func (c *dbCache) lazy() {
	if c.have == nil {
		c.have = make(map[DBConfig]*dbr.Connection, 1)
	}
}

func (c *dbCache) lookupOk(key DBConfig) (*dbr.Connection, bool) {
	c.RLock()
	defer c.RUnlock()

	if c.have == nil {
		return nil, false
	}

	got, ok := c.have[key]
	return got, ok
}

func (c *dbCache) put(key DBConfig, toPut *dbr.Connection) {
	c.Lock()
	defer c.Unlock()
	c.lazy()
	c.have[key] = toPut
}

type cache struct {
	dbCache dbCache
}

func (c *cache) dbFor(conf DBConfig) (DB, error) {
	if conn, ok := c.dbCache.lookupOk(conf); ok {
		if err := conn.Ping(); err != nil {
			return DB{}, err
		}

		return DB{conn.NewSession(nil)}, nil
	}

	conn, err := dbr.Open(conf.Driver, conf.DSN, nil)
	if err != nil {
		return DB{}, err
	}

	if err = conn.Ping(); err != nil {
		return DB{}, err
	}

	c.dbCache.put(conf, conn)
	return DB{conn.NewSession(nil)}, nil
}
