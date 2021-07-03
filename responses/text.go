package responses

import (
	"github.com/bwmarrin/discordgo"
)

type TextMessage struct {
	message string
}

func (tm *TextMessage) ToDiscordMessage() *discordgo.MessageSend {
	return &discordgo.MessageSend{Content: tm.message}
}

func NewTextMessage(msg string) *TextMessage {
	return &TextMessage{message: msg}
}