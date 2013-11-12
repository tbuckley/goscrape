package registry

import (
	"net/url"
)

// SimpleRegistry objects store whether or not URLs have been registered
// (which can mean visited, queued, or anything else).
type SimpleRegistry struct {
	registry map[string]bool
}

// NewMultithreadRegistry creates an empty MultithreadRegistry.
func NewSimpleRegistry() Registry {
	r := new(SimpleRegistry)
	r.registry = make(map[string]bool)
	return r
}

// RegisterIfNot registers the URL if it has not already been registered. It
// returns true if the url is successfully registered.
func (r *SimpleRegistry) RegisterIfNot(page *url.URL) bool {
	if r.IsRegistered(page) {
		return false
	}
	r.Register(page)
	return true
}

// IsRegistered returns true only if the URL has been registered.
func (r *SimpleRegistry) IsRegistered(page *url.URL) bool {
	return r.registry[page.String()]
}

// Register registers the URL.
func (r *SimpleRegistry) Register(page *url.URL) {
	r.registry[page.String()] = true
}
