package grpc

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	authdb "github.com/panduputragit/gym/backend/app/auth-service/internal/db/gen"
	"github.com/panduputragit/gym/backend/app/auth-service/internal/token"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authpb "github.com/panduputragit/gym/backend/packages/proto/auth/v1"
)

func (s *AuthServer) AdminLogout(
	ctx context.Context,
	req *authpb.AdminLogoutRequest,
) (*authpb.AdminLogoutResponse, error) {

	payload, err := s.requireAdminToken(ctx)
	if err != nil {
		return nil, err
	}

	adminID, _ := uuid.Parse(payload.UserID)

	err = s.Queries.RevokeAdminSession(ctx, authdb.RevokeAdminSessionParams{
		TokenID: payload.ID,
		RevokedBy: uuid.NullUUID{
			UUID:  adminID,
			Valid: true,
		},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to revoke session")
	}

	return &authpb.AdminLogoutResponse{
		Message: "logged out successfully",
	}, nil
}

func (s *AuthServer) requireAdminToken(ctx context.Context) (*token.Payload, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}

	tokenStr := bearerToken(values[0])

	payload, err := s.Token.VerifyToken(tokenStr)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	if payload.Role != "admin" {
		return nil, status.Error(codes.PermissionDenied, "forbidden")
	}

	session, err := s.Queries.GetAdminSession(ctx, payload.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, status.Error(codes.Unauthenticated, "session not found")
	}
	if err != nil {
		return nil, status.Error(codes.Internal, "database error")
	}

	if session.RevokedAt.Valid {
		return nil, status.Error(codes.Unauthenticated, "session revoked")
	}

	return payload, nil
}

func bearerToken(header string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return ""
	}
	return strings.TrimSpace(strings.TrimPrefix(header, prefix))
}
