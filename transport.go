package main

import (
	"golang.org/x/time/rate"
	"net/http"
	"runtime"
)

type RateLimitTransport struct {
	limiter *rate.Limiter
	http.RoundTripper
}

func NewRateLimitTransport(transport http.RoundTripper, rps uint64) *RateLimitTransport {
	return &RateLimitTransport{limiter: rate.NewLimiter(rate.Limit(rps), int(rps)), RoundTripper: transport}
}

func (t *RateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for !t.limiter.Allow() {
		// return control
		runtime.Gosched()
	}
	return t.RoundTripper.RoundTrip(req)
}
