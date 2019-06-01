package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
)

var (
	VERSION string
	BUILD   string
)

func main() {
	var urlsChan = make(chan *url.URL)
	var out = os.Stdout
	startUrl, rps := parseAgrs()
	parsedUrl, err := url.Parse(startUrl)
	if err != nil {
		panic(err)
	}

	// easy could be replaced by fasthttp
	var doer Doer = &http.Client{Transport: NewRateLimitTransport(http.DefaultTransport, rps)}
	var linkSearcher = &LinkSearcher{urlsChan, parsedUrl}
	crawler := createCrawler(doer, BodyHandlers{linkSearcher.GetLinks})

	withVisitFiltered := filterVisited(urlsChan)
	withExportTo := exportFoundedUrl(withVisitFiltered, &LineWriter{out})
	// start crawler process
	done := crawler.StartLoop(withExportTo)
	// feed root url
	urlsChan <- parsedUrl

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-quit:
		log.Println("Ctrl+C intercepted. Shutdown")
		done()
		close(urlsChan)
	}
}

func createCrawler(doer Doer, bodyHandlers BodyHandlers) *Crawler {
	requestConstructor := newDefaultRequestConstructor()
	crawler := Crawler{
		requestConstructor,
		doer,
		&responseHandle{
			validators: ResponseValidators{
				checkResponseCode(http.StatusOK),
				checkResponseContentType("html"),
			},
			bodyHandlers: bodyHandlers,
		}}
	return &crawler
}

func exportFoundedUrl(inQueue <-chan *url.URL, writer io.Writer) <-chan *url.URL {
	var outChan = make(chan *url.URL)
	go func() {
		for url := range inQueue {
			writer.Write([]byte(url.String()))
			outChan <- url
		}
	}()
	return outChan
}

type LineWriter struct {
	io.Writer
}

func (w *LineWriter) Write(p []byte) (n int, err error) {
	return w.Writer.Write(append(p, '\n'))
}
