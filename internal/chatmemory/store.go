package chatmemory

import (
	"context"
	"encoding/json"
	"slices"
	"time"

	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{
		client: client,
	}
}

const (
	chatMemoryKey         = "chatmemory"
	maxConversationLength = 100
	maxConversationAge    = 24 * time.Hour
)

var emptyMemory = &ChatMemory{
	Conversation: []ConversationMessage{},
}

func (s *Store) GetChatMemory(ctx context.Context) (*ChatMemory, error) {
	raw, err := s.client.Get(ctx, chatMemoryKey).Result()
	if err == redis.Nil {
		return emptyMemory, nil
	}
	if err != nil {
		return emptyMemory, err
	}

	var chatMemory ChatMemory
	err = json.Unmarshal([]byte(raw), &chatMemory)
	if err != nil {
		return emptyMemory, err
	}

	// 일정 시간이 지난 메시지는 삭제
	now := time.Now()
	chatMemory.Conversation = slices.DeleteFunc(chatMemory.Conversation, func(msg ConversationMessage) bool {
		return now.Sub(msg.SentAt) >= maxConversationAge
	})

	return &chatMemory, nil
}
