package handlers

import (
	"net/url"

	"github.com/tbuckley/goscrape"
)

// Null does nothing
func Null(s goscrape.WebScraper, page *url.URL) {}
