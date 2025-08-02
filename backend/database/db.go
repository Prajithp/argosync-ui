package database

// import (
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"time"

// 	"github.com/prajithp/deploy-heirloom/backend/models"
// 	_ "github.com/lib/pq"
// 	_ "github.com/mattn/go-sqlite3"
// )

// // ErrVersionExists is returned when trying to release a version that already exists
// var ErrVersionExists = errors.New("version already exists")

// // DBType represents the type of database
// type DBType string

// const (
// 	// PostgreSQL database
// 	PostgreSQL DBType = "postgres"
// 	// SQLite database
// 	SQLite DBType = "sqlite3"
// )

// // DB represents the database connection
// type DB struct {
// 	*sql.DB
// 	dbType      DBType
// 	maxVersions int
// }

// // Config represents the database configuration
// type Config struct {
// 	Type        DBType
// 	Host        string
// 	Port        int
// 	User        string
// 	Password    string
// 	DBName      string
// 	SSLMode     string
// 	FilePath    string
// 	MaxVersions int // Maximum number of versions to keep per application/environment/region
// }

// // New creates a new database connection
// func New(config Config) (*DB, error) {
// 	var (
// 		db  *sql.DB
// 		err error
// 	)

// 	switch config.Type {
// 	case PostgreSQL:
// 		connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
// 			config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)
// 		db, err = sql.Open("postgres", connStr)
// 	case SQLite:
// 		db, err = sql.Open("sqlite3", config.FilePath)
// 	default:
// 		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
// 	}

// 	if err != nil {
// 		return nil, err
// 	}

// 	// Test the connection
// 	if err := db.Ping(); err != nil {
// 		return nil, err
// 	}

// 	return &DB{db, config.Type, config.MaxVersions}, nil
// }

// // InitSchema initializes the database schema
// func (db *DB) InitSchema(schemaPath string) error {
// 	schema, err := ioutil.ReadFile(schemaPath)
// 	if err != nil {
// 		return fmt.Errorf("error reading schema file: %w", err)
// 	}

// 	_, err = db.Exec(string(schema))
// 	if err != nil {
// 		return fmt.Errorf("error executing schema: %w", err)
// 	}

// 	return nil
// }

// // LoadSampleData loads sample data into the database
// func (db *DB) LoadSampleData(dataPath string) error {
// 	data, err := ioutil.ReadFile(dataPath)
// 	if err != nil {
// 		return fmt.Errorf("error reading sample data file: %w", err)
// 	}

// 	_, err = db.Exec(string(data))
// 	if err != nil {
// 		return fmt.Errorf("error loading sample data: %w", err)
// 	}

// 	return nil
// }

// // InitSQLiteTestDB initializes a SQLite database for testing
// func InitSQLiteTestDB() (*DB, error) {
// 	// Create a temporary directory for the SQLite database
// 	tempDir, err := ioutil.TempDir("", "heirloom-test")
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating temp directory: %w", err)
// 	}

// 	dbPath := filepath.Join(tempDir, "heirloom-test.db")

// 	// Create a new SQLite database
// 	config := Config{
// 		Type:        SQLite,
// 		FilePath:    dbPath,
// 		MaxVersions: 10, // Default to 10 versions for testing
// 	}

// 	db, err := New(config)
// 	if err != nil {
// 		os.RemoveAll(tempDir)
// 		return nil, fmt.Errorf("error creating SQLite database: %w", err)
// 	}

// 	// Initialize the schema
// 	err = db.InitSchema("database/sqlite_schema.sql")
// 	if err != nil {
// 		db.Close()
// 		os.RemoveAll(tempDir)
// 		return nil, fmt.Errorf("error initializing schema: %w", err)
// 	}

// 	// Load sample data
// 	err = db.LoadSampleData("database/sample_data.sql")
// 	if err != nil {
// 		db.Close()
// 		os.RemoveAll(tempDir)
// 		return nil, fmt.Errorf("error loading sample data: %w", err)
// 	}

// 	return db, nil
// }

// // GetApplicationByName retrieves an application by name
// func (db *DB) GetApplicationByName(name string) (*models.Application, error) {
// 	var app models.Application
// 	err := db.QueryRow("SELECT id, name, description, team, repository_url, created_at, updated_at FROM applications WHERE name = $1", name).
// 		Scan(&app.ID, &app.Name, &app.Description, &app.Team, &app.RepositoryURL, &app.CreatedAt, &app.UpdatedAt)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, errors.New("application not found")
// 		}
// 		return nil, err
// 	}
// 	return &app, nil
// }

// // GetEnvironmentByName retrieves an environment by name
// func (db *DB) GetEnvironmentByName(name string) (*models.Environment, error) {
// 	var env models.Environment
// 	err := db.QueryRow("SELECT id, name, description, priority FROM environments WHERE name = $1", name).
// 		Scan(&env.ID, &env.Name, &env.Description, &env.Priority)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, errors.New("environment not found")
// 		}
// 		return nil, err
// 	}
// 	return &env, nil
// }

// // GetRegionByCode retrieves a region by code
// func (db *DB) GetRegionByCode(code string) (*models.Region, error) {
// 	var region models.Region
// 	err := db.QueryRow("SELECT id, code, name, continent, active FROM regions WHERE code = $1", code).
// 		Scan(&region.ID, &region.Code, &region.Name, &region.Continent, &region.Active)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, errors.New("region not found")
// 		}
// 		return nil, err
// 	}
// 	return &region, nil
// }

// // CheckVersionExists checks if a version already exists for an application in a specific environment and region
// func (db *DB) CheckVersionExists(appID, envID, regionID int, version string) (bool, error) {
// 	var exists bool
// 	err := db.QueryRow(`
// 		SELECT EXISTS(
// 			SELECT 1 FROM deployments
// 			WHERE application_id = $1
// 			AND environment_id = $2
// 			AND region_id = $3
// 			AND version = $4
// 		)`, appID, envID, regionID, version).Scan(&exists)

// 	if err != nil {
// 		return false, err
// 	}

// 	return exists, nil
// }

// // GetActiveDeployment retrieves the active deployment for an application in a specific environment and region
// func (db *DB) GetActiveDeployment(appID, envID, regionID int) (*models.Deployment, error) {
// 	var deployment models.Deployment
// 	err := db.QueryRow(`
// 		SELECT id, application_id, environment_id, region_id, version, status, deployed_by, deployed_at, rollback_target_id, metadata
// 		FROM deployments
// 		WHERE application_id = $1 AND environment_id = $2 AND region_id = $3 AND status = 'active'`,
// 		appID, envID, regionID).
// 		Scan(&deployment.ID, &deployment.ApplicationID, &deployment.EnvironmentID, &deployment.RegionID,
// 			&deployment.Version, &deployment.Status, &deployment.DeployedBy, &deployment.DeployedAt,
// 			&deployment.RollbackTargetID, &deployment.Metadata)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, nil // No active deployment found
// 		}
// 		return nil, err
// 	}
// 	return &deployment, nil
// }

// // CreateDeployment creates a new deployment
// func (db *DB) CreateDeployment(deployment *models.Deployment) (int, error) {
// 	var id int
// 	err := db.QueryRow(`
// 		INSERT INTO deployments (application_id, environment_id, region_id, version, status, deployed_by, deployed_at, rollback_target_id, metadata)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
// 		RETURNING id`,
// 		deployment.ApplicationID, deployment.EnvironmentID, deployment.RegionID, deployment.Version,
// 		deployment.Status, deployment.DeployedBy, deployment.DeployedAt, deployment.RollbackTargetID, deployment.Metadata).
// 		Scan(&id)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return id, nil
// }

// // UpdateDeploymentStatus updates the status of a deployment
// func (db *DB) UpdateDeploymentStatus(id int, status string) error {
// 	_, err := db.Exec("UPDATE deployments SET status = $1 WHERE id = $2", status, id)
// 	return err
// }

// // CreateDeploymentHistory creates a new deployment history record
// func (db *DB) CreateDeploymentHistory(history *models.DeploymentHistory) error {
// 	_, err := db.Exec(`
// 		INSERT INTO deployment_history (deployment_id, action, performed_by, performed_at, details)
// 		VALUES ($1, $2, $3, $4, $5)`,
// 		history.DeploymentID, history.Action, history.PerformedBy, history.PerformedAt, history.Details)
// 	return err
// }

// // CleanupOldVersions removes old versions exceeding the maximum limit
// func (db *DB) CleanupOldVersions(appID, envID, regionID int, maxVersions int) error {
// 	// Get all deployments for the application/environment/region
// 	deployments, err := db.GetDeploymentHistory(appID, envID, regionID)
// 	if err != nil {
// 		return fmt.Errorf("error getting deployment history: %w", err)
// 	}

// 	// If we have fewer deployments than the maximum, no cleanup needed
// 	if len(deployments) <= maxVersions {
// 		return nil
// 	}

// 	// Sort deployments by deployed_at in descending order (newest first)
// 	// This is already done in GetDeploymentHistory

// 	// Keep track of active deployments (we don't want to delete these)
// 	activeDeploymentIDs := make(map[int]bool)
// 	for _, d := range deployments {
// 		if d.Status == "active" {
// 			activeDeploymentIDs[d.ID] = true
// 		}
// 	}

// 	// Identify deployments to delete (older than maxVersions, excluding active ones)
// 	deploymentsToDelete := make([]int, 0)
// 	keptCount := 0

// 	for _, d := range deployments {
// 		if activeDeploymentIDs[d.ID] {
// 			// Always keep active deployments
// 			continue
// 		}

// 		keptCount++
// 		if keptCount > maxVersions {
// 			deploymentsToDelete = append(deploymentsToDelete, d.ID)
// 		}
// 	}

// 	// Delete old deployments
// 	for _, id := range deploymentsToDelete {
// 		// First delete related deployment history records
// 		_, err := db.Exec("DELETE FROM deployment_history WHERE deployment_id = $1", id)
// 		if err != nil {
// 			return fmt.Errorf("error deleting deployment history: %w", err)
// 		}

// 		// Then delete the deployment
// 		_, err = db.Exec("DELETE FROM deployments WHERE id = $1", id)
// 		if err != nil {
// 			return fmt.Errorf("error deleting deployment: %w", err)
// 		}
// 	}

// 	return nil
// }

// // GetDeploymentHistory retrieves the deployment history for an application in a specific environment and region
// func (db *DB) GetDeploymentHistory(appID, envID, regionID int) ([]*models.Deployment, error) {
// 	rows, err := db.Query(`
// 		SELECT id, application_id, environment_id, region_id, version, status, deployed_by, deployed_at, rollback_target_id, metadata
// 		FROM deployments
// 		WHERE application_id = $1 AND environment_id = $2 AND region_id = $3
// 		ORDER BY deployed_at DESC`,
// 		appID, envID, regionID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var deployments []*models.Deployment
// 	for rows.Next() {
// 		var deployment models.Deployment
// 		err := rows.Scan(&deployment.ID, &deployment.ApplicationID, &deployment.EnvironmentID, &deployment.RegionID,
// 			&deployment.Version, &deployment.Status, &deployment.DeployedBy, &deployment.DeployedAt,
// 			&deployment.RollbackTargetID, &deployment.Metadata)
// 		if err != nil {
// 			return nil, err
// 		}
// 		deployments = append(deployments, &deployment)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
// 	return deployments, nil
// }

// // Release creates a new deployment and updates the status of the previous active deployment
// func (db *DB) Release(req *models.ReleaseRequest) (*models.Deployment, error) {
// 	// Start a transaction
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	// Get application
// 	app, err := db.GetApplicationByName(req.Application)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting application: %w", err)
// 	}

// 	// Get environment
// 	env, err := db.GetEnvironmentByName(req.Environment)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting environment: %w", err)
// 	}

// 	// Get region
// 	region, err := db.GetRegionByCode(req.Region)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting region: %w", err)
// 	}

// 	// Check if the version already exists
// 	exists, err := db.CheckVersionExists(app.ID, env.ID, region.ID, req.Version)
// 	if err != nil {
// 		return nil, fmt.Errorf("error checking if version exists: %w", err)
// 	}

// 	if exists {
// 		return nil, ErrVersionExists
// 	}

// 	// Get current active deployment
// 	currentDeployment, err := db.GetActiveDeployment(app.ID, env.ID, region.ID)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting active deployment: %w", err)
// 	}

// 	// Create new deployment
// 	newDeployment := &models.Deployment{
// 		ApplicationID: app.ID,
// 		EnvironmentID: env.ID,
// 		RegionID:      region.ID,
// 		Version:       req.Version,
// 		Status:        "active",
// 		DeployedBy:    req.DeployedBy,
// 		DeployedAt:    time.Now(),
// 		Metadata:      models.JSONField([]byte(`{}`)),
// 	}

// 	// If there's a current active deployment, set it as the rollback target
// 	if currentDeployment != nil {
// 		// Update the status of the current active deployment
// 		if err := db.UpdateDeploymentStatus(currentDeployment.ID, "inactive"); err != nil {
// 			return nil, fmt.Errorf("error updating deployment status: %w", err)
// 		}

// 		// Set the rollback target
// 		newDeployment.RollbackTargetID = sql.NullInt64{Int64: int64(currentDeployment.ID), Valid: true}
// 	}

// 	// Create the new deployment
// 	id, err := db.CreateDeployment(newDeployment)
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating deployment: %w", err)
// 	}
// 	newDeployment.ID = id

// 	// Clean up old versions if maxVersions is set
// 	if db.maxVersions > 0 {
// 		err = db.CleanupOldVersions(app.ID, env.ID, region.ID, db.maxVersions)
// 		if err != nil {
// 			// Log the error but don't fail the deployment
// 			fmt.Printf("Warning: Failed to clean up old versions: %v\n", err)
// 		}
// 	}

// 	// Create deployment history
// 	history := &models.DeploymentHistory{
// 		DeploymentID: id,
// 		Action:       "deploy",
// 		PerformedBy:  req.DeployedBy,
// 		PerformedAt:  time.Now(),
// 		Details:      models.JSONField([]byte(`{}`)),
// 	}
// 	if err := db.CreateDeploymentHistory(history); err != nil {
// 		return nil, fmt.Errorf("error creating deployment history: %w", err)
// 	}

// 	// Commit the transaction
// 	if err := tx.Commit(); err != nil {
// 		return nil, err
// 	}

// 	return newDeployment, nil
// }

// // Rollback rolls back to a previous deployment
// func (db *DB) Rollback(req *models.RollbackRequest) (*models.Deployment, error) {
// 	// Start a transaction
// 	tx, err := db.Begin()
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer func() {
// 		if err != nil {
// 			tx.Rollback()
// 		}
// 	}()

// 	// Get application
// 	app, err := db.GetApplicationByName(req.Application)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting application: %w", err)
// 	}

// 	// Get environment
// 	env, err := db.GetEnvironmentByName(req.Environment)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting environment: %w", err)
// 	}

// 	// Get region
// 	region, err := db.GetRegionByCode(req.Region)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting region: %w", err)
// 	}

// 	// Check if the version already exists
// 	exists, err := db.CheckVersionExists(app.ID, env.ID, region.ID, req.Version)
// 	if err != nil {
// 		return nil, fmt.Errorf("error checking if version exists: %w", err)
// 	}

// 	if exists {
// 		return nil, ErrVersionExists
// 	}

// 	// Get current active deployment
// 	currentDeployment, err := db.GetActiveDeployment(app.ID, env.ID, region.ID)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting active deployment: %w", err)
// 	}
// 	if currentDeployment == nil {
// 		return nil, errors.New("no active deployment found to rollback")
// 	}

// 	// Get deployment history
// 	deployments, err := db.GetDeploymentHistory(app.ID, env.ID, region.ID)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting deployment history: %w", err)
// 	}
// 	if len(deployments) <= 1 {
// 		return nil, errors.New("no previous deployments found to rollback to")
// 	}

// 	// Find the target deployment to rollback to
// 	var targetDeployment *models.Deployment
// 	if req.Version != "" {
// 		// Find deployment with the specified version
// 		for _, d := range deployments {
// 			if d.Version == req.Version && d.ID != currentDeployment.ID {
// 				targetDeployment = d
// 				break
// 			}
// 		}
// 		if targetDeployment == nil {
// 			return nil, fmt.Errorf("no deployment found with version %s", req.Version)
// 		}
// 	} else if currentDeployment.RollbackTargetID.Valid {
// 		// Use the rollback target from the current deployment
// 		for _, d := range deployments {
// 			if d.ID == int(currentDeployment.RollbackTargetID.Int64) {
// 				targetDeployment = d
// 				break
// 			}
// 		}
// 		if targetDeployment == nil {
// 			return nil, errors.New("rollback target not found")
// 		}
// 	} else {
// 		// Use the previous deployment
// 		for _, d := range deployments {
// 			if d.ID != currentDeployment.ID {
// 				targetDeployment = d
// 				break
// 			}
// 		}
// 	}

// 	// Update the status of the current active deployment
// 	if err := db.UpdateDeploymentStatus(currentDeployment.ID, "inactive"); err != nil {
// 		return nil, fmt.Errorf("error updating deployment status: %w", err)
// 	}

// 	// Update the status of the target deployment
// 	if err := db.UpdateDeploymentStatus(targetDeployment.ID, "active"); err != nil {
// 		return nil, fmt.Errorf("error updating deployment status: %w", err)
// 	}

// 	// Create deployment history
// 	history := &models.DeploymentHistory{
// 		DeploymentID: targetDeployment.ID,
// 		Action:       "rollback",
// 		PerformedBy:  req.DeployedBy,
// 		PerformedAt:  time.Now(),
// 		Details:      models.JSONField([]byte(`{}`)),
// 	}
// 	if err := db.CreateDeploymentHistory(history); err != nil {
// 		return nil, fmt.Errorf("error creating deployment history: %w", err)
// 	}

// 	// Commit the transaction
// 	if err := tx.Commit(); err != nil {
// 		return nil, err
// 	}

// 	return targetDeployment, nil
// }
