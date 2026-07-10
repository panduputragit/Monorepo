package grpc

import (
	"context"
	"database/sql"
	"errors"
	"time"

	authdb "github.com/panduputragit/gym/backend/app/auth-service/internal/db/gen"
	"github.com/panduputragit/gym/backend/app/auth-service/internal/token"
	authpb "github.com/panduputragit/gym/backend/packages/proto/auth/v1"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	authpb.UnimplementedAuthServiceServer

	Queries *authdb.Queries
	Token   *token.Maker
}

func (s *AuthServer) AdminLogin(
	ctx context.Context,
	req *authpb.AdminLoginRequest,
) (*authpb.AdminLoginResponse, error) {

	admin, err := s.Queries.GetAdminByEmail(ctx, req.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to fetch admin")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(admin.PasswordHash),
		[]byte(req.Password),
	); err != nil {
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	const duration = 8 * time.Hour

	tokenStr, payload, err := s.Token.CreateToken(
		admin.ID.String(),
		admin.Email,
		"admin",
		duration,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create token")
	}

	err = s.Queries.CreateAdminSession(ctx, authdb.CreateAdminSessionParams{
		AdminID:   admin.ID,
		TokenID:   payload.ID,
		ExpiresAt: payload.ExpiredAt,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to save session")
	}

	return &authpb.AdminLoginResponse{
		AccessToken: tokenStr,
		ExpiresIn:   int64(duration.Seconds()),
		Role:        "admin",
	}, nil
}
