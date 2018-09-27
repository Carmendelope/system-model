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
}
