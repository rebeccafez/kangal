package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rebeccafez/kangal/internal/config"
	"github.com/rebeccafez/kangal/internal/telegram"
)

func main() {
	cfg := config.ConfigFromEnv()
	bot := telegram.NewBot(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sig
		cancel()
	}()

	bot.Run(ctx)
}
