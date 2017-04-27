package chuper

import (
	"net/url"

	"github.com/PuerkitoBio/fetchbot"
)

type Queue struct {
	*fetchbot.Queue
}

type Enqueuer interface {
	Enqueue(string, string, string, int, int) error

	EnqueueWithBasicAuth(string, string, string, int, int, string, string) error
}

func (q *Queue) Enqueue(method, URL, sourceURL string, depth, retries int) error {
	u, err := url.Parse(URL)
	if err != nil {
		return err
	}
	s, err := url.Parse(sourceURL)
	if err != nil {
		return err
	}

	cmd := &Cmd{&fetchbot.Cmd{U: u, M: method}, s, depth, retries}
	if err = q.Send(cmd); err != nil {
		return err
	}
	return nil
}

func (q *Queue) EnqueueWithBasicAuth(method string, URL string, sourceURL string, depth, retries int, user string, password string) error {
	if user == "" && password == "" {
		return q.Enqueue(method, URL, sourceURL, depth, retries)
	}

	u, err := url.Parse(URL)
	if err != nil {
		return err
	}
	s, err := url.Parse(sourceURL)
	if err != nil {
		return err
	}

	cmd := &CmdBasicAuth{&fetchbot.Cmd{U: u, M: method}, s, depth, retries, user, password}
	if err = q.Send(cmd); err != nil {
		return err
	}
	return nil
}
