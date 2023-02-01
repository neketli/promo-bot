package main

import (
	"log"

	"promo-bot/config"
	promobot "promo-bot/internal/promo-bot"
)

const configPath = "./config/config.yml"

func main() {
	config, err := config.New(configPath)
	if err != nil {
		log.Fatal("Can't setup config: ", err)
	}
	promobot.Start(config)
}
