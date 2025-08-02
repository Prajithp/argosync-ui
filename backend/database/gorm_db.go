package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/prajithp/deploy-heirloom/backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB represents the database connection
type DB struct {
	*gorm.DB
	maxVersions int
}

// Config represents the database configuration
type Config struct {
	Type        DBType
	Host        string
	Port        int
	User        string
	Password    string
	DBName      string
	SSLMode     string
	FilePath    string
	MaxVersions int // Maximum number of versions to keep per application/environment/region
}

// DBType represents the type of database
type DBType string

const (
	// PostgreSQL database
	PostgreSQL DBType = "postgres"
	// SQLite database
	SQLite DBType = "sqlite3"
)

// ErrVersionExists is returned when trying to release a version that already exists
var ErrVersionExists = errors.New("version already exists")

// New creates a new database connection
func New(config Config) (*DB, error) {
	var db *gorm.DB
	var err error

	switch config.Type {
	case PostgreSQL:
		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	case SQLite:
		db, err = gorm.Open(sqlite.Open(config.FilePath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}

	if err != nil {
		return nil, err
	}

	return &DB{db, config.MaxVersions}, nil
}

// MigrateSchema initializes the database schema using auto migration
func (db *DB) MigrateSchema() error {
	// Auto migrate the schema
	err := db.AutoMigrate(
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
func (db *DB) GetApplicationByName(name string) (*models.Application, error) {
	var app models.Application

	// Use raw SQL query to avoid GORM's automatic addition of deleted_at IS NULL
	if err := db.Raw("SELECT * FROM applications WHERE name = ? LIMIT 1", name).Scan(&app).Error; err != nil {
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
func (db *DB) CreateApplicationIfNotExists(name string) (*models.Application, error) {
	// Try to get the application first
	app, err := db.GetApplicationByName(name)
	if err == nil {
		// Application already exists
		return app, nil
	}

	// Create the application
	newApp := &models.Application{
		Name: name,
	}

	if err := db.Create(newApp).Error; err != nil {
		return nil, fmt.Errorf("error creating application: %w", err)
	}

	return newApp, nil
}

// GetEnvironmentByName retrieves an environment by name
func (db *DB) GetEnvironmentByName(name string) (*models.Environment, error) {
	var env models.Environment

	// Use raw SQL query to avoid GORM's automatic addition of deleted_at IS NULL
	if err := db.Raw("SELECT * FROM environments WHERE name = ? LIMIT 1", name).Scan(&env).Error; err != nil {
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
func (db *DB) CreateEnvironmentIfNotExists(name string) (*models.Environment, error) {
	// Try to get the environment first
	env, err := db.GetEnvironmentByName(name)
	if err == nil {
		// Environment already exists
		return env, nil
	}

	// Create the environment
	newEnv := &models.Environment{
		Name: name,
	}

	if err := db.Create(newEnv).Error; err != nil {
		return nil, fmt.Errorf("error creating environment: %w", err)
	}

	return newEnv, nil
}

// GetRegionByCode retrieves a region by code
func (db *DB) GetRegionByCode(code string) (*models.Region, error) {
	var region models.Region

	// Use raw SQL query to avoid GORM's automatic addition of deleted_at IS NULL
	if err := db.Raw("SELECT * FROM regions WHERE code = ? LIMIT 1", code).Scan(&region).Error; err != nil {
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
func (db *DB) CreateRegionIfNotExists(code string, name string) (*models.Region, error) {
	// Try to get the region first
	region, err := db.GetRegionByCode(code)
	if err == nil {
		// Region already exists
		return region, nil
	}

	// Create the region
	newRegion := &models.Region{
		Code: code,
		Name: name,
	}

	if err := db.Create(newRegion).Error; err != nil {
		return nil, fmt.Errorf("error creating region: %w", err)
	}

	return newRegion, nil
}

// CheckVersionExists checks if a version already exists for an application in a specific environment and region
func (db *DB) CheckVersionExists(appID, envID, regionID uint, version string) (bool, error) {
	var count int64

	// Use raw SQL query to avoid GORM's automatic addition of deleted_at IS NULL
	err := db.Raw(`
		SELECT COUNT(*) FROM deployments
		WHERE application_id = ? AND environment_id = ? AND region_id = ? AND version = ?`,
		appID, envID, regionID, version).Scan(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetActiveDeployment retrieves the active deployment for an application in a specific environment and region
func (db *DB) GetActiveDeployment(appID, envID, regionID uint) (*models.Deployment, error) {
	var deployment models.Deployment

	// Use raw SQL query to avoid GORM's automatic addition of deleted_at IS NULL
	err := db.Raw(`
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
func (db *DB) CleanupOldVersions(appID, envID, regionID uint, maxVersions int) error {
	// Get all deployments for the application/environment/region using raw SQL
	var deployments []models.Deployment
	if err := db.Raw(`
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

		if err := db.Exec(query).Error; err != nil {
			return fmt.Errorf("error deleting deployments: %w", err)
		}
	}

	return nil
}

// Release creates a new deployment and updates the status of the previous active deployment
func (db *DB) Release(req *models.ReleaseRequest) (*models.Deployment, error) {
	// Start a transaction
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create application if it doesn't exist
	app, err := db.CreateApplicationIfNotExists(req.Application)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating application: %w", err)
	}

	// Create environment if it doesn't exist
	// Default priority is 1
	env, err := db.CreateEnvironmentIfNotExists(req.Environment)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating environment: %w", err)
	}

	// Create region if it doesn't exist
	// Default name is the same as the code
	region, err := db.CreateRegionIfNotExists(req.Region, req.Region)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating region: %w", err)
	}

	// Check if the version already exists
	exists, err := db.CheckVersionExists(app.ID, env.ID, region.ID, req.Version)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error checking if version exists: %w", err)
	}

	if exists {
		tx.Rollback()
		return nil, ErrVersionExists
	}

	// Get current active deployment
	currentDeployment, err := db.GetActiveDeployment(app.ID, env.ID, region.ID)
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
	if db.maxVersions > 0 {
		if err := db.CleanupOldVersions(app.ID, env.ID, region.ID, db.maxVersions); err != nil {
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
func (db *DB) Rollback(req *models.RollbackRequest) (*models.Deployment, error) {
	// Start a transaction
	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create application if it doesn't exist
	app, err := db.CreateApplicationIfNotExists(req.Application)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating application: %w", err)
	}

	// Create environment if it doesn't exist
	// Default priority is 1
	env, err := db.CreateEnvironmentIfNotExists(req.Environment)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating environment: %w", err)
	}

	// Create region if it doesn't exist
	// Default name is the same as the code
	region, err := db.CreateRegionIfNotExists(req.Region, req.Region)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("error creating region: %w", err)
	}

	// Get current active deployment
	currentDeployment, err := db.GetActiveDeployment(app.ID, env.ID, region.ID)
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
	} else if currentDeployment.RollbackTargetID != nil {
		// Use the rollback target from the current deployment
		for i := range deployments {
			if deployments[i].ID == *currentDeployment.RollbackTargetID {
				targetDeployment = &deployments[i]
				break
			}
		}
		if targetDeployment == nil {
			tx.Rollback()
			return nil, errors.New("rollback target not found")
		}
	} else {
		// Use the previous deployment
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
func (db *DB) GetAllDeployments(limit int) ([]models.FrontendDeployment, error) {
	// Use a raw SQL query with window functions to limit the number of versions per application/environment/region
	rows, err := db.Raw(`
		WITH ranked_deployments AS (
			SELECT
				a.name as application_name,
				e.name as environment,
				r.code as region_code,
				d.version,
				d.deployed_at,
				d.status,
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
			status
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
		)

		if err := rows.Scan(
			&appName,
			&envName,
			&regionCode,
			&version,
			&deployedAt,
			&status,
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
		}

		deployments = append(deployments, deployment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return deployments, nil
}

// InitSQLiteTestDB initializes a SQLite database for testing
func InitSQLiteTestDB() (*DB, error) {
	// Create a new SQLite database in memory
	config := Config{
		Type:        SQLite,
		FilePath:    ":memory:",
		MaxVersions: 10, // Default to 10 versions for testing
	}

	db, err := New(config)
	if err != nil {
		return nil, fmt.Errorf("error creating SQLite database: %w", err)
	}

	// Migrate the schema
	if err := db.MigrateSchema(); err != nil {
		return nil, fmt.Errorf("error migrating schema: %w", err)
	}
	return db, nil
}

// Query executes a raw SQL query and returns the result
func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	// Use the underlying GORM DB to execute the raw SQL query
	return db.Raw(query, args...).Rows()
}

// GetDeploymentHistory retrieves the deployment history for an application in a specific environment and region
func (db *DB) GetDeploymentHistory(appID, envID, regionID uint) ([]*models.Deployment, error) {
	var deployments []*models.Deployment

	// Use raw SQL query to avoid GORM's automatic addition of deleted_at IS NULL
	rows, err := db.Raw(`
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
		if err := db.ScanRows(rows, &deployment); err != nil {
			return nil, err
		}
		deployments = append(deployments, &deployment)
	}

	return deployments, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	// GORM v2 doesn't have a Close method, but we need this for compatibility
	// Get the underlying SQL DB
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}

	// Close the SQL DB
	return sqlDB.Close()
}
