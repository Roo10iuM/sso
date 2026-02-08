package ssoservice

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

const (
	prefix = "ssoservice"
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
	const op = prefix + ".Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("username", login),
	)

	log.Info("attempting to login user")
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
	// TODO
	return uuid.New().String(), nil
}
