package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/webdonalds/discord-bot/background"
	"github.com/webdonalds/discord-bot/repositories"
)

const (
	molluHelpMsg = "사용법: !몰루 카페\n\n사용법: !몰루 카페 출석\n\n사용법: !몰루 알림\n\n사용법: !몰루 알림 <on/off>\n예시: !몰루 알림 on"

	notifySettingOn  = "on"
	notifySettingOff = "off"
)

type MolluCommand struct {
	repo repositories.MolluRepository
}

func NewMolluCommand(repo repositories.MolluRepository) Command {
	return &MolluCommand{repo: repo}
}

func (*MolluCommand) CommandTexts() []string {
	return []string{"몰루"}
}

func (cmd *MolluCommand) Execute(args []string, msg *discordgo.MessageCreate) (string, background.Watcher, error) {
	if !cmd.isValidArgs(args) {
		return molluHelpMsg, nil, nil
	}

	ctx := context.Background()
	info, err := cmd.repo.GetByID(ctx, msg.Author.Mention())
	if err == repositories.MolluNotFound {
		info = repositories.MolluInfo{
			ID:            msg.Author.Mention(),
			NotifySetting: false,
			CafeLastVisit: nil,
			LastNotify:    nil,
		}
	} else if err != nil {
		return "", nil, err
	}

	if args[0] == "카페" {
		msg, err := cmd.handleCafeCommand(ctx, info, args[1:])
		return msg, nil, err
	} else if args[0] == "알림" {
		msg, err := cmd.handleNotifyCommand(ctx, info, args[1:])
		return msg, nil, err
	}
	return helpMsg, nil, nil
}

func (cmd *MolluCommand) isValidArgs(args []string) bool {
	if len(args) != 1 && len(args) != 2 {
		return false
	}

	if args[0] == "카페" {
		return len(args) == 1 || args[1] == "출석"
	}
	return args[0] == "알림" && (len(args) == 1 || args[1] == notifySettingOn || args[1] == notifySettingOff)
}

func (cmd *MolluCommand) handleCafeCommand(ctx context.Context, info repositories.MolluInfo, args []string) (string, error) {
	if len(args) == 0 {
		if info.CafeLastVisit == nil {
			return "카페에 출석한 기록이 없습니다.", nil
		} else {
			return fmt.Sprintf("최근 카페 출석 시각: %s", info.CafeLastVisit.Format("2006-01-02 15:04:05")), nil
		}
	}

	currentTime := time.Now()
	info.CafeLastVisit = &currentTime
	info.LastNotify = nil
	err := cmd.repo.Save(ctx, info)
	if err != nil {
		return "", err
	}
	return "정상적으로 출석처리 되었습니다.", nil
}

func (cmd *MolluCommand) handleNotifyCommand(ctx context.Context, info repositories.MolluInfo, args []string) (string, error) {
	reply := ""
	if len(args) != 0 {
		info.NotifySetting = (args[0] == notifySettingOn)
		err := cmd.repo.Save(ctx, info)
		if err != nil {
			return "", err
		}

		reply = "정상적으로 설정되었습니다.\n\n"
	}

	var notifySettingStr string
	if info.NotifySetting {
		notifySettingStr = notifySettingOn
	} else {
		notifySettingStr = notifySettingOff
	}
	reply += "현재 알림 상태: " + notifySettingStr

	return reply, nil
}
