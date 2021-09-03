package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type response interface {
	ApplyTo(msg *tgbotapi.Message) tgbotapi.Chattable
}

type textResponse struct {
	text string
}

func (r *textResponse) ApplyTo(msg *tgbotapi.Message) tgbotapi.Chattable {
	respMsg := tgbotapi.NewMessage(msg.Chat.ID, r.text)
	respMsg.ReplyToMessageID = msg.MessageID
	respMsg.DisableNotification = true
	return respMsg
}

type stickerResponse struct {
	sticker tgbotapi.Sticker
}

func (r *stickerResponse) ApplyTo(msg *tgbotapi.Message) tgbotapi.Chattable {
	respMsg := tgbotapi.NewStickerShare(msg.Chat.ID, r.sticker.FileID)
	respMsg.ReplyToMessageID = msg.MessageID
	respMsg.DisableNotification = true
	return respMsg
}
