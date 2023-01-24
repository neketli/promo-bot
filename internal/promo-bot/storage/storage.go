package storage

import "context"

type Storage interface {
	Save(ctx context.Context, u *User) error
	Remove(ctx context.Context, u *User) error
	Get(ctx context.Context, chatID int) (User, error)
	GetCount(ctx context.Context) (int, error)
}

type User struct {
	ChatID   int
	UserName string
}
