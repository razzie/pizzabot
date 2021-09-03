package main

import (
	_ "embed"
	"strings"
)

//go:embed pizza_facts.txt
var pizzaFacts string
var PizzaFacts = strings.Split(pizzaFacts, "\n")

func NewPizzaBot(token string) *Bot {
	return Must(NewBot("pizzabot", token)).
		WithCommands("pizza").
		WithStickerSet("pizzabot").
		WithLines(PizzaFacts...).
		Done().
		WithKeywords("pizza", "pizzák", "pizzát", "pizzákat", "pizzás").
		WithStickerSet("pizzabox").
		WithLines(PizzaFacts...).
		Done()
}
