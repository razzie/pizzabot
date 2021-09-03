package main

type PizzaBot struct {
	StickerBot
}

func NewPizzaBot(token string) (*PizzaBot, error) {
	bot, err := NewStickerBot(token, "pizzabot", "pizza")
	if err != nil {
		return nil, err
	}

	return &PizzaBot{
		StickerBot: *bot,
	}, nil
}
