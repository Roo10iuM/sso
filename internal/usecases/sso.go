package usecases

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type sso struct {
	log *slog.Logger
}

func NewSso(log slog.Logger) sso {
	return sso{log: &log}
}

func (s *sso) Login(
	ctx context.Context,
	login string,
	password string,
	appUUID string,
) (token string, err error) {
	// TODO
	s.log.Info("start login")
	return login, nil
}

func (s *sso) RegisterNewUser(
	ctx context.Context,
	username string,
	email *string,
	password string,
) (userUUID string, err error) {
	// TODO
	s.log.Info("start register")
	return uuid.New().String(), nil
}
