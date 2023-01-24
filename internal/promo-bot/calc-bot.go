package calcbot

import (
	"context"
	"log"
	"promo-bot/config"
	tgClient "promo-bot/internal/promo-bot/clients/telegram"
	event_consumer "promo-bot/internal/promo-bot/consumer/event-consumer"
	"promo-bot/internal/promo-bot/events/telegram"
	"promo-bot/internal/promo-bot/server"
	"promo-bot/internal/promo-bot/storage/sqlite"
)

const (
	batchSize   = 100
	storagePath = "./data/sqlite/storage.db"
)

func Start(config *config.Config) {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		log.Fatal("Service can't connect db: ", err)
	}

	if err := storage.Init(context.TODO()); err != nil {
		log.Fatal("Service can't init storage: ", err)
	}

	eventsHandler := telegram.New(tgClient.New(config.TG.TgHost, config.TG.TgToken), storage)

	log.Printf("Service has been started")

	go func() {
		if err := server.Start(config, storage); err != nil {
			log.Fatal("Service's server stopped: ", err)
		}
	}()

	consumer := event_consumer.New(eventsHandler, eventsHandler, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("Service stopped: ", err)
	}
}
