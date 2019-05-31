package main

import (
	"net/url"
)

func filterVisited(inQueue <- chan *url.URL) <-chan *url.URL {
	var visited = make(map[string]struct{})
	var outChan = make(chan *url.URL)
	go func() {
		for url := range inQueue {
			if _, ok := visited[url.String()]; ok {
				continue
			}
			visited[url.String()] = struct{}{}
			outChan <- url
		}
	}()
	return outChan
}
