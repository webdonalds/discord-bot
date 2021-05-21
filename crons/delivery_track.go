package crons

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/avast/retry-go"
	"github.com/hellodhlyn/delivery-tracker"
	log "github.com/sirupsen/logrus"

	"github.com/webdonalds/discord-bot/repositories"
)

type DeliveryTrackCron struct {
	repo          repositories.DeliveryTrackRepository
	trackerClient deliverytracker.Client
}

func NewDeliveryTrackCron(repo repositories.DeliveryTrackRepository, trackClient deliverytracker.Client) Cron {
	return &DeliveryTrackCron{repo: repo, trackerClient: trackClient}
}

func (cron DeliveryTrackCron) Pattern() string {
	return "@every 1m"
}

func (cron DeliveryTrackCron) ChannelID() string {
	return os.Getenv("BOT_CHANNEL_ID")
}

func (cron DeliveryTrackCron) Execute() string {
	tracks, err := cron.repo.ListAllShouldExecute(context.Background())
	if err != nil {
		log.Errorf("failed to get delivery tracks: %v", err)
		return ""
	}

	msg := ""
	for _, trackData := range tracks {
		var track *deliverytracker.Track
		err = retry.Do(func() error {
			track, err = cron.trackerClient.GetTrack(trackData.CarrierID, trackData.TrackID)
			return err
		})
		if err != nil {
			log.Errorf("failed to fetch track info (carrierID: %s, trackID: %s)\n%v", trackData.CarrierID, trackData.TrackID, err)
			continue
		} else if track.State != nil && track.State.ID == "delivered" {
			if trackData.ItemName == "" {
				msg += fmt.Sprintf("%s 배송이 완료되었습니다.\n\n", trackData.Mention)
			} else {
				msg += fmt.Sprintf("%s [%s] 배송이 완료되었습니다.\n\n", trackData.Mention, trackData.ItemName)
			}
			continue
		}

		var lastTimestamp *time.Time
		if len(track.Progresses) > 0 {
			// 새로운 배송 추적 데이터가 있는 경우 메시지 발송
			lastProgress := *track.Progresses[0]
			for _, p := range track.Progresses {
				if p.Time.After(*lastProgress.Time) {
					lastProgress = *p
				}
			}
			lastTimestamp = lastProgress.Time

			if trackData.LastTimestamp == nil || lastProgress.Time.After(*trackData.LastTimestamp) {
				minuteAgo := int64(time.Now().Sub(*lastProgress.Time) / time.Minute)

				msg += fmt.Sprintf("%s\n[배송 정보]\n", trackData.Mention)
				if trackData.ItemName != "" {
					msg += fmt.Sprintf("물품명 : %s\n", trackData.ItemName)
				}
				msg += fmt.Sprintf(
					"운송장 : %s %s\n배송 현황 : %s\n현재 위치 : %s\n\n%s (업데이트 : %d분 전)\n\n",
					track.Carrier.Name, trackData.TrackID, lastProgress.Status.Text, lastProgress.Location.Name, lastProgress.Description, minuteAgo,
				)
			}
		}

		runAt := time.Now().Add(20 * time.Minute)
		err = cron.repo.Append(context.Background(), &repositories.DeliveryTrack{
			Mention: trackData.Mention, CarrierID: trackData.CarrierID, TrackID: trackData.TrackID, ItemName: trackData.ItemName, LastTimestamp: lastTimestamp,
		}, &runAt)
		if err != nil {
			log.Errorf("failed to append track data: %v", err)
		}
	}

	return msg
}
