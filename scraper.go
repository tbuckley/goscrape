package goscrape

import (
	"expvar"
	"log"
	"net/url"
	"regexp"
	"time"

	"github.com/tbuckley/goscrape/queue"
	"github.com/tbuckley/goscrape/registry"
)

var (
	varHandled   = expvar.NewInt("handled")
	varUnhandled = expvar.NewInt("unhandled")
	varEnqueued  = expvar.NewInt("enqueued")
)

type Scraper struct {
	lowqueue  queue.Queue
	medqueue  queue.Queue
	highqueue queue.Queue
	registry  registry.Registry
	handlers  *PatternHandler
	pagechan  chan *url.URL
}

func NewScraper() WebScraper {
	return &Scraper{
		lowqueue:  queue.NewMultithreadQueue("lowPriority"),
		medqueue:  queue.NewMultithreadQueue("medPriority"),
		highqueue: queue.NewMultithreadQueue("highPriority"),
		registry:  registry.NewMultithreadRegistry(),
		handlers:  &PatternHandler{},
		pagechan:  make(chan *url.URL),
	}
}

// Enqueue adds a URL to be scraped. If the URL has already been added to the
// queue, it will be ignored. URLs will be visited in the order they are added.
func (s *Scraper) Enqueue(page *url.URL) {
	// @TODO(tbuckley) Rewrite the URL
	if s.registry.RegisterIfNot(page) {
		_, priority, ok := s.handlers.GetHandler(page)
		if !ok {
			return
		}
		switch priority {
		case LowPriority:
			s.lowqueue.Push(page)
		case MediumPriority:
			s.lowqueue.Push(page)
		case HighPriority:
			s.highqueue.Push(page)
		}
		varEnqueued.Add(1)
	}
}

// AddHandler registers a handler for the given URL pattern. Patterns are
// tested in the order they are added, and the handler corresponding to the
// first successful match is chosen.
func (s *Scraper) AddHandler(pattern *regexp.Regexp, handler Handler) {
	s.handlers.Register(pattern, handler, MediumPriority)
}

//
func (s *Scraper) AddHandlerPriority(pattern *regexp.Regexp, handler Handler, priority int) {
	s.handlers.Register(pattern, handler, priority)
}

func (s *Scraper) handle(page *url.URL) {
	handler, _, ok := s.handlers.GetHandler(page)
	if ok {
		varHandled.Add(1)
		log.Printf("HANDLED: %s\n", page.String())
		handler(s, page)
	} else {
		log.Printf("UNHANDLED: %s\n", page.String())
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

	queue := queue.NewComboQueue(s.highqueue, s.medqueue, s.lowqueue)
	for {
		page, ok := queue.Pop()
		if ok {
			s.pagechan <- page
		} else {
			time.Sleep(500 * time.Millisecond)
		}
	}
}
