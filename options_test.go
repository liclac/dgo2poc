package dgo2poc

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithContentType(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/", nil)
	assert.NoError(t, err)
	WithContentType("text/plain;charset=utf-8")(&ReqOptions{Request: req})
	assert.Equal(t, "text/plain;charset=utf-8", req.Header.Get("Content-Type"))
}

func TestWithUserAgent(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/", nil)
	assert.NoError(t, err)
	WithUserAgent("test user agent")(&ReqOptions{Request: req})
	assert.Equal(t, "test user agent", req.Header.Get("User-Agent"))
}
