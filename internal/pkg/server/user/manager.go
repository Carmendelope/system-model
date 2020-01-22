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

package user

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/nalej/system-model/internal/pkg/provider/user"
	"github.com/rs/zerolog/log"
)

// Manager structure with the required providers for user operations.
type Manager struct {
	OrgProvider  organization.Provider
	UserProvider user.Provider
}

// NewManager creates a Manager using a set of providers.
func NewManager(orgProvider organization.Provider, userProvider user.Provider) Manager {
	return Manager{orgProvider, userProvider}
}

// AddUser adds a new user to a given organization.
func (m *Manager) AddUser(addUserRequest *grpc_user_go.AddUserRequest) (*entities.User, derrors.Error) {
	exists, err := m.OrgProvider.Exists(addUserRequest.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("not found organizationID").WithParams(addUserRequest.OrganizationId)
	}
	toAdd := entities.NewUserFromGRPC(addUserRequest)
	err = m.UserProvider.Add(*toAdd)
	if err != nil {
		return nil, err
	}
	err = m.OrgProvider.AddUser(toAdd.OrganizationId, toAdd.Email)
	if err != nil {
		return nil, err
	}
	return toAdd, nil
}

// AddUser adds a new user to a given organization.
func (m *Manager) UpdateUser(request *grpc_user_go.UpdateUserRequest) derrors.Error {
	exists, err := m.OrgProvider.Exists(request.OrganizationId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("not found organizationID").WithParams(request.OrganizationId)
	}

	usr, err := m.UserProvider.Get(request.Email)
	if err != nil {
		return err
	}

	usr.ApplyUpdate(request)

	err = m.UserProvider.Update(*usr)
	if err != nil {
		return err
	}
	return nil
}

// GetUser returns an existing user.
func (m *Manager) GetUser(userID *grpc_user_go.UserId) (*entities.User, derrors.Error) {
	exists, err := m.OrgProvider.Exists(userID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(userID.OrganizationId)
	}

	exists, err = m.OrgProvider.UserExists(userID.OrganizationId, userID.Email)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("userID").WithParams(userID.OrganizationId, userID.Email)
	}
	return m.UserProvider.Get(userID.Email)
}

// GetUsers retrieves the list of users of a given organization.
func (m *Manager) GetUsers(organizationID *grpc_organization_go.OrganizationId) ([]entities.User, derrors.Error) {
	exists, err := m.OrgProvider.Exists(organizationID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(organizationID.OrganizationId)
	}
	users, err := m.OrgProvider.ListUsers(organizationID.OrganizationId)
	if err != nil {
		return nil, err
	}
	result := make([]entities.User, 0)
	for _, email := range users {
		toAdd, err := m.UserProvider.Get(email)
		if err != nil {
			return nil, err
		}
		result = append(result, *toAdd)
		log.Debug().Interface("result", result).Msg("get users result")
	}
	return result, nil
}

// RemoveUser removes a given user from an organization.
func (m *Manager) RemoveUser(removeRequest *grpc_user_go.RemoveUserRequest) derrors.Error {
	exists, err := m.OrgProvider.Exists(removeRequest.OrganizationId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organizationID").WithParams(removeRequest.OrganizationId)
	}

	exists, err = m.OrgProvider.UserExists(removeRequest.OrganizationId, removeRequest.Email)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("userID").WithParams(removeRequest.OrganizationId, removeRequest.Email)
	}

	err = m.OrgProvider.DeleteUser(removeRequest.OrganizationId, removeRequest.Email)
	if err != nil {
		return err
	}

	err = m.UserProvider.Remove(removeRequest.Email)
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("Error removing user. Rollback!")
		rollbackError := m.OrgProvider.AddUser(removeRequest.OrganizationId, removeRequest.Email)
		if rollbackError != nil {
			log.Error().Str("trace", conversions.ToDerror(rollbackError).DebugReport()).Msg("error in Rollback")
		}
	}
	return err

}
