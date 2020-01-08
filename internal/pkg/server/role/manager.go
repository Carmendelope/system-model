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
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-role-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/nalej/system-model/internal/pkg/provider/role"
	"github.com/rs/zerolog/log"
)

// Manager structure with the required providers for role operations.
type Manager struct {
	OrgProvider  organization.Provider
	RoleProvider role.Provider
}

// NewManager creates a Manager using a set of providers.
func NewManager(orgProvider organization.Provider, roleProvider role.Provider) Manager {
	return Manager{orgProvider, roleProvider}
}

// AddRole adds a new role to a given organization.
func (m *Manager) AddRole(addRoleRequest *grpc_role_go.AddRoleRequest) (*entities.Role, derrors.Error) {
	exists, err := m.OrgProvider.Exists(addRoleRequest.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(addRoleRequest.OrganizationId)
	}
	toAdd := entities.NewRoleFromGRPC(addRoleRequest)
	err = m.RoleProvider.Add(*toAdd)
	if err != nil {
		return nil, err
	}
	err = m.OrgProvider.AddRole(toAdd.OrganizationId, toAdd.RoleId)
	if err != nil {
		return nil, err
	}
	return toAdd, nil
}

// GetRole returns an existing role.
func (m *Manager) GetRole(roleID *grpc_role_go.RoleId) (*entities.Role, derrors.Error) {
	exists, err := m.OrgProvider.Exists(roleID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(roleID.OrganizationId)
	}

	exists, err = m.OrgProvider.RoleExists(roleID.OrganizationId, roleID.RoleId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("roleID").WithParams(roleID.OrganizationId, roleID.RoleId)
	}
	return m.RoleProvider.Get(roleID.RoleId)
}

// ListRoles retrieves the list of roles of a given organization.
func (m *Manager) ListRoles(organizationID *grpc_organization_go.OrganizationId) ([]entities.Role, derrors.Error) {
	exists, err := m.OrgProvider.Exists(organizationID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(organizationID.OrganizationId)
	}
	roles, err := m.OrgProvider.ListRoles(organizationID.OrganizationId)
	if err != nil {
		return nil, err
	}
	result := make([]entities.Role, 0)
	for _, rID := range roles {
		toAdd, err := m.RoleProvider.Get(rID)
		if err != nil {
			return nil, err
		}
		result = append(result, *toAdd)
	}
	return result, nil
}

// RemoveRole removes a given role from an organization.
func (m *Manager) RemoveRole(removeRoleRequest *grpc_role_go.RemoveRoleRequest) derrors.Error {
	exists, err := m.OrgProvider.Exists(removeRoleRequest.OrganizationId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organizationID").WithParams(removeRoleRequest.OrganizationId)
	}

	exists, err = m.OrgProvider.RoleExists(removeRoleRequest.OrganizationId, removeRoleRequest.RoleId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("roleID").WithParams(removeRoleRequest.OrganizationId, removeRoleRequest.RoleId)
	}

	err = m.OrgProvider.DeleteRole(removeRoleRequest.OrganizationId, removeRoleRequest.RoleId)
	if err != nil {
		return err
	}
	err = m.RoleProvider.Remove(removeRoleRequest.RoleId)
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("Error removing role. Rollback!")
		rollbackError := m.OrgProvider.AddRole(removeRoleRequest.OrganizationId, removeRoleRequest.RoleId)
		if rollbackError != nil {
			log.Error().Str("trace", conversions.ToDerror(rollbackError).DebugReport()).
				Str("removeRoleRequest.OrganizationId", removeRoleRequest.OrganizationId).
				Str("removeRoleRequest.RoleId", removeRoleRequest.RoleId).
				Msg("error in Rollback")
		}
	}
	return err
}
