package bot

import (
	"context"
	"log/slog"

	tgbot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Echo mirrors any text message back to the originating chat.
// Preserves the original echo behavior of the Java GuessThePassBot.consume.
func Echo(ctx context.Context, b *tgbot.Bot, update *models.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	msg := update.Message

	var userID int64
	if msg.From != nil {
		userID = msg.From.ID
	}
	slog.Debug("received message",
		"text", msg.Text,
		"userId", userID,
		"chatId", msg.Chat.ID)

	if _, err := b.SendMessage(ctx, &tgbot.SendMessageParams{
		ChatID: msg.Chat.ID,
		Text:   msg.Text,
	}); err != nil {
		slog.Error("send message failed", "err", err)
	}
}
