package grpc

import (
	"context"
	"fmt"

	authpb "github.com/panduputragit/gym/backend/packages/proto/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RegisterRoutes(grpcServer *grpc.Server, authGRPCAddr string) error {
	conn, err := grpc.NewClient(authGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("dial auth service: %w", err)
	}

	client := authpb.NewAuthServiceClient(conn)
	authpb.RegisterAuthServiceServer(grpcServer, &proxyServer{authClient: client})

	return nil
}

type proxyServer struct {
	authpb.UnimplementedAuthServiceServer
	authClient authpb.AuthServiceClient
}

func (p *proxyServer) AdminLogin(ctx context.Context, req *authpb.AdminLoginRequest) (*authpb.AdminLoginResponse, error) {
	return p.authClient.AdminLogin(ctx, req)
}
