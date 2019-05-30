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

	// parametrizedDescriptor indexed by AppInstanceID
	parametrizedDescriptor map[string] entities.ParametrizedDescriptor

	instanceParameters map[string] []entities.InstanceParameter

	appEntryPoints map[string] entities.AppEndpoint
	appEntryPointsByName map[string][]*entities.AppEndpoint

	appZtNetworks map[string]map[string]entities.AppZtNetwork

}

func NewMockupOrganizationProvider() * MockupApplicationProvider {
	return &MockupApplicationProvider{
		appDescriptors:make(map[string]entities.AppDescriptor, 0),
		appInstances: make(map[string]entities.AppInstance, 0),
		appEntryPoints:make(map[string]entities.AppEndpoint, 0),
		instanceParameters: make(map[string][]entities.InstanceParameter, 0),
		parametrizedDescriptor: make(map[string]entities.ParametrizedDescriptor, 0),
		appEntryPointsByName: make(map[string][]*entities.AppEndpoint, 0),
		appZtNetworks: make(map[string]map[string]entities.AppZtNetwork,0),
	}
}

// Clear cleans the contents of the mockup.
func (m * MockupApplicationProvider) Clear()  derrors.Error{
	m.Lock()
	defer m.Unlock()

	m.appDescriptors = make(map[string] entities.AppDescriptor, 0)
	m.appInstances = make(map[string] entities.AppInstance, 0)
	m.appEntryPoints = make(map[string]entities.AppEndpoint, 0)
	m.parametrizedDescriptor = make(map[string]entities.ParametrizedDescriptor, 0)

	m.appEntryPointsByName = make(map[string][]*entities.AppEndpoint, 0)
	m.appZtNetworks = make(map[string]map[string]entities.AppZtNetwork,0)

	m.instanceParameters = make(map[string][]entities.InstanceParameter, 0)

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

func (m *MockupApplicationProvider) unsafeExistsParamDesc(instanceID string) bool {
	_, exists := m.parametrizedDescriptor[instanceID]
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

func (m * MockupApplicationProvider) GetDescriptorParameters(appDescriptorID string) ([]entities.Parameter, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	d, e := m.appDescriptors[appDescriptorID]

	if !e {
		return nil, derrors.NewNotFoundError("descriptor").WithParams(appDescriptorID)
	}
	if d.Parameters == nil {
		d.Parameters = make([]entities.Parameter, 0)
	}
	return d.Parameters, nil
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

// -- Instance parameters
// AddInstanceParameters adds deploy parameters of an instance in the system
func (m * MockupApplicationProvider)AddInstanceParameters (appInstanceID string, parameters []entities.InstanceParameter) derrors.Error{
	m.Lock()
	defer m.Unlock()

	_, exists := m.instanceParameters[appInstanceID]

	if exists {
		return derrors.NewAlreadyExistsError("parameters").WithParams(appInstanceID)
	}

	m.instanceParameters[appInstanceID] = parameters

	return nil
}
// GetInstanceParameters retrieves the params of an instance
func (m * MockupApplicationProvider) GetInstanceParameters (appInstanceID string) ([]entities.InstanceParameter, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	params, exists := m.instanceParameters[appInstanceID]

	if !exists {
		params := make ([]entities.InstanceParameter, 0)
		return params, nil
	}
	return params, nil
}


// DeleteInstanceParameters removes the params of an instance
func (m * MockupApplicationProvider)DeleteInstanceParameters (appInstanceID string) derrors.Error {
	m.Lock()
	defer m.Unlock()

	delete (m.instanceParameters, appInstanceID)

	return nil
}

// AddParametrizedDescriptor adds a new parametrized descriptor to the system.
func (m * MockupApplicationProvider)AddParametrizedDescriptor(descriptor entities.ParametrizedDescriptor) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExistsParamDesc(descriptor.AppInstanceId){
		m.parametrizedDescriptor[descriptor.AppInstanceId] = descriptor

		return nil
	}
	return derrors.NewAlreadyExistsError(descriptor.AppInstanceId)
}

// GetParametrizedDescriptor retrieves a parametrized descriptor
func (m * MockupApplicationProvider) GetParametrizedDescriptor(appInstanceID string) (*entities.ParametrizedDescriptor, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	i, e := m.parametrizedDescriptor[appInstanceID]
	if !e {
		return nil, derrors.NewNotFoundError("parametrized descriptor").WithParams(appInstanceID)
	}
	return &i, nil
}

// ParametrizedDescriptorExists checks if a parametrized descriptor exists on the system.
func (m * MockupApplicationProvider)	ParametrizedDescriptorExists (appInstanceID string) (*bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	exists := m.unsafeExistsParamDesc(appInstanceID)

	return &exists, nil
}

// DeleteParametrizedDescriptor removes a parametrized Descriptor from the system
func (m * MockupApplicationProvider) DeleteParametrizedDescriptor (appInstanceID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExistsParamDesc(appInstanceID) {
		return derrors.NewNotFoundError("parametrized descriptor").WithParams(appInstanceID)
	}
	delete(m.parametrizedDescriptor, appInstanceID)
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
		return nil, derrors.NewNotFoundError("appEndPoint").WithParams(fqdn)
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
// AppZtNetwork functions

func (m *MockupApplicationProvider) AddAppZtNetwork(ztNetwork entities.AppZtNetwork) derrors.Error {
	m.Lock()
	defer m.Unlock()

	_, foundOrg := m.appZtNetworks[ztNetwork.OrganizationId]
	if !foundOrg {
		m.appZtNetworks[ztNetwork.OrganizationId] = map[string]entities.AppZtNetwork{ztNetwork.AppInstanceId:ztNetwork}
	} else {
		m.appZtNetworks[ztNetwork.OrganizationId][ztNetwork.AppInstanceId] = ztNetwork
	}
	
	return nil
}


func (m *MockupApplicationProvider) RemoveAppZtNetwork(organizationID string, appInstanceID string) derrors.Error {
	m.Lock()
	defer m.Unlock()

	_, foundOrg := m.appZtNetworks[organizationID]
	if !foundOrg {
		return derrors.NewNotFoundError("non existing organization")
	}
	_,foundAppInstance := m.appZtNetworks[organizationID][appInstanceID]
	if !foundAppInstance {
		return derrors.NewNotFoundError("not existing application instance")
	}
	delete(m.appZtNetworks[organizationID],appInstanceID)
	if len(m.appZtNetworks[organizationID])==0{
		delete(m.appZtNetworks,organizationID)
	}

	return nil
}

func (m *MockupApplicationProvider) GetAppZtNetwork(organizationID string, appInstanceID string) (*entities.AppZtNetwork, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	_, foundOrg := m.appZtNetworks[organizationID]
	if !foundOrg {
		return nil, derrors.NewNotFoundError("non existing organization")
	}
	toReturn,foundAppInstance := m.appZtNetworks[organizationID][appInstanceID]
	if !foundAppInstance {
		return nil,derrors.NewNotFoundError("not existing application instance")
	}

	return &toReturn, nil
}


func (m *MockupApplicationProvider) AddAppZtNetworkMember(member entities.AppZtNetworkMember) (*entities.AppZtNetworkMember, derrors.Error) {
	return nil, derrors.NewUnimplementedError("AddAppZtNetworkMember not implemented yet")
}

func (m *MockupApplicationProvider) RemoveAppZtNetworkMember(organizationId string, appInstanceId string, serviceGroupInstanceId string, serviceInstance string) derrors.Error {
	return derrors.NewUnimplementedError("RemoveAppZtNetworkMember not implemented yet")
}

func (m *MockupApplicationProvider) RemoveCompleteAppZtNetworkMemberNet(organizationId string, appInstanceId string, networkId string) derrors.Error {
	return derrors.NewUnimplementedError("RemoveCompleteAppZtNetworkMemberNet not implemented yet")
}