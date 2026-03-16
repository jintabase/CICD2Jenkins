package logic

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"cicd2jenkins/internal/apperrors"
	"cicd2jenkins/internal/model"
	"cicd2jenkins/internal/repo"
)

type AuthLogic struct {
	users     repo.UserRepository
	jwtSecret []byte
	tokenTTL  time.Duration
}

type LoginResult struct {
	Token string     `json:"token"`
	User  model.User `json:"user"`
}

type Claims struct {
	Role     model.Role `json:"role"`
	Username string     `json:"username"`
	jwt.RegisteredClaims
}

func NewAuthLogic(users repo.UserRepository, secret string, tokenTTL time.Duration) *AuthLogic {
	return &AuthLogic{
		users:     users,
		jwtSecret: []byte(secret),
		tokenTTL:  tokenTTL,
	}
}

func (l *AuthLogic) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || password == "" {
		return nil, apperrors.ErrBadRequest
	}

	user, err := l.users.FindByUsername(ctx, username)
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
			ExpiresAt: jwt.NewNumericDate(now.Add(l.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(l.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}

	return &LoginResult{
		Token: signed,
		User:  *user,
	}, nil
}

func (l *AuthLogic) ParseToken(tokenString string) (model.Actor, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return model.Actor{}, apperrors.ErrUnauthorized
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return l.jwtSecret, nil
	})
	if err != nil {
		return model.Actor{}, apperrors.ErrUnauthorized
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return model.Actor{}, apperrors.ErrUnauthorized
	}

	return model.Actor{
		UserID:   claims.Subject,
		Username: claims.Username,
		Role:     claims.Role,
	}, nil
}
