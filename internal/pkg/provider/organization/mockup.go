/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package organization

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

type MockupOrganizationProvider struct {
	sync.Mutex
	organizations map[string] entities.Organization
}

func NewMockupOrganizationProvider() * MockupOrganizationProvider {
	return &MockupOrganizationProvider{
		organizations:make(map[string]entities.Organization, 0),
	}
}

func (m *MockupOrganizationProvider) unsafeExists(organizationID string) bool {
	_, exists := m.organizations[organizationID]
	return exists
}

func (m *MockupOrganizationProvider) Add(org entities.Organization) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(org.ID){
		m.organizations[org.ID] = org
		return nil
	}
	return derrors.NewAlreadyExistsError(org.ID)
}

func (m *MockupOrganizationProvider) Exists(organizationID string) bool {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExists(organizationID)
}

func (m *MockupOrganizationProvider) Get(organizationID string) (*entities.Organization, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	org, exists := m.organizations[organizationID]
	if exists {
		return &org, nil
	}
	return nil, derrors.NewNotFoundError(organizationID)
}


