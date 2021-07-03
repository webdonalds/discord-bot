package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/webdonalds/discord-bot/background"
	"github.com/webdonalds/discord-bot/responses"
)

type Command interface {
	CommandTexts() []string
	Execute(args []string, msg *discordgo.MessageCreate) (responses.ResponseMessage, background.Watcher, error)
}
