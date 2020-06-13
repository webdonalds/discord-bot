package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/webdonalds/discord-bot/commands"
)

type Bot struct {
	sess *discordgo.Session
	cmds map[string]commands.Command
}

func NewBot(token string) (*Bot, error) {
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		sess: sess,
		cmds: map[string]commands.Command{},
	}, nil
}

func (bot *Bot) AddCommand(cmd commands.Command) {
	for _, cmdText := range cmd.CommandTexts() {
		bot.cmds[cmdText] = cmd
	}
}

func (bot *Bot) NewMessageChannel(channelID, mention string) chan<- string {
	msgChan := make(chan string)
	go func() {
		for msg := range msgChan {
			content := mention + " " + msg
			_, _ = bot.sess.ChannelMessageSend(channelID, content)
		}
	}()
	return msgChan
}

func (bot *Bot) Listen() error {
	bot.sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if (m.Author.ID == s.State.User.ID) || (!strings.HasPrefix(m.Content, "!")) {
			return
		}

		texts := strings.Split(m.Content, " ")
		cmdText := strings.Replace(texts[0], "!", "", 1)
		for text, cmd := range bot.cmds {
			if cmdText == text {
				_ = cmd.Execute(texts[1:], bot.NewMessageChannel(m.ChannelID, m.Author.Mention()), m)
				break
			}
		}
	})

	err := bot.sess.Open()
	if err != nil {
		return err
	}

	fmt.Println("Bot has been started.")
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sigChan

	fmt.Println("Terminating...")
	return bot.sess.Close()
}
