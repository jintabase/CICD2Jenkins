package service

import (
	"context"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/domain"
)

type stubUserRepository struct {
	user *domain.User
	err  error
}

func (s stubUserRepository) FindByUsername(context.Context, string) (*domain.User, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.user, nil
}

func TestAuthServiceLoginSuccess(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	svc := NewAuthService(stubUserRepository{
		user: &domain.User{
			ID:           "u-1",
			Username:     "admin",
			PasswordHash: string(passwordHash),
			Role:         domain.RoleSuperAdmin,
		},
	}, "test-secret", time.Hour)

	result, err := svc.Login(context.Background(), "admin", "secret123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Token == "" {
		t.Fatal("expected token to be generated")
	}

	actor, err := svc.ParseToken(result.Token)
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	if actor.Role != domain.RoleSuperAdmin {
		t.Fatalf("expected role %s, got %s", domain.RoleSuperAdmin, actor.Role)
	}
}

func TestAuthServiceLoginInvalidPassword(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	svc := NewAuthService(stubUserRepository{
		user: &domain.User{
			ID:           "u-1",
			Username:     "admin",
			PasswordHash: string(passwordHash),
			Role:         domain.RoleSuperAdmin,
		},
	}, "test-secret", time.Hour)

	_, err = svc.Login(context.Background(), "admin", "wrong")
	if err != apperrors.ErrInvalidCredentials {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}
