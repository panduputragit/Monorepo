package auth

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type Service struct{}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type ValidateTokenResponse struct {
	Valid  bool   `json:"valid"`
	UserID string `json:"user_id,omitempty"`
	Email  string `json:"email,omitempty"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return nil, ErrInvalidCredentials
	}

	return &LoginResponse{
		AccessToken:  "placeholder_token",
		RefreshToken: "placeholder_refresh",
		ExpiresIn:    3600,
	}, nil
}

func (s *Service) ValidateToken(ctx context.Context, token string) (*ValidateTokenResponse, error) {
	if strings.TrimSpace(token) == "" {
		return &ValidateTokenResponse{Valid: false}, nil
	}

	return &ValidateTokenResponse{
		Valid:  true,
		UserID: "placeholder_user_id",
		Email:  "placeholder@example.com",
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*RefreshResponse, error) {
	if strings.TrimSpace(refreshToken) == "" {
		return nil, ErrInvalidCredentials
	}

	return &RefreshResponse{
		AccessToken: "new_placeholder_token",
		ExpiresIn:   3600,
	}, nil
}
