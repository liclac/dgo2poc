package dgo2poc

// Do not add anything other than event structs here.
// If you change this file, remember to run `go generate` to generate event handlers.
//go:generate go run tools/gen_events/main.go -in wsevents.go -out wsevents_gen.go

import (
	"github.com/bwmarrin/discordgo"
)

type Ready struct {
	Version   int    `json:"v"`
	SessionID string `json:"session_id"`
}

type GuildCreate struct {
	discordgo.Guild // borrowing this definition for a bit
}
