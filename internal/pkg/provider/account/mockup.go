/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package account

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

type MockupAccountProvider struct {
	// Mutex for managing mockup access.
	sync.Mutex
	// accounts with a map of assets indexed by accountID.
	accounts     map[string]entities.Account
	accountNames map[string]bool
}

func NewMockupAccountProvider() *MockupAccountProvider {
	return &MockupAccountProvider{
		accounts:     make(map[string]entities.Account, 0),
		accountNames: make(map[string]bool, 0),
	}
}

func (m *MockupAccountProvider) unsafeExists(accountID string) bool {
	_, exists := m.accounts[accountID]
	return exists
}

// Add a new account to the system.
func (m *MockupAccountProvider) Add(account entities.Account) derrors.Error {
	m.Lock()
	defer m.Unlock()

	if !m.unsafeExists(account.AccountId) {
		m.accounts[account.AccountId] = account
		m.accountNames[account.Name] = true
		return nil
	}
	return derrors.NewAlreadyExistsError(account.AccountId)
}

// Update the information of an account.
func (m *MockupAccountProvider) Update(account entities.Account) derrors.Error {
	m.Lock()
	defer m.Unlock()

	if !m.unsafeExists(account.AccountId) {
		return derrors.NewNotFoundError(account.AccountId)
	}
	delete(m.accountNames, m.accounts[account.AccountId].Name)
	m.accounts[account.AccountId] = account
	m.accountNames[account.Name] = true
	return nil
}

// Exists checks if an account exists on the system.
func (m *MockupAccountProvider) Exists(accountID string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	return m.unsafeExists(accountID), nil
}

func (m *MockupAccountProvider) ExistsByName(accountName string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	_, exists := m.accountNames[accountName]
	return exists, nil

}

// Get an account.
func (m *MockupAccountProvider) Get(accountID string) (*entities.Account, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	asset, exists := m.accounts[accountID]
	if exists {
		return &asset, nil
	}
	return nil, derrors.NewNotFoundError(accountID)
}

func (m *MockupAccountProvider) List() ([]entities.Account, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	list := make([]entities.Account, 0)
	for _, account := range m.accounts {
		list = append(list, account)
	}
	return list, nil
}

// Remove an account
func (m *MockupAccountProvider) Remove(accountID string) derrors.Error {
	m.Lock()
	defer m.Unlock()

	if !m.unsafeExists(accountID) {
		return derrors.NewNotFoundError(accountID)
	}
	delete(m.accountNames, m.accounts[accountID].Name)
	delete(m.accounts, accountID)
	return nil
}

// Clear all accounts
func (m *MockupAccountProvider) Clear() derrors.Error {
	m.Lock()
	defer m.Unlock()
	m.accounts = make(map[string]entities.Account, 0)
	m.accountNames = make(map[string]bool, 0)
	return nil
}
