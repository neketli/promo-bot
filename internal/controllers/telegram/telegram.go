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
			isAdmin, err := b.Repository.IsUserExists(context.TODO(), update.Message.Chat.UserName)
			if err != nil {
				log.Print("ERROR: can't check admin: ", err)
			}
			if isAdmin {
				fmt.Print("msg from ", update.Message.Chat.UserName)
				b.adminControls(update, updates)
			} else {
				// pass
			}
		} else if update.CallbackQuery != nil {
			switch update.CallbackQuery.Data {
			case DeletePostCommand:
				b.deletePost(update)
			case DeleteAdminCommand:
				b.deleteAdmin(update)
			}
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

func (b *TgBot) SendMessageReply(id int64, messageId int, text string) error {
	msg := tgbotapi.NewMessage(id, text)
	msg.ReplyToMessageID = messageId

	if _, err := b.Bot.Send(msg); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}

func (b *TgBot) SendMessageWithKeyboard(id int64, keyboard tgbotapi.ReplyKeyboardMarkup, text string) error {
	msg := tgbotapi.NewMessage(id, text)
	msg.ReplyMarkup = keyboard

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

func (b *TgBot) adminControls(update tgbotapi.Update, updates tgbotapi.UpdatesChannel) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "admin")
	msg.ReplyToMessageID = update.Message.MessageID

	switch update.Message.Text {
	case StartCommand:
		b.SendMessageWithKeyboard(update.Message.Chat.ID, adminKeyboard, msgHello)
	case HelpCommand:
		b.SendMessageWithKeyboard(update.Message.Chat.ID, adminKeyboard, msgHelp)
	case msgTriggerList:
		b.getAllPosts(update)
	case msgTriggerCreate:
		b.enterPost(update, updates)
	case msgAdminList:
		b.getAdminList(update)
	case msgAdminCreate:
		b.createAdmin(update, updates)
	}
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
			msg.ReplyMarkup = adminKeyboard
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
			err := b.Repository.CreatePost(context.TODO(), entity.Post{
				Trigger:     trigger,
				Description: description,
			})
			if err != nil {
				return fmt.Errorf("can't add post: %w", err)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msgTriggerCreateSuccess)
			msg.ReplyMarkup = adminKeyboard
			if _, err := b.Bot.Send(msg); err != nil {
				return fmt.Errorf("can't send message: %w", err)
			}
			return nil
		}
	}
	return nil
}

func (b *TgBot) createAdmin(update tgbotapi.Update, updates tgbotapi.UpdatesChannel) error {
	if err := b.SendMessageWithKeyboard(update.Message.Chat.ID, cancelKeyboard, msgEnterAdminName); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.Text == msgCancel {
			if err := b.SendMessageWithKeyboard(update.Message.Chat.ID, adminKeyboard, msgCanceled); err != nil {
				return fmt.Errorf("can't send message: %w", err)
			}
			return nil
		}

		err := b.Repository.CreateUser(context.TODO(), entity.User{
			UserName: strings.TrimLeft(update.Message.Text, "@"),
		})
		if err != nil {
			return fmt.Errorf("can't add post: %w", err)
		}

		if err := b.SendMessageWithKeyboard(update.Message.Chat.ID, adminKeyboard, msgAdminCreateSuccess); err != nil {
			return fmt.Errorf("can't send message: %w", err)
		}
		return nil
	}
	return nil
}

func (b *TgBot) getAdminList(update tgbotapi.Update) {
	res, err := b.Repository.GetUsers(context.TODO())
	if err != nil {
		log.Print("ERROR: can't get posts: ", err)
	}
	for _, user := range res {
		text := fmt.Sprintf("ID:%d %s|%d", user.ID, user.UserName, user.ChatID)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyMarkup = deleteAdminInlineKeyboard
		if _, err := b.Bot.Send(msg); err != nil {
			log.Fatal("ERROR: failed sending message: ", err)
		}
	}
}

func (b *TgBot) deleteAdmin(update tgbotapi.Update) error {
	res := strings.Split(strings.Split(update.CallbackQuery.Message.Text, " ")[0], ":")[1]

	userId, err := strconv.Atoi(res)
	if err != nil {
		return fmt.Errorf("can't get delete message: %w", err)
	}
	b.Repository.RemoveUser(context.TODO(), userId)

	deleteMessageConfig := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID)
	deleteMessageAlertCallbacl := tgbotapi.NewCallback(update.CallbackQuery.ID, msgDeleteSuccess)
	if _, err := b.Bot.Request(deleteMessageAlertCallbacl); err != nil {
		log.Fatal("ERROR: failed deliting message: ", err)
	}
	if _, err := b.Bot.Request(deleteMessageConfig); err != nil {
		log.Fatal("ERROR: failed deliting message: ", err)
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
		msg.ReplyMarkup = deletePostInlineKeyboard
		if _, err := b.Bot.Send(msg); err != nil {
			log.Fatal("ERROR: failed sending message: ", err)

		}
	}
}

func (b *TgBot) deletePost(update tgbotapi.Update) error {
	res := strings.Split(strings.Split(update.CallbackQuery.Message.Text, " ")[0], ":")[1]

	postId, err := strconv.Atoi(res)
	if err != nil {
		return fmt.Errorf("can't get delete message: %w", err)
	}
	b.Repository.RemovePost(context.TODO(), postId)

	deleteMessageConfig := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID,
		update.CallbackQuery.Message.MessageID)
	deleteMessageAlertCallbacl := tgbotapi.NewCallback(update.CallbackQuery.ID, msgDeleteSuccess)
	if _, err := b.Bot.Request(deleteMessageAlertCallbacl); err != nil {
		log.Fatal("ERROR: failed deliting message: ", err)
	}
	if _, err := b.Bot.Request(deleteMessageConfig); err != nil {
		log.Fatal("ERROR: failed deliting message: ", err)
	}
	return nil
}
