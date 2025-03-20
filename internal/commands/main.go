package commands

import (
	"context"

	"github.com/webdonalds/discord-bot/internal/actions"
)

// Command 는 명령어를 정의하는 인터페이스입니다.
type Command interface {
	Commands() []string
	Execute(ctx context.Context, args []string) (actions.Action, error)
}
