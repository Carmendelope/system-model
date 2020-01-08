/*
 * Copyright 2020 Nalej
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
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/asset"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
)

type Manager struct {
	OrgProvider   organization.Provider
	AssetProvider asset.Provider
}

func NewManager(orgProvider organization.Provider, assetProvider asset.Provider) Manager {
	return Manager{orgProvider, assetProvider}
}

// Add a new asset to the system.
func (m *Manager) Add(addRequest *grpc_inventory_go.AddAssetRequest) (*entities.Asset, derrors.Error) {
	exists, err := m.OrgProvider.Exists(addRequest.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(addRequest.OrganizationId)
	}

	toAdd := entities.NewAssetFromGRPC(addRequest)
	err = m.AssetProvider.Add(*toAdd)
	if err != nil {
		return nil, err
	}
	return toAdd, nil
}

func (m *Manager) Get(assetID *grpc_inventory_go.AssetId) (*entities.Asset, derrors.Error) {
	exists, err := m.OrgProvider.Exists(assetID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(assetID.OrganizationId)
	}
	asset, err := m.AssetProvider.Get(assetID.AssetId)
	if err != nil {
		return nil, err
	}
	if asset.OrganizationId != assetID.OrganizationId {
		return nil, derrors.NewInternalError("asset organization_id does not match").WithParams(assetID.OrganizationId, assetID.AssetId)
	}
	return asset, nil
}

// List the assets of an organization.
func (m *Manager) List(organizationID *grpc_organization_go.OrganizationId) ([]entities.Asset, derrors.Error) {
	exists, err := m.OrgProvider.Exists(organizationID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(organizationID.OrganizationId)
	}
	groups, err := m.AssetProvider.List(organizationID.OrganizationId)
	if err != nil {
		return nil, err
	}
	return groups, nil
}

// Remove a given assets from an organization.
func (m *Manager) Remove(assetID *grpc_inventory_go.AssetId) derrors.Error {
	exists, err := m.OrgProvider.Exists(assetID.OrganizationId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organizationID").WithParams(assetID.OrganizationId)
	}
	asset, err := m.AssetProvider.Get(assetID.AssetId)
	if err != nil {
		return err
	}
	if asset.OrganizationId != assetID.OrganizationId {
		return derrors.NewNotFoundError("organization_id & asset_id").WithParams(assetID.OrganizationId, assetID.AssetId)
	}
	return m.AssetProvider.Remove(assetID.AssetId)
}

// Update the information of an asset.
func (m *Manager) Update(updateRequest *grpc_inventory_go.UpdateAssetRequest) (*entities.Asset, derrors.Error) {
	asset, err := m.AssetProvider.Get(updateRequest.AssetId)
	if err != nil {
		return nil, err
	}
	if asset.OrganizationId != updateRequest.OrganizationId {
		return nil, derrors.NewNotFoundError("organization_id & asset_id").WithParams(updateRequest.OrganizationId, updateRequest.AssetId)
	}
	asset.ApplyUpdate(updateRequest)
	err = m.AssetProvider.Update(*asset)
	if err != nil {
		return nil, err
	}
	return asset, nil
}

func (m *Manager) ListControllerAssets(edgeControllerId *grpc_inventory_go.EdgeControllerId) ([]entities.Asset, derrors.Error) {
	exists, err := m.OrgProvider.Exists(edgeControllerId.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(edgeControllerId.OrganizationId)
	}
	groups, err := m.AssetProvider.ListControllerAssets(edgeControllerId.EdgeControllerId)
	if err != nil {
		return nil, err
	}
	return groups, nil
}
