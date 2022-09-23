package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	acceptHeader = "application/vnd.github+json"
	apiEndpoint  = "https://api.github.com"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
}

type (
	PackageVersionMetadataContainer struct {
		Tags []string `json:"tags"`
	}

	PackageVersionMetadata struct {
		PackageType string                           `json:"package_type"`
		Container   *PackageVersionMetadataContainer `json:"container"`
	}

	GetPackageVersionResponse struct {
		ID       int64                   `json:"id"`
		Metadata *PackageVersionMetadata `json:"metadata"`
	}
)

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		apiKey:     apiKey,
	}
}

func (c *Client) GetPackageVersions(ctx context.Context) ([]GetPackageVersionResponse, error) {
	path := "/orgs/webdonalds/packages/container/discord-bot/versions"
	req, _ := http.NewRequestWithContext(ctx, "GET", apiEndpoint+path, nil)
	req.Header.Set("Accept", acceptHeader)
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected server response: %d", res.StatusCode)
	}

	var resBody []GetPackageVersionResponse
	_ = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return nil, err
	}
	return resBody, nil
}

func (c *Client) Authenticated() bool {
	return c.apiKey != ""
}
