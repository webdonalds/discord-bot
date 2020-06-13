package commands

import "github.com/bwmarrin/discordgo"

type HelpCommand struct{}

func NewHelpCommand() Command {
	return &HelpCommand{}
}

func (*HelpCommand) CommandTexts() []string {
	return []string{"help", "도움"}
}

func (*HelpCommand) Execute(_ []string, msgChan chan<- string, _ *discordgo.MessageCreate) error {
	msgChan <- "\n!help\n" +
		"!도움\n\n" +
		"!택배 <택배사> <운송장번호>"
	return nil
}
