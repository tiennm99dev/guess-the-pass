# Journal — Migrate guess-the-pass: Java → Go

**Date:** 2026-06-13
**Commit:** `de24baa` (main)
**Plan:** `plans/260613-1910-migrate-to-go/`

## What

Full cutover of the Telegram echo bot from Java 21 / Gradle / `org.telegram:telegrambots` to Go. Behavior preserved 1:1 (long-polling echo of any text message back to the same chat). No new game logic.

## Key changes

- **Stack:** Go 1.23 (`go.mod`), `github.com/go-telegram/bot` v1.21.0, stdlib `log/slog` (debug).
- **Layout:** `main.go` + `internal/bot/handler.go` (echo) + `internal/config/config.go` (env load).
- **Container:** Dockerfile rebuilt as static `CGO_ENABLED=0` binary on alpine + ca-certificates (needed for Telegram HTTPS).
- **Removed:** all `src/`, `gradle/`, `gradlew*`, `*.gradle.kts`, logback.xml. Dropped obsolete Compose `version` key.

## Decisions

- Library `go-telegram/bot` chosen over telebot / go-telegram-bot-api: zero-dep, actively maintained, native long-polling. `*bot.Bot` is both poller and send client — collapses Java's separate `GuestThePassClient` singleton.
- `go` directive pinned to 1.23 (not the installed 1.26) for broader build compatibility; Docker builder matched to `golang:1.23-alpine`.

## Notable

- Go port is **strictly safer** than the original: `models.Message.From` is a nilable pointer and is guarded, whereas the Java `getFrom().getId()` would NPE on senderless messages (channel posts). Also added fail-fast `BOT_TOKEN` validation the Java code lacked.
- Verified: `go build` / `vet` / `gofmt` / `go test` (config) / `docker compose config` all clean. Code review returned DONE, parity confirmed against Java via git.

## Unverified (manual gate)

- Live Telegram echo and `docker compose up --build` — require a real `BOT_TOKEN` + network; not run in this session.
