package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"promo-bot/internal/repository"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

// New initial repository function that returns repository object
func New(path string) (*Repository, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("can't open db: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("can't connect to db: %w", err)
	}
	return &Repository{db: db}, nil
}

// CreateUser create new admin user
func (s *Repository) CreateUser(ctx context.Context, user *repository.User) error {
	query := `INSERT INTO users (login, password, chat_id) VALUES (?,?,?)`

	if _, err := s.db.ExecContext(ctx, query, user.Login, user.Password, user.ChatID); err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}
	return nil
}

func (s *Repository) RemoveUser(ctx context.Context, user *repository.User) error {
	query := `DELETE FROM users WHERE login = ? AND password = ? AND chat_id = ?`

	if _, err := s.db.ExecContext(ctx, query, user.Login, user.Password, user.ChatID); err != nil {
		return fmt.Errorf("can't remove user: %w", err)
	}
	return nil
}

func (s *Repository) IsUserExists(ctx context.Context, user *repository.User) (bool, error) {
	query := `SELECT (login, password, chat_id) FROM users WHERE chat_id = ?`
	err := s.db.QueryRowContext(ctx, query, user.ChatID).Scan()
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("can't get user: %w", err)
	}
	return true, nil
}

func (s *Repository) CreatePost(ctx context.Context, post *repository.Post) error {
	query := `INSERT INTO promocodes (trigger, description) VALUES (?,?)`

	if _, err := s.db.ExecContext(ctx, query, post.Trigger, post.Description); err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}
	return nil
}

func (s *Repository) RemovePost(ctx context.Context, post *repository.Post) error {
	query := `DELETE FROM promocodes WHERE id = ?`

	if _, err := s.db.ExecContext(ctx, query, post.ID); err != nil {
		return fmt.Errorf("can't remove user: %w", err)
	}
	return nil
}

func (s *Repository) GetPosts(ctx context.Context) ([]repository.Post, error) {
	query := `SELECT * FROM promocodes`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return []repository.Post{}, fmt.Errorf("can't get posts: %w", err)
	}
	defer rows.Close()
	posts := make([]repository.Post, 0)

	for rows.Next() {
		var id int
		var trigger string
		var description string
		if err := rows.Scan(&id, &trigger, &description); err != nil {
			return []repository.Post{}, fmt.Errorf("can't get posts: %w", err)
		}
		posts = append(posts, repository.Post{
			ID:          id,
			Trigger:     trigger,
			Description: description,
		})
	}

	rerr := rows.Close()
	if rerr != nil {
		return []repository.Post{}, fmt.Errorf("can't get posts: %w", err)
	}

	if err := rows.Err(); err != nil {
		return []repository.Post{}, fmt.Errorf("can't get posts: %w", err)
	}

	return posts, nil
}

func (s *Repository) GetRandomPost(ctx context.Context) (repository.Post, error) {
	query := `SELECT * FROM promocoders ORDER BY RAND() LIMIT 1`
	var id int
	var trigger string
	var description string
	err := s.db.QueryRowContext(ctx, query).Scan(&id, &trigger, &description)
	if err == sql.ErrNoRows {
		return repository.Post{}, nil
	}
	if err != nil {
		return repository.Post{}, fmt.Errorf("can't get user: %w", err)
	}
	return repository.Post{
		ID:          id,
		Trigger:     trigger,
		Description: description,
	}, nil
}
