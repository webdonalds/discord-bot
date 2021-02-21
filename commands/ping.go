package commands

import (
	"github.com/bwmarrin/discordgo"

	"github.com/webdonalds/discord-bot/background"
)

type PingCommand struct{}

func NewPingCommand() Command {
	return &PingCommand{}
}

func (*PingCommand) CommandTexts() []string {
	return []string{"ping"}
}

func (*PingCommand) ExpectedArgsLen() (int, int) {
	return 0, 0
}

func (*PingCommand) Execute(_ []string, _ *discordgo.MessageCreate) (string, background.Watcher, error) {
	return "pong", nil, nil
}
