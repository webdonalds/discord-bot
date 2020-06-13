package commands

import "github.com/bwmarrin/discordgo"

type Command interface {
	CommandTexts() []string
	Execute(args []string, msg *discordgo.MessageCreate) (string, error)
}
