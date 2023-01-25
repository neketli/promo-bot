package telegram

const msgHelp = `Working with me is very easy - just send me a simple math expression and I will quickly calculate it for you

However, remember that I graduated from only 3 classes, so the list of operations I support is not large.

Available operations: a + b, a - b, a * b, a / b, where a, b are integer or real numbers
You can also view this help again using the command /help`

const msgHello = "Привет 👋, я бот с промокодами вот что я умею: \n" + msgHelp

const (
	msgUnknown              = "Я вас не понимаю 😢"
	msgError                = "Произошла какая-то внутренняя ошибка 😢"
	msgEnterTrigger         = "Введите ключевое слово"
	msgEnterDescription     = "Введите сообщение, которым бот будет отвечать"
	msgTriggerList          = "Список текущих триггеров"
	msgTriggerCreate        = "Добавить триггер"
	msgTriggerCreateSuccess = "Успешно добавлено!"
	msgCancel               = "Отменить"
)
