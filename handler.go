package goscrape

import (
	"net/url"
	"regexp"
	"strings"
)

const (
	HighPriority = iota
	MediumPriority
	LowPriority
)

type Handler func(WebScraper, *url.URL)

type patternHandlerPair struct {
	pattern  *regexp.Regexp
	handler  Handler
	priority int
}

type PatternHandler struct {
	handlers []patternHandlerPair
}

// Register associates a handler with the given pattern. Handlers are given
// the opportunity to handle URLs in the order they were registered.
func (s *PatternHandler) Register(pattern *regexp.Regexp, handler Handler, priority int) {
	s.handlers = append(s.handlers, patternHandlerPair{
		pattern:  pattern,
		handler:  handler,
		priority: priority,
	})
}

// Handle tries to pass the page off to the first registered handler that
// matches, returning true only if a handler exists.
func (h *PatternHandler) GetHandler(page *url.URL) (Handler, int, bool) {
	for _, ph := range h.handlers {
		lowerPage := strings.ToLower(page.String())
		if ph.pattern.MatchString(lowerPage) {
			return ph.handler, ph.priority, true
		}
	}
	return nil, 0, false
}
