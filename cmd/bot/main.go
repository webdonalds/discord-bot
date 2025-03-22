package main

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/webdonalds/discord-bot/internal/chatmemory"
	"github.com/webdonalds/discord-bot/internal/commands"
)

func main() {
	bot, err := NewBot(os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Username: os.Getenv("REDIS_USERNAME"),
		Password: os.Getenv("REDIS_PASSWORD"),
	})
	chatMemoryStore := chatmemory.NewStore(redisClient)

	bot.AddCommand(commands.NewPingCommand())
	bot.AddCommand(commands.NewChatCommand(os.Getenv("OPENAI_API_KEY"), chatMemoryStore))

	log.Fatalf("Failed to listen: %v", bot.Listen())
}
