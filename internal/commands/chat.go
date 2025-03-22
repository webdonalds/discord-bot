package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
	"github.com/webdonalds/discord-bot/internal/actions"
)

type ChatCommand struct {
	client openai.Client
}

func NewChatCommand(apiKey string) *ChatCommand {
	return &ChatCommand{
		client: openai.NewClient(),
	}
}

func (c *ChatCommand) Commands() []string {
	return []string{"chat"}
}

func (c *ChatCommand) Execute(ctx context.Context, args []string) (actions.Action, error) {
	if len(args) == 0 {
		return &actions.Reply{Message: "문장을 입력해주세요"}, nil
	}

	prompt := strings.Join(args, " ")

	messageStream := make(chan string)
	replyStream := &actions.ReplyStream{
		MessageStream: messageStream,
	}

	// Start the reply stream in a goroutine
	go func() {
		defer close(messageStream)

		stream := c.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.DeveloperMessage("당신은 메신저에서 대화중이니 가능한 간결하게 답변하세요."),
				openai.UserMessage(prompt),
			},
			Model:               "gpt-4o-mini",
			MaxCompletionTokens: openai.Int(200),
		})

		for stream.Next() {
			chunk := stream.Current()
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				messageStream <- chunk.Choices[0].Delta.Content
			}
		}

		if stream.Err() != nil {
			messageStream <- fmt.Sprintf("오류가 발생했어요.\n\n%v", stream.Err())
			return
		}
	}()

	return replyStream, nil
}
