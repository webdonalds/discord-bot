package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"

	"github.com/webdonalds/discord-bot/background"
	"github.com/webdonalds/discord-bot/commands"
	"github.com/webdonalds/discord-bot/crons"
)

type Bot struct {
	sess   *discordgo.Session
	cmds   map[string]commands.Command
	cron   *cron.Cron
	worker *background.Worker
}

func NewBot(token string) (*Bot, error) {
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Seoul")
	cron := cron.New(cron.WithLocation(loc))

	return &Bot{
		sess:   sess,
		cmds:   map[string]commands.Command{},
		cron:   cron,
		worker: background.NewWorker(sess),
	}, nil
}

func (bot *Bot) AddCommand(cmd commands.Command) {
	for _, cmdText := range cmd.CommandTexts() {
		bot.cmds[cmdText] = cmd
	}
}

func (bot *Bot) AddCron(c crons.Cron) {
	_, _ = bot.cron.AddFunc(c.Pattern(), func() {
		msg := c.Execute()
		if msg != "" {
			_, _ = bot.sess.ChannelMessageSend(c.ChannelID(), msg)
		}
	})
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
				args := texts[1:]
				if len(args) != cmd.ExpectedArgsLen() {
					_, _ = s.ChannelMessageSend(m.ChannelID, "파라미터 개수가 올바르지 않습니다.")
					break
				}

				msg, watcher, err := cmd.Execute(args, m)
				if err != nil {
					fmt.Printf("%v\n", err)
					msg = "오류가 발생했습니다. 서버 로그을 확인하세요."
				}
				if msg != "" {
					_, _ = s.ChannelMessageSend(m.ChannelID, msg)
				}
				if watcher != nil {
					bot.worker.AddWatcher(watcher, m.ChannelID, m.Author.Mention())
				}
				break
			}
		}
	})

	bot.cron.Start()
	bot.worker.Start()

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
