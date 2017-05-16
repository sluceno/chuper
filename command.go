package chuper

import (
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/fetchbot"
)

type Command interface {
	URL() *url.URL
	Method() string
	SourceURL() *url.URL
	FinalURL() *url.URL
	Depth() int
	Retries() int
}

type Cmd struct {
	*fetchbot.Cmd
	S *url.URL
	F *url.URL
	D int
	R int
}

func (c *Cmd) SourceURL() *url.URL {
	return c.S
}

func (c *Cmd) FinalURL() *url.URL {
	return c.F
}

func (c *Cmd) Depth() int {
	return c.D
}

func (c *Cmd) Retries() int {
	return c.R
}

type CmdBasicAuth struct {
	Cmd
	user, pass string
}

func (c *CmdBasicAuth) BasicAuth() (string, string) {
	return c.user, c.pass
}

type CmdHeader struct {
	Cmd
	header http.Header
}

func (c *CmdHeader) Header() http.Header {
	return c.header
}
