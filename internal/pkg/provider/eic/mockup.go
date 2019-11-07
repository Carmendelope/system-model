/*
 * Copyright 2019 Nalej
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
 *
 */

package eic

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

type MockupEICProvider struct {
	// Mutex for managing mockup access.
	sync.Mutex
	// Assets with a map of EIC indexed by edgeControllerID.
	controllers map[string]entities.EdgeController
}

func NewMockupEICProvider() *MockupEICProvider {
	return &MockupEICProvider{
		controllers: make(map[string]entities.EdgeController, 0),
	}
}

func (m *MockupEICProvider) unsafeExists(edgeControllerID string) bool {
	_, exists := m.controllers[edgeControllerID]
	return exists
}

func (m *MockupEICProvider) Add(eic entities.EdgeController) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(eic.EdgeControllerId) {
		m.controllers[eic.EdgeControllerId] = eic
		return nil
	}
	return derrors.NewAlreadyExistsError(eic.EdgeControllerId)
}

func (m *MockupEICProvider) Update(eic entities.EdgeController) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(eic.EdgeControllerId) {
		return derrors.NewNotFoundError(eic.EdgeControllerId)
	}
	m.controllers[eic.EdgeControllerId] = eic
	return nil
}

func (m *MockupEICProvider) Exists(edgeControllerID string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExists(edgeControllerID), nil
}

func (m *MockupEICProvider) Get(edgeControllerID string) (*entities.EdgeController, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	eic, exists := m.controllers[edgeControllerID]
	if exists {
		return &eic, nil
	}
	return nil, derrors.NewNotFoundError(edgeControllerID)
}

// List the EIC in a given organization
func (m *MockupEICProvider) List(organizationID string) ([]entities.EdgeController, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	result := make([]entities.EdgeController, 0)
	for _, eic := range m.controllers {
		if eic.OrganizationId == organizationID {
			result = append(result, eic)
		}
	}
	return result, nil
}

func (m *MockupEICProvider) Remove(edgeControllerID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(edgeControllerID) {
		return derrors.NewNotFoundError(edgeControllerID)
	}
	delete(m.controllers, edgeControllerID)
	return nil
}

func (m *MockupEICProvider) Clear() derrors.Error {
	m.Lock()
	m.controllers = make(map[string]entities.EdgeController, 0)
	m.Unlock()
	return nil
}
