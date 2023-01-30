package telegram

import (
	"context"
	"fmt"
	"log"
	"promo-bot/internal/entity"
	"promo-bot/internal/infrastructure/repository/sqlite"
	"strconv"
	"strings"
	"sync"

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
				fmt.Print("DEBUG: message from admin ", update.Message.Chat.UserName)
				b.adminControls(update, updates)
			} else {
				fmt.Print("DEBUG: message from user ", update.Message.Chat.UserName)
				b.userControls(update, updates)
			}
		} else if update.CallbackQuery != nil {
			isAdmin, err := b.Repository.IsUserExists(context.TODO(), update.CallbackQuery.From.UserName)
			if err != nil {
				log.Print("ERROR: can't check admin: ", err)
			}
			if isAdmin {
				switch update.CallbackQuery.Data {
				case DeletePostCommand:
					b.deletePost(update)
				case DeleteAdminCommand:
					b.deleteAdmin(update)
				}
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
	switch update.Message.Text {
	case StartCommand:
		b.initAdmin(update)
		b.SendMessageWithKeyboard(update.Message.Chat.ID, adminKeyboard, msgHello)
	case HelpCommand:
		b.SendMessageWithKeyboard(update.Message.Chat.ID, adminKeyboard, msgHelp)
	case msgTriggerList:
		b.getAllPosts(update)
	case msgTriggerCreate:
		if err := b.enterPost(update, updates); err != nil {
			b.SendError(update.Message.From.ID)
			log.Printf("ERROR: can't create trigger, %s \n", err.Error())
		}
	case msgAdminList:
		b.getAdminList(update)
	case msgAdminCreate:
		if err := b.createAdmin(update, updates); err != nil {
			b.SendError(update.Message.From.ID)
			log.Printf("ERROR: can't create admin, %s \n", err.Error())
		}
	}
}

func (b *TgBot) userControls(update tgbotapi.Update, updates tgbotapi.UpdatesChannel) error {
	switch update.Message.Text {
	case StartCommand:
		b.SendMessageWithKeyboard(update.Message.Chat.ID, defaultKeyboard, msgHello)
	case HelpCommand:
		b.SendMessageWithKeyboard(update.Message.Chat.ID, defaultKeyboard, msgHelp)
	case msgRandomPromo, RandomCommand:
		b.sendRandomCode(update)
	default:
		b.processMessageTriggers(update)
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

func (b *TgBot) initAdmin(update tgbotapi.Update) error {
	user := entity.User{
		UserName: update.Message.Chat.UserName,
		ChatID:   update.Message.Chat.ID,
	}
	if err := b.Repository.UpdateUser(context.TODO(), user); err != nil {
		return fmt.Errorf("can't update user info: %w", err)
	}
	return nil

}

func (b *TgBot) getAdminList(update tgbotapi.Update) {
	res, err := b.Repository.GetUsers(context.TODO())
	if err != nil {
		log.Print("ERROR: can't get posts: ", err)
	}
	for _, user := range res {
		text := fmt.Sprintf("ID:%d | Chat:%d | %s", user.ID, user.ChatID, user.UserName)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
		msg.ReplyMarkup = deleteAdminInlineKeyboard
		if _, err := b.Bot.Send(msg); err != nil {
			log.Fatal("ERROR: failed sending message: ", err)
		}
	}
}

func (b *TgBot) deleteAdmin(update tgbotapi.Update) error {
	uid := strings.Split(strings.Split(update.CallbackQuery.Message.Text, " | ")[0], ":")[1]
	chatId := strings.Split(strings.Split(update.CallbackQuery.Message.Text, " | ")[1], ":")[1]

	userId, err := strconv.Atoi(uid)
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

	if chatId, err := strconv.ParseInt(chatId, 10, 64); err == nil {
		b.sendMessageToDelited(chatId, update.CallbackQuery.From.UserName)
	} else {
		return fmt.Errorf("can't get delete message: %w", err)
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

func (b *TgBot) sendMessageToDelited(chatId int64, fromUser string) error {
	text := msgDeleteReply + " " + fmt.Sprintf(msgByAdmin, fromUser)
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyMarkup = defaultKeyboard
	if _, err := b.Bot.Send(msg); err != nil {
		return fmt.Errorf("failed sending message: %w", err)
	}
	return nil
}

func (b *TgBot) getRandomPromo() (string, error) {
	post, err := b.Repository.GetRandomPost(context.TODO())
	if err != nil {
		return msgError, fmt.Errorf("can't get random post: %w", err)
	}
	text := msgFindPromo + post.Description
	return text, nil

}

func (b *TgBot) sendRandomCode(update tgbotapi.Update) error {
	text, err := b.getRandomPromo()
	if err != nil {
		b.SendError(update.Message.Chat.ID)
		return err
	}
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := b.Bot.Send(msg); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}

func (b *TgBot) processMessageTriggers(update tgbotapi.Update) error {
	var mutex sync.Mutex
	var wg sync.WaitGroup
	triggers, err := b.Repository.GetTriggerList(context.TODO())
	if err != nil {
		b.SendError(update.Message.Chat.ID)
		return err
	}
	msgText := strings.ReplaceAll(strings.ToLower(update.Message.Text), " ", "")
	descriptions := make([]string, 0)
	for _, trigger := range triggers {
		wg.Add(1)
		go func(trigger string) {
			t := strings.ReplaceAll(strings.ToLower(trigger), " ", "")
			if strings.Contains(msgText, t) {
				posts, err := b.Repository.GetPostsByTrigger(context.TODO(), trigger)
				if err != nil {
					return
				}
				for _, post := range posts {
					mutex.Lock()
					descriptions = append(descriptions, post.Description)
					mutex.Unlock()
				}
			}
			wg.Done()
		}(trigger)
	}

	wg.Wait()

	if len(descriptions) == 0 {
		return nil
	}

	text := msgFindPromo

	for _, description := range descriptions {
		text = text + description + "\n"
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyToMessageID = update.Message.MessageID
	if _, err := b.Bot.Send(msg); err != nil {
		return fmt.Errorf("can't send message: %w", err)
	}
	return nil
}
