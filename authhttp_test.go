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

	time.Sleep(time.Second * 2)
}