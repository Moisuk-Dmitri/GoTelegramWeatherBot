package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) handleCommands(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		b.replyWithKeyboard(msg.Chat.ID, "Выберите город:", cityKeyboard())
	case "help":
		b.reply(msg.Chat.ID, `Доступные команды:
		/start`)
	default:
		b.reply(msg.Chat.ID, `Неизвестная команда, для подсказки напишите /help`)
	}
}
