package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/webdonalds/discord-bot/background"
)

type Command interface {
	CommandTexts() []string
	Execute(args []string, msg *discordgo.MessageCreate) (string, background.Watcher, error)
}
