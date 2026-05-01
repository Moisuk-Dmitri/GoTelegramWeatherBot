package bot

import (
	"context"
	"log"
	"main/internal/weather"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramApi interface {
	GetUpdatesChan(config tgbotapi.UpdateConfig) tgbotapi.UpdatesChannel
	StopReceivingUpdates()
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

type Bot struct {
	api     TelegramApi
	weather weather.Service
}

func NewBot(api TelegramApi, service weather.Service) *Bot {
	return &Bot{
		api:     api,
		weather: service,
	}
}

func (b *Bot) Run(ctx context.Context) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			b.api.StopReceivingUpdates()
			return ctx.Err()
		case update := <-updates:
			if update.Message != nil {
				b.handleMessages(ctx, update.Message)
			}
		}
	}
}

func (b *Bot) reply(chatID int64, text string) {
	message := tgbotapi.NewMessage(chatID, text)
	message.ParseMode = "HTML"

	if _, err := b.api.Send(message); err != nil {
		log.Printf("send message error: %v", err)
	}
}
