package commands

import (
	"context"

	"github.com/webdonalds/discord-bot/internal/actions"
)

type PingCommand struct {
}

func NewPingCommand() *PingCommand {
	return &PingCommand{}
}

func (c *PingCommand) Commands() []string {
	return []string{"ping"}
}

func (c *PingCommand) Execute(ctx context.Context, args *CommandArgs) (actions.Action, error) {
	return &actions.Reply{Message: "pong"}, nil
}
