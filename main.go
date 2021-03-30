package main

import (
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/hellodhlyn/delivery-tracker"
	"github.com/kz/discordrus"
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
	daRepo, _ := repositories.NewRedisDevArticleRepository(rdb)

	// Register commands and crons
	bot.AddCommand(commands.NewPingCommand())
	bot.AddCommand(commands.NewDeliveryCommand(dtRepo, trackClient))
	bot.AddCommand(commands.NewHelpCommand())
	bot.AddCommand(commands.NewTimerCommand())
	bot.AddCommand(commands.NewExchangeCommand())
	bot.AddCommand(commands.NewRandomCommand())
	bot.AddCommand(commands.NewAirQualityCommand())

	bot.AddCron(crons.NewBreakingNewsCron())
	bot.AddCron(crons.NewDeliveryTrackCron(dtRepo, trackClient))
	bot.AddCron(crons.NewDevArticleCron(daRepo))

	err = bot.Listen()
	if err != nil {
		panic(err)
	}
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	loc, _ := time.LoadLocation("Asia/Seoul")

	log.AddHook(discordrus.NewHook(
		// Use environment variable for security reasons
		os.Getenv("ERROR_LOG_WEBHOOK_URL"),
		// Set minimum level to DebugLevel to receive all log entries
		log.WarnLevel,
		&discordrus.Opts{
			Username:         "ErrorBot",
			DisableTimestamp: false,                      // Setting this to true will disable timestamps from appearing in the footer
			TimestampFormat:  "Jan 2 15:04:05.00000 MST", // The timestamp takes this format; if it is unset, it will take log' default format
			TimestampLocale:  loc,                        // The timestamp uses this locale; if it is unset, it will use time.Local
		},
	))
}
