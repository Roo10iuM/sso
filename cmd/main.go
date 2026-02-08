package main

import (
	"fmt"
	"net"

	protosso "github.com/roo10ium/sso-protos/gen/go/sso"
	grpc_app "github.com/roo10ium/sso/internal/app/grpc"
	"github.com/roo10ium/sso/internal/usecases"
	"github.com/roo10ium/sso/internal/util/config"
	"github.com/roo10ium/sso/internal/util/log"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.MustLoad()
	log := log.SetupLogger(cfg.Env)

	log.Info(cfg.Env)

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 8011))
	if err != nil {
		log.Error(fmt.Sprintf("failed to listen: %v", err))
	}

	sso := usecases.NewSso(*log)
	app := grpc_app.NewApp(&sso)

	grpcServer := grpc.NewServer()

	protosso.RegisterSSOServer(grpcServer, app)

	grpcServer.Serve(lis)
}
