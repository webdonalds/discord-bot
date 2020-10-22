package responses

import "github.com/bwmarrin/discordgo"

type Response interface {
	ToDiscordMessage() *discordgo.MessageSend
}
