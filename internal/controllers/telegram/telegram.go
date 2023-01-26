package telegram

import (
	"context"
	"fmt"
	"log"
	"promo-bot/internal/entity"
	"promo-bot/internal/infrastructure/repository"
	"promo-bot/internal/infrastructure/repository/sqlite"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	Bot        *tgbotapi.BotAPI
	Repository repository.Repository
}

var defaultKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(msgTriggerList),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(msgTriggerCreate),
	),
)

var cancelKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(msgCancel),
	),
)

func New(token string, repository *sqlite.Repository) (*TgBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return &TgBot{}, fmt.Errorf("error bot api connection: %w", err)
	}
	bot.Debug = true

	return &TgBot{Bot: bot}, nil
}

func (b *TgBot) Start(timout int) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = timout

	updates := b.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		switch update.Message.Text {
		case "/start":
			msg.ReplyMarkup = defaultKeyboard
		case "Список текущих триггеров":
			res, err := b.Repository.GetPosts(context.TODO())
			if err != nil {
				log.Print("ERROR: can't get posts: ", err)
			}
			for _, post := range res {
				text := fmt.Sprintf("[%s] - %s", post.Trigger, post.Description)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				if _, err := b.Bot.Send(msg); err != nil {
					log.Fatal(err)
				}
			}
			// case "Добавить триггер":
			// 	b.enterPost(update)
		}
	}
	return nil
}

func (b *TgBot) SendMessage(id int64, text string) error {
	msg := tgbotapi.NewMessage(id, text)
	if _, err := b.Bot.Send(msg); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}

func (b *TgBot) SendError(id int64) error {
	msg := tgbotapi.NewMessage(id, msgError)
	if _, err := b.Bot.Send(msg); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}

func (b *TgBot) enterPost(update tgbotapi.Update) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.Bot.GetUpdatesChan(u)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgEnterTrigger)
	msg.ReplyMarkup = cancelKeyboard
	if _, err := b.Bot.Send(msg); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	counter := 0

	var trigger, description string

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == "Отмена" {
			b.Bot.StopReceivingUpdates()
			return nil
		}

		switch counter {
		case 0:
			trigger = update.Message.Text
			counter++
		case 1:
			description = update.Message.Text
			counter++
		case 2:
			err := b.Repository.CreatePost(context.TODO(), &entity.Post{
				Trigger:     trigger,
				Description: description,
			})
			if err != nil {
				return fmt.Errorf("can't add post: %w", err)
			}
			b.SendMessage(update.Message.Chat.ID, msgTriggerCreateSuccess)
			b.Bot.StopReceivingUpdates()
		}

	}
	return nil

}
