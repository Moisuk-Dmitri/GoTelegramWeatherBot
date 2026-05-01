package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func cityKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Москва"),
			tgbotapi.NewKeyboardButton("Санкт-Петербург"),
		),
	)
}

func intervalKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Сейчас"),
			tgbotapi.NewKeyboardButton("День"),
			tgbotapi.NewKeyboardButton("Неделя"),
		),
	)
}

func (b *Bot) replyWithKeyboard(chatID int64, text string, keyboard tgbotapi.ReplyKeyboardMarkup) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}

func (b *Bot) replyRemoveKeyboard(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	b.api.Send(msg)
}
