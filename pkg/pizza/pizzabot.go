package pizza

import (
	_ "embed"
	"strings"

	"github.com/razzie/pizzabot/pkg/bot"
)

//go:embed pizza_facts.txt
var pizzaFacts string
var PizzaFacts = strings.Split(pizzaFacts, "\n")

func NewPizzaBot(token string) *bot.Bot {
	return bot.Must(bot.NewBot("pizzabot", token)).
		WithCommands("pizza").
		WithStickerSet("pizzabot").
		WithLines(PizzaFacts...).
		Done().
		WithKeywords("pizza", "pizzas", "pizzák", "pizzát", "pizzákat", "pizzás").
		WithStickerSet("pizzabot").
		WithLines(PizzaFacts...).
		Done()
}
