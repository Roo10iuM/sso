package main

import (
	"context"
	"fmt"

	protosso "github.com/roo10ium/sso-protos/gen/go/sso"
	"github.com/roo10ium/sso/internal/util/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log := log.SetupLogger("local")
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	serverAddr := "localhost:8011"

	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		log.Error(fmt.Sprintf("fatal: %v", err))
		return
	}
	defer conn.Close()

	client := protosso.NewSSOClient(conn)

	ctx := context.Background()
	mail := "d@d.d"
	log.Info("start")
	client.Register(ctx, &protosso.RegisterRequest{Username: "dan", Email: &mail, Password: "12345678"})
}
