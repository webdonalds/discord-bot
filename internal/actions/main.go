package actions

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

// Action 는 명령어 실행 결과를 정의하는 인터페이스입니다.
type Action interface {
	Execute(ctx context.Context, session *discordgo.Session, message *discordgo.MessageCreate) error
}
