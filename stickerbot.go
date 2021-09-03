package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StickerBot struct {
	name     string
	bot      *tgbotapi.BotAPI
	stickers []tgbotapi.Sticker
	keywords []string
	rand     *rand.Rand
	logger   *log.Logger
}

func NewStickerBot(token, stickerSetName string, keywords ...string) (*StickerBot, error) {
	if len(keywords) == 0 {
		return nil, fmt.Errorf("no keywords")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	stickerSet, err := bot.GetStickerSet(tgbotapi.GetStickerSetConfig{Name: stickerSetName})
	if err != nil {
		return nil, err
	}

	if len(stickerSet.Stickers) == 0 {
		return nil, fmt.Errorf("no stickers in stickerset: %s", stickerSetName)
	}

	return &StickerBot{
		name:     stickerSetName,
		bot:      bot,
		stickers: stickerSet.Stickers,
		keywords: keywords,
		rand:     rand.New(rand.NewSource(42)),
		logger:   log.Default(),
	}, nil
}

func (bot *StickerBot) Seed(seed int64) {
	bot.rand.Seed(seed)
}

func (bot *StickerBot) SetLogger(logger *log.Logger) {
	bot.logger = logger
}

func (bot *StickerBot) Run(exit <-chan struct{}) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.bot.GetUpdatesChan(updateConfig)
	for {
		select {
		case update := <-updates:
			bot.handleUpdate(update)
		case <-exit:
			return
		}
	}
}

func (bot *StickerBot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	if !bot.isMsgAboutKeywords(update.Message) {
		return
	}
	sticker := bot.getSticker()
	msg := tgbotapi.NewStickerShare(update.Message.Chat.ID, sticker.FileID)
	msg.ReplyToMessageID = update.Message.MessageID
	msg.DisableNotification = true
	if _, err := bot.bot.Send(msg); err != nil {
		bot.logger.Printf("[%s] %v", bot.name, err)
	}
}

func (bot *StickerBot) getSticker() tgbotapi.Sticker {
	return bot.stickers[bot.rand.Intn(len(bot.stickers))]
}

func (bot *StickerBot) isMsgAboutKeywords(msg *tgbotapi.Message) bool {
	// command
	if strings.HasPrefix(msg.Text, "/") {
		for _, keyword := range bot.keywords {
			if msg.Text == "/"+keyword {
				return true
			}
			if strings.HasPrefix(msg.Text, "/"+keyword+"@") {
				return true
			}
		}
		return false
	}
	// regular message
	for _, word := range strings.Fields(msg.Text) {
		for _, keyword := range bot.keywords {
			if strings.ToLower(word) == keyword {
				return true
			}
		}
	}
	return false
}
