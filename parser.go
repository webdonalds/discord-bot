package main

import (
	"errors"
	"strings"
)

func ParseCommand(text string) (string, []string, error) {
	splits := strings.SplitN(text, " ", 2)
	cmdText := strings.Replace(splits[0], "!", "", 1)

	if len(splits) == 1 {
		return cmdText, []string{}, nil
	}

	cmdArgs := []string{}
	currStr := ""
	currQuote := ' '
	for _, c := range splits[1] {
		if c == currQuote {
			trimmedStr := strings.TrimSpace(currStr)
			if trimmedStr != "" {
				cmdArgs = append(cmdArgs, trimmedStr)
			}
			currStr = ""
			currQuote = ' '
		} else if c == '"' {
			currQuote = '"'
		} else if c == '\'' {
			currQuote = '\''
		} else {
			currStr += string(c)
		}
	}

	if currQuote != ' ' {
		return "", []string{}, errors.New("unexpected end of arguments")
	} else if trimmedStr := strings.TrimSpace(currStr); trimmedStr != "" {
		cmdArgs = append(cmdArgs, trimmedStr)
	}

	return cmdText, cmdArgs, nil
}
