package responses

import "github.com/bwmarrin/discordgo"

type SimpleResponse struct {
	Text string
}

func NewSimpleResponse(text string) *SimpleResponse {
	return &SimpleResponse{Text: text}
}

func (r *SimpleResponse) ToDiscordMessage() *discordgo.MessageSend {
	return &discordgo.MessageSend{Content: r.Text}
}
