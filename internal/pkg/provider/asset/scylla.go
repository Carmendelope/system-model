/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package asset

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/scylladb"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"sync"
)

// AssetTable with the name of the table that stores asset information.
const AssetTable = "Asset"
// AssetTablePK with the name of the primary key for the asset table.
const AssetTablePK = "asset_id"
// AllAssetColumns contains the name of all the columns in the asset table.
var allAssetColumns = []string{"organization_id", "edge_controller_id", "asset_id", "agent_id", "show",
	"created", "labels", "os", "hardware", "storage", "eic_net_ip", "last_alive_timestamp", "last_op_result"}
// AllAssetColumnsNoPK contains the name of all the columns in the asset table except the PK.
var allAssetColumnsNoPK = []string{"organization_id", "edge_controller_id", "agent_id", "show",
	"created", "labels", "os", "hardware", "storage", "eic_net_ip", "last_alive_timestamp", "last_op_result"}

type ScyllaAssetProvider struct {
	scylladb.ScyllaDB
	sync.Mutex
}

func NewScyllaAssetProvider(address string, port int, keyspace string) * ScyllaAssetProvider{
	provider := ScyllaAssetProvider{
		ScyllaDB : scylladb.ScyllaDB{
			Address: address,
			Port : port,
			Keyspace: keyspace,
		},
	}
	provider.Connect()
	return &provider
}

// disconnect from the database
func (sp *ScyllaAssetProvider) Disconnect() {
	sp.Lock()
	defer sp.Unlock()
	sp.ScyllaDB.Disconnect()
}

func (sp *ScyllaAssetProvider) Add(asset entities.Asset) derrors.Error {
	sp.Lock()
	defer sp.Unlock()
	log.Debug().Interface("asset", asset).Msg("provider add asset")
	return sp.UnsafeAdd(AssetTable, AssetTablePK, asset.AssetId, allAssetColumns, asset)
}

func (sp *ScyllaAssetProvider) Update(asset entities.Asset) derrors.Error {
	sp.Lock()
	defer sp.Unlock()
	return sp.UnsafeUpdate(AssetTable, AssetTablePK, asset.AssetId, allAssetColumnsNoPK, asset)
}

func (sp *ScyllaAssetProvider) Exists(assetID string) (bool, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()
	return sp.UnsafeGenericExist(AssetTable, AssetTablePK, assetID)
}

func (sp *ScyllaAssetProvider) Get(assetID string) (*entities.Asset, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()
	var result interface{} = &entities.Asset{}
	err := sp.UnsafeGet(AssetTable, AssetTablePK, assetID, allAssetColumns, &result)
	if err != nil{
		return nil, err
	}
	return result.(*entities.Asset), nil
}

func (sp *ScyllaAssetProvider) List(organizationID string) ([]entities.Asset, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select(AssetTable).Columns(allAssetColumns...).Where(qb.Eq("organization_id")).ToCql()
	q:= gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
	})

	assets := make ([]entities.Asset, 0)
	cqlErr := q.SelectRelease(&assets)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list assets")
	}

	return assets, nil
}

// ListControllerAssets retrieves the assets associated with a given edge controller
func (sp *ScyllaAssetProvider) ListControllerAssets(edgeControllerID string) ([]entities.Asset, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}
	stmt, names := qb.Select(AssetTable).Columns(allAssetColumns...).Where(qb.Eq("edge_controller_id")).ToCql()
	q:= gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"edge_controller_id": edgeControllerID,
	})

	assets := make ([]entities.Asset, 0)
	cqlErr := q.SelectRelease(&assets)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list assets")
	}

	return assets, nil
}


func (sp *ScyllaAssetProvider) Remove(assetID string) derrors.Error {
	sp.Lock()
	defer sp.Unlock()
	return sp.UnsafeRemove(AssetTable, AssetTablePK, assetID)
}

func (sp *ScyllaAssetProvider) Clear() derrors.Error {
	sp.Lock()
	defer sp.Unlock()
	return sp.UnsafeClear([]string{AssetTable})
}
