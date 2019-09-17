/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package application_network

import (
	"github.com/nalej/derrors"
	"github.com/nalej/scylladb-utils/pkg/scylladb"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"sync"
)

const (
	ConnectionInstanceTable = "Connection_Instances"
	ConnectionInstanceIdIx  = "connection_id"

	ConnectionInsanceLinkTable = "Connection_Instance_Links"
)

var (
	ConnectionInstanceColumns = []string{
		"organization_id",
		"connection_id",
		"source_instance_id",
		"source_instance_name",
		"target_instance_id",
		"target_instance_name",
		"inbound_name",
		"outbound_name",
		"outbound_required",
	}
	//ConnectionInstanceLinkColumns = []string{
	//	"organization_id",
	//	"connection_id",
	//	"source_instance_id",
	//	"source_cluster_id",
	//	"target_instance_id",
	//	"target_cluster_id",
	//	"inbound_name",
	//	"outbound_name",
	//}
)

func (sap *ScyllaApplicationNetworkProvider) createConnectionInsancePkMap(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) map[string]interface{} {
	return map[string]interface{}{
		"organization_id":    organizationId,
		"source_instance_id": sourceInstanceId,
		"target_instance_id": targetInstanceId,
		"inbound_name":       inboundName,
		"outbound_name":      outboundName,
	}
}

func (sap *ScyllaApplicationNetworkProvider) createConnectionInstanceLinkPkMap(connectionId string, sourceClusterId string, targetClusterId string) map[string]interface{} {
	return map[string]interface{}{
		"connection_id":     connectionId,
		"source_cluster_id": sourceClusterId,
		"target_cluster_id": targetClusterId,
	}
}

type ScyllaApplicationNetworkProvider struct {
	sync.Mutex
	scylladb.ScyllaDB
}

func NewScyllaApplicationNetworkProvider(address string, port int, keyspace string) *ScyllaApplicationNetworkProvider {
	provider := ScyllaApplicationNetworkProvider{
		ScyllaDB: scylladb.ScyllaDB{
			Address:  address,
			Port:     port,
			Keyspace: keyspace,
		},
	}
	_ = provider.Connect()
	return &provider
}

func (sap *ScyllaApplicationNetworkProvider) Disconnect() {
	sap.Lock()
	defer sap.Unlock()
	sap.ScyllaDB.Disconnect()
}

func (sap *ScyllaApplicationNetworkProvider) AddConnectionInstance(connectionInstance entities.ConnectionInstance) derrors.Error {
	sap.Lock()
	defer sap.Unlock()
	return sap.UnsafeAdd(ConnectionInstanceTable, ConnectionInstanceIdIx, connectionInstance.ConnectionId, ConnectionInstanceColumns, connectionInstance)
}

func (sap *ScyllaApplicationNetworkProvider) ExistsConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) (bool, derrors.Error) {
	sap.Lock()
	defer sap.Unlock()
	pkComposite := sap.createConnectionInsancePkMap(organizationId, sourceInstanceId, targetInstanceId, inboundName, outboundName)
	return sap.UnsafeGenericCompositeExist(ConnectionInstanceTable, pkComposite)
}

func (sap *ScyllaApplicationNetworkProvider) ExistsConnectionInstanceById(connectionId string) (bool, derrors.Error) {
	sap.Lock()
	defer sap.Unlock()
	return sap.UnsafeGenericExist(ConnectionInstanceTable, ConnectionInstanceIdIx, connectionId)
}

func (sap *ScyllaApplicationNetworkProvider) GetConnectionInstanceById(connectionId string) (*entities.ConnectionInstance, derrors.Error) {
	sap.Lock()
	defer sap.Unlock()
	result := interface{}(&entities.ConnectionInstance{})
	if err := sap.UnsafeGet(ConnectionInstanceTable, ConnectionInstanceIdIx, connectionId, ConnectionInstanceColumns, &result); err != nil {
		return nil, err
	}
	return result.(*entities.ConnectionInstance), nil
}

func (sap *ScyllaApplicationNetworkProvider) GetConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) (*entities.ConnectionInstance, derrors.Error) {
	sap.Lock()
	defer sap.Unlock()
	pkComposite := sap.createConnectionInsancePkMap(organizationId, sourceInstanceId, targetInstanceId, inboundName, outboundName)
	result := interface{}(&entities.ConnectionInstance{})
	if err := sap.UnsafeCompositeGet(ConnectionInstanceTable, pkComposite, ConnectionInstanceColumns, &result); err != nil {
		return nil, err
	}
	return result.(*entities.ConnectionInstance), nil
}

func (sap *ScyllaApplicationNetworkProvider) ListConnectionInstances(organizationId string) ([]entities.ConnectionInstance, derrors.Error) {
	sap.Lock()
	defer sap.Unlock()

	if err := sap.CheckAndConnect(); err != nil {
		return nil, err
	}

	filterColumn := "organization_id"
	stmt, names := qb.Select(ConnectionInstanceTable).Columns(ConnectionInstanceColumns...).Where(qb.Eq(filterColumn)).ToCql()
	q := gocqlx.Query(sap.Session.Query(stmt), names).BindMap(qb.M{
		filterColumn: organizationId,
	})

	connections := make([]entities.ConnectionInstance, 0)
	if qerr := q.SelectRelease(&connections); qerr != nil {
		return nil, derrors.AsError(qerr, "cannot list connection instances")
	}

	return connections, nil
}

func (sap *ScyllaApplicationNetworkProvider) RemoveConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) derrors.Error {
	sap.Lock()
	defer sap.Unlock()
	pkComposite := sap.createConnectionInsancePkMap(organizationId, sourceInstanceId, targetInstanceId, inboundName, outboundName)
	return sap.UnsafeCompositeRemove(ConnectionInstanceTable, pkComposite)
}

// Connection Instance Link
// ------------------------
/*
func (sap *ScyllaApplicationNetworkProvider) AddConnectionInstanceLink(connectionInstanceLink entities.ConnectionInstanceLink) derrors.Error {
	sap.Lock()
	defer sap.Unlock()
	pkComposite := sap.createConnectionInstanceLinkPkMap(connectionInstanceLink.ConnectionId, connectionInstanceLink.SourceClusterId, connectionInstanceLink.TargetClusterId)
	return sap.UnsafeCompositeAdd(ConnectionInsanceLinkTable, pkComposite, ConnectionInstanceLinkColumns, connectionInstanceLink)
}

func (sap *ScyllaApplicationNetworkProvider) ExistsConnectionInstanceLink(connectionId string, sourceClusterId string, targetClusterId string) (bool, derrors.Error) {
	sap.Lock()
	defer sap.Unlock()
	pkComposite := sap.createConnectionInstanceLinkPkMap(connectionId, sourceClusterId, targetClusterId)
	return sap.UnsafeGenericCompositeExist(ConnectionInsanceLinkTable, pkComposite)
}

func (sap *ScyllaApplicationNetworkProvider) GetConnectionInstanceLink(connectionId string, sourceClusterId string, targetClusterId string) (*entities.ConnectionInstanceLink, derrors.Error) {
	sap.Lock()
	defer sap.Unlock()
	pkComposite := sap.createConnectionInstanceLinkPkMap(connectionId, sourceClusterId, targetClusterId)
	result := interface{}(&entities.ConnectionInstanceLink{})
	if err := sap.UnsafeCompositeGet(ConnectionInsanceLinkTable, pkComposite, ConnectionInstanceLinkColumns, &result); err != nil {
		return nil, err
	}
	return result.(*entities.ConnectionInstanceLink), nil
}

func (sap *ScyllaApplicationNetworkProvider) ListConnectionInstanceLinks(connectionId string) ([]entities.ConnectionInstanceLink, derrors.Error) {
	sap.Lock()
	defer sap.Unlock()

	if err := sap.CheckAndConnect(); err != nil {
		return nil, err
	}

	filterColumn := "connection_id"
	stmt, names := qb.Select(ConnectionInsanceLinkTable).Columns(ConnectionInstanceLinkColumns...).Where(qb.Eq(filterColumn)).ToCql()
	q := gocqlx.Query(sap.Session.Query(stmt), names).BindMap(qb.M{
		filterColumn: connectionId,
	})

	connectionInstanceLinks := make([]entities.ConnectionInstanceLink, 0)
	if qerr := q.SelectRelease(&connectionInstanceLinks); qerr != nil {
		return nil, derrors.AsError(qerr, "cannot list connection instance links")
	}
	return connectionInstanceLinks, nil
}

func (sap *ScyllaApplicationNetworkProvider) RemoveConnectionInstanceLinks(connectionId string) derrors.Error {
	sap.Lock()
	defer sap.Unlock()

	if err := sap.CheckAndConnect(); err != nil {
		return err
	}

	filterColumn := "connection_id"
	stmt, names := qb.Delete(ConnectionInsanceLinkTable).Where(qb.Eq(filterColumn)).ToCql()
	q := gocqlx.Query(sap.Session.Query(stmt), names).BindMap(qb.M{
		filterColumn: connectionId,
	})

	if qerr := q.ExecRelease(); qerr != nil {
		return derrors.AsError(qerr, "cannot delete connection instance links")
	}
	return nil
}
*/

func (sap *ScyllaApplicationNetworkProvider) Clear() derrors.Error {
	sap.Lock()
	defer sap.Unlock()

	if err := sap.UnsafeClear([]string{ConnectionInstanceTable, ConnectionInsanceLinkTable}); err != nil {
		return err
	}
	return nil
}
