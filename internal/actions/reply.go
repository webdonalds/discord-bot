package actions

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type Reply struct {
	Message string
}

func (r *Reply) Execute(ctx context.Context, session *discordgo.Session, message *discordgo.MessageCreate) error {
	_, err := session.ChannelMessageSendReply(message.ChannelID, r.Message, &discordgo.MessageReference{
		MessageID: message.ID,
		ChannelID: message.ChannelID,
		GuildID:   message.GuildID,
	})
	return err
}
