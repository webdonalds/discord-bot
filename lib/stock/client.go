package stock

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const apiEndpoint = "https://stockplus.com"

type Client struct {
	httpClient *http.Client
}

type (
	GetSearchResponseAsset struct {
		Type             string `json:"type"`
		Code             string `json:"code"`
		Name             string `json:"name"`
		AssetID          string `json:"assetId"`
		DisplayedSubtype string `json:"displayedSubtype"`
		DisplayedCode    string `json:"displayedCode"`
	}

	GetSearchResponse struct {
		Keyword    string                   `json:"keyword"`
		Assets     []GetSearchResponseAsset `json:"assets"`
		NextCursor string                   `json:"nextCursor"`
	}

	GetRecentSecuritiesResponseRecentSecurity struct {
		AccTradeVolume        int64   `json:"accTradeVolume"`
		Board                 string  `json:"board"`
		Change                string  `json:"change"`
		ChangePrice           float64 `json:"changePrice"`
		ChangePriceRate       float64 `json:"changePriceRate"`
		Code                  string  `json:"code"`
		Currency              string  `json:"currency"`
		Date                  string  `json:"date"`
		DayChartUrl           string  `json:"dayChartUrl"`
		DelayedMinutes        int     `json:"delayedMinutes"`
		DisplayedPrice        float64 `json:"displayedPrice"`
		Eps                   int     `json:"eps"`
		ExchangeCountry       string  `json:"exchangeCountry"`
		ExchangeCountryName   string  `json:"exchangeCountryName"`
		ForeignRatio          string  `json:"foreignRatio"`
		GlobalAccTradePrice   float64 `json:"globalAccTradePrice"`
		High52wPrice          float64 `json:"high52wPrice"`
		HighPrice             float64 `json:"highPrice"`
		Id                    string  `json:"id"`
		IsIndex               bool    `json:"isIndex"`
		IsVi                  bool    `json:"isVi"`
		Low52wPrice           float64 `json:"low52wPrice"`
		LowPrice              float64 `json:"lowPrice"`
		Market                string  `json:"market"`
		MarketCapRank         int     `json:"marketCapRank"`
		MarketName            string  `json:"marketName"`
		MarketWarningMsg      string  `json:"marketWarningMsg"`
		MiniDayChartUrl       string  `json:"miniDayChartUrl"`
		MiniDayGuidedChartUrl string  `json:"miniDayGuidedChartUrl"`
		Name                  string  `json:"name"`
		OpeningPrice          float64 `json:"openingPrice"`
		Per                   float64 `json:"per"`
		PrevClosingPrice      float64 `json:"prevClosingPrice"`
		RegularHoursStatus    string  `json:"regularHoursStatus"`
		SectorName            string  `json:"sectorName"`
		SecurityGroup         string  `json:"securityGroup"`
		IsSecurity            bool    `json:"isSecurity"`
		ShortCode             string  `json:"shortCode"`
		SignedChangePrice     float64 `json:"signedChangePrice"`
		SignedChangeRate      float64 `json:"signedChangeRate"`
		TotalMarketValue      float64 `json:"totalMarketValue"`
		TradePrice            float64 `json:"tradePrice"`
		TradeStrength         float64 `json:"tradeStrength"`
		TradeTime             string  `json:"tradeTime"`
	}

	GetRecentSecuritiesResponse struct {
		RecentSecurities []GetRecentSecuritiesResponseRecentSecurity `json:"recentSecurities"`
	}
)

func (c *Client) GetSearch(ctx context.Context, keyword string) (*GetSearchResponse, error) {
	var res GetSearchResponse
	err := c.get(ctx, "/api/search/assets.json", map[string]string{"keyword": keyword}, &res)
	return &res, err
}

func (c *Client) GetRecentSecurities(ctx context.Context, ids []string) (*GetRecentSecuritiesResponse, error) {
	params := map[string]string{"ids": strings.Join(ids, ",")}

	var res GetRecentSecuritiesResponse
	err := c.get(ctx, "/api/securities.json", params, &res)
	return &res, err
}

func (c *Client) get(ctx context.Context, path string, params map[string]string, resBody interface{}) error {
	req, _ := http.NewRequest("GET", apiEndpoint+path, nil)
	req.WithContext(ctx)

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected server response: %d", res.StatusCode)
	}

	return json.NewDecoder(res.Body).Decode(&resBody)
}

func NewClient() *Client {
	return &Client{httpClient: &http.Client{Timeout: 10 * time.Second}}
}
