# guess-the-pass

Telegram bot game where players guess a secret password — written in Go using the [go-telegram/bot](https://github.com/go-telegram/bot) library.

## Quick start

```bash
# Set environment variable
export BOT_TOKEN=your_telegram_bot_token

go run .
```

Or build a binary:

```bash
go build -o guess-the-pass .
./guess-the-pass
```

Or with Docker:

```bash
docker compose up
```

## Gameplay

The bot receives messages from players in a Telegram chat. Players send their guesses as plain text messages; the bot echoes back responses. The secret password is configured server-side — when a player's message matches, they win.

Sample session:
```
User:  hello
Bot:   hello
User:  secret123
Bot:   secret123
```

> The bot currently mirrors messages (echo mode). Extend `Echo` in `internal/bot/handler.go` to add win-condition logic.

## Configuration

| Variable | Required | Description |
|----------|----------|-------------|
| `BOT_TOKEN` | Yes | Telegram bot token from [@BotFather](https://t.me/BotFather) |

## License

Apache-2.0 — see [LICENSE](LICENSE).
