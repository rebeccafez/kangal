package config

import "time"

type Config struct {
	TelegramToken  string
	OpenAIBaseURL  string
	Model          string
	SystemPrompt   string
	MaxTokens      int
	Temperature    float64
	AllowedUserIDs []int64
	RequestTimeout time.Duration
}

func ConfigFromEnv() Config {
	cfg := Config{
		TelegramToken:  mustEnv("TELEGRAM_TOKEN"),
		OpenAIBaseURL:  getEnv("LLM_BASE_URL", "http://localhost:11434"),
		Model:          getEnv("LLM_MODEL", "llama3"),
		SystemPrompt:   getEnv("LLM_SYSTEM_PROMPT", "You are a helpful assistant"),
		MaxTokens:      envInt("LLM_MAX_TOKENS", 4096),
		Temperature:    envFloat("LLM_TEMPERATURE", 0.7),
		RequestTimeout: time.Duration(envInt("LLM_TIMEOUT_SECONDS", 120)) * time.Second,
	}

	return cfg
}
