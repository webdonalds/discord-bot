package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hellodhlyn/delivery-tracker"
	"github.com/webdonalds/discord-bot/services"
)

type DeliveryCommand struct {
	trackerClient deliverytracker.Client
	follower      *services.DeliveryService
}

func NewDeliveryCommand() Command {
	client, _ := deliverytracker.NewClient()

	follower := services.NewDeliveryFollower()
	follower.Start()

	return &DeliveryCommand{
		trackerClient: client,
		follower:      follower,
	}
}

func (*DeliveryCommand) CommandTexts() []string {
	return []string{"택배"}
}

func (*DeliveryCommand) ExpectedArgsLen() int {
	return 2
}

func (cmd *DeliveryCommand) Execute(args []string, msgChan chan<- string, _ *discordgo.MessageCreate) error {
	carrierName := args[0]
	trackID := args[1]
	return cmd.follower.Track(carrierName, trackID, msgChan, nil)
}
