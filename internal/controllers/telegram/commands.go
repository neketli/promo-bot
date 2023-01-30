package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	StartCommand       = "/start"
	HelpCommand        = "/help"
	DeletePostCommand  = "/delete-post"
	DeleteAdminCommand = "/delete-admin"
)

var defaultKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(HelpCommand),
	),
)

var adminKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(msgTriggerList),
		tgbotapi.NewKeyboardButton(msgTriggerCreate),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(msgAdminList),
		tgbotapi.NewKeyboardButton(msgAdminCreate),
	),
)

var cancelKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(msgCancel),
	),
)

var deletePostInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(msgDelete, DeletePostCommand),
	),
)

var deleteAdminInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(msgDelete, DeleteAdminCommand),
	),
)
