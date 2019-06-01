package main

import (
	"errors"
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type LinkSearcher struct {
	outChan  chan<- *url.URL
	startUrl *url.URL
}

func (s *LinkSearcher) GetLinks(body io.Reader) error {
	htmlDoc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return err
	}

	absoluteUrlFunc, err := s.getAbsoluteUrlFunc(htmlDoc)
	if err != nil {
		return err
	}

	htmlDoc.Find("a[href]").Each(func(_ int, headSelection *goquery.Selection) {
		headSelection.Each(func(_ int, selection *goquery.Selection) {
			val, exits := selection.Attr("href")
			if !exits {
				return
			}
			if newUrl := absoluteUrlFunc(val); newUrl != nil {
				s.outChan <- newUrl
			}
		})
	})
	return nil
}

func (s *LinkSearcher) getAbsoluteUrlFunc(document *goquery.Document) (func(string) *url.URL, error) {
	if !isBaseFound(document) {
		return s.absoluteUrl, nil
	}

	base, err := getBaseValue(document)
	if err != nil {
		return nil, err
	}

	return s.absoluteUrlWithBase(base), nil
}

func (s *LinkSearcher) absoluteUrl(val string) *url.URL {
	return AbsoluteUrl(s.startUrl, nil, val)
}

func (s *LinkSearcher) absoluteUrlWithBase(base *url.URL) func(string) *url.URL {
	return func(val string) *url.URL {
		return AbsoluteUrl(s.startUrl, base, val)
	}
}

func AbsoluteUrl(prev, base *url.URL, path string) *url.URL {
	if strings.HasPrefix(path, "#") {
		return nil
	}
	if base == nil {
		absURL, _ := prev.Parse(path)
		return absURL
	}

	base = prev.ResolveReference(base)

	absURL, err := base.Parse(path)
	if err != nil {
		return nil
	}
	absURL.Fragment = ""
	return absURL
}

func getBaseValue(document *goquery.Document) (*url.URL, error) {
	base, found := document.Find("base[href]").Attr("href")
	if !found {
		return nil, errors.New("base not found")
	}

	return url.Parse(base)

}

func isBaseFound(document *goquery.Document) bool {
	_, found := document.Find("base[href]").Attr("href")
	return found
}
