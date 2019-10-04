/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package application_network

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for the application networking instances.
type Provider interface {
	// AddConnectionInstance Adds a new connection between applications.
	AddConnectionInstance(connectionInstance entities.ConnectionInstance) derrors.Error
	// ExistsConnectionInstance Checks if the connection instance exists on the system.
	ExistsConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) (bool, derrors.Error)
	// GetConnectionInstance Retrieve the connection instance using organizationId, sourceInstanceId, and targetInstanceId.
	GetConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) (*entities.ConnectionInstance, derrors.Error)
	// ListConnectionInstances Lists all the connection instances.
	ListConnectionInstances(organizationId string) ([]entities.ConnectionInstance, derrors.Error)
	// ListInboundConnections retrieve all the connections where instance is the target
	ListInboundConnections(organizationId string, appInstanceId string) ([]entities.ConnectionInstance, derrors.Error)
	// ListOutboundConnections retrieve all the connections where instance is the source
	ListOutboundConnections(organizationId string, appInstanceId string) ([]entities.ConnectionInstance, derrors.Error)
	// UpdateConnectionInstance Updates a connection instance
	UpdateConnectionInstance(connectionInstance entities.ConnectionInstance) derrors.Error
	// RemoveConnectionInstance Removes a connection from the system
	RemoveConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) derrors.Error

	// AddConnectionInstanceLink Adds a new connection between applications.
	AddConnectionInstanceLink(connectionInstanceLink entities.ConnectionInstanceLink) derrors.Error
	// ExistsConnectionInstanceLink Checks if the connection instance exists on the system.
	ExistsConnectionInstanceLink(organizationId string, sourceInstanceId string, targetInstanceId string, sourceClusterId string, targetClusterId string, inboundName string, outboundName string) (bool, derrors.Error)
	// GetConnectionInstanceLink Retrieve the connection instance.
	GetConnectionInstanceLink(organizationId string, sourceInstanceId string, targetInstanceId string, sourceClusterId string, targetClusterId string, inboundName string, outboundName string) (*entities.ConnectionInstanceLink, derrors.Error)
	// ListConnectionInstanceLinks Lists all the connection instance links of one connection instance.
	ListConnectionInstanceLinks(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) ([]entities.ConnectionInstanceLink, derrors.Error)
	// RemoveConnectionInstanceLinks Removes all connection links from a connection instance.
	RemoveConnectionInstanceLinks(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) derrors.Error

	AddZTConnection(ztConnection entities.ZTNetworkConnection) derrors.Error
	ExistsZTConnection(organizationId string, networkId string, appInstanceId string) (bool, derrors.Error)
	GetZTConnection(organizationId string, networkId string, appInstanceId string)(*entities.ZTNetworkConnection, derrors.Error)
	UpdateZTConnection(ztConnection entities.ZTNetworkConnection) derrors.Error
	ListZTConnections(organizationId string, networkId string) ([]entities.ZTNetworkConnection, derrors.Error)
	RemoveZTConnection(organizationId string, networkId string)derrors.Error


	// clear the connections information
	Clear() derrors.Error
}
