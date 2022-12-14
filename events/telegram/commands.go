package telegram

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
	"strings"
	lib "telegram-bot-go/lib/e"
	"telegram-bot-go/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s' ", text, username)

	// add merge collections if user add or delete username
	username, err := checkUsername(chatID, username)
	if err != nil {
		return lib.WrapIfErr("can't check username", err)
	}

	if IsAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		p.sendRandom(chatID, username)
	case HelpCmd:
		p.sendHepl(chatID)
	case StartCmd:
		p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCmd)
	}

	//add page: http://
	//rnd page: /rnd
	//help: /help
	//star: /start: hi + help

	return nil
}

func (p *Processor) namelessUser(chatID int, username string) (err error) {
	defer func() { err = lib.WrapIfErr("can't check username", err) }()
	return p.tg.SendMessage(chatID, msgNamelessUser)
}

func (p *Processor) savePage(chatID int, pageUrl string, username string) (err error) {
	defer func() { err = lib.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageUrl,
		UserName: username,
	}

	isExist, err := p.storage.IsExist(page)
	if err != nil {
		return err
	}

	if isExist {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = lib.WrapIfErr("cant sen random page", err) }()

	page, err := p.storage.PickRandom(username)
	fmt.Print(err)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPage) { // !!!
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPage) {
		return p.tg.SendMessage(chatID, msgNoSavedPage)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)

}

func (p *Processor) sendHepl(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func IsAddCmd(text string) bool {
	return IsUrl(text)
}

// only for links with https:// http://
func IsUrl(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}

func checkUsername(chatID int, username string) (string, error) {
	// Username(chatID)
	if username == "" {
		return Username(chatID)
	}
	return username, nil
}

func Username(chatID int) (string, error) {
	h := sha256.New()
	if _, err := io.WriteString(h, strconv.Itoa(chatID)); err != nil {
		return "", lib.Wrap("can't calculate hash", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
