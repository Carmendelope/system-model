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
		"status",
		"ip_range",
	}
	ConnectionInstanceColumnsNoPK = []string{
		"connection_id",
		"source_instance_name",
		"target_instance_name",
		"outbound_required",
		"status",
		"ip_range",
	}
	ConnectionInstanceLinkColumns = []string{
		"organization_id",
		"connection_id",
		"source_instance_id",
		"source_cluster_id",
		"target_instance_id",
		"target_cluster_id",
		"inbound_name",
		"outbound_name",
		"status",
	}
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

func (sap *ScyllaApplicationNetworkProvider) createConnectionInstanceLinkPkMap(organizationId string, sourceInstanceId string, targetInstanceId string, sourceClusterId string, targetClusterId string, inboundName string, outboundName string) map[string]interface{} {
	return map[string]interface{}{
		"organization_id":    organizationId,
		"source_instance_id": sourceInstanceId,
		"target_instance_id": targetInstanceId,
		"source_cluster_id":  sourceClusterId,
		"target_cluster_id":  targetClusterId,
		"inbound_name":       inboundName,
		"outbound_name":      outboundName,
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
	pkComposite := sap.createConnectionInsancePkMap(connectionInstance.OrganizationId, connectionInstance.SourceInstanceId, connectionInstance.TargetInstanceId, connectionInstance.InboundName, connectionInstance.OutboundName)
	return sap.UnsafeCompositeAdd(ConnectionInstanceTable, pkComposite, ConnectionInstanceColumns, connectionInstance)
}

func (sap *ScyllaApplicationNetworkProvider) UpdateConnectionInstance(connectionInstance entities.ConnectionInstance) derrors.Error {
	sap.Lock()
	defer sap.Unlock()
	pkComposite := sap.createConnectionInsancePkMap(connectionInstance.OrganizationId, connectionInstance.SourceInstanceId, connectionInstance.TargetInstanceId, connectionInstance.InboundName, connectionInstance.OutboundName)
	return sap.UnsafeCompositeUpdate(ConnectionInstanceTable, pkComposite, ConnectionInstanceColumnsNoPK, connectionInstance)
}

func (sap *ScyllaApplicationNetworkProvider) ExistsConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) (bool, derrors.Error) {
	sap.Lock()
	defer sap.Unlock()
	pkComposite := sap.createConnectionInsancePkMap(organizationId, sourceInstanceId, targetInstanceId, inboundName, outboundName)
	return sap.UnsafeGenericCompositeExist(ConnectionInstanceTable, pkComposite)
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

// ListInboundConnections retrieve all the connections where instance is the target
func (sap *ScyllaApplicationNetworkProvider) ListInboundConnections(organizationId string, appInstanceId string) ([]entities.ConnectionInstance, derrors.Error) {
	sap.Lock()
	defer sap.Unlock()

	if err := sap.CheckAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select(ConnectionInstanceTable).Columns(ConnectionInstanceColumns...).Where(qb.Eq("organization_id")).
		Where(qb.Eq("target_instance_id")).ToCql()
	q := gocqlx.Query(sap.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id":    organizationId,
		"target_instance_id": appInstanceId,
	})

	connections := make([]entities.ConnectionInstance, 0)
	cqlErr := q.SelectRelease(&connections)
	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list inbound connections")
	}

	return connections, nil
}

// ListOutboundConnections retrieve all the connections where instance is the source
func (sap *ScyllaApplicationNetworkProvider) ListOutboundConnections(organizationId string, appInstanceId string) ([]entities.ConnectionInstance, derrors.Error) {
	sap.Lock()
	defer sap.Unlock()

	if err := sap.CheckAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select(ConnectionInstanceTable).Columns(ConnectionInstanceColumns...).Where(qb.Eq("organization_id")).
		Where(qb.Eq("source_instance_id")).ToCql()
	q := gocqlx.Query(sap.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id":    organizationId,
		"source_instance_id": appInstanceId,
	})

	connections := make([]entities.ConnectionInstance, 0)
	cqlErr := q.SelectRelease(&connections)
	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list outbound connections")
	}

	return connections, nil
}

// Connection Instance Link
// ------------------------
func (sap *ScyllaApplicationNetworkProvider) AddConnectionInstanceLink(connectionInstanceLink entities.ConnectionInstanceLink) derrors.Error {
	sap.Lock()
	defer sap.Unlock()
	pkComposite := sap.createConnectionInstanceLinkPkMap(
		connectionInstanceLink.OrganizationId,
		connectionInstanceLink.SourceInstanceId,
		connectionInstanceLink.TargetInstanceId,
		connectionInstanceLink.SourceClusterId,
		connectionInstanceLink.TargetClusterId,
		connectionInstanceLink.InboundName,
		connectionInstanceLink.OutboundName,
	)
	return sap.UnsafeCompositeAdd(ConnectionInsanceLinkTable, pkComposite, ConnectionInstanceLinkColumns, connectionInstanceLink)
}

func (sap *ScyllaApplicationNetworkProvider) ExistsConnectionInstanceLink(organizationId string, sourceInstanceId string, targetInstanceId string, sourceClusterId string, targetClusterId string, inboundName string, outboundName string) (bool, derrors.Error) {
	sap.Lock()
	defer sap.Unlock()
	pkComposite := sap.createConnectionInstanceLinkPkMap(organizationId, sourceInstanceId, targetInstanceId, sourceClusterId, targetClusterId, inboundName, outboundName)
	return sap.UnsafeGenericCompositeExist(ConnectionInsanceLinkTable, pkComposite)
}

func (sap *ScyllaApplicationNetworkProvider) GetConnectionInstanceLink(organizationId string, sourceInstanceId string, targetInstanceId string, sourceClusterId string, targetClusterId string, inboundName string, outboundName string) (*entities.ConnectionInstanceLink, derrors.Error) {
	sap.Lock()
	defer sap.Unlock()
	pkComposite := sap.createConnectionInstanceLinkPkMap(organizationId, sourceInstanceId, targetInstanceId, sourceClusterId, targetClusterId, inboundName, outboundName)
	result := interface{}(&entities.ConnectionInstanceLink{})
	if err := sap.UnsafeCompositeGet(ConnectionInsanceLinkTable, pkComposite, ConnectionInstanceLinkColumns, &result); err != nil {
		return nil, err
	}
	return result.(*entities.ConnectionInstanceLink), nil
}

func (sap *ScyllaApplicationNetworkProvider) ListConnectionInstanceLinks(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) ([]entities.ConnectionInstanceLink, derrors.Error) {
	sap.Lock()
	defer sap.Unlock()

	if err := sap.CheckAndConnect(); err != nil {
		return nil, err
	}

	pkMap := sap.createConnectionInsancePkMap(organizationId, sourceInstanceId, targetInstanceId, inboundName, outboundName)
	var whereClause []qb.Cmp
	for column := range pkMap {
		whereClause = append(whereClause, qb.Eq(column))
	}
	stmt, names := qb.Select(ConnectionInsanceLinkTable).Columns(ConnectionInstanceLinkColumns...).Where(whereClause...).ToCql()
	q := gocqlx.Query(sap.Session.Query(stmt), names).BindMap(pkMap)

	connectionInstanceLinks := make([]entities.ConnectionInstanceLink, 0)
	if qerr := q.SelectRelease(&connectionInstanceLinks); qerr != nil {
		return nil, derrors.AsError(qerr, "cannot list connection instance links")
	}
	return connectionInstanceLinks, nil
}

func (sap *ScyllaApplicationNetworkProvider) RemoveConnectionInstanceLinks(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) derrors.Error {
	sap.Lock()
	defer sap.Unlock()

	if err := sap.CheckAndConnect(); err != nil {
		return err
	}

	pkMap := sap.createConnectionInsancePkMap(organizationId, sourceInstanceId, targetInstanceId, inboundName, outboundName)
	var whereClause []qb.Cmp
	for column := range pkMap {
		whereClause = append(whereClause, qb.Eq(column))
	}
	stmt, names := qb.Delete(ConnectionInsanceLinkTable).Where(whereClause...).ToCql()
	q := gocqlx.Query(sap.Session.Query(stmt), names).BindMap(pkMap)

	if qerr := q.ExecRelease(); qerr != nil {
		return derrors.AsError(qerr, "cannot delete connection instance links")
	}
	return nil
}

func (sap *ScyllaApplicationNetworkProvider) Clear() derrors.Error {
	sap.Lock()
	defer sap.Unlock()

	if err := sap.UnsafeClear([]string{ConnectionInstanceTable, ConnectionInsanceLinkTable}); err != nil {
		return err
	}
	return nil
}
