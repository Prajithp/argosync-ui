package repository

import (
	"fmt"
	"sort"
	"time"

	"github.com/Prajithp/argosync/internal/config"
	"github.com/Prajithp/argosync/internal/logger"
	"github.com/Prajithp/argosync/internal/repository/query"
	"github.com/Prajithp/argosync/pkg/models"
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
	default:
		repo, err = initSQLite(cfg)
		if err != nil {
			logger.Fatal().Err(err).Msg("Failed to connect to SQLite database")
			return nil, err
		}
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

	// Initialize query interfaces with the database connection
	query.SetDefault(db)

	// Wrap the repository to implement the updated interface
	return &baseRepository{repo}, nil
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

	// Initialize query interfaces with the database connection
	query.SetDefault(db)

	// Wrap the repository to implement the updated interface
	return &baseRepository{repo}, nil
}

// baseRepository is a wrapper around ReleaseRepository that implements the Repository interface
type baseRepository struct {
	*ReleaseRepository
}

// GetAllApplications returns all applications
func (r *baseRepository) GetAllApplications() ([]*models.Application, error) {
	return r.ReleaseRepository.GetAllApplications()
}

// GetRegionsForApplication returns all regions for a specific application
func (r *baseRepository) GetRegionsForApplication(appID uint) ([]*models.Region, error) {
	return r.ReleaseRepository.GetRegionsForApplication(appID)
}

// GetEnvironmentsForApplicationAndRegion returns all environments for a specific application and region
func (r *baseRepository) GetEnvironmentsForApplicationAndRegion(appID, regionID uint) ([]*models.Environment, error) {
	return r.ReleaseRepository.GetEnvironmentsForApplicationAndRegion(appID, regionID)
}

// GetVersionsForApplicationEnvironmentRegion returns all versions for a specific application, environment, and region
func (r *baseRepository) GetVersionsForApplicationEnvironmentRegion(appID, envID, regionID uint) ([]*models.Deployment, error) {
	return r.ReleaseRepository.GetVersionsForApplicationEnvironmentRegion(appID, envID, regionID)
}

// GetAllDeployments implements the Repository interface with pagination
func (r *baseRepository) GetAllDeployments(page, pageSize int) ([]models.FrontendDeployment, int, error) {
	// Use the query interfaces to fetch deployments with preloaded relations
	deployments, err := query.Deployment.
		Preload(query.Deployment.Application).
		Preload(query.Deployment.Environment).
		Preload(query.Deployment.Region).
		Where(query.Deployment.Status.In("active", "inactive")).
		Order(query.Deployment.DeployedAt.Desc()).
		Find()

	if err != nil {
		return nil, 0, err
	}

	// Group deployments by application, environment, and region
	type groupKey struct {
		AppName    string
		EnvName    string
		RegionCode string
	}
	
	groupedDeployments := make(map[groupKey][]*models.Deployment)
	
	for _, d := range deployments {
		// Check if any of the relations are zero values
		if d.Application.ID == 0 || d.Environment.ID == 0 || d.Region.ID == 0 {
			continue // Skip deployments with missing relations
		}
		
		key := groupKey{
			AppName:    d.Application.Name,
			EnvName:    d.Environment.Name,
			RegionCode: d.Region.Code,
		}
		
		groupedDeployments[key] = append(groupedDeployments[key], d)
	}
	
	// Take the top deployments from each group
	var allResults []models.FrontendDeployment
	
	for key, deps := range groupedDeployments {
		// Sort by deployed_at in descending order (should already be sorted from the query)
		count := 0
		for _, d := range deps {
			if count >= pageSize {
				break
			}
			
			fd := models.FrontendDeployment{
				ApplicationName: key.AppName,
				Environment:     key.EnvName,
				Region:          key.RegionCode,
				Version:         d.Version,
				Timestamp:       d.DeployedAt.Format(time.RFC3339),
				Status:          d.Status,
				DeployedBy:      d.DeployedBy,
			}
			
			allResults = append(allResults, fd)
			count++
		}
	}
	
	// Sort the final result by application, environment, region, and deployed_at
	sort.Slice(allResults, func(i, j int) bool {
		// First sort by application name
		if allResults[i].ApplicationName != allResults[j].ApplicationName {
			return allResults[i].ApplicationName < allResults[j].ApplicationName
		}
		
		// Then by environment
		if allResults[i].Environment != allResults[j].Environment {
			return allResults[i].Environment < allResults[j].Environment
		}
		
		// Then by region
		if allResults[i].Region != allResults[j].Region {
			return allResults[i].Region < allResults[j].Region
		}
		
		// Finally by timestamp (descending)
		return allResults[i].Timestamp > allResults[j].Timestamp
	})
	
	// Calculate total count for pagination
	totalCount := len(allResults)
	
	// Apply pagination
	startIndex := (page - 1) * pageSize
	if startIndex >= totalCount {
		// If page is out of range, return empty result
		return []models.FrontendDeployment{}, totalCount, nil
	}
	
	endIndex := startIndex + pageSize
	if endIndex > totalCount {
		endIndex = totalCount
	}
	
	// Return paginated result
	return allResults[startIndex:endIndex], totalCount, nil
}
