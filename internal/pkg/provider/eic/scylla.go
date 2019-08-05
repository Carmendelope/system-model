/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package eic

import (
	"github.com/nalej/derrors"
	"github.com/nalej/scylladb-utils/pkg/scylladb"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"sync"
)

// ControllerTable with the name of the table that stores controller information.
const ControllerTable = "Controller"
// ControllerTablePK with the name of the primary key for the controller table.
const ControllerTablePK = "edge_controller_id"
// AllControllerColumns contains the name of all the columns in the controller table.
var allControllerColumns = []string{"organization_id", "edge_controller_id", "show",
	"created", "name", "labels", "last_alive_timestamp", "location", "os", "hardware", "storage", "last_op_result"}
// AllControllerColumnsNoPK contains the name of all the columns in the controller table except the PK.
var allControllerColumnsNoPK = []string{"organization_id", "show",
	"created", "name", "labels", "last_alive_timestamp", "location", "os", "hardware", "storage", "last_op_result"}

type ScyllaControllerProvider struct {
	scylladb.ScyllaDB
	sync.Mutex
}

func NewScyllaControllerProvider(address string, port int, keyspace string) * ScyllaControllerProvider{
	provider := ScyllaControllerProvider{
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
func (sp *ScyllaControllerProvider) Disconnect() {
	sp.Lock()
	defer sp.Unlock()
	sp.ScyllaDB.Disconnect()
}

func (sp *ScyllaControllerProvider) Add(eic entities.EdgeController) derrors.Error {
	sp.Lock()
	defer sp.Unlock()
	return sp.UnsafeAdd(ControllerTable, ControllerTablePK, eic.EdgeControllerId, allControllerColumns, eic)
}

func (sp *ScyllaControllerProvider) Update(eic entities.EdgeController) derrors.Error {
	sp.Lock()
	defer sp.Unlock()
	return sp.UnsafeUpdate(ControllerTable, ControllerTablePK, eic.EdgeControllerId, allControllerColumnsNoPK, eic)
}

func (sp *ScyllaControllerProvider) Exists(edgeControllerID string) (bool, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()
	return sp.UnsafeGenericExist(ControllerTable, ControllerTablePK, edgeControllerID)
}

func (sp *ScyllaControllerProvider) Get(edgeControllerID string) (*entities.EdgeController, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()
	var result interface{} = &entities.EdgeController{}
	err := sp.UnsafeGet(ControllerTable, ControllerTablePK, edgeControllerID, allControllerColumns, &result)
	if err != nil{
		return nil, err
	}
	return result.(*entities.EdgeController), nil
}

func (sp *ScyllaControllerProvider) List(organizationID string) ([]entities.EdgeController, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select(ControllerTable).Columns(allControllerColumns...).Where(qb.Eq("organization_id")).ToCql()
	q:= gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
	})

	controllers := make ([]entities.EdgeController, 0)
	cqlErr := q.SelectRelease(&controllers)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list controllers")
	}

	return controllers, nil
}

func (sp *ScyllaControllerProvider) Remove(edgeControllerID string) derrors.Error {
	sp.Lock()
	defer sp.Unlock()
	return sp.UnsafeRemove(ControllerTable, ControllerTablePK, edgeControllerID)
}

func (sp *ScyllaControllerProvider) Clear() derrors.Error {
	sp.Lock()
	defer sp.Unlock()
	return sp.UnsafeClear([]string{ControllerTable})
}