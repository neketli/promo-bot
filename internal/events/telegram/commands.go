package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
)

const (
	StartCommand = "/start"
	HelpCommand  = "/help"
)

func (p *Processor) handle(text string, chatID int, userName string) error {
	text = strings.TrimSpace(text)

	log.Printf("DEBUG: new command '%s' from '%s'", text, userName)

	switch text {
	case HelpCommand:
		return p.tg.SendMessage(chatID, msgHelp)
	case StartCommand:
		p.tg.SendMessage(chatID, msgHello)
		if err := p.addNewUserInfo(context.TODO(), Meta{ChatID: chatID, UserName: userName}); err != nil {
			return fmt.Errorf("can't add new user: %w", err)
		}
		return nil
	default:
		return p.tg.SendMessage(chatID, msgUnknown)
	}
}
