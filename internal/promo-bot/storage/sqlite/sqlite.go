package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"promo-bot/internal/promo-bot/storage"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to db: %w", err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Save(ctx context.Context, user *storage.User) error {
	query := `INSERT INTO users (chat_id, user_name, timestamp) VALUES (?,?,?)`

	if _, err := s.db.ExecContext(ctx, query, user.ChatID, user.UserName, time.Now().Unix()); err != nil {
		return fmt.Errorf("can't add user: %w", err)
	}
	return nil
}
func (s *Storage) Remove(ctx context.Context, user *storage.User) error {
	query := `DELETE FROM users WHERE chat_id = ? AND user_name = ?`
	if _, err := s.db.ExecContext(ctx, query, user.ChatID, user.UserName); err != nil {
		return fmt.Errorf("can't remove user: %w", err)
	}
	return nil
}
func (s *Storage) Get(ctx context.Context, chatID int) (storage.User, error) {
	query := `SELECT user_name FROM users WHERE chat_id = ?`
	var userName string
	err := s.db.QueryRowContext(ctx, query, chatID).Scan(&userName)
	if err == sql.ErrNoRows {
		return storage.User{}, nil
	}
	if err != nil {
		return storage.User{}, fmt.Errorf("can't get user: %w", err)
	}
	return storage.User{
		UserName: userName,
		ChatID:   chatID,
	}, nil
}

func (s *Storage) GetCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM users`
	var result int
	err := s.db.QueryRowContext(ctx, query).Scan(&result)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return -1, fmt.Errorf("can't get users count: %w", err)
	}

	return result, nil
}

func (s *Storage) Init(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS users (chat_id TEXT, user_name TEXT, timestamp TEXT)`
	if _, err := s.db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("can't create table: %w", err)
	}
	return nil
}
