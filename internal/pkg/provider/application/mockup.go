/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package application

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

type MockupApplicationProvider struct {
	sync.Mutex
	appDescriptors map[string] entities.AppDescriptor
	appInstances map[string] entities.AppInstance
}

func NewMockupOrganizationProvider() * MockupApplicationProvider {
	return &MockupApplicationProvider{
		appDescriptors:make(map[string]entities.AppDescriptor, 0),
		appInstances: make(map[string]entities.AppInstance, 0),
	}
}

// Clear cleans the contents of the mockup.
func (m * MockupApplicationProvider) Clear() {
	m.Lock()
	m.appDescriptors = make(map[string] entities.AppDescriptor, 0)
	m.appInstances = make(map[string] entities.AppInstance, 0)
	m.Unlock()
}

func (mockup *MockupApplicationProvider) unsafeExistsAppDesc(descriptorID string) bool {
	_, exists := mockup.appDescriptors[descriptorID]
	return exists
}

func (mockup *MockupApplicationProvider) unsafeExistsAppInst(instanceID string) bool {
	_, exists := mockup.appInstances[instanceID]
	return exists
}

// AddDescriptor adds a new application descriptor to the system.
func (m *MockupApplicationProvider) AddDescriptor(descriptor entities.AppDescriptor) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExistsAppDesc(descriptor.AppDescriptorId){
		m.appDescriptors[descriptor.AppDescriptorId] = descriptor
		return nil
	}
	return derrors.NewAlreadyExistsError(descriptor.AppDescriptorId)
}

// DescriptorExists checks if a given descriptor exists on the system.
func (m *MockupApplicationProvider) DescriptorExists(appDescriptorID string) bool {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsAppDesc(appDescriptorID)
}

// GetDescriptors retrieves an application descriptor.
func (m *MockupApplicationProvider) GetDescriptor(appDescriptorID string) (*entities.AppDescriptor, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	d, e := m.appDescriptors[appDescriptorID]
	if !e {
		return nil, derrors.NewNotFoundError("descriptor").WithParams(appDescriptorID)
	}
	return &d, nil
}

// DeleteDescriptor removes a given descriptor from the system.
func (m * MockupApplicationProvider) DeleteDescriptor(appDescriptorID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExistsAppDesc(appDescriptorID) {
		return derrors.NewNotFoundError("descriptor").WithParams(appDescriptorID)
	}
	delete(m.appDescriptors, appDescriptorID)
	return nil
}

// AddInstance adds a new application instance to the system
func (m *MockupApplicationProvider) AddInstance(instance entities.AppInstance) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExistsAppDesc(instance.AppInstanceId){
		m.appInstances[instance.AppInstanceId] = instance
		return nil
	}
	return derrors.NewAlreadyExistsError(instance.AppDescriptorId)
}

// InstanceExists checks if an application instance exists on the system.
func (m *MockupApplicationProvider) InstanceExists(appInstanceID string) bool {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsAppInst(appInstanceID)
}

// GetInstance retrieves an application instance.
func (m *MockupApplicationProvider) GetInstance(appInstanceID string) (*entities.AppInstance, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	i, e := m.appInstances[appInstanceID]
	if !e {
		return nil, derrors.NewNotFoundError("instance").WithParams(appInstanceID)
	}
	return &i, nil
}

// DeleteInstance removes a given instance from the system.
func (m *MockupApplicationProvider) DeleteInstance(appInstanceID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExistsAppInst(appInstanceID) {
		return derrors.NewNotFoundError("instance").WithParams(appInstanceID)
	}
	delete(m.appInstances, appInstanceID)
	return nil
}

// UpdateInstance updates the information of an instance
func (m *MockupApplicationProvider) UpdateInstance(instance entities.AppInstance) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExistsAppInst(instance.AppInstanceId) {
		return derrors.NewNotFoundError("instance").WithParams(instance.AppInstanceId)
	}
	m.appInstances[instance.AppInstanceId] = instance
	return nil
}




