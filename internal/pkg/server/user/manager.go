/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package user

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-account-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/account"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/nalej/system-model/internal/pkg/provider/user"
	"github.com/rs/zerolog/log"
)

// Manager structure with the required providers for user operations.
type Manager struct {
	OrgProvider  organization.Provider
	UserProvider user.Provider
	AccountProvider account.Provider
}

// NewManager creates a Manager using a set of providers.
func NewManager(orgProvider organization.Provider, userProvider user.Provider, accProvider account.Provider) Manager {
	return Manager{
		OrgProvider:orgProvider,
		UserProvider:userProvider,
		AccountProvider:accProvider,
	}
}

// AddUser adds a new user to a given organization.
func (m *Manager) AddUser(addUserRequest *grpc_user_go.AddUserRequest) (*entities.User, derrors.Error) {

	// check if the organization exists only in case there is no empty
	// OrganizationId is deprecated but we need to keep both versions
	if addUserRequest.OrganizationId != "" {
		exists, err := m.OrgProvider.Exists(addUserRequest.OrganizationId)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, derrors.NewNotFoundError("not found organizationID").WithParams(addUserRequest.OrganizationId)
		}
	}

	toAdd := entities.NewUserFromGRPC(addUserRequest)
	err := m.UserProvider.Add(*toAdd)
	if err != nil {
		return nil, err
	}

	if toAdd.OrganizationId != "" {
		err = m.OrgProvider.AddUser(toAdd.OrganizationId, toAdd.Email)
		if err != nil {
			return nil, err
		}
	}
	return toAdd, nil
}

// AddUser adds a new user to a given organization.
func (m *Manager) UpdateUser(request *grpc_user_go.UpdateUserRequest) derrors.Error {
	// check if the organization exists only in case there is no empty
	// OrganizationId is deprecated but we need to keep both versions
	if request.OrganizationId != "" {
		exists, err := m.OrgProvider.Exists(request.OrganizationId)
		if err != nil {
			return err
		}
		if !exists {
			return derrors.NewNotFoundError("not found organizationID").WithParams(request.OrganizationId)
		}
	}
	usr,err := m.UserProvider.Get(request.Email)
	if err != nil {
		return err
	}

	if request.Name != "" {
		usr.Name=request.Name
	}

	if request.PhotoUrl != "" {
		usr.PhotoUrl = request.PhotoUrl
	}

	err = m.UserProvider.Update(*usr)
	if err != nil {
		return err
	}
	return nil
}

// GetUser returns an existing user.
func (m *Manager) GetUser(userID *grpc_user_go.UserId) (*entities.User, derrors.Error) {
	// check if the organization exists only in case there is no empty
	// OrganizationId is deprecated but we need to keep both versions
	if userID.OrganizationId != "" {
		exists, err := m.OrgProvider.Exists(userID.OrganizationId)
		if err != nil {
			return nil, err
		}
		if ! exists {
			return nil, derrors.NewNotFoundError("organizationID").WithParams(userID.OrganizationId)
		}

		exists, err = m.OrgProvider.UserExists(userID.OrganizationId, userID.Email)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, derrors.NewNotFoundError("userID").WithParams(userID.OrganizationId, userID.Email)
		}
	}
	return m.UserProvider.Get(userID.Email)
}

// GetUsers retrieves the list of users of a given organization.
func (m *Manager) GetUsers(organizationID *grpc_organization_go.OrganizationId) ([]entities.User, derrors.Error) {
	// check if the organization exists only in case there is no empty
	// OrganizationId is deprecated but we need to keep both versions
	if organizationID.OrganizationId != "" {
		exists, err := m.OrgProvider.Exists(organizationID.OrganizationId)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, derrors.NewNotFoundError("organizationID").WithParams(organizationID.OrganizationId)
		}
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
func (m *Manager) RemoveUser(removeRequest *grpc_user_go.RemoveUserRequest) derrors.Error {
	// check if the organization exists only in case there is no empty
	// OrganizationId is deprecated but we need to keep both versions
	if removeRequest.OrganizationId != "" {
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
	}

	err := m.UserProvider.Remove(removeRequest.Email)
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("Error removing user. Rollback!")
		if removeRequest.OrganizationId != "" {
			rollbackError := m.OrgProvider.AddUser(removeRequest.OrganizationId, removeRequest.Email)

			if rollbackError != nil {
				log.Error().Str("trace", conversions.ToDerror(rollbackError).DebugReport()).Msg("error in Rollback")
			}
		}
	}
	return err

}
func (m *Manager) UpdateContactInfo(request *grpc_user_go.UpdateContactInfoRequest)  derrors.Error{

	usr,err := m.UserProvider.Get(request.Email)
	if err != nil {
		return err
	}

	// apply update in contact info
	var contactInfo *entities.UserContactInfo
	if usr.ContactInfo != nil {
		contactInfo = usr.ContactInfo
	}
	if request.FullName != "" {
		contactInfo.FullName = request.FullName
	}
	if request.Address != "" {
		contactInfo.Address = request.Address
	}
	if request.AltEmail != "" {
		contactInfo.AltEmail = request.AltEmail
	}
	if request.CompanyName != "" {
		contactInfo.CompanyName = request.CompanyName
	}
	if request.Title != "" {
		contactInfo.Title = request.Title
	}
	usr.ContactInfo = contactInfo

	err = m.UserProvider.Update(*usr)
	if err != nil {
		return err
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------


func (m *Manager) AddAccountUser(request *grpc_user_go.AddAccountUserRequest) (*entities.AccountUser, derrors.Error){

	// check the user exists
	exists, err := m.UserProvider.Exists(request.Email)
	if err != nil {
		return nil, err
	}
	if ! exists {
		return nil, derrors.NewNotFoundError("User").WithParams(request.Email)
	}

	// check the account exists
	exists, err = m.AccountProvider.Exists(request.AccountId)
	if err != nil {
		return nil, err
	}
	if ! exists {
		return nil, derrors.NewNotFoundError("Account").WithParams(request.AccountId)
	}

	toAdd := entities.NewAccountUserFromGRPC(request)

	// add accountUser
	err = m.UserProvider.AddAccountUser(*toAdd)
	if err != nil {
		return nil, err
	}

	return toAdd, nil
}

func (m *Manager) RemoveAccountUser(accountUserID *grpc_user_go.AccountUserId) derrors.Error{
	// check the user exists
	exists, err := m.UserProvider.Exists(accountUserID.Email)
	if err != nil {
		return  err
	}
	if ! exists {
		return  derrors.NewNotFoundError("User").WithParams(accountUserID.Email)
	}

	// check the account exists
	exists, err = m.AccountProvider.Exists(accountUserID.AccountId)
	if err != nil {
		return  err
	}
	if ! exists {
		return  derrors.NewNotFoundError("Account").WithParams(accountUserID.AccountId)
	}

	// remove accountUser
	return  m.UserProvider.RemoveAccountUser(accountUserID.AccountId, accountUserID.Email)

}

func (m *Manager) UpdateAccountUser(request *grpc_user_go.AccountUserUpdateRequest) (*entities.AccountUser, derrors.Error){
	// check the user exists
	exists, err := m.UserProvider.Exists(request.Email)
	if err != nil {
		return  nil, err
	}
	if ! exists {
		return  nil, derrors.NewNotFoundError("User").WithParams(request.Email)
	}

	// check the account exists
	exists, err = m.AccountProvider.Exists(request.AccountId)
	if err != nil {
		return  nil, err
	}
	if ! exists {
		return  nil, derrors.NewNotFoundError("Account").WithParams(request.AccountId)
	}

	// get the accountUser
	old, err := m.UserProvider.GetAccountUser(request.AccountId, request.Email)
	if err != nil {
		return  nil, err
	}

	// applyUpdate
	if request.UpdateStatus{
		old.Status = entities.UserStatusFromGRPC[request.Status]
	}
	if request.UpdateRoleId {
		old.RoleId = request.RoleId
	}

	err = m.UserProvider.UpdateAccountUser(*old)
	if err != nil {
		return nil, err
	}
	// update accountUser
	return old, nil
}
func (m *Manager) ListAccountsUser(in *grpc_account_go.AccountId) (*grpc_user_go.AccountUserList, derrors.Error){
	// check the user exists

	// check the account exists

	return nil, nil
}
// ---------------------------------------------------------------------------------------------------------------------

/*

// ---------------------------------------------------------------------------------------------------------------------
func (h *Handler) AddAccountUserInvite(in *grpc_user_go.AddAccountInviteRequest) (*grpc_user_go.AccountUserInvite, derrors.Error){return nil, nil}
func (h *Handler) GetAccountUserInvite(in *grpc_user_go.AccountUserInviteId) (*grpc_user_go.AccountUserInvite, derrors.Error){return nil, nil}
func (h *Handler) RemoveAccountUserInvite( in *grpc_user_go.AccountUserInviteId) (*grpc_common_go.Success, derrors.Error){return nil, nil}
func (h *Handler) ListAccountUserInvites(in *grpc_user_go.UserId) (*grpc_user_go.AccountInviteList, derrors.Error){return nil, nil}
// ---------------------------------------------------------------------------------------------------------------------
func (h *Handler) AddProjectUser(in *grpc_user_go.AddProjectUserRequest) (*grpc_user_go.ProjectUser, derrors.Error){return nil, nil}
func (h *Handler) RemoveProjectUser(in *grpc_user_go.ProjectUserId)  derrors.Error {return  nil}
func (h *Handler) UpdateProjectUser(in *grpc_user_go.ProjectUserUpdateRequest) (*grpc_user_go.ProjectUser, derrors.Error){return nil, nil}
func (h *Handler) ListProjectsUser(in *grpc_project_go.ProjectId) (*grpc_user_go.ProjectUserList, derrors.Error){return nil, nil}
// ---------------------------------------------------------------------------------------------------------------------
func (h *Handler) AddProjectUserInvite(in *grpc_user_go.AddProjectInviteRequest) (*grpc_user_go.ProjectUserInvite, derrors.Error){return nil, nil}
func (h *Handler) GetProjectUserInvite(in *grpc_user_go.ProjectUserInviteId) (*grpc_user_go.ProjectUserInvite, derrors.Error){return nil, nil}
func (h *Handler) RemoveProjectUserInvite(in *grpc_user_go.ProjectUserInviteId) derrors.Error {return  nil}
func (h *Handler) ListProjectUserInvites(in *grpc_user_go.UserId) (*grpc_user_go.ProjectInviteList, derrors.Error){return nil, nil}
*/