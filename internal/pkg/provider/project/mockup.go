/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
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
	// projectsNames map indexed by projectID#name to know if a name of a project already exists
	projectNames map[string]bool
}

func NewMockupProjectProvider() *MockupProjectProvider {
	return &MockupProjectProvider{
		projects:         make(map[string]entities.Project, 0),
		account_projects: make(map[string][]string, 0),
		projectNames:     make(map[string]bool, 0),
	}
}

func (m *MockupProjectProvider) getPK(accountID string, projectID string) string {
	return fmt.Sprintf("%s%s", accountID, projectID)
}
func (m *MockupProjectProvider) getNameKey(accountID string, name string) string {
	return fmt.Sprintf("%s%s", accountID, name)
}

func (m *MockupProjectProvider) unsafeExists(key string) bool {
	_, exists := m.projects[key]
	return exists
}

// Add a new project to the system.
func (m *MockupProjectProvider) Add(project entities.Project) derrors.Error {
	m.Lock()
	defer m.Unlock()

	key := m.getPK(project.OwnerAccountId, project.ProjectId)

	if !m.unsafeExists(key) {
		m.projects[key] = project
		m.projectNames[m.getNameKey(project.OwnerAccountId, project.Name)] = true

		// add into account_projects
		account, exists := m.account_projects[project.OwnerAccountId]
		if !exists {
			m.account_projects[project.OwnerAccountId] = []string{project.ProjectId}
		} else {
			m.account_projects[project.OwnerAccountId] = append(account, project.ProjectId)
		}
		return nil
	}
	return derrors.NewAlreadyExistsError("project").WithParams(project.OwnerAccountId, project.ProjectId)
}

// Update the information of a project.
func (m *MockupProjectProvider) Update(project entities.Project) derrors.Error {
	m.Lock()
	defer m.Unlock()

	key := m.getPK(project.OwnerAccountId, project.ProjectId)

	if !m.unsafeExists(key) {
		return derrors.NewNotFoundError("project").WithParams(project.OwnerAccountId, project.ProjectId)
	}

	// delete the all entry
	delete(m.projectNames, m.getNameKey(m.projects[key].OwnerAccountId, m.projects[key].Name))
	// add the new one
	m.projectNames[m.getNameKey(project.OwnerAccountId, project.Name)] = true
	m.projects[key] = project
	return nil
}

// Exists checks if a project exists on the system.
func (m *MockupProjectProvider) Exists(accountID string, projectID string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	key := m.getPK(accountID, projectID)

	return m.unsafeExists(key), nil
}

// check if there is a project in the account with the received name
func (m *MockupProjectProvider) ExistsByName(accountID string, name string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	key := m.getNameKey(accountID, name)

	_, exists := m.projectNames[key]

	return exists, nil
}

// Get a project.
func (m *MockupProjectProvider) Get(accountID string, projectID string) (*entities.Project, derrors.Error) {
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
func (m *MockupProjectProvider) Remove(accountID string, projectID string) derrors.Error {
	m.Lock()
	defer m.Unlock()

	key := m.getPK(accountID, projectID)

	if !m.unsafeExists(key) {
		return derrors.NewNotFoundError("project").WithParams(accountID, projectID)
	}
	delete(m.projectNames, m.getNameKey(m.projects[key].OwnerAccountId, m.projects[key].Name))

	delete(m.projects, key)

	// delete the project from m.account_projects map
	account, exists := m.account_projects[accountID]
	if exists {
		newAccounts := make([]string, 0)
		for _, projectId := range account {
			if projectId != projectID {
				newAccounts = append(newAccounts, projectId)
			}
		}
		m.account_projects[accountID] = newAccounts
	}

	return nil
}

func (m *MockupProjectProvider) ListAccountProjects(accountID string) ([]entities.Project, derrors.Error) {

	res := make([]entities.Project, 0)
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
func (m *MockupProjectProvider) Clear() derrors.Error {
	m.Lock()
	defer m.Unlock()
	m.projects = make(map[string]entities.Project, 0)
	m.account_projects = make(map[string][]string, 0)
	m.projectNames = make(map[string]bool, 0)
	return nil
}
