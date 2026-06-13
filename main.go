package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	tgbot "github.com/go-telegram/bot"

	"github.com/tiennm99/guess-the-pass/internal/bot"
	"github.com/tiennm99/guess-the-pass/internal/config"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	b, err := tgbot.New(cfg.BotToken, tgbot.WithDefaultHandler(bot.Echo))
	if err != nil {
		slog.Error("bot init failed", "err", err)
		os.Exit(1)
	}

	slog.Info("bot starting")
	b.Start(ctx) // blocks until ctx cancelled (long polling)
}
