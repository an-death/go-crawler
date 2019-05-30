package main

import (
	"golang.org/x/time/rate"
	"net/http"
	"time"
)

type RequestError string

func (err RequestError) Error() string {
	return string(err)
}

const OutOfLimit RequestError = "Request rate out of limit"

type RateLimitTransport struct {
	limiter *rate.Limiter
	http.RoundTripper
}

func NewRateLimitTransport(transport http.RoundTripper, rps uint64) *RateLimitTransport {
	return &RateLimitTransport{limiter: rate.NewLimiter(rate.Limit(rps), int(rps*2)), RoundTripper: transport}
}

func (t *RateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for !t.limiter.Allow() {
		time.Sleep(1 * time.Second)
	}
	return t.RoundTripper.RoundTrip(req)
}
