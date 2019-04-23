/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
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

func NewMockupEICProvider() * MockupEICProvider{
	return &MockupEICProvider{
		controllers: make(map[string]entities.EdgeController, 0),
	}
}

func (m*MockupEICProvider) unsafeExists(edgeControllerID string) bool{
	_, exists := m.controllers[edgeControllerID]
	return exists
}

func (m *MockupEICProvider) Add(eic entities.EdgeController) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(eic.EdgeControllerId){
		m.controllers[eic.EdgeControllerId] = eic
		return nil
	}
	return derrors.NewAlreadyExistsError(eic.EdgeControllerId)
}

func (m *MockupEICProvider) Update(eic entities.EdgeController) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(eic.EdgeControllerId){
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
func (m *MockupEICProvider) List(organizationID string) ([]entities.EdgeController, derrors.Error){
	m.Lock()
	defer m.Unlock()
	result := make([]entities.EdgeController, 0)
	for _, eic := range m.controllers{
		if eic.OrganizationId == organizationID{
			result = append(result, eic)
		}
	}
	return result, nil
}

func (m *MockupEICProvider) Remove(edgeControllerID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(edgeControllerID){
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


