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
	DescriptorExists(appDescriptorID string) bool

	// Delete descriptor removes a given descriptor from the system.
	DeleteDescriptor(appDescriptorID string) derrors.Error

	// AddInstance adds a new application instance to the system
	AddInstance(instance entities.AppInstance) derrors.Error

	// InstanceExists checks if an application instance exists on the system.
	InstanceExists(appInstanceID string) bool

	// GetInstance retrieves an application instance.
	GetInstance(appInstanceID string) (* entities.AppInstance, derrors.Error)

	// DeleteInstance removes a given instance from the system.
	DeleteInstance(appInstanceID string) derrors.Error

	// Update status of this instance
	UpdateInstance(instance entities.AppInstance) derrors.Error

}
