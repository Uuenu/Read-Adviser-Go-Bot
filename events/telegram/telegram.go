package telegram

import "telegram-bot-go/clients/telegram"

type Processor struct {
	tg     *telegram.Client
	offset string
	// storage
}

func New(client *telegram.Client) {

}
