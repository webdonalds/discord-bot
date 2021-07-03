package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/webdonalds/discord-bot/responses"

	"github.com/webdonalds/discord-bot/background"
)

type PingCommand struct{}

func NewPingCommand() Command {
	return &PingCommand{}
}

func (*PingCommand) CommandTexts() []string {
	return []string{"ping"}
}

func (*PingCommand) Execute(_ []string, _ *discordgo.MessageCreate) (responses.ResponseMessage, background.Watcher, error) {
	return responses.NewTextMessage("pong"), nil, nil
}
