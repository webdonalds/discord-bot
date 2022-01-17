package crons

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/webdonalds/discord-bot/repositories"
)

type MolluCron struct {
	repo repositories.MolluRepository
}

func NewMolluCron(repo repositories.MolluRepository) Cron {
	return &MolluCron{repo: repo}
}

func (cron MolluCron) Pattern() string {
	return "@every 1m"
}

func (cron MolluCron) ChannelID() string {
	return os.Getenv("MOLLU_CHANNEL_ID")
}

func (cron MolluCron) Execute() string {
	ctx := context.Background()

	molluInfoList, err := cron.repo.ListAll(ctx)
	if err != nil {
		log.Errorf("failed to get mollu info list: ", err)
		return ""
	}

	msg := ""
	currentTime := time.Now()
	for i, info := range molluInfoList {
		if !info.NotifySetting || !cron.isCafeInitialized(info, currentTime) {
			continue
		}

		info.IsNotified = true
		info.CafeLastVisit = &currentTime
		molluInfoList[i] = info
	}

	err = cron.repo.SaveAll(ctx, molluInfoList)
	if err != nil {
		log.Errorf("failed to save mollu info list")
		return ""
	}
	return msg
}

func (cron MolluCron) isCafeInitialized(info repositories.MolluInfo, currentTime time.Time) bool {
	if info.CafeLastVisit == nil {
		return false
	}

	prevTime := *info.CafeLastVisit
	if currentTime.After(prevTime.Add(time.Hour * 3)) {
		return true
	}

	// 4시간씩 빼서 초기화 시간을 0~12시, 12시~24시로 생각합니다.
	// 이미 위의 if문 때문에 두 시간 사이는 3시간 이하입니다.
	// 초기화 될 수 없는 조건은 두 시간 모두 0~12시에 있거나 모두 12~24시에 있는 경우입니다.
	t1 := prevTime.Add(time.Hour * -4)
	t2 := currentTime.Add(time.Hour * -4)
	return !(isSameDay(t1, t2) && (t2.Hour() < 12 || t1.Hour() >= 12))
}

func isSameDay(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
