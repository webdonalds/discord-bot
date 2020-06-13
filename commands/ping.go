package commands

import "github.com/bwmarrin/discordgo"

type PingCommand struct{}

func NewPingCommand() Command {
	return &PingCommand{}
}

func (*PingCommand) CommandTexts() []string {
	return []string{"ping"}
}

func (*PingCommand) Execute(_ []string, msgChan chan<- string, _ *discordgo.MessageCreate) error {
	msgChan <- "pong"
	return nil
}
