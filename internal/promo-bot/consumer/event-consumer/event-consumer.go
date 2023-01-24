package event_consumer

import (
	"log"
	"promo-bot/internal/promo-bot/events"
	"sync"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("ERROR: %s", err.Error())
			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(time.Second)
			continue
		}
		if err := c.handleEvents(gotEvents); err != nil {
			log.Printf("ERROR: %s", err.Error())
			continue
		}
	}
}

func (c *Consumer) handleEvents(arrEvents []events.Event) error {
	var wg sync.WaitGroup
	for _, event := range arrEvents {
		wg.Add(1)
		go func(event events.Event) {
			defer wg.Done()
			if err := c.processor.Process(event); err != nil {
				log.Printf("ERROR: can't handle event: %s", err.Error())

			}
		}(event)

	}
	wg.Wait()
	return nil
}
