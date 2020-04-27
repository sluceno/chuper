package chuper

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Processor interface {
	Process(Context, *http.Response, []byte, *goquery.Document) bool
}

type ProcessorFunc func(Context, *http.Response, []byte, *goquery.Document) bool

func (p ProcessorFunc) Process(ctx Context, res *http.Response, body []byte, doc *goquery.Document) bool {
	return p(ctx, body, doc)
}
