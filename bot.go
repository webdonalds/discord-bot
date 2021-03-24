package main

import (
	"os"
	"os/signal"
	"regexp"
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
	sess     *discordgo.Session
	cmds     map[string]commands.Command
	cmdRegex *regexp.Regexp
	cron     *cron.Cron
}

func NewBot(token string) (*Bot, error) {
	sess, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Seoul")
	cron := cron.New(cron.WithLocation(loc))

	cmdRegex, _ := regexp.Compile("[^\\s\"']+|\"([^\"]*)\"|'([^']*)'")

	return &Bot{
		sess:     sess,
		cmds:     map[string]commands.Command{},
		cmdRegex: cmdRegex,
		cron:     cron,
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

		splits := strings.SplitN(m.Content, " ", 2)
		cmdText := strings.Replace(splits[0], "!", "", 1)

		cmdArgs := []string{}
		if len(splits) == 2 {
			for _, word := range bot.cmdRegex.FindAllString(splits[1], -1) {
				cmdArgs = append(cmdArgs, strings.ReplaceAll(strings.ReplaceAll(word, "\"", ""), "'", ""))
			}
		}

		for text, cmd := range bot.cmds {
			if cmdText == text {
				msg, _, err := cmd.Execute(cmdArgs, m)
				if err != nil {
					log.Error(err)
					msg = "오류가 발생했습니다. 서버 로그을 확인하세요."
				}
				if msg != "" {
					_, _ = s.ChannelMessageSend(m.ChannelID, msg)
				}
				break
			}
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
