package waqi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const apiEndpoint = "https://api.waqi.info"

type Client struct {
	httpClient *http.Client
	apiToken   string
}

type (
	GetCityFeedResponseDataIAQIValue struct {
		V int64 `json:"v"`
	}

	GetCityFeedResponseDataIAQI struct {
		PM25 GetCityFeedResponseDataIAQIValue `json:"pm25"`
		PM10 GetCityFeedResponseDataIAQIValue `json:"pm10"`
	}

	GetCityFeedResponseData struct {
		Idx  int64                       `json:"idx"`
		AQI  int64                       `json:"aqi"`
		IAQI GetCityFeedResponseDataIAQI `json:"iaqi"`
	}

	GetCityFeedResponse struct {
		Status string                  `json:"status"`
		Data   GetCityFeedResponseData `json:"data"`
	}
)

func (c *Client) GetCityFeed(ctx context.Context, cityName string) (*GetCityFeedResponse, error) {
	var res GetCityFeedResponse
	err := c.get(ctx, "/feed/"+cityName+"/", nil, &res)
	return &res, err
}

func (c *Client) get(ctx context.Context, path string, params map[string]string, resBody interface{}) error {
	req, _ := http.NewRequest("GET", apiEndpoint+path, nil)
	req.WithContext(ctx)

	q := req.URL.Query()
	q.Add("token", c.apiToken)
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

func NewClient(apiToken string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiToken:   apiToken,
	}
}
