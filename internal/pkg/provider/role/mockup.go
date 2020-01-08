/*
 * Copyright 2020 Nalej
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
 */

package role

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

type MockupRoleProvider struct {
	sync.Mutex
	// Role indexed by role identifier.
	roles map[string]entities.Role
}

func NewMockupRoleProvider() *MockupRoleProvider {
	return &MockupRoleProvider{
		roles: make(map[string]entities.Role, 0),
	}
}

func (m *MockupRoleProvider) unsafeExists(roleID string) bool {
	_, exists := m.roles[roleID]
	return exists
}

// Add a new role to the system.
func (m *MockupRoleProvider) Add(role entities.Role) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(role.RoleId) {
		m.roles[role.RoleId] = role
		return nil
	}
	return derrors.NewAlreadyExistsError(role.RoleId)
}

// Update an existing role in the system
func (m *MockupRoleProvider) Update(role entities.Role) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(role.RoleId) {
		return derrors.NewNotFoundError(role.RoleId)
	}
	m.roles[role.RoleId] = role
	return nil
}

// Exists checks if a role exists on the system.
func (m *MockupRoleProvider) Exists(roleID string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExists(roleID), nil
}

// Get a role.
func (m *MockupRoleProvider) Get(roleID string) (*entities.Role, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	role, exists := m.roles[roleID]
	if exists {
		return &role, nil
	}
	return nil, derrors.NewNotFoundError(roleID)
}

// Remove a role
func (m *MockupRoleProvider) Remove(roleID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(roleID) {
		return derrors.NewNotFoundError(roleID)
	}
	delete(m.roles, roleID)
	return nil
}

// Clear cleans the contents of the mockup.
func (m *MockupRoleProvider) Clear() derrors.Error {
	m.Lock()
	m.roles = make(map[string]entities.Role, 0)
	m.Unlock()
	return nil
}
