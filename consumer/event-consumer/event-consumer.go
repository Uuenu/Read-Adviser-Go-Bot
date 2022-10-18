package eventconsumer

import (
	"log"
	"telegram-bot-go/events"
	"time"
)

type Consumer struct {
	fetcher    events.Fetcher
	proccessor events.Processor
	batchSize  int
}

func New(fetcher events.Fetcher, proccessor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:    fetcher,
		proccessor: proccessor,
		batchSize:  batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consume: %s", err.Error())
			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)

			continue
		}
	}
}

/*
	1. Потеря событий: ретраи, возвращение в хранилище, фоллбэк (локальный файл/внутри программы/db), подтверждение
	2. Обработка все пачки: остнавливаться после первой ошибки, после n ошибок и т д
	3. Параллельная обработка:
*/
func (c Consumer) handleEvents(events []events.Event) error {
	//sync.WaitGroup{}
	for _, evenet := range events {
		log.Printf("got new events: %s", evenet.Text)

		if err := c.proccessor.Process(evenet); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}
	return nil
}
