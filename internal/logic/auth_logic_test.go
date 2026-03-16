package logic

import (
	"context"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/model"
)

type stubUserRepository struct {
	user *model.User
	err  error
}

func (s stubUserRepository) FindByUsername(context.Context, string) (*model.User, error) {
	if s.err != nil {
		return nil, s.err
	}
	return s.user, nil
}

func TestAuthLogicLoginSuccess(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	svc := NewAuthLogic(stubUserRepository{
		user: &model.User{
			ID:           "u-1",
			Username:     "admin",
			PasswordHash: string(passwordHash),
			Role:         model.RoleSuperAdmin,
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
	if actor.Role != model.RoleSuperAdmin {
		t.Fatalf("expected role %s, got %s", model.RoleSuperAdmin, actor.Role)
	}
}

func TestAuthLogicLoginInvalidPassword(t *testing.T) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	svc := NewAuthLogic(stubUserRepository{
		user: &model.User{
			ID:           "u-1",
			Username:     "admin",
			PasswordHash: string(passwordHash),
			Role:         model.RoleSuperAdmin,
		},
	}, "test-secret", time.Hour)

	_, err = svc.Login(context.Background(), "admin", "wrong")
	if err != apperrors.ErrInvalidCredentials {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}
