package dgo2poc

import (
	"github.com/bwmarrin/discordgo"
)

// ...Message type definition would go here...

// Options for Client.ChannelMessageSend().
type SendOpt func(send *discordgo.MessageSend)

// Attach a file with a message.
func SendWithEmbed(embed *discordgo.MessageEmbed) SendOpt {
	return SendOpt(func(send *discordgo.MessageSend) {
		send.Embed = embed
	})
}
