package goscrape

import (
	"net/url"
	"regexp"
)

type WebScraper interface {
	Enqueue(*url.URL)
	AddHandler(*regexp.Regexp, Handler)
	AddHandlerPriority(*regexp.Regexp, Handler, int)
	Start()
}
