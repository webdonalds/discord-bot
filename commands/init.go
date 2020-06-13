package commands

import "github.com/bwmarrin/discordgo"

type Command interface {
	CommandTexts() []string
	Execute(args []string, msgChan chan<- string, msg *discordgo.MessageCreate) error
}
