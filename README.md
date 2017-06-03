[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/b3ntly/go-authhttp/master/LICENSE.txt)
[![Build Status](https://travis-ci.org/b3ntly/go-authhttp.svg?branch=master)](https://travis-ci.org/b3ntly/go-authhttp) 
[![Coverage Status](https://coveralls.io/repos/github/b3ntly/go-authhttp/badge.svg?branch=master)](https://coveralls.io/github/b3ntly/go-authhttp?branch=master) 
[![GoDoc](https://godoc.org/github.com/b3ntly/go-authhttp?status.svg)](https://godoc.org/github.com/b3ntly/go-authhttp)

### Authhttp

Provides a Golang *http.Client with sane defaults that will reuse a basic auth header
in every request.

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
	"time"
	"encoding/json"
)

type HTTPBinResponse struct {
	Authenticated bool `json:"authenticated"`
	User string `json:"user"`
}

func TestWithBasicAuth(t *testing.T) {
	t.Run("WithBasicAuth will return a client whose future requests satisfy the Basic Auth spec", func(t *testing.T){
		const user = "Benjamin"
		const password = "Jones"

		client := authhttp.WithBasicAuth(user, password)
		req, err := http.NewRequest("GET", fmt.Sprintf("https://httpbin.org/basic-auth/%v/%v", user, password), bytes.NewReader([]byte("hello")))

		if err != nil {
			require.Nil(t, err)
		}

		resp, err := client.Do(req)

		if err != nil {
			require.Nil(t, err)
		}

		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			require.Nil(t, err)
		}

		response := &HTTPBinResponse{}
		err = json.Unmarshal(contents, response)

		require.Nil(t, err)
		require.Equal(t, true, response.Authenticated)
		require.Equal(t, user, response.User)
	})

	t.Run("WithBasicAuth with improper parameters will not satisfy the Basic Auth spec", func(t *testing.T){
		const user = "ben"
		const password = "jones"
		const wrongPassword = "wrong"

		client := authhttp.WithBasicAuth(user, wrongPassword)
		req, err := http.NewRequest("GET", fmt.Sprintf("https://httpbin.org/basic-auth/%v/%v", user, password), bytes.NewReader([]byte("hello")))

		if err != nil {
			require.Nil(t, err)
		}

		resp, err := client.Do(req)

		if err != nil {
			require.Nil(t, err)
		}

		require.Equal(t, 401, resp.StatusCode)
	})
}
```
