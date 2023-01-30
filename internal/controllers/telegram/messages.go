package telegram

const msgHelp = `Введи ключевое слово и я выдам тебе промокод`

const msgHello = "Привет 👋, я бот с промокодами вот что я умею: \n" + msgHelp

const (
	msgUnknown              = "Я вас не понимаю 😢"
	msgError                = "Произошла какая-то внутренняя ошибка 😢"
	msgEnterTrigger         = "Введите ключевое слово"
	msgEnterDescription     = "Введите сообщение, которым бот будет отвечать"
	msgEnterAdminName       = "Введите ник пользователя телеграм, которого вы хотели бы сделать администратором (@ник без собаки и пробелов)"
	msgTriggerList          = "Вывести список текущих триггеров"
	msgTriggerCreate        = "Добавить триггер"
	msgTriggerCreateSuccess = "Успешно добавлено!"
	msgAdminList            = "Вывести список администраторов"
	msgAdminCreate          = "Добавить ник администратора"
	msgAdminCreateSuccess   = "Успешно добавлен!"
	msgCancel               = "Отменить"
	msgCanceled             = "Отменено"
	msgDelete               = "Удалить ❌"
	msgDeleteSuccess        = "Удаление прошло успешно"
	msgDeleteReply          = "Вы были исключены из списка администраторов"

	msgFindPromo   = "Нашёл кое-что для вас 🔍:\n"
	msgRandomPromo = "Получить случайный промокод"

	msgByAdmin = "администратором: %s"
	msgByUser  = "пользователем: %s"
)
