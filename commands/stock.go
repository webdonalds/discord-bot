package commands

import (
	"context"
	"errors"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math"
	"strings"

	"github.com/avast/retry-go"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"

	"github.com/webdonalds/discord-bot/background"
	"github.com/webdonalds/discord-bot/lib/stock"
)

type StockCommand struct {
	stockClient *stock.Client
}

func NewStockCommand() *StockCommand {
	return &StockCommand{
		stockClient: stock.NewClient(),
	}
}

func (*StockCommand) CommandTexts() []string {
	return []string{"주가"}
}

func (c *StockCommand) Execute(args []string, _ *discordgo.MessageCreate) (string, background.Watcher, error) {
	if len(args) == 0 {
		return "종목명을 1개 이상 입력하세요.", nil, nil
	}

	ctx := context.Background()
	msgs := []string{}
	for _, arg := range args {
		eachMsg, err := c.executeEach(ctx, arg)
		if err != nil {
			log.Errorf("failed to get recent security\n%v", err)
		} else {
			msgs = append(msgs, eachMsg)
		}
	}

	return strings.Join(msgs, "\n"), nil, nil
}

func (c *StockCommand) executeEach(ctx context.Context, stockName string) (string, error) {
	var searchResult *stock.GetSearchResponse
	if err := retry.Do(func() (err error) {
		searchResult, err = c.stockClient.GetSearch(ctx, stockName)
		return
	}); err != nil {
		return "", err
	}

	if len(searchResult.Assets) == 0 {
		return "", errors.New("해당하는 종목을 찾을 수 없습니다")
	}

	var securityResult *stock.GetRecentSecuritiesResponse
	if err := retry.Do(func() (err error) {
		securityResult, err = c.stockClient.GetRecentSecurities(ctx, []string{searchResult.Assets[0].AssetID})
		return
	}); err != nil {
		return "", err
	}

	security := securityResult.RecentSecurities[0]
	upOrDown := "－"
	if security.ChangePrice > 0 {
		upOrDown = "▲"
	} else if security.ChangePrice < 0 {
		upOrDown = "▼"
	}

	p := message.NewPrinter(language.English)
	return p.Sprintf("**%s**(%s) %.0f / %s%.0f (%.2f%%)",
		security.Name, security.ShortCode, security.TradePrice,
		upOrDown, math.Abs(security.ChangePrice), security.ChangePriceRate*100,
	), nil
}
