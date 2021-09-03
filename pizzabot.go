package main

import (
	"flag"
	"log"
	"math/rand"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	token := flag.String("token", "", "Telegram bot API token")
	flag.Parse()

	rand.Seed(42)

	bot, err := tgbotapi.NewBotAPI(*token)
	if err != nil {
		log.Fatal(err)
	}

	//bot.Debug = true
	pizzaStickerSet, err := bot.GetStickerSet(tgbotapi.GetStickerSetConfig{Name: "pizzabot"})
	if err != nil {
		log.Fatal(err)
	}
	pizzaStickers := pizzaStickerSet.Stickers

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		isAboutPizza := false
		for _, word := range strings.Fields(update.Message.Text) {
			if strings.ToLower(word) == "pizza" {
				isAboutPizza = true
				break
			}
		}
		if !isAboutPizza {
			continue
		}

		pizza := pizzaStickers[rand.Intn(len(pizzaStickers))]
		msg := tgbotapi.NewStickerShare(update.Message.Chat.ID, pizza.FileID)
		msg.ReplyToMessageID = update.Message.MessageID
		msg.DisableNotification = true
		if _, err := bot.Send(msg); err != nil {
			log.Println(err)
		}
	}
}
