package goscrape

import (
	"expvar"
	"net/url"
	"regexp"

	"github.com/tbuckley/goscrape/queue"
	"github.com/tbuckley/goscrape/registry"
)

var (
	varHandled   = expvar.NewInt("handled")
	varUnhandled = expvar.NewInt("unhandled")
	varEnqueued  = expvar.NewInt("enqueued")
)

type Scraper struct {
	queue    queue.AsyncPrioQueue
	registry registry.Registry
	handlers *PatternHandler
	pagechan chan *url.URL
}

func NewScraper() WebScraper {
	q := queue.NewPrioQueue("scraperqueue")
	return &Scraper{
		queue:    queue.NewMultithreadPrioQueue(q),
		registry: registry.NewMultithreadRegistry(),
		handlers: &PatternHandler{},
		pagechan: make(chan *url.URL),
	}
}

// Enqueue adds a URL to be scraped. If the URL has already been added to the
// queue, it will be ignored. URLs will be visited in the order they are added.
func (s *Scraper) Enqueue(page *url.URL) {
	// @TODO(tbuckley) Rewrite the URL
	if s.registry.RegisterIfNot(page) {
		_, prio, ok := s.handlers.GetHandler(page)
		if ok {
			s.queue.PushPriority(page, prio)
			varEnqueued.Add(1)
		}
	}
}

// AddHandler registers a handler for the given URL pattern. Patterns are
// tested in the order they are added, and the handler corresponding to the
// first successful match is chosen.
func (s *Scraper) AddHandler(pattern *regexp.Regexp, handler Handler) {
	s.handlers.Register(pattern, handler, MediumPriority)
}

//
func (s *Scraper) AddHandlerPriority(pattern *regexp.Regexp, handler Handler, priority uint) {
	s.handlers.Register(pattern, handler, priority)
}

func (s *Scraper) handle(page *url.URL) {
	handler, _, ok := s.handlers.GetHandler(page)
	if ok {
		varHandled.Add(1)
		handler(s, page)
	} else {
		varUnhandled.Add(1)
	}
}

func (s *Scraper) handlerDaemon() {
	for {
		page, ok := <-s.pagechan
		if ok {
			// Got a page, handle it
			s.handle(page)
		} else {
			// Channel was closed
			return
		}
	}
}

func (s *Scraper) Start() {
	for i := 0; i < 20; i++ {
		go s.handlerDaemon()
	}

	// @TODO(tbuckley) figure out how to determine completion and exit gracefully
	// Do you need to tell how many workers are free?
	for {
		s.pagechan <- s.queue.PopBlock()
	}
}
