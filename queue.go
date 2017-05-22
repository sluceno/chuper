package chuper

import (
	"net/http"

	"github.com/PuerkitoBio/fetchbot"
)

type Queue struct {
	*fetchbot.Queue
}

type Enqueuer interface {
	Enqueue(string, string, string, int, int) error
	EnqueueWithBasicAuth(string, string, string, int, int, string, string) error
	EnqueueWithHeader(string, string, string, int, int, http.Header) error
}

func (q *Queue) Enqueue(method, URL, sourceURL string, depth, retries int) error {
	cmd, err := NewCmd(method, URL, sourceURL, depth, retries)
	if err != nil {
		return err
	}

	if err = q.Send(&cmd); err != nil {
		return err
	}
	return nil
}

func (q *Queue) EnqueueWithBasicAuth(method string, URL string, sourceURL string, depth, retries int, user string, password string) error {
	cmd, err := NewCmdBasiAuth(method, URL, sourceURL, depth, retries, user, password)
	if err != nil {
		return err
	}

	if err = q.Send(&cmd); err != nil {
		return err
	}

	return nil
}

func (q *Queue) EnqueueWithHeader(method string, URL string, sourceURL string, depth, retries int, header http.Header) error {
	cmd, err := NewCmdHeader(method, URL, sourceURL, depth, retries, header)
	if err != nil {
		return err
	}

	if err = q.Send(&cmd); err != nil {
		return err
	}

	return nil
}
