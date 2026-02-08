package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/roo10ium/sso/internal/domain/models"
	"github.com/roo10ium/sso/internal/storage"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
)

const prefix = "sqlite"

type SQLiteStorage struct {
	db *sql.DB
}

func New(storagePath string) (*SQLiteStorage, error) {
	const op = prefix + ".New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) SaveUser(ctx context.Context, username string, email *string, passHash []byte) (uuid.UUID, error) {
	switch email {
	case nil:
		return s.SaveUserWithoutEmail(ctx, username, passHash)
	default:
		return s.SaveUserWithEmail(ctx, username, *email, passHash)
	}
}

func (s *SQLiteStorage) SaveUserWithoutEmail(ctx context.Context, username string, passHash []byte) (uuid.UUID, error) {
	const op = prefix + ".SaveUserWithoutEmail"

	stmt, err := s.db.Prepare("INSERT INTO users(id, username, pass_hash) VALUES(?, ?, ?)")
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	id := uuid.New()
	_, err = stmt.ExecContext(ctx, id, username, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return uuid.Nil, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *SQLiteStorage) SaveUserWithEmail(ctx context.Context, username string, email string, passHash []byte) (uuid.UUID, error) {
	const op = prefix + ".SaveUserWithEmail"

	stmt, err := s.db.Prepare("INSERT INTO users(id, username, email, pass_hash) VALUES(?, ?, ?, ?)")
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	id := uuid.New()
	_, err = stmt.ExecContext(ctx, id, username, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return uuid.Nil, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *SQLiteStorage) GetUser(ctx context.Context, username string) (models.User, error) {
	const op = prefix + ".GetUser"

	stmt, err := s.db.Prepare("SELECT id, username, email, pass_hash FROM users WHERE username = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, username)

	var user models.User
	err = row.Scan(&user.ID, &user.Username, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *SQLiteStorage) GetApp(ctx context.Context, id int) (models.App, error) {
	const op = prefix + ".GetApp"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
