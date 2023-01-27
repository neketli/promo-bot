package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"promo-bot/internal/entity"

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

func (s *Repository) GetUsers(ctx context.Context) ([]entity.User, error) {
	query := `SELECT * FROM users`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return []entity.User{}, fmt.Errorf("can't get users: %w", err)
	}
	defer rows.Close()
	users := make([]entity.User, 0)

	for rows.Next() {
		var id int
		var user_name string
		var chat_id int64
		if err := rows.Scan(&id, &user_name, &chat_id); err != nil {
			return []entity.User{}, fmt.Errorf("can't get users: %w", err)
		}
		users = append(users, entity.User{
			ID:       id,
			UserName: user_name,
			ChatID:   chat_id,
		})
	}

	rerr := rows.Close()
	if rerr != nil {
		return []entity.User{}, fmt.Errorf("can't get users: %w", err)
	}

	if err := rows.Err(); err != nil {
		return []entity.User{}, fmt.Errorf("can't get users: %w", err)
	}

	return users, nil
}

// CreateUser create new admin user
func (s *Repository) CreateUser(ctx context.Context, user entity.User) error {
	query := `INSERT INTO users (user_name, chat_id) VALUES (?,?)`

	if _, err := s.db.ExecContext(ctx, query, user.UserName, user.ChatID); err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}
	return nil
}

func (s *Repository) RemoveUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = ?`

	if _, err := s.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("can't remove user: %w", err)
	}
	return nil
}

func (s *Repository) IsUserExists(ctx context.Context, userName string) (bool, error) {
	query := `SELECT * FROM users WHERE user_name = ? LIMIT 1`
	err := s.db.QueryRowContext(ctx, query, userName).Err()
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("can't get user: %w", err)
	}
	return true, nil
}

func (s *Repository) CreatePost(ctx context.Context, post entity.Post) error {
	query := `INSERT INTO promocodes (trigger, description) VALUES (?,?)`

	if _, err := s.db.ExecContext(ctx, query, post.Trigger, post.Description); err != nil {
		return fmt.Errorf("can't create user: %w", err)
	}
	return nil
}

func (s *Repository) RemovePost(ctx context.Context, id int) error {
	query := `DELETE FROM promocodes WHERE id = ?`

	if _, err := s.db.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("can't remove user: %w", err)
	}
	return nil
}

func (s *Repository) GetPosts(ctx context.Context) ([]entity.Post, error) {
	query := `SELECT * FROM promocodes`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return []entity.Post{}, fmt.Errorf("can't get posts: %w", err)
	}
	defer rows.Close()
	posts := make([]entity.Post, 0)

	for rows.Next() {
		var id int
		var trigger string
		var description string
		if err := rows.Scan(&id, &trigger, &description); err != nil {
			return []entity.Post{}, fmt.Errorf("can't get posts: %w", err)
		}
		posts = append(posts, entity.Post{
			ID:          id,
			Trigger:     trigger,
			Description: description,
		})
	}

	rerr := rows.Close()
	if rerr != nil {
		return []entity.Post{}, fmt.Errorf("can't get posts: %w", err)
	}

	if err := rows.Err(); err != nil {
		return []entity.Post{}, fmt.Errorf("can't get posts: %w", err)
	}

	return posts, nil
}

func (s *Repository) GetRandomPost(ctx context.Context) (entity.Post, error) {
	query := `SELECT * FROM promocoders ORDER BY RAND() LIMIT 1`
	var id int
	var trigger string
	var description string
	err := s.db.QueryRowContext(ctx, query).Scan(&id, &trigger, &description)
	if err == sql.ErrNoRows {
		return entity.Post{}, nil
	}
	if err != nil {
		return entity.Post{}, fmt.Errorf("can't get user: %w", err)
	}
	return entity.Post{
		ID:          id,
		Trigger:     trigger,
		Description: description,
	}, nil
}
