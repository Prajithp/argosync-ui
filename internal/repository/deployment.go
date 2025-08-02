package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/Prajithp/argosync/pkg/models"
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
	var app models.Application

	// Use raw SQL query to avoid GORM's automatic addition of deleted_at IS NULL
	if err := r.db.Raw("SELECT * FROM applications WHERE name = ? LIMIT 1", name).Scan(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("application not found")
		}
		return nil, err
	}

	if app.ID == 0 {
		return nil, errors.New("application not found")
	}

	return &app, nil
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

	if err := r.db.Create(newApp).Error; err != nil {
		return nil, fmt.Errorf("error creating application: %w", err)
	}

	return newApp, nil
}

// GetEnvironmentByName retrieves an environment by name
func (r *ReleaseRepository) GetEnvironmentByName(name string) (*models.Environment, error) {
	var env models.Environment

	// Use raw SQL query to avoid GORM's automatic addition of deleted_at IS NULL
	if err := r.db.Raw("SELECT * FROM environments WHERE name = ? LIMIT 1", name).Scan(&env).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("environment not found")
		}
		return nil, err
	}

	if env.ID == 0 {
		return nil, errors.New("environment not found")
	}

	return &env, nil
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

	if err := r.db.Create(newEnv).Error; err != nil {
		return nil, fmt.Errorf("error creating environment: %w", err)
	}

	return newEnv, nil
}

// GetRegionByCode retrieves a region by code
func (r *ReleaseRepository) GetRegionByCode(code string) (*models.Region, error) {
	var region models.Region

	// Use raw SQL query to avoid GORM's automatic addition of deleted_at IS NULL
	if err := r.db.Raw("SELECT * FROM regions WHERE code = ? LIMIT 1", code).Scan(&region).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("region not found")
		}
		return nil, err
	}

	if region.ID == 0 {
		return nil, errors.New("region not found")
	}

	return &region, nil
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

	if err := r.db.Create(newRegion).Error; err != nil {
		return nil, fmt.Errorf("error creating region: %w", err)
	}

	return newRegion, nil
}

// CheckVersionExists checks if a version already exists for an application in a specific environment and region
func (r *ReleaseRepository) CheckVersionExists(appID, envID, regionID uint, version string) (bool, error) {
	var count int64

	// Use raw SQL query to avoid GORM's automatic addition of deleted_at IS NULL
	err := r.db.Raw(`
		SELECT COUNT(*) FROM deployments
		WHERE application_id = ? AND environment_id = ? AND region_id = ? AND version = ?`,
		appID, envID, regionID, version).Scan(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetActiveDeployment retrieves the active deployment for an application in a specific environment and region
func (r *ReleaseRepository) GetActiveDeployment(appID, envID, regionID uint) (*models.Deployment, error) {
	var deployment models.Deployment

	// Use raw SQL query to avoid GORM's automatic addition of deleted_at IS NULL
	err := r.db.Raw(`
		SELECT * FROM deployments
		WHERE application_id = ? AND environment_id = ? AND region_id = ? AND status = ?
		LIMIT 1`,
		appID, envID, regionID, "active").Scan(&deployment).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No active deployment found
		}
		return nil, err
	}

	if deployment.ID == 0 {
		return nil, nil // No active deployment found
	}

	return &deployment, nil
}

// CleanupOldVersions removes old versions exceeding the maximum limit
func (r *ReleaseRepository) CleanupOldVersions(appID, envID, regionID uint, maxVersions int) error {
	// Get all deployments for the application/environment/region using raw SQL
	var deployments []models.Deployment
	if err := r.db.Raw(`
		SELECT * FROM deployments
		WHERE application_id = ? AND environment_id = ? AND region_id = ?
		ORDER BY deployed_at DESC`,
		appID, envID, regionID).Scan(&deployments).Error; err != nil {
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

	// Delete old deployments using raw SQL
	if len(deploymentsToDelete) > 0 {
		query := "DELETE FROM deployments WHERE id IN ("
		for i, id := range deploymentsToDelete {
			if i > 0 {
				query += ","
			}
			query += fmt.Sprintf("%d", id)
		}
		query += ")"

		if err := r.db.Exec(query).Error; err != nil {
			return fmt.Errorf("error deleting deployments: %w", err)
		}
	}

	return nil
}

// Release creates a new deployment and updates the status of the previous active deployment
func (r *ReleaseRepository) Release(req *models.ReleaseRequest) (*models.Deployment, error) {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create application if it doesn't exist
	app, err := r.CreateApplicationIfNotExists(req.Application)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating application: %w", err)
	}

	// Create environment if it doesn't exist
	env, err := r.CreateEnvironmentIfNotExists(req.Environment)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating environment: %w", err)
	}

	// Create region if it doesn't exist
	region, err := r.CreateRegionIfNotExists(req.Region, req.Region)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating region: %w", err)
	}

	// Check if the version already exists
	exists, err := r.CheckVersionExists(app.ID, env.ID, region.ID, req.Version)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error checking if version exists: %w", err)
	}

	if exists {
		tx.Rollback()
		return nil, ErrVersionExists
	}

	// Get current active deployment
	currentDeployment, err := r.GetActiveDeployment(app.ID, env.ID, region.ID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error getting active deployment: %w", err)
	}

	// Create new deployment
	newDeployment := &models.Deployment{
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
		if err := tx.Save(currentDeployment).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error updating deployment status: %w", err)
		}

		// Set the rollback target
		newDeployment.RollbackTargetID = &currentDeployment.ID
	}

	// Create the new deployment
	if err := tx.Create(newDeployment).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating deployment: %w", err)
	}

	// Clean up old versions if maxVersions is set
	if r.maxVersions > 0 {
		if err := r.CleanupOldVersions(app.ID, env.ID, region.ID, r.maxVersions); err != nil {
			// Log the error but don't fail the deployment
			fmt.Printf("Warning: Failed to clean up old versions: %v\n", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return newDeployment, nil
}

// Rollback rolls back to a previous deployment
func (r *ReleaseRepository) Rollback(req *models.RollbackRequest) (*models.Deployment, error) {
	// Start a transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	app, err := r.GetApplicationByName(req.Application)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error getting application: %w", err)
	}

	env, err := r.GetEnvironmentByName(req.Environment)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error getting environment: %w", err)
	}

	region, err := r.GetRegionByCode(req.Region)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error getting region: %w", err)
	}
	// Get current active deployment
	currentDeployment, err := r.GetActiveDeployment(app.ID, env.ID, region.ID)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error getting active deployment: %w", err)
	}
	if currentDeployment == nil {
		tx.Rollback()
		return nil, errors.New("no active deployment found to rollback")
	}

	// Get deployment history
	var deployments []models.Deployment
	if err := tx.Where("application_id = ? AND environment_id = ? AND region_id = ?",
		app.ID, env.ID, region.ID).Order("deployed_at DESC").Find(&deployments).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error getting deployment history: %w", err)
	}

	if len(deployments) <= 1 {
		tx.Rollback()
		return nil, errors.New("no previous deployments found to rollback to")
	}

	// Find the target deployment to rollback to
	var targetDeployment *models.Deployment
	if req.Version != "" {
		// Find deployment with the specified version
		for i := range deployments {
			if deployments[i].Version == req.Version && deployments[i].ID != currentDeployment.ID {
				targetDeployment = &deployments[i]
				break
			}
		}
		if targetDeployment == nil {
			tx.Rollback()
			return nil, fmt.Errorf("no deployment found with version %s", req.Version)
		}
	} else {
		for i := range deployments {
			if deployments[i].ID != currentDeployment.ID {
				targetDeployment = &deployments[i]
				break
			}
		}
	}

	// Update the status of the current active deployment
	currentDeployment.Status = "inactive"
	if err := tx.Save(currentDeployment).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error updating deployment status: %w", err)
	}

	// Update the status of the target deployment
	targetDeployment.Status = "active"
	if err := tx.Save(targetDeployment).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error updating deployment status: %w", err)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return targetDeployment, nil
}

// GetAllDeployments returns all deployments in a format compatible with the frontend
func (r *ReleaseRepository) GetAllDeployments(limit int) ([]models.FrontendDeployment, error) {
	// Use a raw SQL query with window functions to limit the number of versions per application/environment/region
	rows, err := r.db.Raw(`
		WITH ranked_deployments AS (
			SELECT
				a.name as application_name,
				e.name as environment,
				r.code as region_code,
				d.version,
				d.deployed_at,
				d.status,
				d.deployed_by,
				ROW_NUMBER() OVER (
					PARTITION BY a.name, e.name, r.code 
					ORDER BY d.deployed_at DESC
				) as row_num
			FROM deployments d
			JOIN applications a ON d.application_id = a.id
			JOIN environments e ON d.environment_id = e.id
			JOIN regions r ON d.region_id = r.id
			WHERE d.status IN ('active', 'inactive')
		)
		SELECT
			application_name,
			environment,
			region_code,
			version,
			deployed_at,
			status,
			deployed_by
		FROM ranked_deployments
		WHERE row_num <= ?
		ORDER BY application_name, environment, region_code, deployed_at DESC
	`, limit).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deployments []models.FrontendDeployment
	for rows.Next() {
		var (
			appName    string
			envName    string
			regionCode string
			version    string
			deployedAt time.Time
			status     string
			deployedBy string
		)

		if err := rows.Scan(
			&appName,
			&envName,
			&regionCode,
			&version,
			&deployedAt,
			&status,
			&deployedBy,
		); err != nil {
			return nil, err
		}

		deployment := models.FrontendDeployment{
			ApplicationName: appName,
			Environment:     envName,
			Region:          regionCode,
			Version:         version,
			Timestamp:       deployedAt.Format(time.RFC3339),
			Status:          status,
			DeployedBy:      deployedBy,
		}

		deployments = append(deployments, deployment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deployments, nil
}

// GetDeploymentHistory retrieves the deployment history for an application in a specific environment and region
func (r *ReleaseRepository) GetDeploymentHistory(appID, envID, regionID uint) ([]*models.Deployment, error) {
	var deployments []*models.Deployment

	// Use raw SQL query to avoid GORM's automatic addition of deleted_at IS NULL
	rows, err := r.db.Raw(`
		SELECT * FROM deployments
		WHERE application_id = ? AND environment_id = ? AND region_id = ?
		ORDER BY deployed_at DESC`,
		appID, envID, regionID).Rows()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var deployment models.Deployment
		if err := r.db.ScanRows(rows, &deployment); err != nil {
			return nil, err
		}
		deployments = append(deployments, &deployment)
	}

	return deployments, nil
}

// Query executes a raw SQL query and returns the result
func (r *ReleaseRepository) Query(query string, args ...interface{}) (interface{}, error) {
	// Use the underlying GORM DB to execute the raw SQL query
	rows, err := r.db.Raw(query, args...).Rows()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Close closes the database connection
func (r *ReleaseRepository) Close() error {
	// GORM v2 doesn't have a Close method, but we need this for compatibility
	// Get the underlying SQL DB
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}

	// Close the SQL DB
	return sqlDB.Close()
}
