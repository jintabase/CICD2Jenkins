package gormrepo

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	if err := r.db.WithContext(ctx).
		Where("LOWER(username) = ?", strings.ToLower(strings.TrimSpace(username))).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("find user by username: %w", err)
	}
	return &user, nil
}
