package authhttp

import (
	"net/http"
	"net"
	"runtime"
	"time"
)

// CustomRoundTripper provides custom configuration of *http.Client(s) for specific providers. Because each provider may require
// various amounts of extensibility this is abstracted at the provider level.
type CustomRoundTripper struct {
	http.Transport
	Username string
	Password string
}

func (rt *CustomRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(rt.Username, rt.Password)
	return rt.Transport.RoundTrip(req)
}

func getHTTPTransport() *http.Transport {
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

func WithBasicAuth(username string, password string) *http.Client {
	//transport := cleanhttp.DefaultPooledTransport()
	return &http.Client {
		Transport: &CustomRoundTripper{
			Transport: *getHTTPTransport(),
			Username: username,
			Password: password,
		},
	}
}