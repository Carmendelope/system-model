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

	// Clear descriptors and instances
	Clear() derrors.Error

	// AddAppEntryPoint adds a new entry point to the system
	AddAppEntryPoint (appEntryPoint entities.AppEndpoint) derrors.Error

	// GetAppEntryPointByFQDN ()
	GetAppEntryPointByFQDN(fqdn string) ([]*entities.AppEndpoint, derrors.Error)

	// DeleteAppEndpoints removes all the endpoint of an instance
	DeleteAppEndpoints(organizationID string, appInstanceID string) derrors.Error

	// AddAppZtNetwork adds a new zerotier network to an existing application instance
	AddAppZtNetwork(network entities.AppZtNetwork) derrors.Error

	// RemoveAppZtNetwork removes any zt network belonging to an application instance
	RemoveAppZtNetwork(organizationID string, appInstanceID string) derrors.Error

	// GetAppZtNetwork get the zt network
	GetAppZtNetwork(organizationId string, appInstanceId string) (*entities.AppZtNetwork, derrors.Error)

}
