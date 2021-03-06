package commands

import (
	"fmt"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/webdonalds/discord-bot/background"
)

const timerHelpMsg = "사용법: !타이머 <시간> <메시지>\n예시: !타이머 9h 퇴근시간"

type TimerCommand struct {
	re *regexp.Regexp
}

func NewTimerCommand() Command {
	return &TimerCommand{
		re: regexp.MustCompile(`(?P<num>\d+)(?P<unit>[hmsHMS])`),
	}
}

func (*TimerCommand) CommandTexts() []string {
	return []string{"timer", "타이머"}
}

func (c *TimerCommand) Execute(args []string, _ *discordgo.MessageCreate) (string, background.Watcher, error) {
	if len(args) != 2 {
		return timerHelpMsg, nil, nil
	}

	duration, err := time.ParseDuration(args[0])
	if err != nil {
		return "시간 패턴을 파싱할 수 없습니다.", nil, nil
	}

	kst, _ := time.LoadLocation("Asia/Seoul")
	runAt := time.Now().Add(duration)
	msg := fmt.Sprintf("%s에 알림이 발송됩니다.", runAt.In(kst).Format("2006-01-02 15:04:05"))
	watcher := background.NewTimerWatcher(&runAt, args[1])
	return msg, watcher, nil
}
