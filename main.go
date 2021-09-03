package main

import (
	"flag"
	"log"
)

func main() {
	token := flag.String("token", "", "Telegram bot API token")
	flag.Parse()

	bot, err := NewPizzaBot(*token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Run()
}
