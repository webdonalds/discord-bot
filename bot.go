package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	"github.com/webdonalds/discord-bot/commands"
	"github.com/webdonalds/discord-bot/crons"
)

type Bot struct {
	sess *discordgo.Session
	cmds map[string]commands.Command
	cron *cron.Cron
}

func NewBot(token string) (*Bot, error) {
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Seoul")
	cron := cron.New(cron.WithLocation(loc))

	return &Bot{
		sess: sess,
		cmds: map[string]commands.Command{},
		cron: cron,
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

		resMsg := ""
		cmdText, cmdArgs, err := ParseCommand(m.Content)
		if err != nil {
			log.Error(err)
			resMsg = "명령어 파싱에 실패했습니다."
		} else {
			for text, cmd := range bot.cmds {
				if cmdText == text {
					resMsg, _, err = cmd.Execute(cmdArgs, m)
					if err != nil {
						log.Error(err)
						resMsg = "오류가 발생했습니다. 서버 로그을 확인하세요."
					}
					break
				}
			}
		}

		if resMsg != "" {
			_, _ = s.ChannelMessageSend(m.ChannelID, resMsg)
		}
	})

	bot.cron.Start()

	err := bot.sess.Open()
	if err != nil {
		return err
	}

	log.Info("Bot has been started.")
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sigChan

	log.Warn("Terminating...")
	return bot.sess.Close()
}
