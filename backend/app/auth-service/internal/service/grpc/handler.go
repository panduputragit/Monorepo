package grpc

import (
	authdb "github.com/panduputragit/gym/backend/app/auth-service/internal/db/gen"
	"github.com/panduputragit/gym/backend/app/auth-service/internal/token"
	authpb "github.com/panduputragit/gym/backend/packages/proto/auth/v1"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer

	Queries *authdb.Queries
	Token   *token.Maker
}
