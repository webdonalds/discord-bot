package background

import (
	"fmt"
	"time"

	"github.com/hellodhlyn/delivery-tracker"
)

type Watcher interface {
	ShouldRun() bool
	Run() (string, Watcher, error)
}

type TimerWatcher struct {
	runAt   *time.Time
	message string
}

func NewTimerWatcher(runAt *time.Time, msg string) Watcher {
	return &TimerWatcher{runAt: runAt, message: msg}
}

func (w *TimerWatcher) ShouldRun() bool {
	return time.Now().After(*w.runAt)
}

func (w *TimerWatcher) Run() (string, Watcher, error) {
	return w.message, nil, nil
}

type DeliveryWatcher struct {
	runAt      *time.Time
	carrierID  string
	trackID    string
	timeCursor *time.Time

	trackerClient deliverytracker.Client
}

func NewDeliveryWatcher(carrierID, trackID string, timeCursor *time.Time, client deliverytracker.Client) Watcher {
	runAt := time.Now().Add(20 * time.Minute)
	return &DeliveryWatcher{runAt: &runAt, carrierID: carrierID, trackID: trackID, timeCursor: timeCursor, trackerClient: client}
}

func (w *DeliveryWatcher) ShouldRun() bool {
	return time.Now().After(*w.runAt)
}

func (w *DeliveryWatcher) Run() (string, Watcher, error) {
	track, err := w.trackerClient.GetTrack(w.carrierID, w.trackID)
	if err != nil {
		return "", nil, err
	} else if track.State != nil && track.State.ID == "delivered" {
		return "배송이 완료되었습니다.", nil, nil
	}

	var timeCursor *time.Time
	var trackMsg string
	if len(track.Progresses) > 0 && (w.timeCursor == nil || track.Progresses[len(track.Progresses)-1].Time.After(*w.timeCursor)) {
		progress := track.Progresses[len(track.Progresses)-1]

		timeAgo := int64(time.Now().Sub(*progress.Time) / time.Minute)
		trackMsg = fmt.Sprintf(
			"\n[배송 정보]\n운송장 : %s %s\n배송 현황 : %s\n현재 위치 : %s\n\n%s (업데이트 : %d분 전)\n\n",
			track.Carrier.Name, w.trackID, progress.Status.Text, progress.Location.Name, progress.Description, timeAgo,
		)

		timeCursor = progress.Time
	}

	newWatcher := NewDeliveryWatcher(w.carrierID, w.trackID, timeCursor, w.trackerClient)
	return trackMsg, newWatcher, nil
}
