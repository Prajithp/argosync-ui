package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prajithp/deploy-heirloom/backend/database"
	"github.com/prajithp/deploy-heirloom/backend/models"
)

// Handler represents the API handler
type Handler struct {
	DB *database.DB
}

// NewHandler creates a new API handler
func NewHandler(db *database.DB) *Handler {
	return &Handler{DB: db}
}

// Release handles the release endpoint
func (h *Handler) Release(c echo.Context) error {
	// Parse request
	req := new(models.ReleaseRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	// Validate request
	if req.Application == "" || req.Environment == "" || req.Region == "" || req.Version == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields")
	}

	// Set default deployed_by if not provided
	if req.DeployedBy == "" {
		req.DeployedBy = "system"
	}

	// Process the release
	deployment, err := h.DB.Release(req)
	if err != nil {
		// Check if the error is because the version already exists
		if errors.Is(err, database.ErrVersionExists) {
			return echo.NewHTTPError(http.StatusConflict, "Version already exists for this application, environment, and region")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, deployment)
}

// Rollback handles the rollback endpoint
func (h *Handler) Rollback(c echo.Context) error {
	// Parse request
	req := new(models.RollbackRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	// Validate request
	if req.Application == "" || req.Environment == "" || req.Region == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields")
	}

	// Set default deployed_by if not provided
	if req.DeployedBy == "" {
		req.DeployedBy = "system"
	}

	// Process the rollback
	deployment, err := h.DB.Rollback(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, deployment)
}

// GetActiveDeployments handles the get active deployments endpoint
func (h *Handler) GetActiveDeployments(c echo.Context) error {
	// Get application name from query parameter
	appName := c.QueryParam("application")

	// Get application
	app, err := h.DB.GetApplicationByName(appName)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Application not found")
	}

	// Get active deployments for the application
	rows, err := h.DB.Query(`
		SELECT
			a.name as application_name,
			e.name as environment,
			r.code as region_code,
			r.name as region_name,
			d.version,
			d.deployed_at,
			d.deployed_by
		FROM deployments d
		JOIN applications a ON d.application_id = a.id
		JOIN environments e ON d.environment_id = e.id
		JOIN regions r ON d.region_id = r.id
		WHERE d.application_id = $1 AND d.status = 'active'
		ORDER BY r.code
	`, app.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var deployments []*models.ActiveDeployment
	for rows.Next() {
		var deployment models.ActiveDeployment
		err := rows.Scan(
			&deployment.ApplicationName,
			&deployment.Environment,
			&deployment.RegionCode,
			&deployment.RegionName,
			&deployment.Version,
			&deployment.DeployedAt,
			&deployment.DeployedBy,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		deployments = append(deployments, &deployment)
	}
	if err := rows.Err(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, deployments)
}

// GetAllDeployments handles the get all deployments endpoint for the frontend
func (h *Handler) GetAllDeployments(c echo.Context) error {
	// Get limit parameter from query, default to 10 if not provided
	limitStr := c.QueryParam("limit")
	limit := 10 // Default limit
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	// Query all deployments with their related information
	// Use window functions to limit the number of versions per application/environment/region
	rows, err := h.DB.Query(`
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
		WHERE row_num <= $1
		ORDER BY application_name, environment, region_code, deployed_at DESC
	`, limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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

		err := rows.Scan(
			&appName,
			&envName,
			&regionCode,
			&version,
			&deployedAt,
			&status,
			&deployedBy,
		)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, deployments)
}

// GetDeploymentHistory handles the get deployment history endpoint
func (h *Handler) GetDeploymentHistory(c echo.Context) error {
	// Get parameters from query
	appName := c.QueryParam("application")
	envName := c.QueryParam("environment")
	regionCode := c.QueryParam("region")

	// Validate parameters
	if appName == "" || envName == "" || regionCode == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required query parameters")
	}

	// Get application
	app, err := h.DB.GetApplicationByName(appName)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Application not found")
	}

	// Get environment
	env, err := h.DB.GetEnvironmentByName(envName)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Environment not found")
	}

	// Get region
	region, err := h.DB.GetRegionByCode(regionCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Region not found")
	}

	// Get deployment history
	deployments, err := h.DB.GetDeploymentHistory(app.ID, env.ID, region.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, deployments)
}
