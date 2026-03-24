package main

import (
	"fmt"
	"context"

	"github.com/rebeccafez/kangal/internal/config"
	"github.com/rebeccafez/kangal/internal/oaiclient"
)

func main() {
	cfg := config.ConfigFromEnv()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	history := []oaiclient.Message{ oaiclient.Message{Role: "system", Content: cfg.SystemPrompt}, oaiclient.Message{Role: "user", Content: "What continent do kangaroos live on?"} }

	resp, err := oaiclient.CallLLM(ctx, cfg, history)

	if err != nil {
	}

	fmt.Printf("%v", resp)
}
