package commands

import (
	"github.com/bwmarrin/discordgo"
	"math/rand"

	"github.com/webdonalds/discord-bot/background"
)

type RandomCommand struct{}

func NewRandomCommand() Command {
	return &RandomCommand{}
}

func (*RandomCommand) CommandTexts() []string {
	return []string{"랜덤", "뽑기"}
}

func (*RandomCommand) Execute(args []string, _ *discordgo.MessageCreate) (string, background.Watcher, error) {
	if len(args) == 0 {
		return "최소 하나 이상의 항목을 제시해야합니다.", nil, nil
	}
	return args[rand.Intn(len(args))], nil, nil
}
