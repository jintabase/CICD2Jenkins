package config

import (
	"net"
	"os"
	"strings"
	"time"

	"cicd2jenkins/internal/domain"
)

type Config struct {
	AppName       string
	Server        ServerConfig
	Auth          AuthConfig
	Elasticsearch ElasticsearchConfig
	SeedUsers     []SeedUser
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

type ElasticsearchConfig struct {
	Addresses      []string
	Username       string
	Password       string
	Index          string
	RequestTimeout time.Duration
}

type SeedUser struct {
	Username string
	Password string
	Role     domain.Role
}

func Load() Config {
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
		Elasticsearch: ElasticsearchConfig{
			Addresses:      splitCSV(envOrDefault("ES_ADDRESSES", "http://localhost:9200")),
			Username:       strings.TrimSpace(os.Getenv("ES_USERNAME")),
			Password:       strings.TrimSpace(os.Getenv("ES_PASSWORD")),
			Index:          envOrDefault("ES_INDEX", "blog_articles"),
			RequestTimeout: durationOrDefault("ES_REQUEST_TIMEOUT", 5*time.Second),
		},
		SeedUsers: []SeedUser{
			{
				Username: envOrDefault("SUPER_ADMIN_USERNAME", "admin"),
				Password: envOrDefault("SUPER_ADMIN_PASSWORD", "Admin@123456"),
				Role:     domain.RoleSuperAdmin,
			},
			{
				Username: envOrDefault("READER_USERNAME", "reader"),
				Password: envOrDefault("READER_PASSWORD", "Reader@123456"),
				Role:     domain.RoleUser,
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

func splitCSV(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}

	parts := strings.Split(value, ",")
	addresses := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			addresses = append(addresses, trimmed)
		}
	}
	return addresses
}
