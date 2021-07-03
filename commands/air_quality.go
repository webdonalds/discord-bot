package commands

import (
	"context"
	"fmt"
	"github.com/webdonalds/discord-bot/responses"
	"os"

	"github.com/bwmarrin/discordgo"

	"github.com/webdonalds/discord-bot/background"
	"github.com/webdonalds/discord-bot/lib/waqi"
)

type AirQualityCommand struct {
	waqiClient *waqi.Client
}

func NewAirQualityCommand() Command {
	apiToken := os.Getenv("WAQI_API_TOKEN")
	waqiClient := waqi.NewClient(apiToken)
	return &AirQualityCommand{
		waqiClient: waqiClient,
	}
}

func (c *AirQualityCommand) CommandTexts() []string {
	return []string{"미세먼지"}
}

func (c *AirQualityCommand) Execute(_ []string, _ *discordgo.MessageCreate) (responses.ResponseMessage, background.Watcher, error) {
	res, err := c.waqiClient.GetCityFeed(context.Background(), "seoul")
	if err != nil {
		return nil, nil, err
	}

	pm10 := res.Data.IAQI.PM10.V
	pm10Label := ""
	if pm10 < 30 {
		pm10Label = "좋음"
	} else if pm10 < 80 {
		pm10Label = "보통"
	} else if pm10 < 150 {
		pm10Label = "나쁨"
	} else {
		pm10Label = "최악"
	}

	pm25 := res.Data.IAQI.PM25.V
	pm25Label := ""
	if pm25 < 15 {
		pm25Label = "좋음"
	} else if pm25 < 35 {
		pm25Label = "보통"
	} else if pm25 < 75 {
		pm25Label = "나쁨"
	} else {
		pm25Label = "최악"
	}

	msg := fmt.Sprintf("현재 서울의 대기 정보\n미세먼지 : %d (%s)\n초미세먼지 : %d (%s)", pm10, pm10Label, pm25, pm25Label)

	return responses.NewTextMessage(msg), nil, nil
}
