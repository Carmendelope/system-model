/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package user

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for application
type Provider interface {
	// Add a new user to the system.
	Add(user entities.User) derrors.Error
	// Update an existing user in the system
	Update(user entities.User) derrors.Error
	// Exists checks if a user exists on the system.
	Exists(email string) (bool, derrors.Error)
	// Get a user.
	Get(email string) (* entities.User, derrors.Error)
	// Remove a user.
	Remove(email string) derrors.Error

	AddAccountUser(accUser entities.AccountUser) derrors.Error
	UpdateAccountUser(accUser entities.AccountUser) derrors.Error
	RemoveAccountUser(accountID string, email string) derrors.Error
	GetAccountUser(accountID string, email string) (*entities.AccountUser, derrors.Error)
	// TODO: change to ListAccountUser(AccountID string)
	ListAccountUser(email string) ([]entities.AccountUser, derrors.Error)

	AddAccountUserInvite(accUser entities.AccountUserInvite) derrors.Error
	GetAccountUserInvite(accountID string, email string) (*entities.AccountUserInvite, derrors.Error)
	RemoveAccountUserInvite(accountID string, email string) derrors.Error
	ListAccountUserInvites(email string) ([]entities.AccountUserInvite, derrors.Error)

	AddProjectUser(projUser entities.ProjectUser) derrors.Error
	UpdateProjectUser(projUser entities.ProjectUser) derrors.Error
	RemoveProjectUser(accountID string, projectID string, email string) derrors.Error
	ListProjectUser(accountID string, projectID string) ([]entities.ProjectUser, derrors.Error)

	AddProjectUserInvite(invite entities.ProjectUserInvite) derrors.Error
	GetProjectUserInvite(accountID string, projectID string, email string) (*entities.ProjectUserInvite, derrors.Error)
	RemoveProjectUserInvite(accountID string, projectID string, email string) derrors.Error
	ListProjectUserInvites(email string) ([]entities.ProjectUserInvite, derrors.Error)


	// Clear
	Clear() derrors.Error

}
