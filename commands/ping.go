package commands

import (
	"github.com/bwmarrin/discordgo"

	"github.com/webdonalds/discord-bot/background"
	"github.com/webdonalds/discord-bot/responses"
)

type PingCommand struct{}

func NewPingCommand() Command {
	return &PingCommand{}
}

func (*PingCommand) CommandTexts() []string {
	return []string{"ping"}
}

func (*PingCommand) ExpectedArgsLen() int {
	return 0
}

func (*PingCommand) Execute(args []string, msg *discordgo.MessageCreate) (responses.Response, background.Watcher, error) {
	return responses.NewSimpleResponse("pong"), nil, nil
}
