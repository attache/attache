package view

import "sync"

// A Cache provides a common view cache type for driver implementations.
// It is safe to share between goroutines.
type Cache struct {
	lock sync.RWMutex
	list map[string]View
}

func (c *Cache) lazy() {
	// lazy must only be called when c.lock is locked for writing
	if c.list == nil {
		c.list = map[string]View{}
	}
}

// Put stores the View under the given name
func (c *Cache) Put(name string, view View) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.lazy()
	c.list[name] = view
}

// Get retrieves the View stored under the given name, or returns the None no-op View.
// The returned view will always be non-nil.
func (c *Cache) Get(name string) View {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.list == nil {
		return None
	}

	if v := c.list[name]; v != nil {
		return v
	}

	return None
}
