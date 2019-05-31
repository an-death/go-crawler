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

func NewRateLimitTransport(transport http.RoundTripper, rps int) *RateLimitTransport {
	var limit rate.Limit
	if rps <= 0 {
		limit = rate.Inf
	} else {
		limit = rate.Limit(rps)
	}

	return &RateLimitTransport{limiter: rate.NewLimiter(limit, int(rps)), RoundTripper: transport}
}

func (t *RateLimitTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for !t.limiter.Allow() {
		// return control
		runtime.Gosched()
	}
	return t.RoundTripper.RoundTrip(req)
}
