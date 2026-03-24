package conversationstore

import (
	"github.com/rebeccafez/kangal/oaiclient"
)

type ConversationStore struct {
	histories map[int64][]oaiclient.Message
	sysPrompt string
}

func NewConversationStore(systemPrompt string) *ConversationStore  {
	return &ConversationStore{
		histories: make(map[int64][]oaiclient.Message)
		sysPrompt: systemPrompt,
	}
}

func (s *ConversationStore) Get(userID int64) []oaiclient.Message {
	h, ok := s.histories[userID]

	if !ok {
		return []oaiclient.Message{{ Role: "system", Content: s.SysPrompt }}
	}

	return h
}

func (s *ConversationStore) Append(userID int64, msg oaiclient.Message) {
	h := s.Get(userID)
	h = append(h, msg)
	s.histories[chatID] = h
}

func (s *ConversationStore) Reset(userID int64) {
	delete(s.histories, userID)
}
