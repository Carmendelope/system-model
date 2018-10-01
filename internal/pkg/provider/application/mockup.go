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

func (m *MockupApplicationProvider) AddDescriptor(descriptor entities.AppDescriptor) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExistsAppDesc(descriptor.AppDescriptorId){
		m.appDescriptors[descriptor.AppDescriptorId] = descriptor
		return nil
	}
	return derrors.NewAlreadyExistsError(descriptor.AppDescriptorId)
}

func (m *MockupApplicationProvider) DescriptorExists(appDescriptorID string) bool {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsAppDesc(appDescriptorID)
}

func (m *MockupApplicationProvider) GetDescriptor(appDescriptorID string) (*entities.AppDescriptor, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	d, e := m.appDescriptors[appDescriptorID]
	if !e {
		return nil, derrors.NewNotFoundError("descriptor").WithParams(appDescriptorID)
	}
	return &d, nil
}

// Delete descriptor removes a given descriptor from the system.
func (m * MockupApplicationProvider) DeleteDescriptor(appDescriptorID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExistsAppDesc(appDescriptorID) {
		return derrors.NewNotFoundError("descriptor").WithParams(appDescriptorID)
	}
	delete(m.appDescriptors, appDescriptorID)
	return nil
}

func (m *MockupApplicationProvider) AddInstance(instance entities.AppInstance) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExistsAppDesc(instance.AppInstanceId){
		m.appInstances[instance.AppInstanceId] = instance
		return nil
	}
	return derrors.NewAlreadyExistsError(instance.AppDescriptorId)
}

func (m *MockupApplicationProvider) InstanceExists(appInstanceID string) bool {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsAppInst(appInstanceID)
}

func (m *MockupApplicationProvider) GetInstance(appInstanceID string) (*entities.AppInstance, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	i, e := m.appInstances[appInstanceID]
	if !e {
		return nil, derrors.NewNotFoundError("instance").WithParams(appInstanceID)
	}
	return &i, nil
}

func (m *MockupApplicationProvider) DeleteInstance(appInstanceID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExistsAppInst(appInstanceID) {
		return derrors.NewNotFoundError("instance").WithParams(appInstanceID)
	}
	delete(m.appInstances, appInstanceID)
	return nil
}






