/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package application

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for application
type Provider interface {
	// AddDescriptor adds a new application descriptor to the system.
	AddDescriptor(descriptor entities.AppDescriptor) derrors.Error

	// GetDescriptors retrieves an application descriptor.
	GetDescriptor(appDescriptorID string) (* entities.AppDescriptor, derrors.Error)

	// DescriptorExists checks if a given descriptor exists on the system.
	DescriptorExists(appDescriptorID string) (bool, derrors.Error)

	// UpdateDescriptor updates the information of an application descriptor.
	UpdateDescriptor(descriptor entities.AppDescriptor) derrors.Error

	// DeleteDescriptor removes a given descriptor from the system.
	DeleteDescriptor(appDescriptorID string) derrors.Error

	// GetDescriptorParameters retrieves the params of a descriptor
	GetDescriptorParameters(appDescriptorID string) ([]entities.Parameter, derrors.Error)

	// AddInstance adds a new application instance to the system
	AddInstance(instance entities.AppInstance) derrors.Error

	// InstanceExists checks if an application instance exists on the system.
	InstanceExists(appInstanceID string) (bool, derrors.Error)

	// GetInstance retrieves an application instance.
	GetInstance(appInstanceID string) (* entities.AppInstance, derrors.Error)

	// DeleteInstance removes a given instance from the system.
	DeleteInstance(appInstanceID string) derrors.Error

	// UpdateInstance updates the information of an instance
	UpdateInstance(instance entities.AppInstance) derrors.Error

	// AddInstanceParameters adds deploy parameters of an instance in the system
	AddInstanceParameters (appInstanceID string, parameters []entities.InstanceParameter) derrors.Error

	// GetInstanceParameters retrieves the params of an instance
	GetInstanceParameters (appInstanceID string) ([]entities.InstanceParameter, derrors.Error)

	// DeleteInstanceParameters removes the params of an instance
	DeleteInstanceParameters (appInstanceID string) derrors.Error

	// AddParametrizedDescriptor adds a new parametrized descriptor to the system.
	AddParametrizedDescriptor(descriptor entities.ParametrizedDescriptor) derrors.Error

	// GetParametrizedDescriptor retrieves a parametrized descriptor
	GetParametrizedDescriptor(appInstanceID string) (*entities.ParametrizedDescriptor, derrors.Error)

	// ParametrizedDescriptorExists checks if a parametrized descriptor exists on the system.
	ParametrizedDescriptorExists (appInstanceID string) (*bool, derrors.Error)

	// DeleteParametrizedDescriptor removes a parametrized Descriptor from the system
	DeleteParametrizedDescriptor (appInstanceID string) derrors.Error

	// Clear descriptors and instances
	Clear() derrors.Error

	// AddAppEndPoint adds a new entry point to the system
	AddAppEndpoint (appEntryPoint entities.AppEndpoint) derrors.Error

	// GetAppEndPointByFQDN ()
	GetAppEndpointByFQDN(fqdn string) ([]*entities.AppEndpoint, derrors.Error)

	// DeleteAppEndpoints removes all the endpoint of an instance
	DeleteAppEndpoints(organizationID string, appInstanceID string) derrors.Error

	GetAppEndpointList(organizationID string , appInstanceId string, serviceGroupInstanceID string) ([]*entities.AppEndpoint, derrors.Error)

	// AddAppZtNetwork adds a new zerotier network to an existing application instance
	AddAppZtNetwork(network entities.AppZtNetwork) derrors.Error

	// RemoveAppZtNetwork removes any zt network belonging to an application instance
	RemoveAppZtNetwork(organizationID string, appInstanceID string) derrors.Error

	// GetAppZtNetwork get the zt network
	GetAppZtNetwork(organizationId string, appInstanceId string) (*entities.AppZtNetwork, derrors.Error)

	// AddZtNetworkMember add a new member for an existing zt network
	AddAppZtNetworkMember(member entities.AppZtNetworkMember) (*entities.AppZtNetworkMember, derrors.Error)

	// RemoveZtNetworkMember remove an existing member for a zt network
	RemoveAppZtNetworkMember(organizationId string, appInstanceId string, serviceGroupInstanceId string, serviceInstance string) derrors.Error

	// RemoveCompleteAppZtNetworkMemberNet removes all the entries stored under the same network id
	RemoveCompleteAppZtNetworkMemberNet(organizationId string, appInstanceId string, networkId string) derrors.Error

}
