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
	// GetConnectionInstance Retrieve the connection instance using organizationId, sourceInstanceId, targetInstanceId, inboundName, and outboundName.
	GetConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) (*entities.ConnectionInstance, derrors.Error)
	// GetConnectionByZtNetworkId Retrieve the connection instance using organizationId, and ztNetworkId
	GetConnectionByZtNetworkId(ztNetworkId string) ([]entities.ConnectionInstance, derrors.Error)
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

	// AddZTConnection adds a new ZTConnection
	AddZTConnection(ztConnection entities.ZTNetworkConnection) derrors.Error
	// ExistsZTConnection checks if a ztConnection exists
	ExistsZTConnection(organizationId string, networkId string, appInstanceId string, serviceId string, clusterId string) (bool, derrors.Error)
	// GetZTConnection retrieve the ztConnection using organizationId, networkId, appInstanceId, serviceId and clusterId
	GetZTConnection(organizationId string, networkId string, appInstanceId string, serviceId string, clusterId string) (*entities.ZTNetworkConnection, derrors.Error)
	// UpdateZTConnection updates a ztConnection
	UpdateZTConnection(ztConnection entities.ZTNetworkConnection) derrors.Error
	// ListZTConnections retrieve all the zt connections of a zero tier network
	ListZTConnections(organizationId string, networkId string) ([]entities.ZTNetworkConnection, derrors.Error)
	// RemoveZTConnection removes one zt connections
	RemoveZTConnection(organizationId string, networkId string, appInstanceId string, serviceId string, clusterId string) derrors.Error
	// RemoveZTConnectionByNetworkId removes all the zt connections of a zero tier network
	RemoveZTConnectionByNetworkId(organizationId string, networkId string) derrors.Error

	// clear the connections information
	Clear() derrors.Error
}
