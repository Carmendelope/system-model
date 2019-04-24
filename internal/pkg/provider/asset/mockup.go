/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package asset

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

type MockupAssetProvider struct {
	// Mutex for managing mockup access.
	sync.Mutex
	// Assets with a map of assets indexed by assetID.
	assets map[string]entities.Asset
}

func NewMockupAssetProvider() * MockupAssetProvider{
	return &MockupAssetProvider{
		assets: make(map[string]entities.Asset, 0),
	}
}

func (m*MockupAssetProvider) unsafeExists(assetID string) bool{
	_, exists := m.assets[assetID]
	return exists
}

func (m * MockupAssetProvider) Add(asset entities.Asset) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(asset.AssetId){
		m.assets[asset.AssetId] = asset
		return nil
	}
	return derrors.NewAlreadyExistsError(asset.AssetId)
}

func (m * MockupAssetProvider) Update(asset entities.Asset) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(asset.AssetId){
		return derrors.NewNotFoundError(asset.AssetId)
	}
	m.assets[asset.AssetId] = asset
	return nil
}

func (m * MockupAssetProvider) Exists(assetID string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExists(assetID), nil
}

func (m * MockupAssetProvider) Get(assetID string) (*entities.Asset, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	asset, exists := m.assets[assetID]
	if exists {
		return &asset, nil
	}
	return nil, derrors.NewNotFoundError(assetID)
}

// List the EIC in a given organization
func (m *MockupAssetProvider) List(organizationID string) ([]entities.Asset, derrors.Error){
	m.Lock()
	defer m.Unlock()
	result := make([]entities.Asset, 0)
	for _, a := range m.assets{
		if a.OrganizationId == organizationID{
			result = append(result, a)
		}
	}
	return result, nil
}

func (m * MockupAssetProvider) Remove(assetID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(assetID){
		return derrors.NewNotFoundError(assetID)
	}
	delete(m.assets, assetID)
	return nil
}

func (m * MockupAssetProvider) Clear() derrors.Error {
	m.Lock()
	m.assets = make(map[string]entities.Asset, 0)
	m.Unlock()
	return nil
}


