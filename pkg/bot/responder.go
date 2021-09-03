package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Responder struct {
	bot       *Bot
	responses []response
}

func (rr *Responder) WithLines(lines ...string) *Responder {
	responses := make([]response, len(lines))
	for i, line := range lines {
		responses[i] = &textResponse{text: line}
	}
	rr.responses = append(rr.responses, responses...)
	return rr
}

func (rr *Responder) WithStickerSet(stickerSetName string) *Responder {
	stickerSet, err := rr.bot.api.GetStickerSet(tgbotapi.GetStickerSetConfig{Name: stickerSetName})
	if err != nil {
		rr.bot.log("%v", err)
		return rr
	}
	responses := make([]response, len(stickerSet.Stickers))
	for i, sticker := range stickerSet.Stickers {
		responses[i] = &stickerResponse{sticker: sticker}
	}
	rr.responses = append(rr.responses, responses...)
	return rr
}

func (rr *Responder) Done() *Bot {
	return rr.bot
}

func (rr *Responder) respond(msg *tgbotapi.Message) tgbotapi.Chattable {
	if len(rr.responses) == 0 {
		rr.bot.log("trying to respond with no available responses")
		return nil
	}
	response := rr.responses[rr.bot.rand.Intn(len(rr.responses))]
	return response.ApplyTo(msg)
}
