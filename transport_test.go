package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

type mockRateLimitRoundTrip struct {
	first  time.Time
	second time.Time
	called bool
}

func (m *mockRateLimitRoundTrip) RoundTrip(*http.Request) (*http.Response, error) {
	if !m.called {
		m.called = true
		m.first = time.Now()
		return nil, nil
	}
	m.second = time.Now()
	return nil, nil
}

func TestRateLimitTransport_RoundTrip(t *testing.T) {
	var requestPerSec int = 1
	mockTransport := mockRateLimitRoundTrip{}
	limitedTransport := NewRateLimitTransport(&mockTransport, requestPerSec)
	limitedTransport.RoundTrip(nil)
	limitedTransport.RoundTrip(nil)
	delta := 5 * time.Millisecond // in case of -race request a lot of time
	assert.InDelta(t, 1*time.Second, mockTransport.second.Sub(mockTransport.first), float64(delta))
}
