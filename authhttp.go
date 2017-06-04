package authhttp

import (
	"net/http"
	"net"
	"runtime"
	"time"
)

type roundTripFunc func(*http.Request)(*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request)(*http.Response, error){
	return fn(req)
}

// transportOption represents a transport-level option for an http.RoundTripper.
type TransportOption func(http.RoundTripper) http.RoundTripper


func WithBasicAuth(clientID, secret string) TransportOption {
	return func(rt http.RoundTripper) http.RoundTripper {
		return roundTripFunc(func(req *http.Request) (*http.Response, error) {
			req.SetBasicAuth(clientID, secret)
			return rt.RoundTrip(req)
		})
	}
}

func WithHeader(key, value string) TransportOption {
	return func(rt http.RoundTripper) http.RoundTripper {
		return roundTripFunc(func(req *http.Request) (*http.Response, error) {
			req.Header.Set(key, value)
			return rt.RoundTrip(req)
		})
	}
}

// newBaseRoundTripper provides sane defaults for an http.RoundTripper.
func newBaseRoundTripper() http.RoundTripper {
	return &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
	}
}

// newHTTPClient constructs a new HTTP client with the specified transport-level options.
func NewHTTPClient(opts ...TransportOption) *http.Client {
	rt := newBaseRoundTripper()
	for _, opt := range opts {
		rt = opt(rt)
	}
	return &http.Client{Transport: rt}
}

