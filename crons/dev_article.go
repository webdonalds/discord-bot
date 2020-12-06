package crons

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/webdonalds/discord-bot/repositories"
)

const crawlURL = "https://awesome-devblog-v2-beta-backend.vercel.app/api/korean/people/feeds/rank?day=14&size=14"

type DevArticle struct {
	ID          string    `json:"_id"`
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	ImageURL    string    `json:"imgUrl"`
	Tags        []string  `json:"tags"`
	Count       int64     `json:"count"`
}

type DevArticleCron struct {
	repo repositories.DevArticleRepository
}

func NewDevArticleCron(repo repositories.DevArticleRepository) Cron {
	return &DevArticleCron{repo: repo}
}

func (cron *DevArticleCron) Pattern() string {
	return "0 10 * * *" // 10AM KST
}

func (cron *DevArticleCron) ChannelID() string {
	return os.Getenv("DEV_ARTICLE_CHANNEL_ID")
}

func (cron *DevArticleCron) Execute() string {
	ctx := context.Background()
	res, err := http.Get(crawlURL)
	if err != nil {
		log.Errorf("failed to read articles: %v", err)
		return ""
	}

	defer func() { _ = res.Body.Close() }()
	var articles []DevArticle
	err = json.NewDecoder(res.Body).Decode(&articles)
	if err != nil {
		log.Errorf("failed to parse articles: %v", err)
		return ""
	}

	readArticleIDs, err := cron.repo.ListAllReadArticleID(ctx)
	if err != nil {
		log.Errorf("failed to list read articles: %v", err)
		return ""
	}

	for _, article := range articles {
		if cron.containID(readArticleIDs, article.ID) {
			continue
		}

		_ = cron.repo.AddReadArticleID(ctx, article.ID)
		return fmt.Sprintf("**%s** (%s)\n%s...\n%s", article.Title, article.Author, article.Description, article.Link)
	}
	return "최근 14일 이내 새로운 포스트가 없습니다."
}

func (cron *DevArticleCron) containID(articleIDs []string, articleID string) bool {
	for _, i := range articleIDs {
		if i == articleID {
			return true
		}
	}
	return false
}
