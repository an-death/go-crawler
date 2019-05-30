package main

import (
	"net/http"
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
