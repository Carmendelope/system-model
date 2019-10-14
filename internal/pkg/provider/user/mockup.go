/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package user

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

type MockupUserProvider struct {
	sync.Mutex
	// Users indexed by user email.
	users map[string]entities.User
}

func NewMockupUserProvider() *MockupUserProvider {
	return &MockupUserProvider{
		users: make(map[string]entities.User, 0),
	}
}

func (m *MockupUserProvider) unsafeExists(email string) bool {
	_, exists := m.users[email]
	return exists
}

// Add a new user to the system.
func (m *MockupUserProvider) Add(user entities.User) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(user.Email) {
		m.users[user.Email] = user
		return nil
	}
	return derrors.NewAlreadyExistsError(user.Email)
}

// Update an existing user in the system
func (m *MockupUserProvider) Update(user entities.User) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(user.Email) {
		return derrors.NewNotFoundError(user.Email)
	}
	m.users[user.Email] = user
	return nil
}

// Exists checks if a user exists on the system.
func (m *MockupUserProvider) Exists(email string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExists(email), nil
}

// Get a user.
func (m *MockupUserProvider) Get(email string) (*entities.User, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	user, exists := m.users[email]
	if exists {
		return &user, nil
	}
	return nil, derrors.NewNotFoundError(email)
}

// Remove a user.
func (m *MockupUserProvider) Remove(email string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(email) {
		return derrors.NewNotFoundError(email)
	}
	delete(m.users, email)
	return nil
}

// Clear cleans the contents of the mockup.
func (m *MockupUserProvider) Clear() derrors.Error {
	m.Lock()
	m.users = make(map[string]entities.User, 0)
	m.Unlock()
	return nil
}
