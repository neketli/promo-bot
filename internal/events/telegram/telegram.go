package telegram

import (
	"context"
	"errors"
	"fmt"
	"promo-bot/internal/promo-bot/events"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserName string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, fmt.Errorf("can't fetch events: %w", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))
	for _, value := range updates {
		res = append(res, event(value))
	}

	p.offset = updates[len(updates)-1].ID + 1
	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return fmt.Errorf("can't handle message: %w", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return fmt.Errorf("can't handle message: %w", err)
	}
	if err := p.handle(event.Text, meta.ChatID, meta.UserName); err != nil {
		return fmt.Errorf("can't handle message: %w", err)
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	if res, ok := event.Meta.(Meta); ok {
		return res, nil
	}

	return Meta{}, ErrUnknownMetaType
}

func event(upd telegram.Update) events.Event {
	updateType := fetchType(upd)
	res := events.Event{
		Type: updateType,
		Text: fetchText(upd),
	}

	if updateType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			UserName: upd.Message.User.UserName,
		}
	}

	return res
}

func fetchText(upd telegram.Update) string {
	if upd.Message != nil {
		return upd.Message.Text
	}
	return ""
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message != nil {
		return events.Message
	}
	return events.Unknown
}

func (p *Processor) addNewUserInfo(ctx context.Context, userMeta Meta) error {
	user, err := p.storage.Get(ctx, userMeta.ChatID)
	if err != nil {
		return fmt.Errorf("can't get user info :%w", err)
	}
	if user.UserName != "" {
		return nil
	}
	if err := p.storage.Save(ctx, (*storage.User)(&userMeta)); err != nil {
		return err
	}
	return nil
}
