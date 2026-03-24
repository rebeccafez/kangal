package main

import (
	"fmt"
	"context"

	"github.com/rebeccafez/kangal/internal/config"
	"github.com/rebeccafez/kangal/internal/oaiclient"
	"github.com/rebeccafez/kangal/internal/conversationstore"
)

func main() {
	cfg := config.ConfigFromEnv()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	store := conversationstore.NewConversationStore(cfg.SystemPrompt)

	store.Append(1, oaiclient.Message{Role: "user", Content: "What continent do capybaras live on?"})

	history := store.Get(1)
	resp, err := oaiclient.CallLLM(ctx, cfg, history)

	if err != nil {
	}

	fmt.Printf("%v", resp)
}
