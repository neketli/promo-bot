package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

const (
	StartCommand = "/start"
	HelpCommand  = "/help"
)

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

var editInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(msgDelete, "/delete"),
	),
)
