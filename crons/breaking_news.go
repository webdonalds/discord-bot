package crons

import (
	"fmt"
	"os"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type BreakingNewsCron struct {
	twitterClient *twitter.Client
	lastTweetID   int64
}

func NewBreakingNewsCron() Cron {
	return &BreakingNewsCron{
		twitterClient: newTwitterClient(),
	}
}

func newTwitterClient() *twitter.Client {
	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	return twitter.NewClient(config.Client(oauth1.NoContext, token))
}

func (*BreakingNewsCron) Pattern() string {
	return "@every 1m"
}

func (*BreakingNewsCron) ChannelID() string {
	return os.Getenv("BREAKING_NEWS_CHANNEL_ID")
}

func (cron *BreakingNewsCron) Execute() string {
	tweets, _, err := cron.twitterClient.Timelines.UserTimeline(&twitter.UserTimelineParams{
		UserID:          147451838,
		SinceID:         cron.lastTweetID,
		IncludeRetweets: twitter.Bool(false),
		ExcludeReplies:  twitter.Bool(true),
	})
	if err != nil {
		fmt.Printf("%v\n", err)
		return ""
	}

	if len(tweets) > 0 {
		cron.lastTweetID = tweets[0].ID
	}

	texts := []string{}
	for _, tweet := range tweets {
		if strings.Contains(tweet.Text, "속보") || strings.Contains(tweet.Text, "1보") {
			texts = append(texts, tweet.Text)
		}
	}

	return strings.Join(texts, "\n\n")
}
