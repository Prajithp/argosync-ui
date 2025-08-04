package repository

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/Prajithp/argosync/internal/repository/query"
	"github.com/Prajithp/argosync/pkg/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// ErrVersionExists is returned when trying to release a version that already exists

// ErrVersionExists is returned when trying to release a version that already exists
var ErrVersionExists = errors.New("version already exists")

// ReleaseRepository implements the Repository interface using GORM
type ReleaseRepository struct {
	db          *gorm.DB
	maxVersions int
}

// NewReleaseRepository creates a new ReleaseRepository
func NewReleaseRepository(db *gorm.DB, maxVersions int) (*ReleaseRepository, error) {
	return &ReleaseRepository{
		db:          db,
		maxVersions: maxVersions,
	}, nil
}

// Close closes the database connection
func (r *ReleaseRepository) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// MigrateSchema initializes the database schema using auto migration
func (r *ReleaseRepository) MigrateSchema() error {
	// Auto migrate the schema
	err := r.db.AutoMigrate(
		&models.Application{},
		&models.Environment{},
		&models.Region{},
		&models.Deployment{},
	)
	if err != nil {
		return fmt.Errorf("error migrating schema: %w", err)
	}

	return nil
}

// GetApplicationByName retrieves an application by name
func (r *ReleaseRepository) GetApplicationByName(name string) (*models.Application, error) {
	// Use the generated query interface
	app, err := query.Application.Where(query.Application.Name.Eq(name)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("application not found")
		}
		return nil, err
	}

	return app, nil
}

// CreateApplicationIfNotExists creates an application if it doesn't exist
func (r *ReleaseRepository) CreateApplicationIfNotExists(name string) (*models.Application, error) {
	// Try to get the application first
	app, err := r.GetApplicationByName(name)
	if err == nil {
		// Application already exists
		return app, nil
	}

	// Create the application
	newApp := &models.Application{
		Name: name,
	}

	// Use the generated query interface
	err = query.Application.Create(newApp)
	if err != nil {
		return nil, fmt.Errorf("error creating application: %w", err)
	}

	return newApp, nil
}

// GetEnvironmentByName retrieves an environment by name
func (r *ReleaseRepository) GetEnvironmentByName(name string) (*models.Environment, error) {
	// Use the generated query interface
	env, err := query.Environment.Where(query.Environment.Name.Eq(name)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("environment not found")
		}
		return nil, err
	}

	return env, nil
}

// CreateEnvironmentIfNotExists creates an environment if it doesn't exist
func (r *ReleaseRepository) CreateEnvironmentIfNotExists(name string) (*models.Environment, error) {
	// Try to get the environment first
	env, err := r.GetEnvironmentByName(name)
	if err == nil {
		// Environment already exists
		return env, nil
	}

	// Create the environment
	newEnv := &models.Environment{
		Name: name,
	}

	// Use the generated query interface
	err = query.Environment.Create(newEnv)
	if err != nil {
		return nil, fmt.Errorf("error creating environment: %w", err)
	}

	return newEnv, nil
}

// GetRegionByCode retrieves a region by code
func (r *ReleaseRepository) GetRegionByCode(code string) (*models.Region, error) {
	// Use the generated query interface
	region, err := query.Region.Where(query.Region.Code.Eq(code)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("region not found")
		}
		return nil, err
	}

	return region, nil
}

// CreateRegionIfNotExists creates a region if it doesn't exist
func (r *ReleaseRepository) CreateRegionIfNotExists(code string, name string) (*models.Region, error) {
	// Try to get the region first
	region, err := r.GetRegionByCode(code)
	if err == nil {
		// Region already exists
		return region, nil
	}

	// Create the region
	newRegion := &models.Region{
		Code: code,
		Name: name,
	}

	// Use the generated query interface
	err = query.Region.Create(newRegion)
	if err != nil {
		return nil, fmt.Errorf("error creating region: %w", err)
	}

	return newRegion, nil
}

// CheckVersionExists checks if a version already exists for an application in a specific environment and region
func (r *ReleaseRepository) CheckVersionExists(appID, envID, regionID uint, version string) (bool, error) {
	// Use the generated query interface
	count, err := query.Deployment.Where(
		query.Deployment.ApplicationID.Eq(appID),
		query.Deployment.EnvironmentID.Eq(envID),
		query.Deployment.RegionID.Eq(regionID),
		query.Deployment.Version.Eq(version),
	).Count()

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetActiveDeployment retrieves the active deployment for an application in a specific environment and region
func (r *ReleaseRepository) GetActiveDeployment(appID, envID, regionID uint) (*models.Deployment, error) {
	// Use the generated query interface
	deployment, err := query.Deployment.Where(
		query.Deployment.ApplicationID.Eq(appID),
		query.Deployment.EnvironmentID.Eq(envID),
		query.Deployment.RegionID.Eq(regionID),
		query.Deployment.Status.Eq("active"),
	).First()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No active deployment found
		}
		return nil, err
	}

	return deployment, nil
}

// CleanupOldVersions removes old versions exceeding the maximum limit
func (r *ReleaseRepository) CleanupOldVersions(appID, envID, regionID uint, maxVersions int) error {
	// Get all deployments for the application/environment/region using the generated query interface
	deployments, err := query.Deployment.Where(
		query.Deployment.ApplicationID.Eq(appID),
		query.Deployment.EnvironmentID.Eq(envID),
		query.Deployment.RegionID.Eq(regionID),
	).Order(query.Deployment.DeployedAt.Desc()).Find()

	if err != nil {
		return fmt.Errorf("error getting deployments: %w", err)
	}

	// If we have fewer deployments than the maximum, no cleanup needed
	if len(deployments) <= maxVersions {
		return nil
	}

	// Keep track of active deployments (we don't want to delete these)
	activeDeploymentIDs := make(map[uint]bool)
	for _, d := range deployments {
		if d.Status == "active" {
			activeDeploymentIDs[d.ID] = true
		}
	}

	// Identify deployments to delete (older than maxVersions, excluding active ones)
	deploymentsToDelete := make([]uint, 0)
	keptCount := 0

	for _, d := range deployments {
		if activeDeploymentIDs[d.ID] {
			// Always keep active deployments
			continue
		}

		keptCount++
		if keptCount > maxVersions {
			deploymentsToDelete = append(deploymentsToDelete, d.ID)
		}
	}

	// Delete old deployments using the generated query interface
	if len(deploymentsToDelete) > 0 {
		_, err := query.Deployment.Where(query.Deployment.ID.In(deploymentsToDelete...)).Delete()
		if err != nil {
			return fmt.Errorf("error deleting deployments: %w", err)
		}
	}

	return nil
}

// Release creates a new deployment and updates the status of the previous active deployment
func (r *ReleaseRepository) Release(req *models.ReleaseRequest) (*models.Deployment, error) {
	// Use the query interface's transaction support
	var newDeployment *models.Deployment
	err := query.Q.Transaction(func(tx *query.Query) error {
		// Create application if it doesn't exist
		app, err := r.CreateApplicationIfNotExists(req.Application)
		if err != nil {
			return fmt.Errorf("error creating application: %w", err)
		}

		// Create environment if it doesn't exist
		env, err := r.CreateEnvironmentIfNotExists(req.Environment)
		if err != nil {
			return fmt.Errorf("error creating environment: %w", err)
		}

		// Create region if it doesn't exist
		region, err := r.CreateRegionIfNotExists(req.Region, req.Region)
		if err != nil {
			return fmt.Errorf("error creating region: %w", err)
		}

		// Check if the version already exists
		exists, err := r.CheckVersionExists(app.ID, env.ID, region.ID, req.Version)
		if err != nil {
			return fmt.Errorf("error checking if version exists: %w", err)
		}

		if exists {
			return ErrVersionExists
		}

		// Get current active deployment
		currentDeployment, err := r.GetActiveDeployment(app.ID, env.ID, region.ID)
		if err != nil {
			return fmt.Errorf("error getting active deployment: %w", err)
		}

		// Create new deployment
		newDeployment = &models.Deployment{
			ApplicationID: app.ID,
			EnvironmentID: env.ID,
			RegionID:      region.ID,
			Version:       req.Version,
			Status:        "active",
			DeployedBy:    req.DeployedBy,
			DeployedAt:    time.Now(),
		}

		// If there's a current active deployment, set it as the rollback target and update its status
		if currentDeployment != nil {
			// Update the status of the current active deployment
			currentDeployment.Status = "inactive"
			err = tx.Deployment.Save(currentDeployment)
			if err != nil {
				return fmt.Errorf("error updating deployment status: %w", err)
			}

			// Set the rollback target
			newDeployment.RollbackTargetID = &currentDeployment.ID
		}

		// Create the new deployment
		err = tx.Deployment.Create(newDeployment)
		if err != nil {
			return fmt.Errorf("error creating deployment: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Clean up old versions if maxVersions is set
	// This is done outside the transaction to avoid locking issues
	if r.maxVersions > 0 {
		if err := r.CleanupOldVersions(newDeployment.ApplicationID, newDeployment.EnvironmentID, newDeployment.RegionID, r.maxVersions); err != nil {
			// Log the error but don't fail the deployment
			fmt.Printf("Warning: Failed to clean up old versions: %v\n", err)
		}
	}

	return newDeployment, nil
}

// Rollback rolls back to a previous deployment
func (r *ReleaseRepository) Rollback(req *models.RollbackRequest) (*models.Deployment, error) {
	var targetDeployment *models.Deployment

	// Use the query interface's transaction support
	err := query.Q.Transaction(func(tx *query.Query) error {
		app, err := r.GetApplicationByName(req.Application)
		if err != nil {
			return fmt.Errorf("error getting application: %w", err)
		}

		env, err := r.GetEnvironmentByName(req.Environment)
		if err != nil {
			return fmt.Errorf("error getting environment: %w", err)
		}

		region, err := r.GetRegionByCode(req.Region)
		if err != nil {
			return fmt.Errorf("error getting region: %w", err)
		}

		// Get current active deployment
		currentDeployment, err := r.GetActiveDeployment(app.ID, env.ID, region.ID)
		if err != nil {
			return fmt.Errorf("error getting active deployment: %w", err)
		}
		if currentDeployment == nil {
			return errors.New("no active deployment found to rollback")
		}
		// Get deployment history using the query interface
		deployments, err := tx.Deployment.Where(
			tx.Deployment.ApplicationID.Eq(app.ID),
			tx.Deployment.EnvironmentID.Eq(env.ID),
			tx.Deployment.RegionID.Eq(region.ID),
		).Order(tx.Deployment.DeployedAt.Desc()).Find()
		if err != nil {
			return fmt.Errorf("error getting deployment history: %w", err)
		}

		if len(deployments) <= 1 {
			return errors.New("no previous deployments found to rollback to")
		}

		var foundDeployment *models.Deployment
		if req.Version != "" {
			for i := range deployments {
				log.Info().Str("version", req.Version).Any("deployment", deployments[i]).Any("currentDeployment", currentDeployment).Msg("Checking deployment for rollback")
				if deployments[i].Version == req.Version && deployments[i].ID != currentDeployment.ID {
					foundDeployment = deployments[i]
					break
				}
			}

			if foundDeployment == nil {
				return fmt.Errorf("no deployment found with version %s", req.Version)
			}
		} else {
			for i := range deployments {
				if deployments[i].ID != currentDeployment.ID {
					foundDeployment = deployments[i]
					break
				}
			}
		}

		// Assign the found deployment to targetDeployment
		targetDeployment = foundDeployment

		// Update the status of the current active deployment
		currentDeployment.Status = "inactive"
		err = tx.Deployment.Save(currentDeployment)
		if err != nil {
			return fmt.Errorf("error updating deployment status: %w", err)
		}

		// Update the status of the target deployment
		targetDeployment.Status = "active"
		targetDeployment.DeployedBy = req.DeployedBy
		err = tx.Deployment.Save(targetDeployment)
		if err != nil {
			return fmt.Errorf("error updating deployment status: %w", err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return targetDeployment, nil
}

// GetAllDeployments returns all deployments in a format compatible with the frontend
// This method is kept for backward compatibility but is no longer used directly
// The baseRepository wrapper implements the new interface with pagination
func (r *ReleaseRepository) GetAllDeployments(limit int) ([]models.FrontendDeployment, error) {
	// Use the query interfaces to fetch deployments with preloaded relations
	deployments, err := query.Deployment.
		Preload(query.Deployment.Application).
		Preload(query.Deployment.Environment).
		Preload(query.Deployment.Region).
		Where(query.Deployment.Status.In("active", "inactive")).
		Order(query.Deployment.DeployedAt.Desc()).
		Find()

	if err != nil {
		return nil, err
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

	// Take the top 'limit' deployments from each group
	var result []models.FrontendDeployment

	for key, deps := range groupedDeployments {
		// Sort by deployed_at in descending order (should already be sorted from the query)
		count := 0
		for _, d := range deps {
			if count >= limit {
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

			result = append(result, fd)
			count++
		}
	}

	// Sort the final result by application, environment, region, and deployed_at
	sort.Slice(result, func(i, j int) bool {
		// First sort by application name
		if result[i].ApplicationName != result[j].ApplicationName {
			return result[i].ApplicationName < result[j].ApplicationName
		}

		// Then by environment
		if result[i].Environment != result[j].Environment {
			return result[i].Environment < result[j].Environment
		}

		// Then by region
		if result[i].Region != result[j].Region {
			return result[i].Region < result[j].Region
		}

		// Finally by timestamp (descending)
		return result[i].Timestamp > result[j].Timestamp
	})

	return result, nil
}

// GetDeploymentHistory retrieves the deployment history for an application in a specific environment and region
func (r *ReleaseRepository) GetDeploymentHistory(appID, envID, regionID uint) ([]*models.Deployment, error) {
	// Use the generated query interface
	deployments, err := query.Deployment.Where(
		query.Deployment.ApplicationID.Eq(appID),
		query.Deployment.EnvironmentID.Eq(envID),
		query.Deployment.RegionID.Eq(regionID),
	).Order(query.Deployment.DeployedAt.Desc()).Find()

	if err != nil {
		return nil, err
	}

	// Convert []*models.Deployment to []*models.Deployment (no conversion needed, but keeping the function signature)
	result := make([]*models.Deployment, len(deployments))
	for i, d := range deployments {
		result[i] = d
	}

	return result, nil
}

// Query executes a raw SQL query and returns the result
func (r *ReleaseRepository) Query(queryStr string, args ...interface{}) (interface{}, error) {
	// Use the underlying GORM DB to execute the raw SQL query
	rows, err := r.db.Raw(queryStr, args...).Rows()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// GetAllApplications returns all applications
func (r *ReleaseRepository) GetAllApplications() ([]*models.Application, error) {
	// Use the query interfaces to fetch all applications
	apps, err := query.Application.Find()
	if err != nil {
		return nil, err
	}
	return apps, nil
}

// GetRegionsForApplication returns all regions for a specific application
func (r *ReleaseRepository) GetRegionsForApplication(appID uint) ([]*models.Region, error) {
	// Use the query interfaces to fetch regions for a specific application
	// First get all deployments for the application
	deployments, err := query.Deployment.
		Where(query.Deployment.ApplicationID.Eq(appID)).
		Find()
	if err != nil {
		return nil, err
	}

	// Extract unique region IDs
	regionIDs := make(map[uint]bool)
	for _, d := range deployments {
		regionIDs[d.RegionID] = true
	}

	// Convert to slice
	ids := make([]uint, 0, len(regionIDs))
	for id := range regionIDs {
		ids = append(ids, id)
	}

	// If no regions found, return empty slice
	if len(ids) == 0 {
		return []*models.Region{}, nil
	}

	// Fetch regions by IDs
	regions, err := query.Region.
		Where(query.Region.ID.In(ids...)).
		Find()
	if err != nil {
		return nil, err
	}

	return regions, nil
}

// GetEnvironmentsForApplicationAndRegion returns all environments for a specific application and region
func (r *ReleaseRepository) GetEnvironmentsForApplicationAndRegion(appID, regionID uint) ([]*models.Environment, error) {
	// Use the query interfaces to fetch environments for a specific application and region
	// First get all deployments for the application and region
	deployments, err := query.Deployment.
		Where(
			query.Deployment.ApplicationID.Eq(appID),
			query.Deployment.RegionID.Eq(regionID),
		).
		Find()
	if err != nil {
		return nil, err
	}

	// Extract unique environment IDs
	envIDs := make(map[uint]bool)
	for _, d := range deployments {
		envIDs[d.EnvironmentID] = true
	}

	// Convert to slice
	ids := make([]uint, 0, len(envIDs))
	for id := range envIDs {
		ids = append(ids, id)
	}

	// If no environments found, return empty slice
	if len(ids) == 0 {
		return []*models.Environment{}, nil
	}

	// Fetch environments by IDs
	environments, err := query.Environment.
		Where(query.Environment.ID.In(ids...)).
		Find()
	if err != nil {
		return nil, err
	}

	return environments, nil
}

// GetVersionsForApplicationEnvironmentRegion returns all versions for a specific application, environment, and region
func (r *ReleaseRepository) GetVersionsForApplicationEnvironmentRegion(appID, envID, regionID uint) ([]*models.Deployment, error) {
	// Use the query interfaces to fetch versions for a specific application, environment, and region
	deployments, err := query.Deployment.
		Where(
			query.Deployment.ApplicationID.Eq(appID),
			query.Deployment.EnvironmentID.Eq(envID),
			query.Deployment.RegionID.Eq(regionID),
		).
		Order(query.Deployment.DeployedAt.Desc()).
		Find()
	if err != nil {
		return nil, err
	}

	return deployments, nil
}
