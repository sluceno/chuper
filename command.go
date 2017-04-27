package chuper

import (
	"net/url"

	"github.com/PuerkitoBio/fetchbot"
)

type Command interface {
	URL() *url.URL
	Method() string
	SourceURL() *url.URL
	Depth() int
	Retries() int
}

type Cmd struct {
	*fetchbot.Cmd
	S *url.URL
	D int
	R int
}

func (c *Cmd) SourceURL() *url.URL {
	return c.S
}

func (c *Cmd) Depth() int {
	return c.D
}

func (c *Cmd) Retries() int {
	return c.R
}

type CmdBasicAuth struct {
	*fetchbot.Cmd
	S          *url.URL
	D          int
	R          int
	user, pass string
}

func (c *CmdBasicAuth) SourceURL() *url.URL {
	return c.S
}

func (c *CmdBasicAuth) Depth() int {
	return c.D
}

func (c *CmdBasicAuth) Retries() int {
	return c.R
}

func (c *CmdBasicAuth) BasicAuth() (string, string) {
	return c.user, c.pass
}
