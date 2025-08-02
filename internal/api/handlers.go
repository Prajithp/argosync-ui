package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/Prajithp/argosync/internal/repository"
	"github.com/Prajithp/argosync/internal/service"
	"github.com/Prajithp/argosync/pkg/models"
)

// Handler represents the API handler
type Handler struct {
	Service *service.DeploymentService
}

// NewHandler creates a new API handler
func NewHandler(svc *service.DeploymentService) *Handler {
	return &Handler{Service: svc}
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

	// Process the release
	deployment, err := h.Service.Release(req)
	if err != nil {
		// Check if the error is because the version already exists
		if errors.Is(err, repository.ErrVersionExists) {
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

	// Process the rollback
	deployment, err := h.Service.Rollback(req)
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
	app, err := h.Service.GetApplicationByName(appName)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Application not found")
	}

	// Get environment
	envName := c.QueryParam("environment")
	env, err := h.Service.GetEnvironmentByName(envName)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Environment not found")
	}

	// Get region
	regionCode := c.QueryParam("region")
	region, err := h.Service.GetRegionByCode(regionCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Region not found")
	}

	// Get deployment history
	deployments, err := h.Service.GetDeploymentHistory(app.ID, env.ID, region.ID)
	if err != nil {
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

	// Get all deployments
	deployments, err := h.Service.GetAllDeployments(limit)
	if err != nil {
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
	app, err := h.Service.GetApplicationByName(appName)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Application not found")
	}

	// Get environment
	env, err := h.Service.GetEnvironmentByName(envName)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Environment not found")
	}

	// Get region
	region, err := h.Service.GetRegionByCode(regionCode)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Region not found")
	}

	// Get deployment history
	deployments, err := h.Service.GetDeploymentHistory(app.ID, env.ID, region.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, deployments)
}
