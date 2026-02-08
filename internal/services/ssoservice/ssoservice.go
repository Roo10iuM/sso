package ssoservice

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/roo10ium/sso/internal/domain/models"
	"github.com/roo10ium/sso/internal/storage"
)

const (
	prefix = "ssoservice"
)

type sso struct {
	log     *slog.Logger
	storage Storage
}

type Storage interface {
	UserSaver
	UserProvider
	AppProvider
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		username string,
		email *string,
		passHash []byte,
	) (uid uuid.UUID, err error)
}

type UserProvider interface {
	GetUser(ctx context.Context, email string) (models.User, error)
}

type AppProvider interface {
	GetApp(ctx context.Context, appID int) (models.App, error)
}

func NewSSO(log *slog.Logger, storage Storage) *sso {
	return &sso{log: log, storage: storage}
}

func (s *sso) Login(
	ctx context.Context,
	login string,
	password string,
	appUUID string,
) (token string, err error) {
	const op = prefix + ".Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("username", login),
	)

	log.Info("attempting to login user")

	user, err := s.storage.GetUser(ctx, login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			s.log.Warn("user not found", slog.Any("error", err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		s.log.Error("failed to get user", slog.Any("error", err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.log.Info("invalid credentials", slog.Any("error", err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// TODO
	return login, nil
}

func (s *sso) RegisterNewUser(
	ctx context.Context,
	username string,
	email *string,
	password string,
) (userUUID string, err error) {
	const op = prefix + ".RegisterNewUser"

	log := s.log.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	if email != nil {
		log = log.With("email", *email)
	}

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.Any("error", err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	id, err := s.storage.SaveUser(ctx, username, email, passHash)
	if err != nil {
		log.Error("failed to save user", slog.Any("error", err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id.String(), nil
}
