/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package eic

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for application
type Provider interface {
	// Add a new edge controller to the system.
	Add(eic entities.EdgeController) derrors.Error
	// Update the information of an edge controller.
	Update(eic entities.EdgeController) derrors.Error
	// Exists checks if an EIC exists on the system.
	Exists(edgeControllerID string) (bool, derrors.Error)
	// Get an EIC.
	Get(edgeControllerID string) (*entities.EdgeController, derrors.Error)
	// List the EIC in a given organization
	List(organizationID string) ([]entities.EdgeController, derrors.Error)
	// Remove an EIC
	Remove(edgeControllerID string) derrors.Error
	// Clear all assets
	Clear() derrors.Error
}
