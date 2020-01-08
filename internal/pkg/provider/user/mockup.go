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
