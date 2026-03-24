package oaiclient

type Message struct {
	Role string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest {
	Model string `json:"model"`
	Messages []Message `json:"Messages"`
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


