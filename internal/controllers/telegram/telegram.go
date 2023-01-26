package telegram

import (
	"context"
	"fmt"
	"log"
	"promo-bot/internal/entity"
	"promo-bot/internal/infrastructure/repository/sqlite"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TgBot struct {
	Bot        *tgbotapi.BotAPI
	Repository *sqlite.Repository
}

func New(token string, repository *sqlite.Repository) (*TgBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return &TgBot{}, fmt.Errorf("error bot api connection: %w", err)
	}
	bot.Debug = true

	return &TgBot{Bot: bot, Repository: repository}, nil
}

func (b *TgBot) Start(timout int) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = timout

	updates := b.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID

			switch update.Message.Text {
			case "/start":
				msg.ReplyMarkup = defaultKeyboard
			case "Список текущих триггеров":
				b.getAllPosts(update)
			case "Добавить триггер":
				b.enterPost(update, updates)
			}
		} else if update.CallbackQuery != nil {
			b.deletePost(update)

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

func (b *TgBot) enterPost(update tgbotapi.Update, updates tgbotapi.UpdatesChannel) error {

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

		if update.Message.Text == msgCancel {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgCanceled)
			msg.ReplyMarkup = defaultKeyboard
			if _, err := b.Bot.Send(msg); err != nil {
				return fmt.Errorf("can't send message: %w", err)
			}
			return nil
		}

		switch counter {
		case 0:
			trigger = update.Message.Text
			counter++
			b.SendMessage(update.Message.Chat.ID, msgEnterDescription)
		case 1:
			description = update.Message.Text
			err := b.Repository.CreatePost(context.TODO(), &entity.Post{
				Trigger:     trigger,
				Description: description,
			})
			if err != nil {
				return fmt.Errorf("can't add post: %w", err)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgTriggerCreateSuccess)
			msg.ReplyMarkup = defaultKeyboard
			if _, err := b.Bot.Send(msg); err != nil {
				return fmt.Errorf("can't send message: %w", err)
			}
			return nil
		}
	}
	return nil
}

func (b *TgBot) getAllPosts(update tgbotapi.Update) {
	res, err := b.Repository.GetPosts(context.TODO())
	if err != nil {
		log.Print("ERROR: can't get posts: ", err)
	}
	for _, post := range res {
		text := fmt.Sprintf("ID:%d [%s] - %s", post.ID, post.Trigger, post.Description)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyMarkup = editInlineKeyboard
		if _, err := b.Bot.Send(msg); err != nil {
			log.Fatal(err)
		}
	}
}

func (b *TgBot) deletePost(update tgbotapi.Update) error {
	log.Printf("!!%s", update.CallbackQuery.Message.Text)
	res := strings.Split(strings.Split(update.CallbackQuery.Message.Text, " ")[0], ":")[1]
	log.Printf("!!%s", res)

	postId, err := strconv.Atoi(res)
	if err != nil {
		return fmt.Errorf("can't get delete message: %w", err)
	}
	b.Repository.RemovePost(context.TODO(), postId)

	deleteMessageConfig := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID)
	deleteMessageAlertCallbacl := tgbotapi.NewCallback(update.CallbackQuery.ID, msgDeleteSuccess)
	if _, err := b.Bot.Request(deleteMessageAlertCallbacl); err != nil {
		panic(err)
	}
	if _, err := b.Bot.Request(deleteMessageConfig); err != nil {
		panic(err)
	}
	return nil
}
