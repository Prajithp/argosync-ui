package repository

import (
	"fmt"

	zerologgorm "github.com/go-mods/zerolog-gorm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitSQLiteTestDB initializes a SQLite database for testing
func InitSQLiteTestDB() (Repository, error) {
	// Open the SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: &zerologgorm.GormLogger{
			FieldsExclude: []string{zerologgorm.DurationFieldName, zerologgorm.FileFieldName},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error opening SQLite database: %w", err)
	}

	repo, err := NewReleaseRepository(db, 10)
	if err != nil {
		return nil, fmt.Errorf("error creating SQLite database: %w", err)
	}

	// Migrate the schema
	if err := repo.MigrateSchema(); err != nil {
		return nil, fmt.Errorf("error migrating schema: %w", err)
	}

	return repo, nil
}
