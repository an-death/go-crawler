package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

func GetEnvOrDefault(envName, bakoff string) string {
	if value, ok := os.LookupEnv(envName); ok {
		return value
	}
	return bakoff
}

func parseAgrs() (string, uint64) {
	var startUrlStr string
	var rps uint64
	flag.Uint64Var(&rps, "rps", 10, "requests per second limit")
	flag.StringVar(&startUrlStr, "url", "", "define url for crawler")
	flag.Parse()
	if startUrlStr == "" {
		panic("No host available")
	}
	return startUrlStr, rps
}

func isHTMLContentType(resp *http.Response) bool {
	return strings.Contains(strings.ToLower(resp.Header.Get("Content-Type")), "html")
}

func AbsoluteUrl(base, prev *url.URL, u string) string {
	if strings.HasPrefix(u, "#") {
		return ""
	}
	if base == nil {
		base = prev
	}
	if !base.IsAbs() {
		base.Host = prev.Host
	}
	absURL, err := base.Parse(u)
	if err != nil {
		return ""
	}
	absURL.Fragment = ""
	if absURL.Scheme == "//" || absURL.Scheme == "" {
		absURL.Scheme = prev.Scheme
	}
	return absURL.String()
}

func main() {
	var urlsChan = make(chan string, 2)
	var out = os.Stdout
	startUrl, rps := parseAgrs()
	parsedUrl, err := checkUrl(startUrl)
	if err != nil {
		panic(err)
	}

	client := http.Client{Transport: NewRateLimitTransport(http.DefaultTransport, rps)}
	crawler := Crawler{checkUrl, prepareRequest, client.Do, checkResponse, getLinks(parsedUrl, urlsChan)}
	err = crawler(startUrl)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(out, startUrl)
	var visited = make(map[string]struct{})
	var group sync.WaitGroup
	for newUrl := range urlsChan {
		if _, ok := visited[newUrl]; ok {
			continue
		}
		visited[newUrl] = struct{}{}
		go func(newUrl string) {
			group.Add(1)
			err := crawler(newUrl)
			if err != nil {
				fmt.Fprintf(out, "ERROR: while request \"%s\" %s\n", newUrl, err)
			} else {
				fmt.Fprintln(out, newUrl)
			}
			group.Done()
		}(newUrl)
	}
}

func prepareRequest(parsed *url.URL) *http.Request {
	return &http.Request{
		Method:     "GET",
		URL:        parsed,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     http.Header{"User-Agent": []string{"Test user agent"}, "Accept": []string{"*/*"}},
		Body:       ioutil.NopCloser(&bytes.Buffer{}),
		Host:       parsed.Host,
	}
}

func checkUrl(rawUrl string) (*url.URL, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return nil, err
	}
	return parsedUrl, nil
}
func checkResponse(response *http.Response) error {
	if response.StatusCode != 200 || !isHTMLContentType(response) {
		return errors.New("not valid response")
	}
	return nil
}

func getLinks(startUrl *url.URL, out chan string) func(io.Reader) {
	return func(body io.Reader) {

		var base *url.URL
		htmlDoc, err := goquery.NewDocumentFromReader(body)
		if err != nil {
			return
		}
		if href, found := htmlDoc.Find("base[href]").Attr("href"); found {
			base, _ = url.Parse(href)
		}
		htmlDoc.Find("a[href]").Each(func(_ int, selection *goquery.Selection) {
			for _, n := range selection.Nodes {
				for _, a := range n.Attr {
					if a.Key == "href" {
						newUrl := AbsoluteUrl(base, startUrl, a.Val)
						if newUrl != "" {
							out <- newUrl
						}
					}
				}
			}
		})
	}
}
