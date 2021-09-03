package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func launchBot(token string, exit <-chan struct{}) {
	bot, err := NewPizzaBot(token)
	if err != nil {
		log.Fatal(err)
	}

	go bot.Run(exit)
}

func waitForSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-sigChan
}

func main() {
	exitChan := make(chan struct{})

	tokens := os.Args[1:]
	for _, token := range tokens {
		launchBot(token, exitChan)
	}
	log.Printf("Launched %d bots", len(tokens))

	waitForSignal()
	close(exitChan)
}
