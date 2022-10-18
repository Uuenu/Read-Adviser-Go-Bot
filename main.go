package main

import (
	"flag"
	"log"

	tgClient "telegram-bot-go/clients/telegram"
	eventconsumer "telegram-bot-go/consumer/event-consumer"
	"telegram-bot-go/events/telegram"
	"telegram-bot-go/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {

	eventsProccessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		files.New(storagePath),
	)

	log.Print("service started")

	consumer := eventconsumer.New(eventsProccessor, eventsProccessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal()
	}

}

func mustToken() string {

	token := flag.String(
		"token-bot-token",
		"",
		"token to access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
