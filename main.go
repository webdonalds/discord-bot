package main

import (
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/hellodhlyn/delivery-tracker"
	log "github.com/sirupsen/logrus"

	"github.com/webdonalds/discord-bot/commands"
	"github.com/webdonalds/discord-bot/crons"
	"github.com/webdonalds/discord-bot/repositories"
)

func getEnvOrDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func main() {
	// Initialize connections
	rdb := redis.NewClient(&redis.Options{
		Addr:     getEnvOrDefault("REDIS_URL", "127.0.0.1:6379"),
		Password: getEnvOrDefault("REDIS_PASSWORD", ""),
	})

	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	bot, err := NewBot(botToken)
	if err != nil {
		panic(err)
	}

	trackClient, _ := deliverytracker.NewClient()

	// Initialize repositories
	dtRepo, _ := repositories.NewRedisDeliveryTrackRepository(rdb)

	// Register commands and crons
	bot.AddCommand(commands.NewPingCommand())
	bot.AddCommand(commands.NewDeliveryCommand(dtRepo, trackClient))
	bot.AddCommand(commands.NewHelpCommand())
	bot.AddCommand(commands.NewTimerCommand())
	bot.AddCommand(commands.NewExchangeCommand())

	bot.AddCron(crons.NewBreakingNewsCron())
	bot.AddCron(crons.NewDeliveryTrackCron(dtRepo, trackClient))

	err = bot.Listen()
	if err != nil {
		panic(err)
	}
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}
