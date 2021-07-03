package responses

import (
	"github.com/bwmarrin/discordgo"
)

type EmbedMessage struct {
	embed *discordgo.MessageEmbed
}

func (em *EmbedMessage) ToDiscordMessage() *discordgo.MessageSend {
	return &discordgo.MessageSend{Embed: em.embed}
}

func NewEmbedMessage(e *discordgo.MessageEmbed) *EmbedMessage {
	return &EmbedMessage{embed: e}
}