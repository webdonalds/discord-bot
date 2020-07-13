package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hellodhlyn/delivery-tracker"

	"github.com/webdonalds/discord-bot/background"
)

type DeliveryCommand struct {
	trackerClient deliverytracker.Client
}

func NewDeliveryCommand() Command {
	client, _ := deliverytracker.NewClient()
	return &DeliveryCommand{trackerClient: client}
}

func (*DeliveryCommand) CommandTexts() []string {
	return []string{"택배"}
}

func (*DeliveryCommand) ExpectedArgsLen() int {
	return 2
}

func (cmd *DeliveryCommand) Execute(args []string, _ *discordgo.MessageCreate) (string, background.Watcher, error) {
	carrierName := args[0]
	trackID := args[1]

	carriers, err := cmd.trackerClient.FindCarriersByName(carrierName)
	if len(carriers) != 1 {
		return "해당하는 이름의 택배사를 찾을 수 없습니다.", nil, nil
	} else if err != nil {
		return "", nil, err
	}

	carrierID := carriers[0].ID
	track, err := cmd.trackerClient.GetTrack(carrierID, trackID)
	if err != nil {
		return "", nil, err
	} else if track == nil {
		return "올바르지 않은 운송장이거나, 택배사에서 아직 물건을 인수하지 않았습니다.", nil, nil
	} else if track.State != nil && track.State.ID == "delivered" {
		return "이미 배송이 완료되었습니다.", nil, nil
	}

	var timeCursor *time.Time
	var trackMsg string
	if len(track.Progresses) > 0 {
		progress := track.Progresses[len(track.Progresses)-1]

		timeAgo := int64(time.Now().Sub(*progress.Time) / time.Minute)
		trackMsg = fmt.Sprintf(
			"[배송 정보]\n운송장 : %s %s\n배송 현황 : %s\n현재 위치 : %s\n\n%s (업데이트 : %d분 전)\n\n",
			carriers[0].Name, trackID, progress.Status.Text, progress.Location.Name, progress.Description, timeAgo,
		)

		timeCursor = progress.Time
	}

	watcher := background.NewDeliveryWatcher(carrierID, trackID, timeCursor, cmd.trackerClient)
	return fmt.Sprintf("%s배송 상태에 변경이 있을 시 30분 간격으로 알림을 발송합니다.", trackMsg), watcher, nil
}
