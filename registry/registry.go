package registry

import (
	"net/url"
)

// Registry objects are used to keep track of URLs that have been seen,
// visited, etc.
type Registry interface {
	// IsRegistered returns true only if the URL has been registered.
	IsRegistered(*url.URL) bool

	// Register registers the URL.
	Register(*url.URL)

	// RegisterIfNot is an atomic operation that registers the URL if it has not
	// already been registered. It returns true if the url is successfully
	// registered.
	RegisterIfNot(*url.URL) bool
}
