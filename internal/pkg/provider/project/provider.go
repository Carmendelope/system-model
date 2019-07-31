/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package project

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for project
type Provider interface {
	// Add a new project to the system.
	Add(project entities.Project) derrors.Error
	// Update the information of a project.
	Update(project entities.Project) derrors.Error
	// Exists checks if a project exists on the system.
	Exists(accountID string, projectID string) (bool, derrors.Error)
	// Get a project.
	Get(accountID string, projectID string) (*entities.Project, derrors.Error)
	// Remove a project
	Remove(accountID string, projectID string) derrors.Error
	// List all the projects of an account
	ListAccountProjects(accountID string) ([]entities.Project, derrors.Error)
	// Clear all projects
	Clear() derrors.Error
}

