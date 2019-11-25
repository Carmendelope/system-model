/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
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
	// ListControllerAssets retrieves the assets associated with a given edge controller
	ListControllerAssets(edgeControllerID string) ([]entities.Asset, derrors.Error)
	// Get an asset.
	Get(assetID string) (*entities.Asset, derrors.Error)
	// Remove an asset
	Remove(assetID string) derrors.Error
	// Clear all assets
	Clear() derrors.Error
}
