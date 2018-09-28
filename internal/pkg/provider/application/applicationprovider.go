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

	// DescriptorExists checks if an application descriptor exists on the system.
	DescriptorExists(organizationID string, appDescriptorID string) bool

	// ListDescriptors obtains a list of the application descriptors associated with a given organization.
	ListDescriptors(organizationID string) ([] entities.AppDescriptor, derrors.Error)

	// GetDescriptors retrieves an application descriptor.
	GetDescriptor(organizationID string, appDescriptorID string) (* entities.AppDescriptor, derrors.Error)

	// AddInstance adds a new application instance to the system
	AddInstance(instance entities.AppInstance) derrors.Error

	// InstanceExists checks if an application instance exists on the system.
	InstanceExists(organizationID string, appInstanceID string) bool

	// ListInstances obtains a list of the application instances associated with a given organization.
	ListInstances(organizationID string) ([]entities.AppInstance, derrors.Error)

	// GetInstance retrieves an application instance.
	GetInstance(organizationID string, appInstanceID string) (* entities.AppInstance, derrors.Error)

}
