package attache

import (
	"sync"

	viewDriver "github.com/attache/attache/drivers/view"
)

// A View is the interface implemented by a renderable type
type View = viewDriver.View

var _ View = viewDriver.None

// A ViewCache is a read-only view of a set of cached views. Implementations
// must be safe for concurrent use.
type ViewCache interface {
	// Get will return a valid view if one is stored under the given name,
	// otherwise it should return a valid no-op view.
	Get(name string) View
}

var viewCacheCache = &struct {
	sync.RWMutex
	list map[ViewConfig]*viewDriver.Cache
}{
	list: map[ViewConfig]*viewDriver.Cache{},
}

// ViewCacheFor will return the cached ViewCache or a
// newly initialized ViewCache for the given ViewConfig.
// Any error encountered while initializing the new ViewCahce
// is returned.
func ViewCacheFor(conf ViewConfig) (ViewCache, error) {
	if found := getCachedViewCache(conf); found != nil {
		return found, nil
	}

	return initViewCache(conf)
}

// ViewCacheRefresh will re-initialize the ViewCache for the given ViewConfig.
func ViewCacheRefresh(conf ViewConfig) error {
	_, err := initViewCache(conf)
	return err
}

func initViewCache(conf ViewConfig) (*viewDriver.Cache, error) {
	cache, err := viewDriver.DriverInit(conf.Driver, conf.Root)
	if err != nil {
		return nil, err
	}

	viewCacheCache.Lock()
	defer viewCacheCache.Unlock()
	viewCacheCache.list[conf] = cache
	return cache, nil
}

func getCachedViewCache(conf ViewConfig) *viewDriver.Cache {
	viewCacheCache.RLock()
	defer viewCacheCache.RUnlock()
	return viewCacheCache.list[conf]
}
