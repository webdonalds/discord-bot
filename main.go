package main

import (
	"os"

	"github.com/webdonalds/discord-bot/commands"
	"github.com/webdonalds/discord-bot/crons"
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
	bot.AddCommand(commands.NewTimerCommand())
	bot.AddCommand(commands.NewExchangeCommand())

	bot.AddCron(crons.NewBreakingNewsCron())

	err = bot.Listen()
	if err != nil {
		panic(err)
	}
}
