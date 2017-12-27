package dgo2poc

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClientRequest(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(t, "/test", req.URL.Path)
			assert.Equal(t, "Bot hi", req.Header.Get("Authorization"))
			assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
			assert.Equal(t, UserAgent, req.Header.Get("User-Agent"))

			data, err := ioutil.ReadAll(req.Body)
			assert.NoError(t, err)
			assert.Equal(t, "{}", string(data))
			req.Body.Close()

			_, _ = rw.Write([]byte("hi"))
		}))
		defer srv.Close()
		cl := NewClient(BotToken("hi"))
		data, err := cl.Request(context.Background(), "GET", srv.URL+"/test", []byte("{}"))
		assert.NoError(t, err)
		assert.Equal(t, "hi", string(data))
	})
	t.Run("Opts", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(t, "/test", req.URL.Path)
			assert.Equal(t, "Bot hi", req.Header.Get("Authorization"))
			assert.Equal(t, "text/plain", req.Header.Get("Content-Type"))
			assert.Equal(t, "test user agent", req.Header.Get("User-Agent"))

			data, err := ioutil.ReadAll(req.Body)
			assert.NoError(t, err)
			assert.Equal(t, "hi", string(data))
			req.Body.Close()

			_, _ = rw.Write([]byte("hi"))
		}))
		defer srv.Close()
		cl := NewClient(BotToken("hi"), WithUserAgent("test user agent"))
		data, err := cl.Request(context.Background(), "GET", srv.URL+"/test", []byte("hi"), WithContentType("text/plain"))
		assert.NoError(t, err)
		assert.Equal(t, "hi", string(data))
	})

	t.Run("Error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(400)
			rw.Write([]byte(`{"code":1234,"message":"everything is broken"}`))
		}))
		defer srv.Close()
		cl := NewClient(BotToken("hi"))
		_, err := cl.Request(context.Background(), "GET", srv.URL, nil)
		assert.EqualError(t, err, "400: everything is broken")

		t.Run("Plaintext", func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(400)
				rw.Write([]byte(`aaaa`))
			}))
			defer srv.Close()
			cl := NewClient(BotToken("hi"))
			_, err := cl.Request(context.Background(), "GET", srv.URL, nil)
			assert.EqualError(t, err, "400: aaaa")
		})
	})
}

func TestClientRequestJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte(`{"id":"1234"}`))
	}))
	defer srv.Close()

	var obj struct {
		ID string `json:"id"`
	}
	cl := NewClient(BotToken("hi"))
	assert.NoError(t, cl.RequestJSON(context.Background(), "GET", srv.URL, nil, &obj))

	t.Run("Error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(`{"id":"1234`))
		}))
		defer srv.Close()

		var obj struct {
			ID string `json:"id"`
		}
		cl := NewClient(BotToken("hi"))
		assert.EqualError(t, cl.RequestJSON(context.Background(), "GET", srv.URL, nil, &obj), "unexpected end of JSON input")
	})
}

func TestClientUser(t *testing.T) {
	requireTestToken(t)

	cl := NewClient(TestToken)

	t.Run("@me", func(t *testing.T) {
		user, err := cl.User(context.Background(), "@me")
		require.NoError(t, err)
		require.NotNil(t, user)
		assert.NotEmpty(t, user.ID)
		assert.NotEmpty(t, user.Username)
		assert.NotEmpty(t, user.Discriminator)
		assert.True(t, user.Bot)

		t.Run("By ID", func(t *testing.T) {
			user2, err := cl.User(context.Background(), user.ID)
			require.NoError(t, err)
			require.NotNil(t, user2)
			assert.Equal(t, user.ID, user2.ID)
			assert.Equal(t, user.Username, user2.Username)
			assert.Equal(t, user.Discriminator, user2.Discriminator)
		})
	})
}
