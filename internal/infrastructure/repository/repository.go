package repository

import (
	"context"
	"promo-bot/internal/entity"
)

type Repository interface {
	UserRepository
	PostRepository
}

type UserRepository interface {
	CreateUser(ctx context.Context, u *entity.User) error
	RemoveUser(ctx context.Context, u *entity.User) error
	IsUserExists(ctx context.Context, user *entity.User) (bool, error)
}

type PostRepository interface {
	CreatePost(ctx context.Context, p *entity.Post) error
	RemovePost(ctx context.Context, id int) error
	GetPosts(ctx context.Context) ([]entity.Post, error)
	GetRandomPost(ctx context.Context) (entity.Post, error)
}
