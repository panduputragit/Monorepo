package route

import (
	"database/sql"

	authdb "github.com/panduputragit/gym/backend/app/auth-service/internal/db/gen"
	authgrpc "github.com/panduputragit/gym/backend/app/auth-service/internal/service/grpc"
	"github.com/panduputragit/gym/backend/app/auth-service/internal/token"
	authpb "github.com/panduputragit/gym/backend/packages/proto/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func RegisterRoutes(grpcServer *grpc.Server, db *sql.DB, tokenMaker *token.Maker) error {
	authpb.RegisterAuthServiceServer(grpcServer, &authgrpc.AuthServer{
		Queries: authdb.New(db),
		Token:   tokenMaker,
	})
	reflection.Register(grpcServer)
	return nil
}
