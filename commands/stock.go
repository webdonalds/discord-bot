package commands

import (
	"context"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math"

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
	if len(args) != 1 {
		return "종목명을 입력하세요.", nil, nil
	}

	stockName := args[0]
	ctx := context.Background()

	var searchResult *stock.GetSearchResponse
	if err := retry.Do(func() (err error) {
		searchResult, err = c.stockClient.GetSearch(ctx, stockName)
		return
	}); err != nil {
		log.Errorf("failed to search asset\n%v", err)
		return "종목 정보 검색에 실패했습니다.", nil, nil
	}

	if len(searchResult.Assets) == 0 {
		return "해당하는 종목을 찾을 수 없습니다.", nil, nil
	}

	var securityResult *stock.GetRecentSecuritiesResponse
	if err := retry.Do(func() (err error) {
		securityResult, err = c.stockClient.GetRecentSecurities(ctx, []string{searchResult.Assets[0].AssetID})
		return
	}); err != nil {
		log.Errorf("failed to get recent security\n%v", err)
		return "종목 정보 조회에 실패했습니다.", nil, nil
	}

	security := securityResult.RecentSecurities[0]
	upOrDown := "－"
	if security.ChangePrice > 0 {
		upOrDown = "▲"
	} else if security.ChangePrice < 0 {
		upOrDown = "▼"
	}

	p := message.NewPrinter(language.English)
	msg := p.Sprintf("**%s**(%s) %.0f / %s%.0f (%.2f%%)",
		security.Name, security.ShortCode, security.TradePrice,
		upOrDown, math.Abs(security.ChangePrice), security.ChangePriceRate*100,
	)
	return msg, nil, nil
}
