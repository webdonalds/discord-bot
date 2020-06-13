package commands

import "github.com/bwmarrin/discordgo"

type PingCommand struct{}

func NewPingCommand() Command {
	return &PingCommand{}
}

func (*PingCommand) CommandTexts() []string {
	return []string{"ping"}
}

func (*PingCommand) Execute(_ []string, _ *discordgo.MessageCreate) (string, error) {
	return "pong", nil
}
