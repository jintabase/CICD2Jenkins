package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/domain"
	"cicd2jenkins/internal/repository"
)

type AuthService struct {
	users     repository.UserRepository
	jwtSecret []byte
	tokenTTL  time.Duration
}

type LoginResult struct {
	Token string      `json:"token"`
	User  domain.User `json:"user"`
}

type Claims struct {
	Role     domain.Role `json:"role"`
	Username string      `json:"username"`
	jwt.RegisteredClaims
}

func NewAuthService(users repository.UserRepository, secret string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		users:     users,
		jwtSecret: []byte(secret),
		tokenTTL:  tokenTTL,
	}
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || password == "" {
		return nil, apperrors.ErrBadRequest
	}

	user, err := s.users.FindByUsername(ctx, username)
	if err != nil {
		if err == apperrors.ErrNotFound {
			return nil, apperrors.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, apperrors.ErrInvalidCredentials
	}

	now := time.Now().UTC()
	claims := Claims{
		Role:     user.Role,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}

	return &LoginResult{
		Token: signed,
		User:  *user,
	}, nil
}

func (s *AuthService) ParseToken(tokenString string) (domain.Actor, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return domain.Actor{}, apperrors.ErrUnauthorized
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return domain.Actor{}, apperrors.ErrUnauthorized
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return domain.Actor{}, apperrors.ErrUnauthorized
	}

	return domain.Actor{
		UserID:   claims.Subject,
		Username: claims.Username,
		Role:     claims.Role,
	}, nil
}
