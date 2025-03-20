package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/webdonalds/discord-bot/internal/actions"
	"github.com/webdonalds/discord-bot/internal/commands"
)

type Bot struct {
	session  *discordgo.Session
	commands map[string]commands.Command
}

func NewBot(token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		session:  session,
		commands: make(map[string]commands.Command),
	}, nil
}

func (bot *Bot) AddCommand(command commands.Command) {
	for _, cmd := range command.Commands() {
		bot.commands[cmd] = command
	}
}

func (bot *Bot) Listen() error {
	bot.session.AddHandler(func(session *discordgo.Session, message *discordgo.MessageCreate) {
		if message.Author.Bot || !strings.HasPrefix(message.Content, "!") {
			return
		}

		command := strings.TrimPrefix(message.Content, "!")
		args := strings.Split(command, " ")
		if cmd, ok := bot.commands[args[0]]; ok {
			action, err := cmd.Execute(context.Background(), args[1:])
			if err != nil {
				action = &actions.Reply{Message: fmt.Sprintf("명령어 실행 중 오류가 발생했습니다:\n\n%v", err)}
			}

			if err := action.Execute(context.Background(), session, message); err != nil {
				session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("결과 처리 중 오류가 발생했습니다:\n\n%v", err))
			}
		}
	})

	err := bot.session.Open()
	if err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Println("Bot is running...")
	<-sigChan

	log.Println("Bot is shutting down...")
	return bot.session.Close()
}
