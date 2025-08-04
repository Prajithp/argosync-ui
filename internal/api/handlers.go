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
	// Get page parameter from query, default to 1 if not provided
	pageStr := c.QueryParam("page")
	page := 1 // Default page
	if pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	// Get page size parameter from query, default to 10 if not provided
	pageSizeStr := c.QueryParam("pageSize")
	pageSize := 10 // Default page size
	if pageSizeStr != "" {
		if parsedPageSize, err := strconv.Atoi(pageSizeStr); err == nil && parsedPageSize > 0 {
			pageSize = parsedPageSize
		}
	}

	// Get all deployments with pagination
	deployments, totalCount, err := h.Service.GetAllDeployments(page, pageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Create response with pagination metadata
	response := map[string]interface{}{
		"deployments": deployments,
		"pagination": map[string]interface{}{
			"page":       page,
			"pageSize":   pageSize,
			"totalCount": totalCount,
			"totalPages": (totalCount + pageSize - 1) / pageSize, // Ceiling division
		},
	}

	return c.JSON(http.StatusOK, response)
}

// GetAllApplications handles the get all applications endpoint
func (h *Handler) GetAllApplications(c echo.Context) error {
	// Get all applications
	applications, err := h.Service.GetAllApplications()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, applications)
}

// GetRegionsForApplication handles the get regions for application endpoint
func (h *Handler) GetRegionsForApplication(c echo.Context) error {
	// Get application ID from path parameter
	appIDStr := c.Param("appID")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid application ID")
	}

	// Get regions for application
	regions, err := h.Service.GetRegionsForApplication(uint(appID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, regions)
}

// GetEnvironmentsForApplicationAndRegion handles the get environments for application and region endpoint
func (h *Handler) GetEnvironmentsForApplicationAndRegion(c echo.Context) error {
	// Get application ID from path parameter
	appIDStr := c.Param("appID")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid application ID")
	}

	// Get region ID from path parameter
	regionIDStr := c.Param("regionID")
	regionID, err := strconv.ParseUint(regionIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid region ID")
	}

	// Get environments for application and region
	environments, err := h.Service.GetEnvironmentsForApplicationAndRegion(uint(appID), uint(regionID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, environments)
}

// GetVersionsForApplicationEnvironmentRegion handles the get versions for application, environment, and region endpoint
func (h *Handler) GetVersionsForApplicationEnvironmentRegion(c echo.Context) error {
	// Get application ID from path parameter
	appIDStr := c.Param("appID")
	appID, err := strconv.ParseUint(appIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid application ID")
	}

	// Get environment ID from path parameter
	envIDStr := c.Param("envID")
	envID, err := strconv.ParseUint(envIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid environment ID")
	}

	// Get region ID from path parameter
	regionIDStr := c.Param("regionID")
	regionID, err := strconv.ParseUint(regionIDStr, 10, 32)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid region ID")
	}

	// Get versions for application, environment, and region
	versions, err := h.Service.GetVersionsForApplicationEnvironmentRegion(uint(appID), uint(envID), uint(regionID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, versions)
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
