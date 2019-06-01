package main

import (
	"net/http"
	"net/url"
	"sync"
)

type RequestConstructor interface {
	CreateRequest(rawUlr string) (*http.Request, error)
}

type Doer interface {
	Do(r *http.Request) (*http.Response, error)
}

type ResponseHandler interface {
	HandleResponse(response *http.Response) error
}

type Crawler struct {
	RequestConstructor
	Doer
	ResponseHandler
}

func (c *Crawler) Fetch(rawUrl string) error {

	req, err := c.RequestConstructor.CreateRequest(rawUrl)
	if err != nil {
		return err

	}
	resp, err := c.Doer.Do(req)
	if err != nil {
		return err
	}
	return c.ResponseHandler.HandleResponse(resp)
}

func (c *Crawler) StartLoop(inQueue <-chan *url.URL) (done func()) {
	var group sync.WaitGroup
	doneChan := make(chan struct{})
	done = func() {
		doneChan <- struct{}{}
		<-doneChan
	}
	go func() {
	Loop:
		for {
			select {
			case newUrl := <-inQueue:
				c.asyncFetch(&group, newUrl)
			case <-doneChan:
				break Loop
			}
		}

		group.Wait()
		doneChan <- struct{}{}
	}()
	return done
}

func (c Crawler) asyncFetch(group *sync.WaitGroup, url *url.URL) {
	group.Add(1)
	go func() {
		c.Fetch(url.String())
		group.Done()
	}()
}
