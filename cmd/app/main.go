package main

import (
	"log"

	"promo-bot/config"
	promobot "promo-bot/internal/promo-bot"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.Fatal("Can't setup config: ", err)
	}
	promobot.Start(config)
}
