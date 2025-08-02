package service

import (
	"github.com/Prajithp/argosync/internal/repository"
	"github.com/Prajithp/argosync/pkg/models"
)

// DeploymentService handles business logic for deployments
type DeploymentService struct {
	repo repository.Repository
}

// NewDeploymentService creates a new DeploymentService
func NewDeploymentService(repo repository.Repository) *DeploymentService {
	return &DeploymentService{
		repo: repo,
	}
}

// MigrateSchema initializes the database schema
func (s *DeploymentService) MigrateSchema() error {
	return s.repo.MigrateSchema()
}

// Release creates a new deployment
func (s *DeploymentService) Release(req *models.ReleaseRequest) (*models.Deployment, error) {
	// Set default deployed_by if not provided
	if req.DeployedBy == "" {
		req.DeployedBy = "system"
	}

	return s.repo.Release(req)
}

// Rollback rolls back to a previous deployment
func (s *DeploymentService) Rollback(req *models.RollbackRequest) (*models.Deployment, error) {
	// Set default deployed_by if not provided
	if req.DeployedBy == "" {
		req.DeployedBy = "system"
	}

	return s.repo.Rollback(req)
}

// GetApplicationByName retrieves an application by name
func (s *DeploymentService) GetApplicationByName(name string) (*models.Application, error) {
	return s.repo.GetApplicationByName(name)
}

// GetEnvironmentByName retrieves an environment by name
func (s *DeploymentService) GetEnvironmentByName(name string) (*models.Environment, error) {
	return s.repo.GetEnvironmentByName(name)
}

// GetRegionByCode retrieves a region by code
func (s *DeploymentService) GetRegionByCode(code string) (*models.Region, error) {
	return s.repo.GetRegionByCode(code)
}

// GetAllDeployments returns all deployments in a format compatible with the frontend
func (s *DeploymentService) GetAllDeployments(limit int) ([]models.FrontendDeployment, error) {
	return s.repo.GetAllDeployments(limit)
}

// GetDeploymentHistory retrieves the deployment history for an application in a specific environment and region
func (s *DeploymentService) GetDeploymentHistory(appID, envID, regionID uint) ([]*models.Deployment, error) {
	return s.repo.GetDeploymentHistory(appID, envID, regionID)
}

// Close closes the database connection
func (s *DeploymentService) Close() error {
	return s.repo.Close()
}
