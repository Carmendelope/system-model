/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package application

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

type MockupApplicationProvider struct {
	sync.Mutex
	appDescriptors map[string] entities.AppDescriptor
	appInstances map[string] entities.AppInstance

	appEntryPoints map[string] entities.AppEndpoint
	appEntryPointsByName map[string][]*entities.AppEndpoint
}

func NewMockupOrganizationProvider() * MockupApplicationProvider {
	return &MockupApplicationProvider{
		appDescriptors:make(map[string]entities.AppDescriptor, 0),
		appInstances: make(map[string]entities.AppInstance, 0),
		appEntryPoints:make(map[string]entities.AppEndpoint, 0),
		appEntryPointsByName: make(map[string][]*entities.AppEndpoint, 0),
	}
}

// Clear cleans the contents of the mockup.
func (m * MockupApplicationProvider) Clear()  derrors.Error{
	m.Lock()
	defer m.Unlock()

	m.appDescriptors = make(map[string] entities.AppDescriptor, 0)
	m.appInstances = make(map[string] entities.AppInstance, 0)
	m.appEntryPoints = make(map[string]entities.AppEndpoint, 0)
	m.appEntryPointsByName = make(map[string][]*entities.AppEndpoint, 0)

	return nil
}

func (m *MockupApplicationProvider) unsafeExistsAppDesc(descriptorID string) bool {
	_, exists := m.appDescriptors[descriptorID]
	return exists
}

func (m *MockupApplicationProvider) unsafeExistsAppInst(instanceID string) bool {
	_, exists := m.appInstances[instanceID]
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
func (m *MockupApplicationProvider) DescriptorExists(appDescriptorID string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsAppDesc(appDescriptorID), nil
}

// UpdateDescriptor updates the information of an application descriptor.
func (m *MockupApplicationProvider) UpdateDescriptor(descriptor entities.AppDescriptor) derrors.Error{
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExistsAppDesc(descriptor.AppDescriptorId){
		return derrors.NewNotFoundError(descriptor.AppDescriptorId)
	}
	m.appDescriptors[descriptor.AppDescriptorId] = descriptor
	return nil
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
func (m *MockupApplicationProvider) InstanceExists(appInstanceID string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsAppInst(appInstanceID), nil
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

func (m*MockupApplicationProvider)getAppEndpointKey(appEntryPoint entities.AppEndpoint) string {
	return fmt.Sprintf("%s-%s-%s-%s-%d", appEntryPoint.OrganizationId, appEntryPoint.AppInstanceId,
		appEntryPoint.ServiceGroupInstanceId, appEntryPoint.ServiceInstanceId, appEntryPoint.Port)
}

// AddAppEntryPoint adds a new entry point to the system
func (m *MockupApplicationProvider)AddAppEndpoint (appEntryPoint entities.AppEndpoint) derrors.Error {
	m.Lock()
	defer m.Unlock()

	key := m.getAppEndpointKey(appEntryPoint)
	m.appEntryPoints[key] = appEntryPoint

	list, exists := m.appEntryPointsByName[appEntryPoint.GlobalFqdn]
	if exists{
		m.appEntryPointsByName[appEntryPoint.GlobalFqdn] = append(list, &appEntryPoint)
	}else {
		m.appEntryPointsByName[appEntryPoint.GlobalFqdn] = []*entities.AppEndpoint{&appEntryPoint}
	}

	return nil
}

// GetAppEntryPointByFQDN ()
func (m *MockupApplicationProvider) GetAppEndpointByFQDN(fqdn string) ([]*entities.AppEndpoint, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	list, exists := m.appEntryPointsByName[fqdn]
	if exists{
		return list, nil
	}else {
		return nil, derrors.NewNotFoundError("appEntryPoint").WithParams(fqdn)
	}
}

func (m *MockupApplicationProvider) DeleteAppEndpoints(organizationID string, appInstanceID string) derrors.Error {
	m.Lock()
	defer m.Unlock()

	for key, endpoint := range m.appEntryPoints{
		if endpoint.OrganizationId == organizationID && endpoint.AppInstanceId == appInstanceID {
			delete (m.appEntryPointsByName, endpoint.GlobalFqdn)
			delete(m.appEntryPoints, key)
		}
	}
	return nil
}

func (m *MockupApplicationProvider) GetAppEndpointList(organizationID string , appInstanceId string,
	serviceGroupInstanceID string) ([]*entities.AppEndpoint, derrors.Error) {

	m.Lock()
	defer m.Unlock()

	list := make ([]*entities.AppEndpoint, 0)
	for _, endpoint := range m.appEntryPoints{
		if endpoint.OrganizationId == organizationID && endpoint.AppInstanceId == appInstanceId &&
			endpoint.ServiceGroupInstanceId == serviceGroupInstanceID {
			list = append (list, &endpoint)
		}
	}
	return list, nil
}

