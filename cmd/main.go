package main

import (
	"github.com/roo10ium/sso/internal/app/grpcapp"
	"github.com/roo10ium/sso/internal/services/ssoservice"
	"github.com/roo10ium/sso/internal/util/config"
	"github.com/roo10ium/sso/internal/util/log"
)

func main() {
	cfg := config.MustLoad()
	log := log.SetupLogger(cfg.Env)

	sso := ssoservice.NewSSO(log)
	app := grpcapp.NewApp(log, sso, cfg.GRPC.Host, cfg.GRPC.Port)

	app.MustRun()
}
