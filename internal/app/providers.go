package app

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"cicd2jenkins/internal/config"
	"cicd2jenkins/internal/logic"
	"cicd2jenkins/internal/model"
	"cicd2jenkins/internal/repo"
	"cicd2jenkins/internal/repo/gormrepo"
)

func provideDatabase(cfg config.Config, seedUsers []model.User) (*gorm.DB, func(), error) {
	dialector, err := provideDialector(cfg)
	if err != nil {
		return nil, nil, err
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("open database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, fmt.Errorf("get raw database handle: %w", err)
	}
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(20)

	cleanup := func() {
		_ = sqlDB.Close()
	}

	if err := db.AutoMigrate(&model.User{}, &model.Article{}); err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("auto migrate schema: %w", err)
	}

	if err := seedUsersToDB(db, seedUsers); err != nil {
		cleanup()
		return nil, nil, fmt.Errorf("seed users: %w", err)
	}

	return db, cleanup, nil
}

func provideSeedUsers(cfg config.Config) ([]model.User, error) {
	users := make([]model.User, 0, len(cfg.SeedUsers))
	now := time.Now().UTC()

	for _, seed := range cfg.SeedUsers {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(seed.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("hash seed user password: %w", err)
		}

		users = append(users, model.User{
			ID:           uuid.NewString(),
			Username:     strings.TrimSpace(seed.Username),
			PasswordHash: string(passwordHash),
			Role:         seed.Role,
			CreatedAt:    now,
			UpdatedAt:    now,
		})
	}

	return users, nil
}

func provideUserRepository(db *gorm.DB) *gormrepo.UserRepository {
	return gormrepo.NewUserRepository(db)
}

func provideArticleRepository(db *gorm.DB) *gormrepo.ArticleRepository {
	return gormrepo.NewArticleRepository(db)
}

func provideAuthLogic(cfg config.Config, users repo.UserRepository) *logic.AuthLogic {
	return logic.NewAuthLogic(users, cfg.Auth.JWTSecret, cfg.Auth.TokenTTL)
}

func provideHTTPServer(cfg config.Config, router http.Handler) *http.Server {
	return &http.Server{
		Addr:         net.JoinHostPort(cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}
}

func provideDialector(cfg config.Config) (gorm.Dialector, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Database.Driver)) {
	case "sqlite":
		return sqlite.Open(cfg.Database.DSN), nil
	case "mysql":
		return mysql.Open(cfg.Database.DSN), nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Database.Driver)
	}
}

func seedUsersToDB(db *gorm.DB, users []model.User) error {
	for _, user := range users {
		entry := user
		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "username"}},
			DoUpdates: clause.Assignments(map[string]any{
				"password_hash": entry.PasswordHash,
				"role":          entry.Role,
				"updated_at":    entry.UpdatedAt,
			}),
		}).Create(&entry).Error; err != nil {
			return err
		}
	}
	return nil
}
