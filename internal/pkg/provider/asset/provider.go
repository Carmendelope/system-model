/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package asset

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for application
type Provider interface {
	// Add a new asset to the system.
	Add(asset entities.Asset) derrors.Error
	// Update the information of an asset.
	Update(asset entities.Asset) derrors.Error
	// Exists checks if an asset exists on the system.
	Exists(assetID string) (bool, derrors.Error)
	// List the assets in a given organization
	List(organizationID string) ([]entities.Asset, derrors.Error)
	// Get an asset.
	Get(assetID string) (* entities.Asset, derrors.Error)
	// Remove an asset
	Remove(assetID string) derrors.Error
	// Clear all assets
	Clear() derrors.Error
}