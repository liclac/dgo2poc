package dgo2poc

import (
	"golang.org/x/oauth2"
)

// OAuth2 Endpoint for authenticating with Discord.
var Endpoint = &oauth2.Endpoint{
	AuthURL:  "https://discordapp.com/api/oauth2/authorize",
	TokenURL: "https://discordapp.com/api/oauth2/token",
}

// Returns a token for authenticating as a regular user.
// For bot accounts, use BotToken().
func UserToken(t string) *oauth2.Token {
	return &oauth2.Token{AccessToken: t, TokenType: "Bearer"}
}

// Returns a token for authenticating as a bot.
// For user accounts, use UserToken().
func BotToken(t string) *oauth2.Token {
	return &oauth2.Token{AccessToken: t, TokenType: "Bot"}
}
