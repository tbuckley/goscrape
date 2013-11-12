package registry

import (
	"net/url"
	"sync"
)

// MultithreadRegistry objects store whether or not URLs have been registered
// (which can mean visited, queued, or anything else). This implementation
// wraps the SimpleRegistry with synchronization primitives to be thread-safe.
type MultithreadRegistry struct {
	*SimpleRegistry
	lock sync.Mutex
}

// NewMultithreadRegistry creates an empty MultithreadRegistry.
func NewMultithreadRegistry() Registry {
	r := new(MultithreadRegistry)
	r.SimpleRegistry = NewSimpleRegistry().(*SimpleRegistry)
	return r
}

// RegisterIfNot is an atomic operation that registers the URL if it has not
// already been registered. It returns true if the url is successfully
// registered.
func (r *MultithreadRegistry) RegisterIfNot(page *url.URL) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.SimpleRegistry.RegisterIfNot(page)
}

// IsRegistered returns true only if the URL has been registered.
func (r *MultithreadRegistry) IsRegistered(page *url.URL) bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.SimpleRegistry.IsRegistered(page)
}

// Register registers the URL.
func (r *MultithreadRegistry) Register(page *url.URL) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.SimpleRegistry.Register(page)
}
