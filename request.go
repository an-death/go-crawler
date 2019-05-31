package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	DefaultMethod = "GET"
	DefaultProto  = "HTTP/1.1"
)

var DefaultHeaders = http.Header{"User-Agent": []string{"Test user agent"}, "Accept": []string{"*/*"}}

type requestConstructor struct {
	requestMethod  string
	requestProto   string
	requestHeaders http.Header
}

func newDefaultRequestConstructor() *requestConstructor {
	return &requestConstructor{
		requestMethod:  DefaultMethod,
		requestProto:   DefaultProto,
		requestHeaders: DefaultHeaders,
	}
}

func (c *requestConstructor) CreateRequest(rawUlr string) (*http.Request, error) {
	parsed, err := url.Parse(rawUlr)
	if err != nil {
		return nil, err
	}

	return prepareRequest(parsed), nil
}

func prepareRequest(parsed *url.URL) *http.Request {
	return &http.Request{
		Method:     DefaultMethod,
		URL:        parsed,
		Proto:      DefaultProto,
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     DefaultHeaders,
		Body:       ioutil.NopCloser(&bytes.Buffer{}),
		Host:       parsed.Host,
	}
}
