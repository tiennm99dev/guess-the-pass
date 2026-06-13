---
phase: 3
title: Build & Docker
status: completed
priority: P1
effort: 30m
dependencies:
  - 2
---

# Phase 3: Build & Docker

## Overview

Replace the Gradle/JDK build and container with a Go build. Delete all Java/Gradle artifacts in one atomic cutover. Rewrite the multi-stage Dockerfile to build a static Go binary on a minimal runtime image.

## Requirements

- Functional: `docker compose up` builds and runs the Go bot; `BOT_TOKEN` injected via env (unchanged compose contract).
- Non-functional: small final image (scratch/alpine), reproducible build, no Java toolchain anywhere.

## Architecture

Multi-stage Docker: `golang:1.26-alpine` builder → `CGO_ENABLED=0 go build` static binary → copied into `alpine` (or `scratch` + ca-certificates) runtime. `ca-certificates` required for HTTPS to the Telegram API.

## Related Code Files

- Modify: `Dockerfile` (Gradle/Corretto → Go multi-stage)
- Modify: `docker-compose.yml` (drop `version` key — obsolete in Compose v2; keep service/env/restart)
- Delete: `build.gradle.kts`, `settings.gradle.kts`, `gradlew`, `gradlew.bat`, `gradle/` dir
- Delete: `src/` (entire Java tree incl. `logback.xml`)
- Modify: `.gitignore` (Java/Gradle ignores → Go: built binary, etc.)

## Implementation Steps

1. Rewrite `Dockerfile`:
   ```dockerfile
   FROM golang:1.26-alpine AS builder
   WORKDIR /app
   COPY go.mod go.sum ./
   RUN go mod download
   COPY . .
   RUN CGO_ENABLED=0 go build -o /app/guess-the-pass .

   FROM alpine:3.20
   RUN apk add --no-cache ca-certificates
   WORKDIR /app
   COPY --from=builder /app/guess-the-pass /app/guess-the-pass
   CMD ["/app/guess-the-pass"]
   ```
2. Edit `docker-compose.yml`: remove the `version: '3.8'` line (deprecated in Compose v2); leave `services.app` build/container_name/environment/restart intact.
3. Delete Java/Gradle files:
   ```bash
   rm -rf src gradle build .gradle
   rm -f build.gradle.kts settings.gradle.kts gradlew gradlew.bat
   ```
4. Update `.gitignore` for Go (drop `build/`, `.gradle/`; add `/guess-the-pass` binary, optional `vendor/`). Keep existing IDE/OS entries.
5. `docker compose build` — confirm image builds.

## Success Criteria

- [ ] `Dockerfile` builds a static Go binary on a minimal runtime with ca-certificates.
- [ ] `docker compose build` succeeds; no Gradle/JDK references remain.
- [ ] All Java/Gradle files deleted; `src/` gone.
- [ ] `.gitignore` reflects a Go project.
- [ ] `docker compose config` is valid with no `version` warning.

## Risk Assessment

- Missing ca-certificates in runtime image → TLS handshake to api.telegram.org fails at runtime. Mitigated by `apk add ca-certificates` (or copy certs if using scratch).
- Deleting `src/` is destructive but intentional (user chose full replacement); the original Java is recoverable from git history.
