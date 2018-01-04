package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/liclac/dgo2poc"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <token>\n", os.Args[0])
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cl := dgo2poc.NewClient(dgo2poc.BotToken(os.Args[1]))
	u, err := cl.User(ctx, "@me")
	if err != nil {
		log.Fatalf("Couldn't get @me: %s\n", err)
	}
	log.Printf("Authenticated as: %s#%s\n", u.Username, u.Discriminator)

	log.Printf("Connecting to WS gateway...")
	ws := dgo2poc.NewWSClient(cl)
	ws.AddHandler(dgo2poc.OnReady(func(ctx context.Context, r *dgo2poc.Ready) {
		log.Printf("Connected! v%d\n", r.Version)
		cancel() // cancel the context to disconnect
	}))
	if err := ws.Run(ctx); err != nil {
		log.Fatalf("websocket error: %s\n", err)
	}
}
