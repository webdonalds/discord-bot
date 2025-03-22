package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/openai/openai-go"
	"github.com/webdonalds/discord-bot/internal/actions"
	"github.com/webdonalds/discord-bot/internal/chatmemory"
)

type ChatCommand struct {
	client      openai.Client
	memoryStore *chatmemory.Store
}

func NewChatCommand(apiKey string, memoryStore *chatmemory.Store) *ChatCommand {
	return &ChatCommand{
		client:      openai.NewClient(),
		memoryStore: memoryStore,
	}
}

func (c *ChatCommand) Commands() []string {
	return []string{"chat"}
}

func (c *ChatCommand) Execute(ctx context.Context, args *CommandArgs) (actions.Action, error) {
	if len(args.Texts) == 0 {
		return &actions.Reply{Message: "문장을 입력해주세요"}, nil
	}

	prompt := strings.Join(args.Texts, " ")
	c.memoryStore.AppendConversationMessage(ctx, "user", prompt, args.Author.Name)

	messageStream := make(chan string)
	replyStream := &actions.ReplyStream{
		MessageStream: messageStream,
	}

	// Start the reply stream in a goroutine
	go func() {
		defer close(messageStream)

		chatMemory, err := c.memoryStore.GetChatMemory(ctx)
		if err != nil {
			messageStream <- fmt.Sprintf("기억을 꺼내오던 중 오류가 발생했어요.\n\n%v", err)
			return
		}

		messages := chatMemory.ToChatCompletionMessages(ctx)
		messages = append(messages, openai.DeveloperMessage("당신은 메신저에서 대화중이니 가능한 간결하게 답변하세요."))

		stream := c.client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
			Messages:            messages,
			Model:               "gpt-4o-mini",
			MaxCompletionTokens: openai.Int(200),
		})

		acc := openai.ChatCompletionAccumulator{}
		for stream.Next() {
			chunk := stream.Current()
			if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
				acc.AddChunk(chunk)
				messageStream <- chunk.Choices[0].Delta.Content
			}
		}

		if stream.Err() != nil {
			messageStream <- fmt.Sprintf("\n\n답변 중 오류가 발생했어요.\n\n%v", stream.Err())
			return
		}

		c.memoryStore.AppendConversationMessage(ctx, "assistant", acc.Choices[0].Message.Content, "Bot")
	}()

	return replyStream, nil
}
