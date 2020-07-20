package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/bwmarrin/discordgo"

	"github.com/webdonalds/discord-bot/background"
)

type ExchangeCommand struct{}

func NewExchangeCommand() Command {
	return &ExchangeCommand{}
}

func (*ExchangeCommand) CommandTexts() []string {
	return []string{"환율"}
}

func (*ExchangeCommand) Execute(args []string, _ *discordgo.MessageCreate) (string, background.Watcher, error) {
	api := "https://earthquake.kr:23490/query/" // 참고: https://jaeheon.kr/12
	url := api + args[0]

	targetPrice, _ := strconv.ParseFloat(args[1], 64)

	req, err := http.Get(url)

	if err != nil {
		return "환율 API 에러: ", nil, err
	}

	defer req.Body.Close()

	data, err := ioutil.ReadAll(req.Body)

	if err != nil {
		return "환율 API Read 에러", nil, err
	}

	var res map[string]interface{}
	jsonErr := json.Unmarshal([]byte(data), &res)

	if jsonErr != nil {
		return "환율 json 처리 에러", nil, jsonErr
	}
	exchageData := res[args[0]].([]interface{})
	price := exchageData[0].(float64) * targetPrice

	msg := fmt.Sprintf("%s: %s to %s", args[0], args[1], strconv.FormatFloat(price, 'f', -1, 32))
	return msg, nil, nil
}
