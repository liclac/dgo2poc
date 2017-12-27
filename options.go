package dgo2poc

import (
	"net/http"
)

type ReqOptions struct {
	Request *http.Request
}

// Options can be passed to Client.Request() to customise requests.
type ReqOption func(opts *ReqOptions)

// Set a request's content type to something other than the default "application/json".
func WithContentType(ct string) ReqOption {
	return ReqOption(func(opts *ReqOptions) {
		opts.Request.Header.Set("Content-Type", ct)
	})
}

// Override the default user agent; for bots, this must begin with: "DiscordBot (url, version)",
// additional information may be appended at the end.
func WithUserAgent(ua string) ReqOption {
	return ReqOption(func(opts *ReqOptions) {
		opts.Request.Header.Set("User-Agent", ua)
	})
}
