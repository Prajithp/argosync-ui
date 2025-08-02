package models

import (
	"time"

	"gorm.io/gorm"
)

// Application represents an application/service
type Application struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null" json:"name"`
	Deployments []Deployment
}

// Environment represents a deployment environment
type Environment struct {
	gorm.Model
	Name string `gorm:"uniqueIndex;not null" json:"name"`
}

// Region represents a deployment region
type Region struct {
	gorm.Model
	Code string `gorm:"uniqueIndex;not null" json:"code"`
	Name string `json:"name"`
}

// Deployment represents a deployment record
type Deployment struct {
	gorm.Model
	ApplicationID    uint        `json:"application_id"`
	Application      Application `gorm:"foreignKey:ApplicationID" json:"-"`
	EnvironmentID    uint        `json:"environment_id"`
	Environment      Environment `gorm:"foreignKey:EnvironmentID" json:"-"`
	RegionID         uint        `json:"region_id"`
	Region           Region      `gorm:"foreignKey:RegionID" json:"-"`
	Version          string      `json:"version"`
	Status           string      `gorm:"default:active" json:"status"`
	DeployedBy       string      `json:"deployed_by"`
	DeployedAt       time.Time   `gorm:"default:CURRENT_TIMESTAMP" json:"deployed_at"`
	RollbackTargetID *uint       `json:"rollback_target_id,omitempty"`
	RollbackTarget   *Deployment `gorm:"foreignKey:RollbackTargetID" json:"-"`
}

// ReleaseRequest represents the request payload for the release endpoint
type ReleaseRequest struct {
	Application string `json:"application" validate:"required"`
	Environment string `json:"environment" validate:"required"`
	Region      string `json:"region" validate:"required"`
	Version     string `json:"version" validate:"required"`
	DeployedBy  string `json:"deployed_by"`
}

// RollbackRequest represents the request payload for the rollback endpoint
type RollbackRequest struct {
	Application string `json:"application" validate:"required"`
	Environment string `json:"environment" validate:"required"`
	Region      string `json:"region" validate:"required"`
	Version     string `json:"version,omitempty"` // Optional, if not provided, will rollback to the previous version
	DeployedBy  string `json:"deployed_by"`
}

// FrontendDeployment represents a deployment in the format expected by the frontend
type FrontendDeployment struct {
	ApplicationName string `json:"applicationName"`
	Environment     string `json:"environment"`
	Region          string `json:"region"`
	Version         string `json:"version"`
	Timestamp       string `json:"timestamp"`
	Status          string `json:"status"`
	DeployedBy      string `json:"deployedBy"`
}

// ActiveDeployment represents an active deployment with additional information
type ActiveDeployment struct {
	ApplicationName string    `json:"applicationName"`
	Environment     string    `json:"environment"`
	RegionCode      string    `json:"regionCode"`
	RegionName      string    `json:"regionName"`
	Version         string    `json:"version"`
	DeployedAt      time.Time `json:"deployedAt"`
	DeployedBy      string    `json:"deployedBy"`
}
