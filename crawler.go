package chuper

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/fetchbot"
	"github.com/PuerkitoBio/goquery"
)

const (
	DefaultCrawlDelay      = 5 * time.Second
	DefaultCrawlPoliteness = false
	DefaultUserAgent       = fetchbot.DefaultUserAgent
)

var (
	DefaultHTTPClient = http.DefaultClient

	DefaultCache = NewMemoryCache()

	DefaultErrorHandler = fetchbot.HandlerFunc(func(ctx *fetchbot.Context, res *http.Response, err error) {
		fmt.Printf("chuper - %s - error: %s %s - %s\n", time.Now().Format(time.RFC3339), ctx.Cmd.Method(), ctx.Cmd.URL(), err)
	})

	DefaultLogHandlerFunc = func(ctx *fetchbot.Context, res *http.Response, err error) {
		if err == nil {
			fmt.Printf("chuper - %s - info: [%d] %s %s - %s\n", time.Now().Format(time.RFC3339), res.StatusCode, ctx.Cmd.Method(), ctx.Cmd.URL(), res.Header.Get("Content-Type"))
		}
	}
)

type Crawler struct {
	CrawlDelay      time.Duration
	CrawlDuration   time.Duration
	CrawlPoliteness bool
	HTTPClient      fetchbot.Doer
	Cache           Cache
	ErrorHandler    fetchbot.Handler
	UserAgent       string
	LogHandlerFunc  func(ctx *fetchbot.Context, res *http.Response, err error)

	mux *fetchbot.Mux
	f   *fetchbot.Fetcher
	q   *fetchbot.Queue
}

// New returns an initialized Crawler.
func New() *Crawler {
	return &Crawler{
		CrawlDelay:      DefaultCrawlDelay,
		CrawlPoliteness: DefaultCrawlPoliteness,
		HTTPClient:      DefaultHTTPClient,
		Cache:           DefaultCache,
		ErrorHandler:    DefaultErrorHandler,
		UserAgent:       DefaultUserAgent,
		LogHandlerFunc:  DefaultLogHandlerFunc,
		mux:             fetchbot.NewMux(),
	}
}

func (c *Crawler) Start() Enqueuer {
	c.mux.HandleErrors(c.ErrorHandler)
	l := newLogHandler(c.mux, c.LogHandlerFunc)

	f := fetchbot.New(l)
	f.CrawlDelay = c.CrawlDelay
	f.DisablePoliteness = !c.CrawlPoliteness
	f.HttpClient = c.HTTPClient
	f.UserAgent = c.UserAgent

	c.f = f
	c.q = c.f.Start()

	if c.CrawlDuration > 0 {
		go func() {
			t := time.After(c.CrawlDuration)
			<-t
			c.q.Close()
		}()
	}

	return &Queue{c.q}
}

func (c *Crawler) Block() {
	c.q.Block()
}

func (c *Crawler) Finish() {
	c.q.Close()
}

type ResponseCriteria struct {
	Method      string
	ContentType string
	Status      int
	MinStatus   int
	MaxStatus   int
	Path        string
	Host        string
}

func (c *Crawler) Match(r *ResponseCriteria) *fetchbot.ResponseMatcher {
	m := c.mux.Response()

	if r.Method != "" {
		m.Method(r.Method)
	}

	if r.ContentType != "" {
		m.ContentType(r.ContentType)
	}

	if r.Status != 0 {
		m.Status(r.Status)
	} else {
		if r.MinStatus != 0 && r.MaxStatus != 0 {
			m.StatusRange(r.MinStatus, r.MaxStatus)
		} else {
			if r.MinStatus != 0 {
				m.Status(r.MinStatus)
			}
			if r.MaxStatus != 0 {
				m.Status(r.MaxStatus)
			}
		}
	}

	if r.Path != "" {
		m.Path(r.Path)
	}

	if r.Host != "" {
		m.Host(r.Host)
	}

	return m
}

func (c *Crawler) Register(rc *ResponseCriteria, procs ...Processor) {
	m := c.Match(rc)
	h := newDocHandler(c.Cache, procs...)
	m.Handler(h)
}

func newLogHandler(wrapped fetchbot.Handler, f func(ctx *fetchbot.Context, res *http.Response, err error)) fetchbot.Handler {
	return fetchbot.HandlerFunc(func(ctx *fetchbot.Context, res *http.Response, err error) {
		f(ctx, res, err)
		wrapped.Handle(ctx, res, err)
	})
}

func newDocHandler(cache Cache, procs ...Processor) fetchbot.Handler {
	return fetchbot.HandlerFunc(func(ctx *fetchbot.Context, res *http.Response, err error) {
		context := &Ctx{ctx, cache}
		doc, err := goquery.NewDocumentFromResponse(res)
		if err != nil {
			fmt.Printf("chuper - %s - error: %s %s - %s\n", time.Now().Format(time.RFC3339), ctx.Cmd.Method(), ctx.Cmd.URL(), err)
			return
		}
		for _, p := range procs {
			ok := p.Process(context, doc)
			if !ok {
				return
			}
		}
	})
}
