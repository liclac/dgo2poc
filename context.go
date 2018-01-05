package dgo2poc

import (
	"context"
)

type ctxKey string

const (
	ctxKeyClient   ctxKey = "client"
	ctxKeyWSClient ctxKey = "wsclient"
)

// Returns the Client for a context. Returns nil if used outside of a handler function.
func GetClient(ctx context.Context) Client {
	cl, _ := ctx.Value(ctxKeyClient).(Client)
	return cl
}

// Adds a Client to a context.
func withClient(ctx context.Context, cl Client) context.Context {
	return context.WithValue(ctx, ctxKeyClient, cl)
}

// Returns the WSClient for a context. Returns nil if used outside of a handler function.
func GetWSClient(ctx context.Context) WSClient {
	ws, _ := ctx.Value(ctxKeyWSClient).(WSClient)
	return ws
}

// Adds a WSClient to a context.
func withWSClient(ctx context.Context, ws WSClient) context.Context {
	return context.WithValue(ctx, ctxKeyWSClient, ws)
}
