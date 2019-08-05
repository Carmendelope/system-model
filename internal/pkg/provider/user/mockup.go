/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package user

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

type MockupUserProvider struct {
	sync.Mutex
	// Users indexed by user email.
	users map[string]entities.User
	// accountUser indexed by email and accountID
	accountUser map[string]map[string]entities.AccountUser
	// accountUserInvite indexed by email and accountID
	accountUserInvite map[string]map[string]entities.AccountUserInvite
	// projectUser indexed by accountID#projectID and email
	projectUser map[string]map[string]entities.ProjectUser
	// accountUserInvite indexed by email and accountID#projectID
	projectUserInvite map[string]map[string]entities.ProjectUserInvite
}

func NewMockupUserProvider() * MockupUserProvider {
	return &MockupUserProvider{
		users: make(map[string]entities.User, 0),
		accountUser: make(map[string]map[string]entities.AccountUser, 0),
		accountUserInvite: make(map[string]map[string]entities.AccountUserInvite, 0),
		projectUser: make (map[string]map[string]entities.ProjectUser, 0),
		projectUserInvite: make (map[string]map[string]entities.ProjectUserInvite, 0),
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func (m * MockupUserProvider) unsafeExists(email string) bool {
	_, exists := m.users[email]
	return exists
}
// Add a new user to the system.
func (m * MockupUserProvider) Add(user entities.User) derrors.Error{
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(user.Email){
		m.users[user.Email] = user
		return nil
	}
	return derrors.NewAlreadyExistsError(user.Email)
}
// Update an existing user in the system
func (m * MockupUserProvider) Update(user entities.User) derrors.Error{
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(user.Email){
		return derrors.NewNotFoundError(user.Email)
	}
	m.users[user.Email] = user
	return nil
}
// Exists checks if a user exists on the system.
func (m * MockupUserProvider) Exists(email string) (bool, derrors.Error){
	m.Lock()
	defer m.Unlock()
	return m.unsafeExists(email), nil
}
// Get a user.
func (m * MockupUserProvider) Get(email string) (* entities.User, derrors.Error){
	m.Lock()
	defer m.Unlock()
	user, exists := m.users[email]
	if exists {
		return &user, nil
	}
	return nil, derrors.NewNotFoundError(email)
}
// Remove a user.
func (m * MockupUserProvider) Remove(email string) derrors.Error{
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(email){
		return derrors.NewNotFoundError(email)
	}
	delete(m.users, email)
	return nil
}
// ---------------------------------------------------------------------------------------------------------------------
func (m * MockupUserProvider) unsafeAccountUserExists(accountID string, email string) bool {
	accounts, exists := m.accountUser[email]
	if ! exists{
		return false
	}
	_, exists = accounts[accountID]
	return exists
}
func (m * MockupUserProvider) AddAccountUser(accUser entities.AccountUser) derrors.Error{
	m.Lock()
	defer m.Unlock()

	// ask if the user exists
	userExists := m.unsafeExists(accUser.Email)
	if ! userExists{
		return derrors.NewNotFoundError("User").WithParams(accUser.Email)
	}

	accounts, exists := m.accountUser[accUser.Email]
	if ! exists{

		m.accountUser[accUser.Email] = map[string]entities.AccountUser{
			accUser.AccountId: accUser,
		}

	} else{
		_, exists = accounts[accUser.AccountId]
		if ! exists{
			accounts[accUser.AccountId] = accUser
			m.accountUser[accUser.Email] = accounts // ????
		} else{
			return derrors.NewAlreadyExistsError("AccountUser").WithParams(accUser.AccountId, accUser.Email)
		}
	}
	return nil
}
func (m * MockupUserProvider) UpdateAccountUser(accUser entities.AccountUser) derrors.Error{
	m.Lock()
	defer m.Unlock()

	accounts, exists := m.accountUser[accUser.Email]
	if ! exists{
		return derrors.NewNotFoundError("AccountUser").WithParams(accUser.AccountId, accUser.Email)
	} else{
		_, exists = accounts[accUser.AccountId]
		if exists{
			m.accountUser[accUser.Email][accUser.AccountId] = accUser
		} else{
			return derrors.NewNotFoundError("AccountUser").WithParams(accUser.AccountId, accUser.Email)
		}
	}
	return nil
}
func (m * MockupUserProvider) RemoveAccountUser(accountID string, email string) derrors.Error{
	m.Lock()
	defer m.Unlock()

	accounts, exists := m.accountUser[email]
	if ! exists{
		return derrors.NewNotFoundError("AccountUser").WithParams(accountID, email)
	} else{
		_, exists = accounts[accountID]
		if ! exists{
			return derrors.NewNotFoundError("AccountUser").WithParams(accountID, email)
		} else{
			delete(m.accountUser[email], accountID)
		}
	}
	return nil
}
func (m * MockupUserProvider) ListAccountUser(email string) ([]entities.AccountUser, derrors.Error){
	m.Lock()
	defer m.Unlock()

	list := make([]entities.AccountUser, 0)

	// ask if the user exists
	userExists := m.unsafeExists(email)
	if ! userExists{
		return list, derrors.NewNotFoundError("User").WithParams(email)
	}

	accounts, exists := m.accountUser[email]
	if ! exists {
		return list, nil
	} else{
		for _, account := range accounts {
			list = append(list, account)
		}
	}
	return list, nil
}
// ---------------------------------------------------------------------------------------------------------------------
func (m * MockupUserProvider) unsafeAccountUserInviteExists(accountID string, email string) bool {
	accounts, exists := m.accountUserInvite[email]
	if ! exists{
		return false
	}
	_, exists = accounts[accountID]
	return exists
}
func (m * MockupUserProvider) AddAccountUserInvite(accUser entities.AccountUserInvite) derrors.Error{
	m.Lock()
	defer m.Unlock()
	// ask if the user exists
	userExists := m.unsafeExists(accUser.Email)
	if ! userExists{
		return derrors.NewNotFoundError("User").WithParams(accUser.Email)
	}

	accounts, exists := m.accountUserInvite[accUser.Email]
	if ! exists{

		m.accountUserInvite[accUser.Email] = map[string]entities.AccountUserInvite{
			accUser.AccountId: accUser,
		}

	} else{
		_, exists = accounts[accUser.AccountId]
		if ! exists{
			accounts[accUser.AccountId] = accUser
			m.accountUserInvite[accUser.Email] = accounts // ????
		} else{
			return derrors.NewAlreadyExistsError("AccountUserInvite").WithParams(accUser.AccountId, accUser.Email)
		}
	}
	return nil
}
func (m * MockupUserProvider) GetAccountUserInvite(accountID string, email string) (*entities.AccountUserInvite, derrors.Error){
	m.Lock()
	defer m.Unlock()

	accounts, exists := m.accountUserInvite[email]
	if exists{
		invite, exists := accounts[accountID]
		if exists{
			return &invite, nil
		} else{
			return nil, derrors.NewNotFoundError("AccountUserInvite").WithParams(accountID, email)
		}
	}
	return nil, derrors.NewNotFoundError("AccountUserInvite").WithParams(accountID, email)

}
func (m * MockupUserProvider) RemoveAccountUserInvite(accountID string, email string) derrors.Error{
	m.Lock()
	defer m.Unlock()

	accounts, exists := m.accountUserInvite[email]
	if ! exists{
		return derrors.NewNotFoundError("AccountUserInvite").WithParams(accountID, email)
	} else{
		_, exists = accounts[accountID]
		if ! exists{
			return derrors.NewNotFoundError("AccountUserInvite").WithParams(accountID, email)
		} else{
			delete(m.accountUserInvite[email], accountID)
		}
	}
	return nil
}
func (m * MockupUserProvider) ListAccountUserInvites(email string) ([]	entities.AccountUserInvite, derrors.Error){
	m.Lock()
	defer m.Unlock()

	list := make([]entities.AccountUserInvite, 0)

	// ask if the user exists
	userExists := m.unsafeExists(email)
	if ! userExists{
		return list, derrors.NewNotFoundError("User").WithParams(email)
	}

	accounts, exists := m.accountUserInvite[email]
	if ! exists {
		return list, nil
	} else{
		for _, account := range accounts {
			list = append(list, account)
		}
	}
	return list, nil
}
// ---------------------------------------------------------------------------------------------------------------------
func (m * MockupUserProvider) getProjectUserKey (accountID string, projectID string) string{
	return fmt.Sprintf("%s#%s", accountID, projectID)
}
func (m * MockupUserProvider) unsafeProjectUserExists(accountID string, projectID string, email string) bool {

	key := m.getProjectUserKey(accountID, projectID)

	users, exists := m.projectUser[key]
	if ! exists{
		return false
	}
	_, exists = users[email]
	return exists
}
func (m * MockupUserProvider) AddProjectUser(projUser entities.ProjectUser) derrors.Error{
	m.Lock()
	defer m.Unlock()

	key := m.getProjectUserKey(projUser.AccountId, projUser.ProjectId)

	// ask if the user exists
	userExists := m.unsafeExists(projUser.Email)
	if ! userExists{
		return derrors.NewNotFoundError("User").WithParams(projUser.Email)
	}

	users, exists := m.projectUser[key]
	if ! exists{

		m.projectUser[key] = map[string]entities.ProjectUser{
			projUser.Email: projUser,
		}

	} else{
		_, exists = users[projUser.Email]
		if ! exists{
			users[projUser.Email] = projUser
			m.projectUser[key] = users // ????
		} else{
			return derrors.NewAlreadyExistsError("ProjectUser").WithParams(projUser.AccountId, projUser.ProjectId, projUser.Email)
		}
	}
	return nil
}
func (m * MockupUserProvider) UpdateProjectUser(projUser entities.ProjectUser) derrors.Error{
	m.Lock()
	defer m.Unlock()

	key := m.getProjectUserKey(projUser.AccountId, projUser.ProjectId)

	users, exists := m.projectUser[key]
	if ! exists{
		return derrors.NewNotFoundError("ProjectUser").WithParams(projUser.AccountId, projUser.ProjectId, projUser.Email)
	} else{
		_, exists = users[projUser.Email]
		if exists{
			m.projectUser[key][projUser.Email] = projUser
		} else{
			return derrors.NewNotFoundError("ProjectUser").WithParams(projUser.AccountId, projUser.ProjectId, projUser.Email)
		}
	}
	return nil
}
func (m * MockupUserProvider) RemoveProjectUser(accountID string, projectID string, email string) derrors.Error{
	m.Lock()
	defer m.Unlock()

	key := m.getProjectUserKey(accountID, projectID)

	users, exists := m.projectUser[key]
	if ! exists{
		return derrors.NewNotFoundError("ProjectUser").WithParams(accountID, projectID, email)
	} else{
		_, exists = users[email]
		if ! exists{
			return derrors.NewNotFoundError("ProjectUser").WithParams(accountID, projectID, email)
		} else{
			delete(m.projectUser[key], email)
		}
	}
	return nil
}
func (m * MockupUserProvider) ListProjectUser(accountID string, projectID string) ([]entities.ProjectUser, derrors.Error){
	m.Lock()
	defer m.Unlock()

	list := make([]entities.ProjectUser, 0)

	key := m.getProjectUserKey(accountID, projectID)

	users, exists := m.projectUser[key]
	if ! exists {
		return list, nil
	} else{
		for _, user := range users {
			list = append(list, user)
		}
	}
	return list, nil
}
// ---------------------------------------------------------------------------------------------------------------------
func (m * MockupUserProvider) unsafeProjectUserInviteExists(accountID string, projectID string, email string) bool {
	invites, exists := m.projectUserInvite[email]
	if ! exists{
		return false
	}
	_, exists = invites[m.getProjectUserKey(accountID, projectID)]
	return exists
}

func (m * MockupUserProvider) AddProjectUserInvite(invite entities.ProjectUserInvite) derrors.Error{
	m.Lock()
	defer m.Unlock()
	// ask if the user exists
	userExists := m.unsafeExists(invite.Email)
	if ! userExists{
		return derrors.NewNotFoundError("User").WithParams(invite.Email)
	}

	invites, exists := m.projectUserInvite[invite.Email]
	if ! exists{

		m.projectUserInvite[invite.Email] = map[string]entities.ProjectUserInvite{
			m.getProjectUserKey(invite.AccountId, invite.ProjectId): invite,
		}

	} else{
		_, exists = invites[m.getProjectUserKey(invite.AccountId, invite.ProjectId)]
		if ! exists{
			invites[m.getProjectUserKey(invite.AccountId, invite.ProjectId)] = invite
			m.projectUserInvite[invite.Email] = invites // ????
		} else{
			return derrors.NewAlreadyExistsError("ProjectUserInvite").WithParams(invite.AccountId, invite.ProjectId, invite.Email)
		}
	}
	return nil
}
func (m * MockupUserProvider) GetProjectUserInvite(accountID string, projectId string, email string) (*entities.ProjectUserInvite, derrors.Error){
	m.Lock()
	defer m.Unlock()

	invites, exists := m.projectUserInvite[email]
	if exists{
		invite, exists := invites[m.getProjectUserKey(accountID, projectId)]
		if exists{
			return &invite, nil
		} else{
			return nil, derrors.NewNotFoundError("ProjectUserInvite").WithParams(accountID, projectId, email)
		}
	}
	return nil, derrors.NewNotFoundError("ProjectUserInvite").WithParams(accountID, projectId, email)
}
func (m * MockupUserProvider) RemoveProjectUserInvite(accountId string, projectId string, email string) derrors.Error{
	m.Lock()
	defer m.Unlock()

	invites, exists := m.projectUserInvite[email]
	if ! exists{
		return derrors.NewNotFoundError("ProjectUserInvite").WithParams(accountId, projectId, email)
	} else{
		_, exists = invites[m.getProjectUserKey(accountId, projectId)]
		if ! exists{
			return derrors.NewNotFoundError("AccountUserInvite").WithParams(accountId, projectId, email)
		} else{
			delete(m.projectUserInvite[email], m.getProjectUserKey(accountId, projectId))
		}
	}
	return nil
}
func (m * MockupUserProvider) ListProjectUserInvites(email string) ([]entities.ProjectUserInvite, derrors.Error){
	m.Lock()
	defer m.Unlock()

	list := make([]entities.ProjectUserInvite, 0)

	// ask if the user exists
	userExists := m.unsafeExists(email)
	if ! userExists{
		return list, derrors.NewNotFoundError("User").WithParams(email)
	}

	invites, exists := m.projectUserInvite[email]
	if ! exists {
		return list, nil
	} else{
		for _, invite := range invites {
			list = append(list, invite)
		}
	}
	return list, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// Clear cleans the contents of the mockup.
func (m * MockupUserProvider) Clear() derrors.Error{
	m.Lock()

	m.users = make(map[string]entities.User, 0)
	m.accountUser = make(map[string]map[string]entities.AccountUser)
	m.accountUserInvite = make(map[string]map[string]entities.AccountUserInvite)
	m.projectUser =  make (map[string]map[string]entities.ProjectUser, 0)
	m.projectUserInvite = make (map[string]map[string]entities.ProjectUserInvite, 0)
	m.Unlock()
	return nil
}


