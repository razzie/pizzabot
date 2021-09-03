package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func launchBot(token string, exit <-chan struct{}, wg *sync.WaitGroup) {
	bot, err := NewPizzaBot(token)
	if err != nil {
		log.Fatal(err)
	}

	wg.Add(1)
	go bot.RunUntil(exit, wg)
}

func waitForSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-sigChan
}

func main() {
	log.SetOutput(os.Stdout)

	var wg sync.WaitGroup
	exitChan := make(chan struct{})
	tokens := os.Args[1:]
	for _, token := range tokens {
		launchBot(token, exitChan, &wg)
	}
	log.Printf("Launched %d bot(s)", len(tokens))

	waitForSignal()
	close(exitChan)
	wg.Wait()
}
