package repository

import "github.com/Prajithp/argosync/pkg/models"

// ApplicationRepository defines methods for application operations
type ApplicationRepository interface {
	BaseRepository

	// Application operations
	GetApplicationByName(name string) (*models.Application, error)
	CreateApplicationIfNotExists(name string) (*models.Application, error)
}

type DeploymentRepository interface {
	BaseRepository

	// Deployment operations
	Release(req *models.ReleaseRequest) (*models.Deployment, error)
	Rollback(req *models.RollbackRequest) (*models.Deployment, error)
	GetActiveDeployment(appID, envID, regionID uint) (*models.Deployment, error)
	CheckVersionExists(appID, envID, regionID uint, version string) (bool, error)
	CleanupOldVersions(appID, envID, regionID uint, maxVersions int) error
	GetAllDeployments(page, pageSize int) ([]models.FrontendDeployment, int, error)
	GetDeploymentHistory(appID, envID, regionID uint) ([]*models.Deployment, error)
	
	// New hierarchical methods
	GetAllApplications() ([]*models.Application, error)
	GetRegionsForApplication(appID uint) ([]*models.Region, error)
	GetEnvironmentsForApplicationAndRegion(appID, regionID uint) ([]*models.Environment, error)
	GetVersionsForApplicationEnvironmentRegion(appID, envID, regionID uint) ([]*models.Deployment, error)
}

type EnvironmentRepository interface {
	BaseRepository

	// Environment operations
	GetEnvironmentByName(name string) (*models.Environment, error)
	CreateEnvironmentIfNotExists(name string) (*models.Environment, error)
}

type RegionRepository interface {
	BaseRepository

	// Region operations
	GetRegionByCode(code string) (*models.Region, error)
	CreateRegionIfNotExists(code string, name string) (*models.Region, error)
}

type Repository interface {
	ApplicationRepository
	EnvironmentRepository
	RegionRepository
	DeploymentRepository

	// Raw query execution
	Query(query string, args ...interface{}) (interface{}, error)
}
