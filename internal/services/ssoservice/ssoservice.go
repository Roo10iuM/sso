package ssoservice

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type sso struct {
	log *slog.Logger
}

func NewSSO(log *slog.Logger) *sso {
	return &sso{log: log}
}

func (s *sso) Login(
	ctx context.Context,
	login string,
	password string,
	appUUID string,
) (token string, err error) {
	// TODO
	s.log.Info("attempting to login user")
	return login, nil
}

func (s *sso) RegisterNewUser(
	ctx context.Context,
	username string,
	email *string,
	password string,
) (userUUID string, err error) {
	// TODO
	s.log.Info("registering user")
	return uuid.New().String(), nil
}
