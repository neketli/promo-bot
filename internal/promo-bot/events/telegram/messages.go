package telegram

const msgHelp = `Working with me is very easy - just send me a simple math expression and I will quickly calculate it for you

However, remember that I graduated from only 3 classes, so the list of operations I support is not large.

Available operations: a + b, a - b, a * b, a / b, where a, b are integer or real numbers
You can also view this help again using the command /help`

const msgHello = "Hello 👋, I'm a calculator bot and here's what I can do: \n" + msgHelp

const (
	msgUnknown = "I'm sorry, I didn't understand your command 😢"
)
