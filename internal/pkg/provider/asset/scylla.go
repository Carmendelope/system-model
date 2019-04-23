/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package asset

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
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
var allAssetColumns = []string{"organization_id", "asset_id", "agent_id", "show",
	"created", "labels", "os", "hardware", "storage", "eic_net_ip"}

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

func (sp *ScyllaAssetProvider) unsafeExists(assetID string) (bool, derrors.Error) {

	var count int

	stmt, names := qb.Select(AssetTable).CountAll().Where(qb.Eq(AssetTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		AssetTablePK: assetID})

	err := q.GetRelease(&count)
	if err != nil {
		if err.Error() == scylladb.RowNotFoundMsg {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if asset exists")
		}
	}

	return count == 1, nil
}


func (sp *ScyllaAssetProvider) Add(asset entities.Asset) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return err
	}
	exists, err := sp.unsafeExists(asset.AssetId)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError(asset.AssetId)
	}

	// insert the cluster instance
	stmt, names := qb.Insert(AssetTable).Columns(allAssetColumns...).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(asset)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add asset")
	}

	return nil
}

func (sp *ScyllaAssetProvider) Update(asset entities.Asset) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.CheckAndConnect()
	if err != nil {
		return err
	}

	exists, err := sp.unsafeExists(asset.AssetId)
	if err != nil {
		return err
	}
	if ! exists {
		return derrors.NewNotFoundError(asset.AssetId)
	}

	// insert the cluster instance
	stmt, names := qb.Update(AssetTable).Set(
		"organization_id", "agent_id", "show",
		"created", "labels", "os", "hardware", "storage", "eic_net_ip").
		Where(qb.Eq(AssetTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(asset)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(err,"cannot update asset")
	}

	return nil
}

func (sp *ScyllaAssetProvider) Exists(assetID string) (bool, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	err := sp.CheckAndConnect()
	if err != nil {
		return false, err
	}
	return sp.unsafeExists(assetID)
}

func (sp *ScyllaAssetProvider) Get(assetID string) (*entities.Asset, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	var asset entities.Asset
	stmt, names := qb.Select(AssetTable).Columns(allAssetColumns...).Where(qb.Eq(AssetTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		AssetTablePK: assetID,
	})

	err := q.GetRelease(&asset)
	if err != nil {
		if err.Error() == scylladb.RowNotFoundMsg {
			return nil, derrors.NewNotFoundError("asset").WithParams(assetID)
		} else {
			return nil, derrors.AsError(err, "cannot get asset")
		}
	}

	return &asset, nil
}

func (sp *ScyllaAssetProvider) Remove(assetID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.CheckAndConnect(); err != nil {
		return err
	}

	// check if the asset exists
	exists, err := sp.unsafeExists(assetID)
	if err != nil {
		return err
	}
	if ! exists {
		return derrors.NewNotFoundError(assetID)
	}

	// delete cluster instance
	stmt, _ := qb.Delete(AssetTable).Where(qb.Eq(AssetTablePK)).ToCql()
	cqlErr := sp.Session.Query(stmt, assetID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot remove asset")
	}
	return nil
}

func (sp *ScyllaAssetProvider) Clear() derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return err
	}

	// delete clusters table
	err := sp.Session.Query("TRUNCATE TABLE asset").Exec()
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("failed to truncate the asset table")
		return derrors.AsError(err, "cannot truncate asset table")
	}

	return nil
}
