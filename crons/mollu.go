package crons

import (
	"context"
	"fmt"
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
		if !info.NotifySetting || !shouldNotify(info, currentTime) {
			continue
		}

		info.LastNotify = &currentTime
		molluInfoList[i] = info
		if len(msg) != 0 {
			msg += "\n\n"
		}
		msg += fmt.Sprintf("%s 선생님 카페 방문 하실 시간입니다.", info.ID)
	}

	if len(msg) != 0 {
		err = cron.repo.SaveAll(ctx, molluInfoList)
		if err != nil {
			log.Errorf("failed to save mollu info list")
			return ""
		}
	}
	return msg
}

// 알람을 받은 이후에 12시간 초기화가 되지 않았다면 3시간 뒤에 한번 발송
// 12시간마다 초기화 되는 알림은 항상 한번 발송
func shouldNotify(info repositories.MolluInfo, currentTime time.Time) bool {
	if info.CafeLastVisit == nil {
		return false
	}

	if info.LastNotify == nil && (!isSameInterval(*info.CafeLastVisit, currentTime) || info.CafeLastVisit.Before(currentTime.Add(time.Hour*-3))) {
		return true
	}
	return info.LastNotify != nil && !isSameInterval(*info.LastNotify, currentTime)
}

// 4시간씩 빼서 초기화 시간을 0:00~11:59, 12:00~23:59로 생각합니다.
func isSameInterval(t1, t2 time.Time) bool {
	if t1.After(t2) {
		t1, t2 = t2, t1
	}

	t1 = t1.Add(time.Hour * -4)
	t2 = t2.Add(time.Hour * -4)
	return isSameDate(t1, t2) && (t2.Hour() < 12 || t1.Hour() >= 12)
}

func isSameDate(t1, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
