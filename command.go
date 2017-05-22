package chuper

import (
	"errors"
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

func NewCmd(method, urlStr, sourceURLStr string, depth, retries int) (Cmd, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return Cmd{}, err
	}

	s, err := url.Parse(sourceURLStr)
	if err != nil {
		return Cmd{}, err
	}

	c := Cmd{
		Cmd: &fetchbot.Cmd{U: u, M: method},
		S:   s,
		D:   depth,
		R:   retries,
	}

	return c, nil
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

func NewCmdBasiAuth(method, urlStr, sourceURLStr string, depth, retries int, user, password string) (CmdBasicAuth, error) {
	var cmd CmdBasicAuth

	if user == "" && password == "" {
		return cmd, errors.New("user or password not provided")
	}

	baseCommand, err := NewCmd(method, urlStr, sourceURLStr, depth, retries)
	if err != nil {
		return cmd, err
	}

	cmd = CmdBasicAuth{
		Cmd:  baseCommand,
		user: user,
		pass: password,
	}

	return cmd, nil
}

func (c *CmdBasicAuth) BasicAuth() (string, string) {
	return c.user, c.pass
}

type CmdHeader struct {
	Cmd
	header http.Header
}

func NewCmdHeader(method, urlStr, sourceURLStr string, depth, retries int, header http.Header) (CmdHeader, error) {
	var cmd CmdHeader
	if header == nil {
		return cmd, errors.New("header not provided")
	}

	baseCommand, err := NewCmd(method, urlStr, sourceURLStr, depth, retries)
	if err != nil {
		return cmd, err
	}

	cmd = CmdHeader{
		Cmd:    baseCommand,
		header: header,
	}

	return cmd, nil
}

func (c *CmdHeader) Header() http.Header {
	return c.header
}
