[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/b3ntly/go-authhttp/master/LICENSE.txt)
[![Build Status](https://travis-ci.org/b3ntly/go-authhttp.svg?branch=master)](https://travis-ci.org/b3ntly/go-authhttp) 
[![Coverage Status](https://coveralls.io/repos/github/b3ntly/go-authhttp/badge.svg?branch=master)](https://coveralls.io/github/b3ntly/go-authhttp?branch=master) 
[![GoDoc](https://godoc.org/github.com/b3ntly/go-authhttp?status.svg)](https://godoc.org/github.com/b3ntly/go-authhttp)

### Authhttp

Provides a Golang *http.Client with sane defaults that will reuse a basic auth header
in every request.

It may also be used to append arbitrary headers to all subsequent requests.

### Usage

```go 
package authhttp_test

import (
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/b3ntly/go-authhttp"
	"net/http"
	"bytes"
	"io/ioutil"
	"fmt"
	"encoding/json"
)

type HTTPBinBasicAuthQuery struct {
	Authenticated bool `json:"authenticated"`
	User string `json:"user"`
}

type HTTBinHeadersQuery struct {
	Headers struct {
		Accept string `json:"Accept"`
		AcceptEncoding string `json:"Accept-Encoding"`
		AcceptLanguage string `json:"Accept-Language"`
		Connection string `json:"Connection"`
		Cookie string `json:"Cookie"`
		Host string `json:"Host"`
		Referer string `json:"Referer"`
		UpgradeInsecureRequests string `json:"Upgrade-Insecure-Requests"`
		UserAgent string `json:"User-Agent"`
	} `json:"headers"`
}

func TestWithBasicAuth(t *testing.T) {
	t.Run("WithBasicAuth will return a client whose future requests satisfy the Basic Auth spec", func(t *testing.T){
		const user = "Benjamin"
		const password = "Jones"

		client := authhttp.NewHTTPClient(authhttp.WithBasicAuth(user, password))
		req, err := http.NewRequest("GET", fmt.Sprintf("https://httpbin.org/basic-auth/%v/%v", user, password), bytes.NewReader([]byte("hello")))

		require.Nil(t, err)

		resp, err := client.Do(req)

		require.Nil(t, err)

		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)

		require.Nil(t, err)

		response := &HTTPBinBasicAuthQuery{}
		err = json.Unmarshal(contents, response)

		require.Nil(t, err)
		require.Equal(t, true, response.Authenticated)
		require.Equal(t, user, response.User)
	})

	t.Run("WithBasicAuth with improper parameters will not satisfy the Basic Auth spec", func(t *testing.T){
		const user = "ben"
		const password = "jones"
		const wrongPassword = "wrong"

		client := authhttp.NewHTTPClient(authhttp.WithBasicAuth(user, wrongPassword))
		req, err := http.NewRequest("GET", fmt.Sprintf("https://httpbin.org/basic-auth/%v/%v", user, password), bytes.NewReader([]byte("hello")))

		require.Nil(t, err)

		resp, err := client.Do(req)

		require.Nil(t, err)

		require.Equal(t, 401, resp.StatusCode)
	})
}

func TestWithHeader(t *testing.T) {
	t.Run("WithHeader will send requests with the desired headers", func(t *testing.T){
		const headerKey = "Accept"
		const headerValue = "application/json"

		client := authhttp.NewHTTPClient(authhttp.WithHeader(headerKey, headerValue))
		req, err := http.NewRequest("GET", "https://httpbin.org/headers", bytes.NewReader([]byte("")))

		require.Nil(t, err)

		resp, err := client.Do(req)

		require.Nil(t, err)
		require.Equal(t, 200, resp.StatusCode)
		defer resp.Body.Close()

		contents, err := ioutil.ReadAll(resp.Body)
		require.Nil(t, err)

		response := &HTTBinHeadersQuery{}
		err = json.Unmarshal(contents, response)

		require.Equal(t, headerValue, response.Headers.Accept)
	})
}
```

## Contributions

Significant contributions from Andrei Tudor Călin from gophers slack channel.