---
phase: 4
title: Verify & docs
status: completed
priority: P2
effort: 30m
dependencies:
  - 3
---

# Phase 4: Verify & docs

## Overview

Prove the port works end-to-end and update README to reflect the Go stack. Optional minimal unit test for the config loader.

## Requirements

- Functional: bot connects to Telegram and echoes a real message; README commands work as written.
- Non-functional: no stale Java/Gradle references in docs.

## Related Code Files

- Modify: `README.md` (Quick start, Configuration, stack description)
- Create (optional): `internal/config/config_test.go`

## Implementation Steps

1. **Static checks:** `go build ./...`, `go vet ./...`, `gofmt -l .` (expect empty output).
2. **Optional unit test** — `internal/config/config_test.go` covering `Load()` empty vs set `BOT_TOKEN` (use `t.Setenv`). Run `go test ./...`.
3. **Live smoke test** (needs real token + network):
   ```bash
   export BOT_TOKEN=<token>
   go run .
   ```
   Send a message to the bot in Telegram; confirm it echoes and a debug log line appears. Ctrl-C → clean shutdown.
4. **Docker smoke test:** `docker compose up --build`; repeat echo check; `docker compose down`.
5. **Update README.md:**
   - Line 3: change "written in Java using the Telegrambots library" → Go + go-telegram/bot.
   - Quick start: replace `./gradlew run` with `go run .` (and `go build` note); keep `docker compose up`.
   - Gameplay note line 32: update `GuessThePassBot.consume()` reference → `internal/bot/handler.go` `Echo`.
   - Configuration table (`BOT_TOKEN`) unchanged — contract preserved.

## Success Criteria

- [ ] `go build`, `go vet`, `gofmt -l .` all clean.
- [ ] (If added) `go test ./...` passes.
- [ ] Live echo verified via `go run .` against a real bot token.
- [ ] `docker compose up --build` runs and echoes.
- [ ] README has zero Java/Gradle references; commands match the Go project.

## Risk Assessment

- Live test requires a valid `BOT_TOKEN` and outbound network; if unavailable in CI/dev, document as a manual gate and rely on build/vet/test for automated confidence.

## Next Steps

- Optional follow-ups (out of scope, not requested): implement the actual guess/win-condition game logic; add CI workflow (`go build`/`test`); add golangci-lint.
