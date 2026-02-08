package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	protosso "github.com/roo10ium/sso-protos/gen/go/sso"
)

type app struct {
	protosso.UnimplementedSSOServer
	sso sso
}

func NewApp(s sso) *app {
	return &app{sso: s}
}

type sso interface {
	Login(
		ctx context.Context,
		login string,
		password string,
		appUUID string,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		username string,
		email *string,
		password string,
	) (userUUID string, err error)
}

func Register(gRPCServer *grpc.Server, sso sso) {
	protosso.RegisterSSOServer(gRPCServer, &app{sso: sso})
}

func (s *app) Login(
	ctx context.Context,
	in *protosso.LoginRequest,
) (*protosso.LoginResponse, error) {
	// TODO
	var login string
	switch in.Login.(type) {
	case *protosso.LoginRequest_Email:
		login = in.GetEmail()
	case *protosso.LoginRequest_Username:
		login = in.GetUsername()
	default:
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	token, err := s.sso.Login(ctx, login, in.Password, in.AppUuid)
	if err != nil {
		return nil, err
	}
	return &protosso.LoginResponse{Token: token}, nil
}

func (s *app) Register(
	ctx context.Context,
	in *protosso.RegisterRequest,
) (*protosso.RegisterResponse, error) {
	// TODO
	uuid, _ := s.sso.RegisterNewUser(ctx, in.Username, in.Email, in.Password)
	return &protosso.RegisterResponse{UserUuid: uuid}, nil
}
