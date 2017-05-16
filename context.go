package chuper

import (
	"net/url"

	"github.com/PuerkitoBio/fetchbot"
	"github.com/Sirupsen/logrus"
)

type Context interface {
	Cache() Cache
	Queue() Enqueuer
	Log(fields map[string]interface{}) *logrus.Entry
	URL() *url.URL
	Method() string
	SourceURL() *url.URL
	FinalURL() *url.URL
	Depth() int
	Retries() int
}

type Ctx struct {
	*fetchbot.Context
	C Cache
	L *logrus.Logger
}

func (c *Ctx) Cache() Cache {
	return c.C
}

func (c *Ctx) Queue() Enqueuer {
	return &Queue{c.Q}
}

func (c *Ctx) Log(fields map[string]interface{}) *logrus.Entry {
	data := logrus.Fields{}
	for k, v := range fields {
		data[k] = v
	}
	return c.L.WithFields(data)
}

func (c *Ctx) URL() *url.URL {
	return c.Cmd.URL()
}

func (c *Ctx) Method() string {
	return c.Cmd.Method()
}

func (c *Ctx) SourceURL() *url.URL {
	switch cmd := c.Cmd.(type) {
	case *Cmd:
		return cmd.SourceURL()
	case *CmdBasicAuth:
		return cmd.SourceURL()
	case *CmdHeader:
		return cmd.SourceURL()
	default:
		return nil
	}
}

func (c *Ctx) Depth() int {
	switch cmd := c.Cmd.(type) {
	case *Cmd:
		return cmd.Depth()
	case *CmdBasicAuth:
		return cmd.Depth()
	case *CmdHeader:
		return cmd.Depth()
	default:
		return 0
	}
}

func (c *Ctx) Retries() int {
	switch cmd := c.Cmd.(type) {
	case *Cmd:
		return cmd.Retries()
	case *CmdBasicAuth:
		return cmd.Retries()
	case *CmdHeader:
		return cmd.Retries()
	default:
		return 0
	}
}

func (c *Ctx) FinalURL() *url.URL {
	switch cmd := c.Cmd.(type) {
	case *Cmd:
		return cmd.FinalURL()
	case *CmdBasicAuth:
		return cmd.FinalURL()
	case *CmdHeader:
		return cmd.FinalURL()
	default:
		return nil
	}
}

func (c *Ctx) SetFinalURL(url *url.URL) {
	switch cmd := c.Cmd.(type) {
	case *Cmd:
		cmd.F = url
	case *CmdBasicAuth:
		cmd.F = url
	case *CmdHeader:
		cmd.F = url
	}
}
