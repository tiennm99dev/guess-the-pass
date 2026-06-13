---
phase: 2
title: Port bot logic
status: completed
priority: P1
effort: 45m
dependencies:
  - 1
---

# Phase 2: Port bot logic

## Overview

Port the echo handler and the bot bootstrap. After this phase the bot runs and mirrors text messages, matching `GuessThePassBot.consume()` exactly.

## Requirements

- Functional: any incoming text message is echoed back to the same chat; non-text/empty updates ignored. Debug log line per received message (text, userId, chatId) matching Java's log statement.
- Non-functional: graceful shutdown on SIGINT/SIGTERM; structured logging via slog at debug level.

## Architecture

```
internal/bot/handler.go    # Echo(ctx, b, update) — the default-handler func
main.go                    # slog setup, config.Load(), bot.New(...), b.Start(ctx)
```

go-telegram/bot model: `bot.New(token, bot.WithDefaultHandler(fn))` then `b.Start(ctx)` runs long polling internally (replaces `TelegramBotsLongPollingApplication`). `*bot.Bot` is both the poller and the send client — collapses Java's separate `GuestThePassClient` singleton.

Handler signature: `func(ctx context.Context, b *bot.Bot, update *models.Update)`.

## Related Code Files

- Create: `internal/bot/handler.go`
- Create: `main.go`

## Implementation Steps

1. Create `internal/bot/handler.go`:
   ```go
   package bot

   import (
       "context"
       "log/slog"

       tgbot "github.com/go-telegram/bot"
       "github.com/go-telegram/bot/models"
   )

   // Echo mirrors any text message back to the originating chat.
   // Mirrors the original GuessThePassBot.consume echo behavior.
   func Echo(ctx context.Context, b *tgbot.Bot, update *models.Update) {
       if update.Message == nil || update.Message.Text == "" {
           return
       }

       msg := update.Message
       slog.Debug("received message",
           "text", msg.Text,
           "userId", msg.From.ID,
           "chatId", msg.Chat.ID)

       if _, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
           ChatID: msg.Chat.ID,
           Text:   msg.Text,
       }); err != nil {
           slog.Error("send message failed", "err", err)
       }
   }
   ```
2. Create `main.go`:
   ```go
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
   ```
3. Verify the exact `WithDefaultHandler` / `SendMessageParams` / `models.Update` field names against the resolved go-telegram/bot version (API is stable but confirm `From.ID`, `Chat.ID`). Use `go doc github.com/go-telegram/bot` if unsure.
4. `go build ./...` and `go vet ./...` — both clean.

## Success Criteria

- [ ] `internal/bot/handler.go` echoes text, ignores nil/empty messages, logs debug line with text/userId/chatId.
- [ ] `main.go` loads config, sets slog debug default, starts long polling, shuts down cleanly on SIGINT/SIGTERM.
- [ ] `go build ./...` and `go vet ./...` exit 0.

## Risk Assessment

- API field-name drift between go-telegram/bot versions → mitigate via `go doc` / compile check in step 3-4.
- `update.Message.From` can be nil for some update types; text messages always have From, but the nil-text guard plus accessing From only after it is reachable in practice. If vet/lint flags, add a `msg.From != nil` guard.
