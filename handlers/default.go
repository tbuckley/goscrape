package handlers

import (
	"log"
	"net/url"
	"strings"

	"github.com/tbuckley/goscrape"
	"github.com/tbuckley/htmlutils"
)

func HrefToUrl(base *url.URL, href string) (*url.URL, bool) {
	if strings.HasPrefix(href, "#") {
		return nil, false
	}
	hrefURL, err := base.Parse(href)
	if err != nil {
		return nil, false
	}
	if hrefURL.Scheme == "javascript" {
		return nil, false
	}
	return hrefURL, true
}

func DefaultPage(s goscrape.WebScraper, page *url.URL, body []byte) {
	// Create a query
	query, err := htmlutils.NewQueryFromPage(body)
	if err != nil {
		log.Printf("Error w/ Query: %s", err)
		return
	}

	// Identify links
	hrefs := query.ElementsByTagName("a").Attr("href")

	// Add links to queue
	for _, href := range hrefs {
		hrefURL, ok := HrefToUrl(page, href)
		if ok {
			s.Enqueue(hrefURL)
		}
	}
}

// Default finds the links on a page and adds them to the scraper's queue.
func Default(s goscrape.WebScraper, page *url.URL) {
	// Load a page
	body, err := htmlutils.FetchPage(page)
	if err != nil {
		log.Printf("Error w/ Query: %s", err)
		return
	}
	DefaultPage(s, page, body)
}
