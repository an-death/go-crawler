package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

func ExampleMainRun() {
	var testPost = 63666
	stopFileServer := startTestFileServer(testPost)
	var urlsChan = make(chan *url.URL)
	parsedUrl, _ := url.Parse(fmt.Sprintf("http://localhost:%v", testPost))
	searcher := &LinkSearcher{urlsChan, parsedUrl}
	crawler := createCrawler(&http.Client{},[]func(io.Reader)error {searcher.GetLinks})
	withVisitFiltered := filterVisited(urlsChan)
	withExportTo := exportFoundedUrl(withVisitFiltered, &LineWriter{os.Stdout})
	done := crawler.StartLoop(withExportTo)
	urlsChan <- parsedUrl
	time.Sleep(100*time.Millisecond)
	defer func() {
		done()
		close(urlsChan)
		stopFileServer()
	}()
	// Output:
	// http://localhost:63666
	// http://localhost:63666/first.html
	// http://localhost:63666/index.html
	// http://localhost:63666/second.html
	// http://localhost:63666/base/third.html
	// http://localhost:63666/base/1.html
	// http://localhost:63666/base/2.html
	// http://localhost:63666/base/3.html
	// http://localhost:63666/base/4.html
	// http://localhost:63666/base/5.html
	// http://localhost:63666/base/6.html
	// http://localhost:63666/base/7.html
	// http://localhost:63666/base/8.html
	// http://localhost:63666/base/9.html
	// http://localhost:63666/base/10.html
}

func startTestFileServer(port int) func() {
	fs := http.FileServer(http.Dir("test_html"))
	server := http.Server{Addr:fmt.Sprintf(":%v", port), Handler:fs}
	go server.ListenAndServe()
	return func() {
		server.Close()
	}
}