package repository

import (
	"fmt"

	"github.com/Prajithp/argosync/internal/repository/query"
	"github.com/Prajithp/argosync/pkg/models"
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

	// Initialize query interfaces
	query.SetDefault(db)

	return &testRepository{repo}, nil
}

// testRepository is a wrapper around ReleaseRepository that implements the Repository interface
type testRepository struct {
	*ReleaseRepository
}

// GetAllDeployments implements the Repository interface with pagination
func (r *testRepository) GetAllDeployments(page, pageSize int) ([]models.FrontendDeployment, int, error) {
	// For testing, just return empty results
	return []models.FrontendDeployment{}, 0, nil
}

// GetAllApplications implements the Repository interface
func (r *testRepository) GetAllApplications() ([]*models.Application, error) {
	// For testing, just return empty results
	return []*models.Application{}, nil
}

// GetRegionsForApplication implements the Repository interface
func (r *testRepository) GetRegionsForApplication(appID uint) ([]*models.Region, error) {
	// For testing, just return empty results
	return []*models.Region{}, nil
}

// GetEnvironmentsForApplicationAndRegion implements the Repository interface
func (r *testRepository) GetEnvironmentsForApplicationAndRegion(appID, regionID uint) ([]*models.Environment, error) {
	// For testing, just return empty results
	return []*models.Environment{}, nil
}

// GetVersionsForApplicationEnvironmentRegion implements the Repository interface
func (r *testRepository) GetVersionsForApplicationEnvironmentRegion(appID, envID, regionID uint) ([]*models.Deployment, error) {
	// For testing, just return empty results
	return []*models.Deployment{}, nil
}
