package responses

import "github.com/bwmarrin/discordgo"

type ResponseMessage interface {
	ToDiscordMessage() *discordgo.MessageSend
}

