package commands

import (
	"github.com/bwmarrin/discordgo"

	"github.com/webdonalds/discord-bot/background"
	"github.com/webdonalds/discord-bot/responses"
)

type HelpCommand struct{}

const helpMsg = `
!help
!도움

!택배 <택배사> <운송장번호>

!타이머 <시간> <메시지>
!timer <시간> <메시지>

!환율 <통화쌍> <금액>
`

func NewHelpCommand() Command {
	return &HelpCommand{}
}

func (*HelpCommand) CommandTexts() []string {
	return []string{"help", "도움"}
}

func (*HelpCommand) ExpectedArgsLen() int {
	return 0
}

func (*HelpCommand) Execute(_ []string, _ *discordgo.MessageCreate) (responses.Response, background.Watcher, error) {
	return responses.NewSimpleResponse(helpMsg), nil, nil
}
