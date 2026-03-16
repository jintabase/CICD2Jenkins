package memory

import (
	"context"
	"strings"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/domain"
)

type UserRepository struct {
	users map[string]domain.User
}

func NewUserRepository(users []domain.User) *UserRepository {
	store := make(map[string]domain.User, len(users))
	for _, user := range users {
		store[strings.ToLower(user.Username)] = user
	}

	return &UserRepository{users: store}
}

func (r *UserRepository) FindByUsername(_ context.Context, username string) (*domain.User, error) {
	user, ok := r.users[strings.ToLower(strings.TrimSpace(username))]
	if !ok {
		return nil, apperrors.ErrNotFound
	}
	return &user, nil
}
