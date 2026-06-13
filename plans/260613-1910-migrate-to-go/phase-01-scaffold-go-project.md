---
phase: 1
title: Scaffold Go project
status: completed
priority: P1
effort: 30m
dependencies: []
---

# Phase 1: Scaffold Go project

## Overview

Stand up the Go module and dependency, plus config loading. No bot wiring yet — just a compilable skeleton.

## Requirements

- Functional: `go build ./...` succeeds; `BOT_TOKEN` read from env with fail-fast on empty.
- Non-functional: idiomatic Go layout, snake_case files, zero dead deps.

## Architecture

```
go.mod                          # module github.com/tiennm99/guess-the-pass, go 1.23
go.sum
internal/config/config.go       # Config struct + Load() reads BOT_TOKEN, validates non-empty
```

Config mirrors `GuestThePassEnv` but adds the validation Java lacked (Java passed a possibly-null token to the client). `Load()` returns `(Config, error)`; empty token → error.

## Related Code Files

- Create: `go.mod`, `go.sum`
- Create: `internal/config/config.go`
- (Java sources remain until Phase 3 cutover — do not delete yet, keeps repo buildable mid-migration is N/A since languages are independent, but defer deletion to one atomic step.)

## Implementation Steps

1. `cd /config/workspace/tiennm99/guess-the-pass && go mod init github.com/tiennm99/guess-the-pass`
2. `go get github.com/go-telegram/bot@latest` (populates go.mod/go.sum; pin the resolved version).
3. Create `internal/config/config.go`:
   ```go
   package config

   import (
       "errors"
       "os"
   )

   // Config holds runtime settings sourced from the environment.
   type Config struct {
       BotToken string
   }

   // Load reads configuration from environment variables.
   // BOT_TOKEN is required; an empty value is a fatal misconfiguration.
   func Load() (Config, error) {
       token := os.Getenv("BOT_TOKEN")
       if token == "" {
           return Config{}, errors.New("BOT_TOKEN is required")
       }
       return Config{BotToken: token}, nil
   }
   ```
4. `go build ./...` — confirm clean compile.

## Success Criteria

- [ ] `go.mod` declares module `github.com/tiennm99/guess-the-pass` and go 1.23+.
- [ ] `go-telegram/bot` present in go.mod with a pinned version, recorded in go.sum.
- [ ] `internal/config/config.go` compiles; `Load()` errors on empty `BOT_TOKEN`.
- [ ] `go build ./...` exits 0.

## Risk Assessment

- Toolchain mismatch: installed Go 1.26 > target 1.23 — fine. Keep go.mod `go` directive at 1.23 for broad compatibility unless go-telegram/bot requires higher.
