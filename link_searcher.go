package main

import (
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
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

	var base *url.URL
	if href, found := htmlDoc.Find("base[href]").Attr("href"); found {
		base, _ = url.Parse(href)
	}

	htmlDoc.Find("a[href]").Each(func(_ int, selection *goquery.Selection) {
		for _, n := range selection.Nodes {
			for a := range filterKey("href", n.Attr) {
				if a.Val == "#" {
					continue
				}

				if newUrl := AbsoluteUrl(s.startUrl, base, a.Val); newUrl != nil {
					s.outChan <- newUrl
				}
			}
		}
	})
	return nil
}

func filterKey(key string, in []html.Attribute) <-chan html.Attribute {
	var out = make(chan html.Attribute)
	go func() {
		for _, a := range in {
			if a.Key == key {
				out <- a
			}
		}
		close(out)
	}()
	return out
}

func AbsoluteUrl(prev, base *url.URL, path string) *url.URL {
	if strings.HasPrefix(path, "#") {
		return nil
	}
	if base == nil {
		absURL,_ := prev.Parse(path)
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
