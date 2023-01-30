package calcbot

import (
	"log"
	"promo-bot/config"
	"promo-bot/internal/controllers/telegram"
	"promo-bot/internal/infrastructure/repository/sqlite"
)

const (
	storagePath = "./data/sqlite/storage.db"
)

func Start(config *config.Config) {
	repository, err := sqlite.New(storagePath)
	if err != nil {
		log.Fatal("Service can't connect db: ", err)
	}

	log.Printf("Service has been started")

	bot, err := telegram.New(config.TG.Token, repository)
	if err != nil {
		log.Fatalf("FATAL: can't create bot, %s", err.Error())
	}
	if err := bot.Start(config.TG.Timeout); err != nil {
		log.Fatalf("FATAL: can't start bot, %s", err.Error())
	}

	bot.Bot.Debug = config.TG.Mode == "debug"

	log.Printf("Authorized on account %s", bot.Bot.Self.UserName)

}
