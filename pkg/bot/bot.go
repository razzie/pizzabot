package bot

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	name     string
	api      *tgbotapi.BotAPI
	rand     *rand.Rand
	logger   *log.Logger
	commands map[string]*Responder
	keywords map[string]*Responder
}

func Must(bot *Bot, err error) *Bot {
	if err != nil {
		panic(err)
	}
	return bot
}

func NewBot(name, token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &Bot{
		name:     name,
		api:      bot,
		rand:     rand.New(rand.NewSource(42)),
		logger:   log.Default(),
		commands: make(map[string]*Responder),
		keywords: make(map[string]*Responder),
	}, nil
}

func (bot *Bot) Seed(seed int64) {
	bot.rand.Seed(seed)
}

func (bot *Bot) SetLogger(logger *log.Logger) {
	bot.logger = logger
}

func (bot *Bot) WithCommands(commands ...string) *Responder {
	resp := &Responder{bot: bot}
	for _, cmd := range commands {
		bot.commands["/"+cmd] = resp
	}
	return resp
}

func (bot *Bot) WithKeywords(keywords ...string) *Responder {
	resp := &Responder{bot: bot}
	for _, keyword := range keywords {
		bot.keywords[keyword] = resp
	}
	return resp
}

func (bot *Bot) RunUntil(exit <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer bot.log("exited")
	bot.log("launched")
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.api.GetUpdatesChan(updateConfig)
	for {
		select {
		case update := <-updates:
			bot.handleUpdate(update)
		case <-exit:
			bot.api.StopReceivingUpdates()
			return
		}
	}
}

func (bot *Bot) Run() {
	bot.log("launched")
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.api.GetUpdatesChan(updateConfig)
	for update := range updates {
		bot.handleUpdate(update)
	}
}

func (bot *Bot) log(format string, args ...interface{}) {
	bot.logger.Printf("[%s] %s", bot.name, fmt.Sprintf(format, args...))
}

func (bot *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}
	if responder := bot.getResponder(update.Message.Text); responder != nil {
		if response := responder.respond(update.Message); response != nil {
			bot.log("responding to %s's message", update.Message.Chat.FirstName)
			if _, err := bot.api.Send(response); err != nil {
				bot.log("sending response failed: %v", err)
			}
		}
	}
}

func (bot *Bot) getResponder(msg string) *Responder {
	// command
	if strings.HasPrefix(msg, "/") {
		for cmd, responder := range bot.commands {
			if len(responder.responses) == 0 {
				continue
			}
			if cmd == msg || strings.HasPrefix(msg, cmd+"@") {
				return responder
			}
		}
		return nil
	}
	// normal message
	for _, word := range strings.Fields(msg) {
		for keyword, responder := range bot.keywords {
			if len(responder.responses) == 0 {
				continue
			}
			if strings.ToLower(word) == keyword {
				return responder
			}
		}
	}
	return nil
}
