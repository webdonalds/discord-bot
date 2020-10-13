package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hellodhlyn/delivery-tracker"

	"github.com/webdonalds/discord-bot/background"
	"github.com/webdonalds/discord-bot/repositories"
)

type DeliveryCommand struct {
	repo          repositories.DeliveryTrackRepository
	trackerClient deliverytracker.Client
}

func NewDeliveryCommand(repo repositories.DeliveryTrackRepository, trackerClient deliverytracker.Client) Command {
	return &DeliveryCommand{repo: repo, trackerClient: trackerClient}
}

func (*DeliveryCommand) CommandTexts() []string {
	return []string{"택배"}
}

func (*DeliveryCommand) ExpectedArgsLen() int {
	return 2
}

func (cmd *DeliveryCommand) Execute(args []string, msg *discordgo.MessageCreate) (string, background.Watcher, error) {
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

	var lastTimestamp *time.Time
	var trackMsg string
	if len(track.Progresses) > 0 {
		progress := track.Progresses[len(track.Progresses)-1]

		timeAgo := int64(time.Now().Sub(*progress.Time) / time.Minute)
		trackMsg = fmt.Sprintf(
			"[배송 정보]\n운송장 : %s %s\n배송 현황 : %s\n현재 위치 : %s\n\n%s (업데이트 : %d분 전)\n\n",
			carriers[0].Name, trackID, progress.Status.Text, progress.Location.Name, progress.Description, timeAgo,
		)

		lastTimestamp = progress.Time
	}

	runAt := time.Now().Add(20 * time.Minute)
	err = cmd.repo.Append(context.Background(), &repositories.DeliveryTrack{
		Mention: msg.Author.Mention(), CarrierID: carrierID, TrackID: trackID, LastTimestamp: lastTimestamp,
	}, &runAt)
	return fmt.Sprintf("%s배송 상태에 변경이 있을 시 20분 간격으로 알림을 발송합니다.", trackMsg), nil, err
}
