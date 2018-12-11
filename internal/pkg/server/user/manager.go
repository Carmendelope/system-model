/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
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
	OrgProvider organization.Provider
	UserProvider user.Provider
}

// NewManager creates a Manager using a set of providers.
func NewManager(orgProvider organization.Provider, userProvider user.Provider) Manager{
	return Manager{orgProvider, userProvider}
}

// AddUser adds a new user to a given organization.
func (m * Manager) AddUser(addUserRequest *grpc_user_go.AddUserRequest) (*entities.User, derrors.Error){
	exists, err := m.OrgProvider.Exists(addUserRequest.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists{
		return nil, derrors.NewNotFoundError("organizationID").WithParams(addUserRequest.OrganizationId)
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

// GetUser returns an existing user.
func (m * Manager) GetUser(userID *grpc_user_go.UserId) (*entities.User, derrors.Error){
	exists, err := m.OrgProvider.Exists(userID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if ! exists{
		return nil, derrors.NewNotFoundError("organizationID").WithParams(userID.OrganizationId)
	}

	exists, err = m.OrgProvider.UserExists(userID.OrganizationId, userID.Email)
	if err != nil {
		return nil, err
	}
	if !exists{
		return nil, derrors.NewNotFoundError("userID").WithParams(userID.OrganizationId, userID.Email)
	}
	return m.UserProvider.Get(userID.Email)
}

// GetUsers retrieves the list of users of a given organization.
func (m * Manager) GetUsers(organizationID *grpc_organization_go.OrganizationId) ([]entities.User, derrors.Error){
	exists, err := m.OrgProvider.Exists(organizationID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists{
		return nil, derrors.NewNotFoundError("organizationID").WithParams(organizationID.OrganizationId)
	}
	users, err := m.OrgProvider.ListUsers(organizationID.OrganizationId)
	if err != nil {
		return nil, err
	}
	result := make([] entities.User, 0)
	for _, email := range users {
		toAdd, err := m.UserProvider.Get(email)
		if err != nil {
			return nil, err
		}
		result = append(result, *toAdd)
	}
	return result, nil
}

// RemoveUser removes a given user from an organization.
func (m * Manager) RemoveUser(removeRequest *grpc_user_go.RemoveUserRequest) derrors.Error {
	exists, err := m.OrgProvider.Exists(removeRequest.OrganizationId)
	if err != nil {
		return err
	}
	if !exists{
		return derrors.NewNotFoundError("organizationID").WithParams(removeRequest.OrganizationId)
	}

	exists, err = m.OrgProvider.UserExists(removeRequest.OrganizationId, removeRequest.Email)
	if err != nil {
		return err
	}
	if !exists{
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

