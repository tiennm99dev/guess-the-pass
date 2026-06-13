---
title: Migrate guess-the-pass from Java to Go
description: >-
  1:1 port of the echo Telegram bot from Java/Gradle/telegrambots to Go using
  go-telegram/bot + slog. Full cutover, Java removed.
status: completed
priority: P2
branch: main
tags:
  - migration
  - go
  - telegram
blockedBy: []
blocks: []
created: '2026-06-13T12:21:05.107Z'
createdBy: 'ck:plan'
source: skill
---

# Migrate guess-the-pass from Java to Go

## Overview

Port the existing Telegram echo bot (4 Java files, ~70 LOC) from Java 21 / Gradle / `org.telegram:telegrambots` to Go. Behavior is preserved exactly: long-polling bot that echoes any text message back to the same chat. No new game/win-condition logic.

**Decisions (user-confirmed):**
- Library: `github.com/go-telegram/bot` (modern, zero-dep, native long-polling)
- Logging: `log/slog` (stdlib structured logging)
- Scope: echo only — exact 1:1 behavioral port
- Java code: replaced entirely (clean cutover; repo becomes a Go project)

**Module path:** `github.com/tiennm99/guess-the-pass` (matches git remote)
**Go toolchain:** 1.26 installed; target `go 1.23` in go.mod (minimum for go-telegram/bot).

### Java → Go mapping

| Java | Go |
|------|-----|
| `Main.java` (bootstrap long-polling app) | `main.go` (wire config + bot, `b.Start(ctx)`) |
| `GuessThePassBot.consume()` (echo handler) | `internal/bot/handler.go` (default handler) |
| `GuestThePassClient` (OkHttp singleton) | (absent — `*bot.Bot` is the client) |
| `GuestThePassEnv.BOT_TOKEN` | `internal/config/config.go` (`Load()` reads env) |
| `logback.xml` (debug root) | slog text handler at debug in `main.go` |
| `build.gradle.kts` / Gradle wrapper | `go.mod` / `go.sum` |

## Phases

| Phase | Name | Status |
|-------|------|--------|
| 1 | [Scaffold Go project](./phase-01-scaffold-go-project.md) | Completed |
| 2 | [Port bot logic](./phase-02-port-bot-logic.md) | Completed |
| 3 | [Build & Docker](./phase-03-build-docker.md) | Completed |
| 4 | [Verify & docs](./phase-04-verify-docs.md) | Completed |

## Dependencies

- Go 1.23+ toolchain (1.26 present).
- `github.com/go-telegram/bot` (fetched via `go get`).
- Network access to Telegram API for live smoke test (Phase 4).
- No cross-plan dependencies.
