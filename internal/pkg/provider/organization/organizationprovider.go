/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package organization

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

type Provider interface {
	// Add a new organization to the system.
	Add(org entities.Organization) derrors.Error
	// Check if an organization exists on the system.
	Exists(organizationID string) bool
	// Get an organization.
	Get(organizationID string) (* entities.Organization, derrors.Error)

	// AddDescriptor adds a new descriptor ID to a given organization.
	AddDescriptor(organizationID string, appDescriptorID string) derrors.Error
	// DescriptorExists checks if an application descriptor exists on the system.
	DescriptorExists(organizationID string, appDescriptorID string) bool
	// ListDescriptors returns the identifiers of the application descriptors associated with an organization.
	ListDescriptors(organizationID string) ([]string, derrors.Error)
	// DeleteDescriptor removes a descriptor from an organization
	DeleteDescriptor(organizationID string, appDescriptorID string) derrors.Error
}
