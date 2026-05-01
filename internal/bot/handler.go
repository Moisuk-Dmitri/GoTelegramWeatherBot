package bot

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleMessages(ctx context.Context, msg *tgbotapi.Message) {
	if msg.IsCommand() {
		b.handleCommands(msg)
		return
	}

	b.handleWeatherMessage(ctx, msg)
}
