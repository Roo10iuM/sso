package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/roo10ium/sso/internal/app/grpcapp"
	"github.com/roo10ium/sso/internal/services/ssoservice"
	"github.com/roo10ium/sso/internal/storage/sqlite"
	"github.com/roo10ium/sso/internal/util/config"
	"github.com/roo10ium/sso/internal/util/log"
)

func main() {
	cfg := config.MustLoad()
	log := log.SetupLogger(cfg.Env)

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		panic(fmt.Errorf("storage not initialized: %w", err))
	}

	sso := ssoservice.NewSSO(log, storage)
	app := grpcapp.NewApp(log, sso, cfg.GRPC.Host, cfg.GRPC.Port)

	go func() {
		app.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.Stop()
	log.Info("Gracefully stopped")
}
