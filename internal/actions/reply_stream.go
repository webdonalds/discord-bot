package actions

import (
	"context"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ReplyStream struct {
	MessageStream chan string
}

func (r *ReplyStream) Execute(ctx context.Context, session *discordgo.Session, message *discordgo.MessageCreate) error {
	// Start typing indicator
	session.ChannelTyping(message.ChannelID)

	var sentMessage *discordgo.Message
	var err error
	isFirstMessage := true
	var messageBuffer strings.Builder

	// Process messages from the stream until it's closed
	for msg := range r.MessageStream {
		session.ChannelTyping(message.ChannelID)

		if isFirstMessage {
			if msg == "" {
				continue
			}

			// First message: Reply to the original message
			sentMessage, err = session.ChannelMessageSendReply(message.ChannelID, msg, &discordgo.MessageReference{
				MessageID: message.ID,
				ChannelID: message.ChannelID,
				GuildID:   message.GuildID,
			})
			if err != nil {
				return err
			}
			messageBuffer.WriteString(msg)
			isFirstMessage = false
		} else {
			// Subsequent messages: Edit the existing message by appending new content
			if sentMessage != nil {
				messageBuffer.WriteString(msg)
				sentMessage, err = session.ChannelMessageEdit(message.ChannelID, sentMessage.ID, messageBuffer.String())
				if err != nil {
					return err
				}
			}
		}
	}

	// Add check mark reaction after streaming is complete
	if sentMessage != nil {
		_ = session.MessageReactionAdd(message.ChannelID, sentMessage.ID, ":kazusaThumbsUp:1138625229229400096")
	}

	return nil
}
