package main

func NewPizzaBot(token string) *Bot {
	return Must(NewBot("pizzabot", token)).
		WithCommands("pizza").
		WithStickerSet("pizzabot").
		Done().
		WithKeywords("pizza").
		WithStickerSet("pizzabox").
		Done()
}
