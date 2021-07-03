package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/webdonalds/discord-bot/responses"

	"github.com/webdonalds/discord-bot/background"
)

type HelpCommand struct{}

func NewHelpCommand() Command {
	return &HelpCommand{}
}

func (*HelpCommand) CommandTexts() []string {
	return []string{"help", "도움"}
}

func (*HelpCommand) Execute(_ []string, _ *discordgo.MessageCreate) (responses.ResponseMessage, background.Watcher, error) {
	e := responses.NewEmbed().
		SetTitle("Webdo Bot Help").
		SetDescription("`!help` 또는 `!도움`을 입력").
		AddField("미세먼지", "`!미세먼지`").
		AddField("택배", "`!택배 <택배사> <운송장번호>`").
		AddField("환율", "`!환율 <통화쌍> <금액>`").
		AddField("랜덤", "`!랜덤 ...` or `!뽑기 ...`").
		AddField("주가", "`!주가 <종목명>`").
		SetColor(0xc0ffee)

	return responses.NewEmbedMessage(e.MessageEmbed), nil, nil
}
