package repository

import (
	"errors"
	"fmt"

	"github.com/Prajithp/argosync/internal/config"
	"github.com/Prajithp/argosync/internal/logger"
	zerologgorm "github.com/go-mods/zerolog-gorm"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBType represents the type of database
type DBType string

const (
	// PostgreSQL database
	PostgreSQL DBType = "postgres"
	// SQLite database
	SQLite DBType = "sqlite3"
)

// BaseRepository defines common methods for all repositories
type BaseRepository interface {
	// Schema operations
	MigrateSchema() error
	Close() error
}

var (
	// ErrUnsupportedDBType is returned when an unsupported database type is specified
	ErrUnsupportedDBType = errors.New("unsupported database type")
)

// InitRepositoryFromConfig initializes a repository from the application config
func InitRepositoryFromConfig(cfg *config.Config) (Repository, error) {
	var repo Repository
	var err error

	switch cfg.DBType {
	case string(PostgreSQL):
		repo, err = initPostgresql(cfg)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to connect to PostgreSQL database")
			return nil, err
		}
	case string(SQLite):
		repo, err = initSQLite(cfg)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to connect to SQLite database")
			return nil, err
		}
	default:
		return nil, ErrUnsupportedDBType
	}

	logger.Info().Msg("Migrating database schema...")
	err = repo.MigrateSchema()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to migrate schema")
		return nil, err
	}

	logger.Info().Msg("Database setup completed successfully")

	return repo, nil
}

func initPostgresql(cfg *config.Config) (Repository, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
		cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: &zerologgorm.GormLogger{
			FieldsExclude: []string{zerologgorm.DurationFieldName, zerologgorm.FileFieldName},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error opening PostgreSQL database: %w", err)
	}

	repo, err := NewReleaseRepository(db, cfg.MaxVersions)
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func initSQLite(cfg *config.Config) (Repository, error) {
	db, err := gorm.Open(sqlite.Open(cfg.SQLitePath), &gorm.Config{
		Logger: &zerologgorm.GormLogger{
			FieldsExclude: []string{zerologgorm.DurationFieldName, zerologgorm.FileFieldName},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error opening SQLite database: %w", err)
	}

	repo, err := NewReleaseRepository(db, cfg.MaxVersions)
	if err != nil {
		return nil, err
	}

	return repo, nil
}
