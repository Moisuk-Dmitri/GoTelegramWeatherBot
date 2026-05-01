package main

import (
	"context"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"

	"main/internal/bot"
	"main/internal/weather"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env load fail")
	}

	weatherService := weather.NewService(os.Getenv("WeatherApiURL"))
	token := os.Getenv("TelegramBotApiKey")

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Authorized on account @%s", api.Self.UserName)

	bot := bot.NewBot(api, weatherService)
	ctx := context.Background()
	if err := bot.Run(ctx); err != nil {
		log.Printf("bot stopped: %v", err)
	}
}
