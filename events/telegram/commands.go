package telegram

import (
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s' ", text, username)

	if IsAddCmd(text) {

	}

	switch text {
	case RndCmd:
	case HelpCmd:
	case StartCmd:

	}

	//add page: http://
	//rnd page: /rnd
	//help: /help
	//star: /start: hi + help

	return nil
}

func IsAddCmd(text string) bool {
	return IsUrl(text)
}

func IsUrl(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
