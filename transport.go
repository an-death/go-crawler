package main

import (
	"golang.org/x/time/rate"
	"net/http"
	"runtime"
)

type RateLimitedTransport struct {
	limiter *rate.Limiter
	http.RoundTripper
}

func NewRateLimitTransport(transport http.RoundTripper, rps int) *RateLimitedTransport {
	var limit rate.Limit
	if rps <= 0 {
		limit = rate.Inf
	} else {
		limit = rate.Limit(rps)
	}

	return &RateLimitedTransport{limiter: rate.NewLimiter(limit, rps), RoundTripper: transport}
}

func (t *RateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.wait()
	return t.RoundTripper.RoundTrip(req)
}

func (t *RateLimitedTransport) wait() {
	for !t.limiter.Allow() {
		// return control
		runtime.Gosched()
	}
}
