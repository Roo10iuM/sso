package grpcapp

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/roo10ium/sso-protos/gen/go/pbsso"
)

type gRPCServer struct {
	pbsso.UnimplementedSSOServer
	sso ssoService
}

type ssoService interface {
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

func (s *gRPCServer) Login(
	ctx context.Context,
	in *pbsso.LoginRequest,
) (*pbsso.LoginResponse, error) {
	// TODO
	var login string
	switch in.Login.(type) {
	case *pbsso.LoginRequest_Email:
		login = in.GetEmail()
	case *pbsso.LoginRequest_Username:
		login = in.GetUsername()
	default:
		return nil, status.Error(codes.InvalidArgument, "login is required")
	}
	token, err := s.sso.Login(ctx, login, in.Password, in.AppUuid)
	if err != nil {
		return nil, err
	}
	return &pbsso.LoginResponse{Token: token}, nil
}

func (s *gRPCServer) Register(
	ctx context.Context,
	in *pbsso.RegisterRequest,
) (*pbsso.RegisterResponse, error) {
	// TODO
	uuid, _ := s.sso.RegisterNewUser(ctx, in.Username, in.Email, in.Password)
	return &pbsso.RegisterResponse{UserUuid: uuid}, nil
}
