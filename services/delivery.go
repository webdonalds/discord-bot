package services

import (
	"fmt"
	"time"

	"github.com/hellodhlyn/delivery-tracker"
)

type DeliveryService struct {
	client     deliverytracker.Client
	followList []*followData
}

type followData struct {
	carrierName string
	trackID     string
	msgChan     chan<- string
	runAt       *time.Time
	timestamp   *time.Time
}

func NewDeliveryFollower() *DeliveryService {
	client, _ := deliverytracker.NewClient()
	return &DeliveryService{
		client:     client,
		followList: []*followData{},
	}
}

func (svc *DeliveryService) Start() {
	go func() {
		for range time.Tick(1 * time.Minute) {
			svc.runCron()
		}
	}()
}

func (svc *DeliveryService) Track(carrierName, trackID string, msgChan chan<- string, timestamp *time.Time) error {
	carriers, err := svc.client.FindCarriersByName(carrierName)
	if err != nil {
		return err
	}
	if len(carriers) == 0 {
		msgChan <- "해당하는 이름의 택배사를 찾을 수 없습니다."
		return nil
	}

	track, err := svc.client.GetTrack(carriers[0].ID, trackID)
	if err != nil {
		return err
	} else if track == nil || track.State == nil || len(track.Progresses) == 0 {
		msgChan <- "올바르지 않은 운송장이거나, 택배사에서 아직 물건을 인수하지 않았습니다."
		return nil
	} else if track.State.ID == "delivered" {
		msgChan <- "배송이 완료되었습니다."
		return nil
	}

	progress := track.Progresses[len(track.Progresses)-1]
	if timestamp == nil || progress.Time.After(*timestamp) {
		timeAgo := int64(time.Now().Sub(*progress.Time) / time.Minute)
		msgChan <- fmt.Sprintf("[배송 정보]\n운송장 : %s %s\n배송 현황 : %s\n현재 위치 : %s\n\n%s (업데이트 : %d분 전)",
			carriers[0].Name, trackID,
			progress.Status.Text, progress.Location.Name, progress.Description, timeAgo)
	}

	runAt := time.Now().Add(30 * time.Minute)
	svc.follow(&followData{
		carrierName: carrierName,
		trackID:     trackID,
		msgChan:     msgChan,
		runAt:       &runAt,
		timestamp:   progress.Time,
	})
	return nil
}

func (svc *DeliveryService) follow(data *followData) {
	svc.followList = append(svc.followList, data)
}

func (svc *DeliveryService) runCron() {
	newList := []*followData{}
	for _, data := range svc.followList {
		_ = svc.Track(data.carrierName, data.trackID, data.msgChan, nil)
	}

	svc.followList = newList
}
