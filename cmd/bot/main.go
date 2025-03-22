package main

import (
	"log"
	"os"

	"github.com/webdonalds/discord-bot/internal/commands"
)

func main() {
	bot, err := NewBot(os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	bot.AddCommand(commands.NewPingCommand())
	bot.AddCommand(commands.NewChatCommand(os.Getenv("OPENAI_API_KEY")))

	log.Fatalf("Failed to listen: %v", bot.Listen())
}
