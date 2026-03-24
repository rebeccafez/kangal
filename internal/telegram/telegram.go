package telegram

import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"
	"context"
	"errors"
	"bytes"
	"log"
	"strings"

	"github.com/rebeccafez/kangal/internal/config"
	"github.com/rebeccafez/kangal/internal/conversationstore"
	"github.com/rebeccafez/kangal/internal/oaiclient"
)
type TelegramUpdate struct {
	UpdateID int `json:"update_id"`
	Message *struct {
		MessageID int `json:"message_id"`
		From *struct {
			ID int64 `json:"id"`
			FirstName string `json:"first_name"`
			Username string `json:"username"`
		} `json:"from"`
		Chat struct {
			ID int64 `json:"id"`
		} `json:"chat"`
		Text string `json:"text"`
	} `json:"message"`
}

type TelegramGetUpdatesResponse struct {
	OK bool `json:"ok"`
	Result []TelegramUpdate `json:"result"`
}

type Bot struct {
	cfg config.Config
	store *conversationstore.ConversationStore
	apiURL string
	client *http.Client
}

func NewBot(cfg config.Config) *Bot {
	return &Bot{
		cfg: cfg,
		store: conversationstore.NewConversationStore(cfg.SystemPrompt),
		apiURL: "https://api.telegram.org/bot" + cfg.TelegramToken,
		client: &http.Client{ Timeout: 30 * time.Second },
	}
}

func (b *Bot) getUpdates(offset int) ([]TelegramUpdate, error) {
	url := fmt.Sprintf("%s/getUpdates?timeout=30&offset=%d", b.apiURL, offset)
	resp, err := b.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result TelegramGetUpdatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.OK {
		return nil, errors.New("telegram getUpdates returned not OK")
	}

	return result.Result, nil
}

func (b *Bot) sendMessage(chatID int64, text string) error {
	chunks := splitMessage(text, 4000)

	for _, chunk := range chunks {
		payload := map[string]interface{}{
			"chat_id": chatID,
			"text": chunk,
			"parse_mode": "Markdown",
		}

		body, err := json.Marshal(payload)

		if err != nil {
			return err
		}

		resp, err := b.client.Post(b.apiURL+"/sendMessage", "application/json", bytes.NewBuffer(body))

		if err != nil {
			return err
		}

		resp.Body.Close()
	}

	return nil
}

func (b *Bot) sendTyping(chatID int64) error {
	payload := map[string]interface{}{"chat_id": chatID, "action": "typing"}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	resp, err := b.client.Post(b.apiURL+"/sendChatAction", "application/json", bytes.NewReader(body))

	if err == nil {
		resp.Body.Close()
	}

	return nil
}

func (b *Bot) isAllowed(userID int64) bool {
	if len(b.cfg.AllowedUserIDs) == 0 {
		return true
	}

	for _, id := range b.cfg.AllowedUserIDs {
		if id == userID {
			return true
		}
	}

	return false
}

func (b *Bot) handleMessage(update TelegramUpdate) {
	msg := update.Message

	if msg == nil || msg.Text == "" {
		return
	}

	chatID := msg.Chat.ID
	userID := int64(0)

	if msg.From != nil {
		userID = msg.From.ID
	}

	if !b.isAllowed(userID) {
		log.Printf("Blocked user %d in chat %d", userID, chatID)
		b.sendMessage(userID, "You are not authorized to use this bot")
		return
	}

	text := strings.TrimSpace(msg.Text)

	switch {
	case text == "/start":
			 b.sendMessage(chatID,	fmt.Sprintf("Hello! I'm connected to %s via LLaMa.cpp. Just send a message to chat.\n\nCommands:\n* /reset - clear conversation history\n* /model - show current model\n* /help - show this help message", b.cfg.Model))
			return
	

		case text == "/help":
			 b.sendMessage(chatID, "Commands:\n* /reset - clear conversation history\n* /model - show current model\n* /help - show this help message")
			return

		case text == "/reset":
			b.store.Reset(chatID)
			 b.sendMessage(chatID, "Conversation history cleared.")
			return

		case text == "/model":
			 b.sendMessage(chatID, fmt.Sprintf("Model: %s", b.cfg.Model))
			return
	}

	b.store.Append(chatID, oaiclient.Message{Role: "user", Content: text})
	b.sendTyping(chatID)

	ctx, cancel := context.WithTimeout(context.Background(), b.cfg.RequestTimeout+5*time.Second)
	defer cancel()

	reply, err := oaiclient.CallLLM(ctx, b.cfg, b.store.Get(chatID))
	if err != nil {
		log.Printf("LLM error for chat %d: %v", chatID, err)
		b.store.Reset(chatID)
		 b.sendMessage(chatID, fmt.Sprintf("Error: %v", err))

		return
	}

	b.store.Append(chatID, oaiclient.Message{Role: "assistant", Content: reply})
	 b.sendMessage(chatID, reply)
}

func (b *Bot) Run(ctx context.Context) {
	log.Printf("Bot started. Model: %s @ %s", b.cfg.Model, b.cfg.OpenAIBaseURL)
	offset := 0

	for {
		select {
		case <- ctx.Done():
			log.Printf("Bot shutting down.")
			return
		default:
		}

		updates, err := b.getUpdates(offset)
		if err != nil {
			log.Printf("getUpdates error: %v - retrying in 5s", err)
			time.Sleep(5 * time.Second)
			continue
		}

		for _, u := range updates {
			offset = u.UpdateID + 1
			go b.handleMessage(u)
		}
	}
}

func splitMessage(text string, maxLen int) []string {
	if len(text) <= maxLen {
		return []string{text}
	}

	var chunks []string

	for len(text) > maxLen {
		chunks = append(chunks, text[:maxLen])
		text = text[maxLen:]
	}

	if text != "" {
		chunks = append(chunks, text)
	}

	return chunks
}
