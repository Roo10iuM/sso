package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	"github.com/roo10ium/sso-protos/gen/go/pbsso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	prefix = "app"
)

type app struct {
	log  *slog.Logger
	s    *grpc.Server
	host string
	port int
}

func NewApp(log *slog.Logger, sso ssoService, host string, port int) *app {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p any) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	validatorOpts := []validator.Option{
		validator.WithOnValidationErrCallback(func(ctx context.Context, err error) {
			log.Error("Validation error", slog.Any("error", err))
		}),
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
		validator.UnaryServerInterceptor(validatorOpts...),
	))

	pbsso.RegisterSSOServer(s, &gRPCServer{sso: sso})

	return &app{log: log, s: s, host: host, port: port}
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (a app) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a app) Run() error {
	const op = prefix + ".Run"
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.host, a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	a.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := a.s.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
