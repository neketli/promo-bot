package repository

import (
	"context"
)

type Repository interface {
	UserRepository
	PostRepository
}

type UserRepository interface {
	CreateUser(ctx context.Context, u *User) error
	RemoveUser(ctx context.Context, u *User) error
	IsUserExists(ctx context.Context, user *User) (bool, error)
}

type PostRepository interface {
	CreatePost(ctx context.Context, p *Post) error
	RemovePost(ctx context.Context, p *Post) error
	GetPosts(ctx context.Context) ([]Post, error)
	GetRandomPost(ctx context.Context) (Post, error)
}

type User struct {
	ID       int
	ChatID   int64
	Login    string
	Password string
}

type Post struct {
	ID          int
	Trigger     string
	Description string
}
