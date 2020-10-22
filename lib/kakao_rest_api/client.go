package kakao_rest_api

import (
	"context"
	"net/http"

	"github.com/google/go-querystring/query"
)

type Client struct {
	apiKey     string
	httpClient *http.Client

	daumEndpoint string
}

func NewClient(apiKey string) *Client {
	httpClient := &http.Client{}
	return &Client{
		apiKey:     apiKey,
		httpClient: httpClient,

		daumEndpoint: "https://dapi.kakao.com",
	}
}

func (client *Client) get(ctx context.Context, url string, params interface{}) (*http.Response, error) {
	q, err := query.Values(params)
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Authorization", "KakaoAK "+client.apiKey)

	return client.httpClient.Do(req)
}
