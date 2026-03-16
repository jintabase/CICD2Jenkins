package config

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"cicd2jenkins/internal/model"
)

type Config struct {
	AppName   string
	Server    ServerConfig
	Auth      AuthConfig
	Database  DatabaseConfig
	SeedUsers []SeedUser
}

type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type AuthConfig struct {
	JWTSecret string
	TokenTTL  time.Duration
}

type DatabaseConfig struct {
	Driver string
	DSN    string
}

type SeedUser struct {
	Username string
	Password string
	Role     model.Role
}

func Load() Config {
	driver := envOrDefault("DB_DRIVER", "sqlite")

	return Config{
		AppName: envOrDefault("APP_NAME", "blog-api"),
		Server: ServerConfig{
			Host:         envOrDefault("APP_HOST", "0.0.0.0"),
			Port:         envOrDefault("APP_PORT", "8080"),
			ReadTimeout:  durationOrDefault("APP_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: durationOrDefault("APP_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  durationOrDefault("APP_IDLE_TIMEOUT", 60*time.Second),
		},
		Auth: AuthConfig{
			JWTSecret: envOrDefault("JWT_SECRET", "replace-this-in-production"),
			TokenTTL:  durationOrDefault("JWT_TTL", 24*time.Hour),
		},
		Database: DatabaseConfig{
			Driver: driver,
			DSN:    databaseDSN(driver),
		},
		SeedUsers: []SeedUser{
			{
				Username: envOrDefault("SUPER_ADMIN_USERNAME", "admin"),
				Password: envOrDefault("SUPER_ADMIN_PASSWORD", "Admin@123456"),
				Role:     model.RoleSuperAdmin,
			},
			{
				Username: envOrDefault("READER_USERNAME", "reader"),
				Password: envOrDefault("READER_PASSWORD", "Reader@123456"),
				Role:     model.RoleUser,
			},
		},
	}
}

func (c ServerConfig) BindAddress() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func durationOrDefault(key string, fallback time.Duration) time.Duration {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func databaseDSN(driver string) string {
	driver = strings.ToLower(strings.TrimSpace(driver))
	key := "DB_DSN"
	if driver == "mysql" {
		return envOrDefault(key, "blog:blog@tcp(127.0.0.1:3306)/blog_api?charset=utf8mb4&parseTime=True&loc=Local")
	}
	return envOrDefault(key, fmt.Sprintf("%s.db", envOrDefault("APP_NAME", "blog-api")))
}
