package chatmemory

import (
	"context"
	"encoding/json"
	"time"

	"github.com/openai/openai-go"
)

type ChatMemory struct {
	Conversation []ConversationMessage `json:"conversation"`
}

type ConversationMessage struct {
	Role    string    `json:"role"`
	Content string    `json:"content"`
	Name    string    `json:"name,omitempty"`
	SentAt  time.Time `json:"sent_at,omitempty"`
}

func (mem *ChatMemory) ToChatCompletionMessages(ctx context.Context) []openai.ChatCompletionMessageParamUnion {
	messages := []openai.ChatCompletionMessageParamUnion{}
	for _, msg := range mem.Conversation {
		if msg.Role == "user" {
			messages = append(messages, openai.ChatCompletionMessageParamUnion{
				OfUser: &openai.ChatCompletionUserMessageParam{
					Name: openai.String(msg.Name),
					Content: openai.ChatCompletionUserMessageParamContentUnion{
						OfString: openai.String(msg.Content),
					},
				},
			})
		} else {
			messages = append(messages, openai.AssistantMessage(msg.Content))
		}
	}
	return messages
}

func (s *Store) AppendConversationMessage(ctx context.Context, role string, content string, name string) error {
	chatMemory, err := s.GetChatMemory(ctx)
	if err != nil {
		return err
	}

	chatMemory.Conversation = append(chatMemory.Conversation, ConversationMessage{
		Role:    role,
		Content: content,
		Name:    name,
		SentAt:  time.Now(),
	})
	if len(chatMemory.Conversation) > maxConversationLength {
		chatMemory.Conversation = chatMemory.Conversation[len(chatMemory.Conversation)-maxConversationLength:]
	}

	raw, err := json.Marshal(chatMemory)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, chatMemoryKey, raw, 0).Err()
}
