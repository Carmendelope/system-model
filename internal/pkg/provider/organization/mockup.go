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
	// organizations contains the organization indexed per organization identifier.
	organizations map[string] entities.Organization
	// Descriptors contains the application descriptors ids per organization.
	descriptors map[string][]string
	// Instances contains the application instances ids per organization.
	instances map[string][]string
}


func NewMockupOrganizationProvider() * MockupOrganizationProvider {
	return &MockupOrganizationProvider{
		organizations:make(map[string]entities.Organization, 0),
		descriptors:make(map[string][]string, 0),
		instances:make(map[string][]string, 0),
	}
}

func (m * MockupOrganizationProvider) Clear() {
	m.Lock()
	m.organizations = make(map[string] entities.Organization, 0)
	m.descriptors = make(map[string] []string, 0)
	m.instances = make(map[string] []string, 0)
	m.Unlock()
}

func (m *MockupOrganizationProvider) unsafeExists(organizationID string) bool {
	_, exists := m.organizations[organizationID]
	return exists
}

func (m *MockupOrganizationProvider) unsafeExistsAppDesc(organizationID string, descriptorID string) bool {
	descriptors, ok := m.descriptors[organizationID]
	if ok {
		for _, descriptor := range descriptors {
			if descriptor == descriptorID {
				return true
			}
		}
		return false
	}
	return false
}

func (m *MockupOrganizationProvider) unsafeExistsAppInst(organizationID string, instanceID string) bool {
	instances, ok := m.instances[organizationID]
	if ok {
		for _, inst := range instances {
			if inst == instanceID {
				return true
			}
		}
		return false
	}
	return false
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

func (m *MockupOrganizationProvider) AddDescriptor(organizationID string, appDescriptorID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExists(organizationID) {
		if !m.unsafeExistsAppDesc(organizationID, appDescriptorID) {
			descriptors, _ := m.descriptors[organizationID]
			m.descriptors[organizationID] = append(descriptors, appDescriptorID)
			return nil
		}
		return derrors.NewAlreadyExistsError("descriptor").WithParams(organizationID, appDescriptorID)
	}
	return derrors.NewNotFoundError("organization").WithParams(organizationID)
}

func (m *MockupOrganizationProvider) DescriptorExists(organizationID string, appDescriptorID string) bool {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsAppDesc(organizationID, appDescriptorID)
}

func (m *MockupOrganizationProvider) ListDescriptors(organizationID string) ([]string, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	if !m.unsafeExists(organizationID) {
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	descriptors, ok := m.descriptors[organizationID]
	if ok {
		return descriptors, nil
	}
	return make([]string, 0), nil
}

func (m *MockupOrganizationProvider) DeleteDescriptor(organizationID string, appDescriptorID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExistsAppDesc(organizationID, appDescriptorID) {
		previous := m.descriptors[organizationID]
		newList := make([] string, 0, len(previous)-1)
		for _, id := range previous {
			if id != appDescriptorID {
				newList = append(newList, id)
			}
		}
		m.descriptors[organizationID] = newList
		return nil
	}
	return derrors.NewNotFoundError("descriptor").WithParams(organizationID, appDescriptorID)
}

func (m *MockupOrganizationProvider) AddInstance(organizationID string, appInstanceID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExists(organizationID) {
		if !m.unsafeExistsAppInst(organizationID, appInstanceID) {
			instances, _ := m.instances[organizationID]
			m.instances[organizationID] = append(instances, appInstanceID)
			return nil
		}
		return derrors.NewAlreadyExistsError("instance").WithParams(organizationID, appInstanceID)
	}
	return derrors.NewNotFoundError("organization").WithParams(organizationID)
}

func (m *MockupOrganizationProvider) InstanceExists(organizationID string, appInstanceID string) bool {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsAppInst(organizationID, appInstanceID)
}

func (m *MockupOrganizationProvider) ListInstances(organizationID string) ([]string, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	if !m.unsafeExists(organizationID) {
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	instances, ok := m.instances[organizationID]
	if ok {
		return instances, nil
	}
	return make([]string, 0), nil
}

func (m *MockupOrganizationProvider) DeleteInstance(organizationID string, appInstanceID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExistsAppDesc(organizationID, appInstanceID) {
		previous := m.instances[organizationID]
		newList := make([] string, 0, len(previous)-1)
		for _, id := range previous {
			if id != appInstanceID {
				newList = append(newList, id)
			}
		}
		m.instances[organizationID] = newList
		return nil
	}
	return derrors.NewNotFoundError("instance").WithParams(organizationID, appInstanceID)
}