package calcbot

import (
	"context"
	"fmt"
	"log"
	"promo-bot/config"
	"promo-bot/internal/repository"
	"promo-bot/internal/repository/sqlite"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	storagePath = "./data/sqlite/storage.db"
)

var keyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Список текущих триггеров"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Добавить триггер"),
	),
)

func Start(config *config.Config) {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		log.Fatal("Service can't connect db: ", err)
	}

	log.Printf("Service has been started")

	bot, err := tgbotapi.NewBotAPI(config.TG.Token)
	if err != nil {
		log.Fatal(err)
	}

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = config.TG.Timeout

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		switch update.Message.Text {
		case "start":
			msg.ReplyMarkup = keyboard
		case "Список текущих триггеров":
			res, err := storage.GetPosts(context.TODO())
			if err != nil {
				log.Print("ERROR: can't get posts: ", err)
			}
			for _, post := range res {
				text := fmt.Sprintf("[%s] - %s", post.Trigger, post.Description)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				if _, err := bot.Send(msg); err != nil {
					log.Fatal(err)
				}
			}
		case "Добавить триггер":
			err := storage.CreatePost(context.TODO(), &repository.Post{
				Trigger:     "Trigger",
				Description: "Description",
			})
			if err != nil {
				log.Print("ERROR: can't add post: ", err)
			}
		}

		if _, err := bot.Send(msg); err != nil {
			log.Fatal(err)
		}
	}
}
