package kakao_rest_api

import (
	"context"
	"encoding/json"
	"time"
)

type SearchMeta struct {
	TotalCount    int64 `json:"total_count"`
	PageableCount int64 `json:"pageable_count"`
	IsEnd         bool  `json:"is_end"`
}

// Web Search API
// See: https://developers.kakao.com/tool/rest-api/open/get/v2-search-web
type SearchWebInputSort string

const (
	SearchWebInputSortAccuracy SearchWebInputSort = "accuracy"
	SearchWebInputSortRecency  SearchWebInputSort = "recency"
)

type SearchWebInput struct {
	Query string             `url:"query"`
	Sort  SearchWebInputSort `url:"sort,omitempty"`
	Page  int                `url:"page,omitempty"`
	Size  int                `url:"size,omitempty"`
}

type SearchWebOutputResult struct {
	Title    string     `json:"title"`
	Contents string     `json:"contents"`
	URL      string     `json:"url"`
	Datetime *time.Time `json:"datetime"`
}

type SearchWebOutput struct {
	Meta      *SearchMeta             `json:"meta"`
	Documents []SearchWebOutputResult `json:"documents"`
}

func (client *Client) SearchWeb(ctx context.Context, input *SearchWebInput) (*SearchWebOutput, error) {
	res, err := client.get(ctx, client.daumEndpoint+"/v2/search/web", input)
	if err != nil {
		return nil, err
	}

	var output SearchWebOutput
	if err = json.NewDecoder(res.Body).Decode(&output); err != nil {
		return nil, err
	}
	return &output, nil
}

// Image Search API
// https://developers.kakao.com/tool/rest-api/open/get/v2-search-image
type SearchImageInputSort string

const (
	SearchImageInputSortAccuracy SearchImageInputSort = "accuracy"
	SearchImageInputSortRecency  SearchImageInputSort = "recency"
)

type SearchImageInput struct {
	Query string               `url:"query"`
	Sort  SearchImageInputSort `url:"sort,omitempty"`
	Page  int                  `url:"page,omitempty"`
	Size  int                  `url:"size,omitempty"`
}

type SearchImageOutputResult struct {
	Collection      string     `json:"collection"`
	ThumbnailURL    string     `json:"thumbnail_url"`
	ImageURL        string     `json:"image_url"`
	Width           int        `json:"width"`
	Height          int        `json:"height"`
	DisplaySitename string     `json:"display_sitename"`
	DocURL          string     `json:"doc_url"`
	DateTime        *time.Time `json:"datetime"`
}

type SearchImageOutput struct {
	Meta      *SearchMeta               `json:"meta"`
	Documents []SearchImageOutputResult `json:"documents"`
}

func (client *Client) SearchImage(ctx context.Context, input *SearchImageInputSort) (*SearchImageOutput, error) {
	res, err := client.get(ctx, client.daumEndpoint+"/v2/search/image", input)
	if err != nil {
		return nil, err
	}

	var output SearchImageOutput
	if err = json.NewDecoder(res.Body).Decode(&output); err != nil {
		return nil, err
	}
	return &output, nil
}
