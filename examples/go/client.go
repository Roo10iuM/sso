package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-faker/faker/v4"

	"github.com/roo10ium/sso-protos/gen/go/pbsso"
	"github.com/roo10ium/sso/internal/util/config"
	"github.com/roo10ium/sso/internal/util/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.MustLoad()
	log := log.SetupLogger("local")

	serverAddr := fmt.Sprintf("%s:%d", cfg.GRPC.Host, cfg.GRPC.Port)

	var opts []grpc.DialOption
	// TODO secure
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		log.Error("NewClient error", slog.Any("error", err))
		return
	}
	defer conn.Close()

	client := pbsso.NewSSOClient(conn)

	ctx := context.Background()
	username := faker.FirstNameMale()
	mail := faker.Email()
	password := faker.Password()
	log.Info("start register", "user", []string{username, mail, password})
	res, err := client.Register(ctx, &pbsso.RegisterRequest{Username: username, Email: &mail, Password: password})
	if err != nil {
		log.Error("Request error", slog.Any("error", err))
	}
	log.Info("Register end", slog.Any("res", res))
}
