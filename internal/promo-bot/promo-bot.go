package calcbot

import (
	"log"
	"promo-bot/config"
	"promo-bot/internal/controllers/telegram"
	"promo-bot/internal/infrastructure/repository/postgres"
)

const (
	storagePath = "./data/sqlite/storage.db"
)

func Start(config *config.Config) {
	repository, err := postgres.New(config.DB.Connection)
	if err != nil {
		log.Fatal("Service can't connect db: ", err)
	}

	log.Printf("Service has been started")

	bot, err := telegram.New(config.TG.Token, repository)
	if err != nil {
		log.Fatalf("FATAL: can't create bot, %s", err.Error())
	}
	if err := bot.Start(config.TG); err != nil {
		log.Fatalf("FATAL: can't start bot, %s", err.Error())
	}

}
