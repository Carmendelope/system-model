/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package project

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-account-go"
	"github.com/nalej/grpc-project-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/account"
	"github.com/nalej/system-model/internal/pkg/provider/project"
)

type Manager struct {
	AccountProvider account.Provider
	ProjectProvider project.Provider
}

func NewManager(accProvider account.Provider, proProvider project.Provider) Manager {
	return Manager{
		AccountProvider: accProvider,
		ProjectProvider: proProvider,
	}
}

// AddProject adds a new project to a given account
func (m *Manager) AddProject(request *grpc_project_go.AddProjectRequest) (*entities.Project, derrors.Error) {

	// check if the account exists
	exists, err := m.AccountProvider.Exists(request.AccountId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("account").WithParams(request.AccountId)
	}

	// check there is no another project with the same name
	exists, err = m.ProjectProvider.ExistsByName(request.AccountId, request.Name)
	if exists {
		return nil, derrors.NewInvalidArgumentError("A Project with that name already exists").WithParams(request.Name)
	}

	toAdd := entities.NewProjectToGRPC(request)
	err = m.ProjectProvider.Add(*toAdd)
	if err != nil {
		return nil, err
	}

	return toAdd, nil
}

// GetProject retrieves a given project
func (m *Manager) GetProject(project *grpc_project_go.ProjectId) (*entities.Project, derrors.Error) {

	exists, err := m.AccountProvider.Exists(project.AccountId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("account").WithParams(project.AccountId)
	}

	return m.ProjectProvider.Get(project.AccountId, project.ProjectId)

}

// RemoveProject removes a given project
func (m *Manager) RemoveProject(project *grpc_project_go.ProjectId) derrors.Error {
	exists, err := m.AccountProvider.Exists(project.AccountId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("account").WithParams(project.AccountId)
	}
	return m.ProjectProvider.Remove(project.AccountId, project.ProjectId)

}

// ListAccountProjects list the projects of a given account
func (m *Manager) ListAccountProjects(project *grpc_account_go.AccountId) ([]entities.Project, derrors.Error) {
	exists, err := m.AccountProvider.Exists(project.AccountId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("account").WithParams(project.AccountId)
	}

	return m.ProjectProvider.ListAccountProjects(project.AccountId)

}

// UpdateProject updates the project information
func (m *Manager) UpdateProject(request *grpc_project_go.UpdateProjectRequest) derrors.Error {
	exists, err := m.AccountProvider.Exists(request.AccountId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("account").WithParams(request.AccountId)
	}

	oldProject, err := m.ProjectProvider.Get(request.AccountId, request.ProjectId)
	if err != nil {
		return err
	}
	// if the name is been changed, check if the ner one already exists
	if request.UpdateName {
		exists, err := m.ProjectProvider.ExistsByName(request.ProjectId, request.Name)
		if err != nil {
			return err
		}
		if exists {
			return derrors.NewInvalidArgumentError("A Project with that name already exists").WithParams(request.Name)
		}
	}

	oldProject.ApplyUpdate(request)

	return m.ProjectProvider.Update(*oldProject)

}
