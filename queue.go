package chuper

import (
	"net/url"

	"github.com/PuerkitoBio/fetchbot"
)

type Queue struct {
	*fetchbot.Queue
}

type Enqueuer interface {
	Enqueue(string, string, string) error

	EnqueueWithBasicAuth(string, string, string, string, string) error
}

func (q *Queue) Enqueue(method, URL, sourceURL string) error {
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}
	s, err := url.Parse(sourceURL)
	if err != nil {
		return err
	}

	cmd := &Cmd{&fetchbot.Cmd{U: u, M: method}, s}
	if err = q.Send(cmd); err != nil {
		return err
	}
	return nil
}

func (q *Queue) EnqueueWithBasicAuth(method, URL, sourceURL, user, password string) error {
	if user == "" && password == "" {
		return q.Enqueue(method, URL, sourceURL)
	}

	u, err := url.Parse(URL)
	if err != nil {
		return err
	}
	s, err := url.Parse(sourceURL)
	if err != nil {
		return err
	}

	cmd := &CmdBasicAuth{&fetchbot.Cmd{U: u, M: method}, s, user, password}
	if err = q.Send(cmd); err != nil {
		return err
	}
	return nil
}