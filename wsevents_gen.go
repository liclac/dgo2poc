package dgo2poc

// -- THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT. --
// Instead, edit wsevents.go and/or tools/gen_events/main.go and re-run 'go generate'.

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/pkg/errors"
)

type wsHandlers struct {
	GuildCreate     []*func(ctx context.Context, ev *GuildCreate)
	GuildCreateLock sync.RWMutex

	Ready     []*func(ctx context.Context, ev *Ready)
	ReadyLock sync.RWMutex
}

type wsHandler func(hls *wsHandlers) func()

func (hls *wsHandlers) Dispatch(ctx context.Context, t string, data []byte) error {
	switch t {
	case "GUILD_CREATE":
		var ev GuildCreate
		if err := json.Unmarshal(data, &ev); err != nil {
			return errors.Wrap(err, t)
		}
		hls.DispatchGuildCreate(ctx, &ev)
	case "READY":
		var ev Ready
		if err := json.Unmarshal(data, &ev); err != nil {
			return errors.Wrap(err, t)
		}
		hls.DispatchReady(ctx, &ev)
	}
	return nil
}

func (hls *wsHandlers) DispatchGuildCreate(ctx context.Context, ev *GuildCreate) {
	hls.GuildCreateLock.RLock()
	fns := hls.GuildCreate
	hls.GuildCreateLock.RUnlock()
	for _, ptr := range fns {
		fn := *ptr
		go fn(ctx, ev)
	}
}

func (hls *wsHandlers) DispatchReady(ctx context.Context, ev *Ready) {
	hls.ReadyLock.RLock()
	fns := hls.Ready
	hls.ReadyLock.RUnlock()
	for _, ptr := range fns {
		fn := *ptr
		go fn(ctx, ev)
	}
}

// Handle a GuildCreate event. See WSClient.AddHandler().
func OnGuildCreate(fn func(ctx context.Context, ev *GuildCreate)) wsHandler {
	return wsHandler(func(hls *wsHandlers) func() {
		hls.GuildCreateLock.Lock()
		hls.GuildCreate = append(hls.GuildCreate, &fn)
		hls.GuildCreateLock.Unlock()
		return func() {
			hls.GuildCreateLock.Lock()
			for i, v := range hls.GuildCreate {
				if v == &fn {
					hls.GuildCreate = append(hls.GuildCreate[:i], hls.GuildCreate[i+1:]...)
				}
			}
			hls.GuildCreateLock.Unlock()
		}
	})
}

// Handle a Ready event. See WSClient.AddHandler().
func OnReady(fn func(ctx context.Context, ev *Ready)) wsHandler {
	return wsHandler(func(hls *wsHandlers) func() {
		hls.ReadyLock.Lock()
		hls.Ready = append(hls.Ready, &fn)
		hls.ReadyLock.Unlock()
		return func() {
			hls.ReadyLock.Lock()
			for i, v := range hls.Ready {
				if v == &fn {
					hls.Ready = append(hls.Ready[:i], hls.Ready[i+1:]...)
				}
			}
			hls.ReadyLock.Unlock()
		}
	})
}
