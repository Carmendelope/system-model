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
 *
 */

package application_history_logs

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

func getCompositePK(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) string {
	return organizationId + sourceInstanceId + targetInstanceId + inboundName + outboundName
}

type MockupApplicationNetworkProvider struct {
	sync.Mutex
	// connectionInstances indexed by composite PK organizationId + sourceInstanceId + targetInstanceId + inboundName + outboundName.
	connectionInstances map[string]*entities.ConnectionInstance
	// connectionInstances indexed by zt_networkId
	connectionInstancesByNetwork map[string][]entities.ConnectionInstance
	// connectionInstanceLinks derived from a connectionInstance.
	connectionInstanceLinks map[string][]entities.ConnectionInstanceLink
	// ztNetworkConnections map of ZTNetworkConnection indexed by organizationId+ztnetworkId+appInstanceId+serviceId+clusterId
	ztNetworkConnections map[string]entities.ZTNetworkConnection
}

// NewMockupApplicationNetworkProvider Create a new mockup provider for the application network domain.
func NewMockupApplicationNetworkProvider() *MockupApplicationNetworkProvider {
	return &MockupApplicationNetworkProvider{
		connectionInstances:          make(map[string]*entities.ConnectionInstance, 0),
		connectionInstancesByNetwork: make(map[string][]entities.ConnectionInstance, 0),
		connectionInstanceLinks:      make(map[string][]entities.ConnectionInstanceLink, 0),
		ztNetworkConnections:         make(map[string]entities.ZTNetworkConnection, 0),
	}
}

// unsafeExistsConnectionInstance Checks the existence of the connection instance without locking the DB.
func (m *MockupApplicationNetworkProvider) unsafeExistsConnectionInstance(compositePK string) bool {
	_, exists := m.connectionInstances[compositePK]
	return exists
}

// unsafeExistsLink Checks the existence of the connection instance link without locking the DB.
func (m *MockupApplicationNetworkProvider) unsafeExistsLink(organizationId string, sourceInstanceId string, targetInstanceId string, sourceClusterId string, targetClusterId string, inboundName string, outboundName string) bool {
	compositePK := getCompositePK(organizationId, sourceInstanceId, targetInstanceId, inboundName, outboundName)
	links, exists := m.connectionInstanceLinks[compositePK]
	if exists {
		for _, link := range links {
			if link.OrganizationId == organizationId &&
				link.SourceInstanceId == sourceInstanceId &&
				link.TargetInstanceId == targetInstanceId &&
				link.SourceClusterId == sourceClusterId &&
				link.TargetClusterId == targetClusterId &&
				link.InboundName == inboundName &&
				link.OutboundName == outboundName {
				return true
			}
		}
	}
	return false
}

// AddConnectionInstance Adds a new ConnectionInstance to the system.
func (m *MockupApplicationNetworkProvider) AddConnectionInstance(toAdd entities.ConnectionInstance) derrors.Error {
	m.Lock()
	defer m.Unlock()
	compositePK := getCompositePK(toAdd.OrganizationId, toAdd.SourceInstanceId, toAdd.TargetInstanceId, toAdd.InboundName, toAdd.OutboundName)
	if m.unsafeExistsConnectionInstance(compositePK) {
		return derrors.NewAlreadyExistsError("connection instance").WithParams(toAdd.ConnectionId)
	}
	m.connectionInstances[compositePK] = &toAdd

	// connectionInstancesByNetwork
	list, exists := m.connectionInstancesByNetwork[toAdd.ZtNetworkId]
	if exists {
		m.connectionInstancesByNetwork[toAdd.ZtNetworkId] = append(list, toAdd)
	} else {
		m.connectionInstancesByNetwork[toAdd.ZtNetworkId] = []entities.ConnectionInstance{toAdd}
	}

	return nil
}

func (m *MockupApplicationNetworkProvider) UpdateConnectionInstance(toUpdate entities.ConnectionInstance) derrors.Error {
	m.Lock()
	defer m.Unlock()
	compositePK := getCompositePK(toUpdate.OrganizationId, toUpdate.SourceInstanceId, toUpdate.TargetInstanceId, toUpdate.InboundName, toUpdate.OutboundName)
	old, exists := m.connectionInstances[compositePK]
	if !exists {
		return derrors.NewNotFoundError(compositePK)
	}

	m.connectionInstances[compositePK] = &toUpdate

	// connectionInstancesByNetwork
	// delete the old entry
	var newList []entities.ConnectionInstance
	list, exists := m.connectionInstancesByNetwork[old.ZtNetworkId]
	delete(m.connectionInstancesByNetwork, old.ZtNetworkId)
	if exists { // it should exist
		for _, conn := range list {
			if conn.OrganizationId != toUpdate.OrganizationId || conn.SourceInstanceId != toUpdate.SourceInstanceId ||
				conn.TargetInstanceId != toUpdate.TargetInstanceId || conn.InboundName != toUpdate.InboundName ||
				conn.OutboundName != toUpdate.OutboundName || conn.ConnectionId != toUpdate.ConnectionId {
				newList = append(newList, conn)
			}
		}
	}
	if len(newList) > 0 {
		m.connectionInstancesByNetwork[old.ZtNetworkId] = newList
	}

	// add the new one
	list, exists = m.connectionInstancesByNetwork[toUpdate.ZtNetworkId]
	if exists {
		m.connectionInstancesByNetwork[toUpdate.ZtNetworkId] = append(list, toUpdate)
	} else {
		m.connectionInstancesByNetwork[toUpdate.ZtNetworkId] = []entities.ConnectionInstance{toUpdate}
	}

	return nil
}

// ExistsConnectionInstance Checks the existence of the connection instance using organizationId, sourceInstanceId, targetInstanceId, inboundName, and outboundName.
func (m *MockupApplicationNetworkProvider) ExistsConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	compositePK := getCompositePK(organizationId, sourceInstanceId, targetInstanceId, inboundName, outboundName)
	return m.unsafeExistsConnectionInstance(compositePK), nil
}

// GetConnectionInstance Retrieves a connection instance using organizationId, sourceInstanceId, targetInstanceId, inboundName, and outboundName.
func (m *MockupApplicationNetworkProvider) GetConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) (*entities.ConnectionInstance, derrors.Error) {
	compositePK := getCompositePK(organizationId, sourceInstanceId, targetInstanceId, inboundName, outboundName)
	return m.GetConnectionInstanceById(compositePK)
}

func (m *MockupApplicationNetworkProvider) GetConnectionByZtNetworkId(ztNetworkId string) ([]entities.ConnectionInstance, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	instance, exists := m.connectionInstancesByNetwork[ztNetworkId]
	if exists {
		return instance, nil
	}
	return nil, derrors.NewNotFoundError(ztNetworkId)
}

// GetConnectionInstanceById Retrieves a connection instance using connectionId.
func (m *MockupApplicationNetworkProvider) GetConnectionInstanceById(connectionId string) (*entities.ConnectionInstance, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	instance, exists := m.connectionInstances[connectionId]
	if exists {
		return instance, nil
	}
	return nil, derrors.NewNotFoundError(connectionId)
}

// ListConnectionInstances Retrieves a list with all the connection instances of an organization using OrganizationID
func (m *MockupApplicationNetworkProvider) ListConnectionInstances(organizationId string) ([]entities.ConnectionInstance, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	ret := make([]entities.ConnectionInstance, 0)
	for _, instance := range m.connectionInstances {
		if instance.OrganizationId == organizationId {
			ret = append(ret, *instance)
		}
	}
	return ret, nil
}

func (m *MockupApplicationNetworkProvider) RemoveConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) derrors.Error {
	compositePK := getCompositePK(organizationId, sourceInstanceId, targetInstanceId, inboundName, outboundName)
	if m.unsafeExistsConnectionInstance(compositePK) {
		instance := m.connectionInstances[compositePK]
		delete(m.connectionInstances, compositePK)
		delete(m.connectionInstancesByNetwork, instance.ZtNetworkId)
		return nil
	}
	return derrors.NewNotFoundError("connectionInstance").WithParams(compositePK)
}

// ListInboundConnections retrieve all the connections where instance is the target
func (m *MockupApplicationNetworkProvider) ListInboundConnections(organizationId string, appInstanceId string) ([]entities.ConnectionInstance, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	ret := make([]entities.ConnectionInstance, 0)
	for _, instance := range m.connectionInstances {
		if instance.OrganizationId == organizationId && instance.TargetInstanceId == appInstanceId {
			ret = append(ret, *instance)
		}
	}
	return ret, nil
}

// ListOutboundConnections retrieve all the connections where instance is the source
func (m *MockupApplicationNetworkProvider) ListOutboundConnections(organizationId string, appInstanceId string) ([]entities.ConnectionInstance, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	ret := make([]entities.ConnectionInstance, 0)
	for _, instance := range m.connectionInstances {
		if instance.OrganizationId == organizationId && instance.SourceInstanceId == appInstanceId {
			ret = append(ret, *instance)
		}
	}
	return ret, nil
}

// Connection Instance Link
// ------------------------

// AddConnectionInstanceLink Inserts a new connection instance link in the DB
func (m *MockupApplicationNetworkProvider) AddConnectionInstanceLink(link entities.ConnectionInstanceLink) derrors.Error {
	m.Lock()
	defer m.Unlock()
	compositePK := getCompositePK(link.OrganizationId, link.SourceInstanceId, link.TargetInstanceId, link.InboundName, link.OutboundName)
	if m.unsafeExistsConnectionInstance(compositePK) {
		if !m.unsafeExistsLink(link.OrganizationId, link.SourceInstanceId, link.TargetInstanceId, link.SourceClusterId, link.TargetClusterId, link.InboundName, link.OutboundName) {
			m.connectionInstanceLinks[compositePK] = append(m.connectionInstanceLinks[compositePK], link)
			return nil
		}
		return derrors.NewAlreadyExistsError("ConnectionInstanceLink").WithParams(link.OrganizationId, link.SourceInstanceId, link.TargetInstanceId, link.SourceClusterId, link.TargetClusterId, link.InboundName, link.OutboundName)
	}
	return derrors.NewNotFoundError("ConnectionInstance").WithParams(link.OrganizationId, link.SourceInstanceId, link.TargetInstanceId, link.InboundName, link.OutboundName)
}

// ExistsConnectionInstanceLink Checks the existence of the connection instance link
func (m *MockupApplicationNetworkProvider) ExistsConnectionInstanceLink(organizationId string, sourceInstanceId string, targetInstanceId string, sourceClusterId string, targetClusterId string, inboundName string, outboundName string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsLink(organizationId, sourceInstanceId, targetInstanceId, sourceClusterId, targetClusterId, inboundName, outboundName), nil
}

// GetConnectionInstanceLink Retrieves a connection instance link
func (m *MockupApplicationNetworkProvider) GetConnectionInstanceLink(organizationId string, sourceInstanceId string, targetInstanceId string, sourceClusterId string, targetClusterId string, inboundName string, outboundName string) (*entities.ConnectionInstanceLink, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	compositePK := getCompositePK(organizationId, sourceInstanceId, targetInstanceId, inboundName, outboundName)
	if m.unsafeExistsConnectionInstance(compositePK) {
		if m.unsafeExistsLink(organizationId, sourceInstanceId, targetInstanceId, sourceClusterId, targetClusterId, inboundName, outboundName) {
			var foundLink entities.ConnectionInstanceLink
			for _, link := range m.connectionInstanceLinks[compositePK] {
				if link.SourceClusterId == sourceClusterId && link.TargetClusterId == targetClusterId {
					foundLink = link
					break
				}
			}
			return &foundLink, nil
		}
		return nil, derrors.NewNotFoundError("ConnectionInstanceLink").WithParams(organizationId, sourceClusterId, targetClusterId)
	}
	return nil, derrors.NewNotFoundError("ConnectionInstance").WithParams(organizationId)
}

// ListConnectionInstanceLinks Retrieves a list with all the links from a connection instance
func (m *MockupApplicationNetworkProvider) ListConnectionInstanceLinks(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) ([]entities.ConnectionInstanceLink, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	compositePK := getCompositePK(organizationId, sourceInstanceId, targetInstanceId, inboundName, outboundName)
	if m.unsafeExistsConnectionInstance(compositePK) {
		return m.connectionInstanceLinks[compositePK], nil
	}
	return nil, derrors.NewNotFoundError("ConnectionInstance").WithParams(organizationId)
}

// RemoveConnectionInstanceLinks Removes all the links from a connection instance
func (m *MockupApplicationNetworkProvider) RemoveConnectionInstanceLinks(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	compositePK := getCompositePK(organizationId, sourceInstanceId, targetInstanceId, inboundName, outboundName)
	if m.unsafeExistsConnectionInstance(compositePK) {
		delete(m.connectionInstanceLinks, compositePK)
		return nil
	}
	return derrors.NewNotFoundError("ConnectionInstanceLinks").WithParams(organizationId)

}

// ------------------ //
// -- ZTConnection -- //
// ------------------ //
func (m *MockupApplicationNetworkProvider) unsafeExistsZTConnection(pk string) bool {
	_, exists := m.ztNetworkConnections[pk]
	return exists
}

func (m *MockupApplicationNetworkProvider) getZTPk(organizationID string, networkId string, appInstanceId string, serviceId string, clusterId string) string {
	return fmt.Sprintf("%s%s%s%s%s", organizationID, networkId, appInstanceId, serviceId, clusterId)
}
func (m *MockupApplicationNetworkProvider) AddZTConnection(ztConnection entities.ZTNetworkConnection) derrors.Error {
	m.Lock()
	defer m.Unlock()
	pk := m.getZTPk(ztConnection.OrganizationId, ztConnection.ZtNetworkId, ztConnection.AppInstanceId, ztConnection.ServiceId, ztConnection.ClusterId)
	if m.unsafeExistsZTConnection(pk) {
		return derrors.NewAlreadyExistsError("ztNetwork")
	}
	m.ztNetworkConnections[pk] = ztConnection
	return nil
}

func (m *MockupApplicationNetworkProvider) ExistsZTConnection(organizationId string, networkId string, appInstanceId string, serviceId string, clusterId string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	pk := m.getZTPk(organizationId, networkId, appInstanceId, serviceId, clusterId)
	return m.unsafeExistsZTConnection(pk), nil
}

func (m *MockupApplicationNetworkProvider) GetZTConnection(organizationId string, networkId string, appInstanceId string, serviceId string, clusterId string) (*entities.ZTNetworkConnection, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	pk := m.getZTPk(organizationId, networkId, appInstanceId, serviceId, clusterId)
	zt, exists := m.ztNetworkConnections[pk]
	if !exists {
		return nil, derrors.NewNotFoundError("ztNetwork")
	}
	return &zt, nil
}

func (m *MockupApplicationNetworkProvider) ListZTConnections(organizationId string, networkId string) ([]entities.ZTNetworkConnection, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	list := make([]entities.ZTNetworkConnection, 0)
	for _, conn := range m.ztNetworkConnections {
		if conn.OrganizationId == organizationId && conn.ZtNetworkId == networkId {
			list = append(list, conn)
		}
	}
	return list, nil
}

func (m *MockupApplicationNetworkProvider) UpdateZTConnection(ztConnection entities.ZTNetworkConnection) derrors.Error {
	m.Lock()
	defer m.Unlock()
	pk := m.getZTPk(ztConnection.OrganizationId, ztConnection.ZtNetworkId, ztConnection.AppInstanceId, ztConnection.ServiceId, ztConnection.ClusterId)
	_, exists := m.ztNetworkConnections[pk]
	if !exists {
		return derrors.NewNotFoundError("ztNetwork")
	}
	m.ztNetworkConnections[pk] = ztConnection

	return nil
}

func (m *MockupApplicationNetworkProvider) RemoveZTConnection(organizationId string, networkId string, appInstanceId string, serviceId string, clusterId string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	pk := m.getZTPk(organizationId, networkId, appInstanceId, serviceId, clusterId)
	_, exists := m.ztNetworkConnections[pk]
	if !exists {
		return derrors.NewNotFoundError("ztNetwork")
	}
	delete(m.ztNetworkConnections, pk)
	return nil
}

func (m *MockupApplicationNetworkProvider) RemoveZTConnectionByNetworkId(organizationId string, networkId string) derrors.Error {
	m.Lock()
	defer m.Unlock()

	list := make([]entities.ZTNetworkConnection, 0)
	for _, conn := range m.ztNetworkConnections {
		if conn.OrganizationId == organizationId && conn.ZtNetworkId == networkId {
			list = append(list, conn)
		}
	}
	if len(list) == 0 {
		return derrors.NewNotFoundError("ztnetworkConnection").WithParams(organizationId, networkId)
	}
	for _, ztConnection := range list {
		pk := m.getZTPk(ztConnection.OrganizationId, ztConnection.ZtNetworkId, ztConnection.AppInstanceId, ztConnection.ServiceId, ztConnection.ClusterId)
		delete(m.ztNetworkConnections, pk)
	}

	return nil
}

func (m *MockupApplicationNetworkProvider) Clear() derrors.Error {
	m.Lock()
	defer m.Unlock()
	m.connectionInstances = make(map[string]*entities.ConnectionInstance, 0)
	m.connectionInstancesByNetwork = make(map[string][]entities.ConnectionInstance, 0)
	m.connectionInstanceLinks = make(map[string][]entities.ConnectionInstanceLink, 0)
	m.ztNetworkConnections = make(map[string]entities.ZTNetworkConnection, 0)
	return nil
}
