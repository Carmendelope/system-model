/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package project

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

type MockupProjectProvider struct {
	// Mutex for managing mockup access.
	sync.Mutex
	// projects with a map of assets indexed by accountID and projectID.
	projects map[string]entities.Project
	// projects by account indexed by accountID
	account_projects map[string][]string
}

func NewMockupProjectProvider() * MockupProjectProvider{
	return &MockupProjectProvider{
		projects: make(map[string]entities.Project, 0),
		account_projects: make (map[string][]string, 0),
	}
}

func (m *MockupProjectProvider)getPK(accountID string, projectID string) string {
	return fmt.Sprintf("%s%s", accountID, projectID)
}

func (m *MockupProjectProvider) unsafeExists(key string) bool{
	_, exists := m.projects[key]
	return exists
}

// Add a new project to the system.
func (m * MockupProjectProvider) Add(project entities.Project) derrors.Error{
	m.Lock()
	defer m.Unlock()

	key := m.getPK(project.OwnerAccountId, project.ProjectId)

	if !m.unsafeExists(key){
		m.projects[key] = project

		// add into account_projects
		account, exists := m.account_projects[project.OwnerAccountId]
		if !exists{
			m.account_projects[project.OwnerAccountId] = []string{project.ProjectId}
		}else{
			m.account_projects[project.OwnerAccountId] = append(account, project.ProjectId)
		}
		return nil
	}
	return derrors.NewAlreadyExistsError("project").WithParams(project.OwnerAccountId, project.ProjectId)
}

// Update the information of a project.
func (m * MockupProjectProvider) Update(project entities.Project) derrors.Error{
	m.Lock()
	defer m.Unlock()

	key := m.getPK(project.OwnerAccountId, project.ProjectId)

	if !m.unsafeExists(key){
		return derrors.NewNotFoundError("project").WithParams(project.OwnerAccountId, project.ProjectId)
	}
	m.projects[key] = project
	return nil
}

// Exists checks if a project exists on the system.
func (m * MockupProjectProvider) Exists(accountID string, projectID string) (bool, derrors.Error){
	m.Lock()
	defer m.Unlock()
	key := m.getPK(accountID, projectID)

	return m.unsafeExists(key), nil
}

// Get a project.
func (m * MockupProjectProvider) Get(accountID string, projectID string) (*entities.Project, derrors.Error){
	m.Lock()
	defer m.Unlock()

	key := m.getPK(accountID, projectID)

	project, exists := m.projects[key]
	if exists {
		return &project, nil
	}
	return nil, derrors.NewNotFoundError("project").WithParams(project.OwnerAccountId, project.ProjectId)
}

// Remove a project
func (m * MockupProjectProvider) Remove(accountID string, projectID string) derrors.Error{
	m.Lock()
	defer m.Unlock()

	key := m.getPK(accountID, projectID)

	if !m.unsafeExists(key){
		return derrors.NewNotFoundError("project").WithParams(accountID, projectID)
	}
	delete(m.projects, key)

	// delete the project from m.account_projects map
	account, exists := m.account_projects[accountID]
	if exists{
		newAccounts := make ([]string, 0)
		for _, projectId := range account{
			if projectId != projectID{
				newAccounts = append(newAccounts, projectId)
			}
		}
		m.account_projects[accountID] = newAccounts
	}

	return nil
}

func (m * MockupProjectProvider) ListAccountProjects(accountID string) ([]entities.Project, derrors.Error) {

	res := make ([]entities.Project, 0)
	accounts, exists := m.account_projects[accountID]
	if exists {
		for _, projectID := range accounts {
			project, err := m.Get(accountID, projectID)
			if err != nil {
				return nil, err
			}
			res = append(res, *project)
		}
	}
	return res, nil
}

// Clear all projects
func (m * MockupProjectProvider) Clear() derrors.Error{
	m.Lock()
	defer m.Unlock()
	m.projects = make(map[string]entities.Project, 0)
	m.account_projects = make (map[string][]string, 0)

	return nil
}