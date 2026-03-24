package oaiclient

import (
	"context"
	"bytes"
	"encoding/json"
	"strings"
	"fmt"
	"net/http"
	"io"
	"errors"

	"github.com/rebeccafez/kangal/internal/config"
)

type Message struct {
	Role string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model string `json:"model"`
	Messages []Message `json:"messages"`
	MaxTokens int `json:"max_tokens,omitempty"`
	Temperature float64 `json:"temperature"`
	Stream bool `json:"Stream"`
}

type ChatResponse struct {
	Choices []struct{
		Message Message `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func CallLLM(ctx context.Context, cfg config.Config, history []Message) (string, error) {
	payload := ChatRequest {
		Model: cfg.Model,
		Messages: history,
		MaxTokens: cfg.MaxTokens,
		Temperature: cfg.Temperature,
		Stream: false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	url := strings.TrimRight(cfg.OpenAIBaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w")
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: cfg.RequestTimeout}
	resp, err := client.Do(req)

	if err != nil {
		return "", fmt.Errorf("http request: %w", err)
	}

	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM returned HTTP %d: %s", resp.StatusCode, string(raw))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(raw, &chatResp); err != nil {
		return "", fmt.Errorf("unmarshal response: %w", err)
	}

	if chatResp.Error != nil {
		return "", fmt.Errorf("LLM error: %s", chatResp.Error.Message)
	}

	if len(chatResp.Choices) == 0 {
		return "", errors.New("LLM returned no choices")
	}

	return strings.TrimSpace(chatResp.Choices[0].Message.Content), nil
}
