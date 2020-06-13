package main

import (
	"github.com/webdonalds/discord-bot/commands"
	"os"
)

func main() {
	botToken := os.Getenv("DISCORD_BOT_TOKEN")
	bot, err := NewBot(botToken)
	if err != nil {
		panic(err)
	}

	bot.AddCommand(commands.NewPingCommand())
	bot.AddCommand(commands.NewDeliveryCommand())
	bot.AddCommand(commands.NewHelpCommand())

	err = bot.Listen()
	if err != nil {
		panic(err)
	}
}
