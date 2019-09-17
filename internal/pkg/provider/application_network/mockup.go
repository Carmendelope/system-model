/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package application_network

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

type MockupApplicationNetworkProvider struct {
	sync.Mutex
	// connectionInstances indexed by composite PK.
	connectionInstancesPkToId map[string]string
	// connectionInstancesById indexed by connectionId.
	connectionInstancesById map[string]entities.ConnectionInstance
	// connectionInstanceLinks derived from a connectionInstance.
	//connectionInstanceLinks map[string][]entities.ConnectionInstanceLink
}

// NewMockupApplicationNetworkProvider Create a new mockup provider for the application network domain.
func NewMockupApplicationNetworkProvider() *MockupApplicationNetworkProvider {
	return &MockupApplicationNetworkProvider{
		connectionInstancesPkToId: make(map[string]string, 0),
		connectionInstancesById:   make(map[string]entities.ConnectionInstance, 0),
		//connectionInstanceLinks:   make(map[string][]entities.ConnectionInstanceLink, 0),
	}
}

// unsafeExistsConnectionInstance Checks the existence of the connection instance without locking the DB.
func (m *MockupApplicationNetworkProvider) unsafeExistsConnectionInstance(connectionId string) bool {
	_, exists := m.connectionInstancesById[connectionId]
	return exists
}

// unsafeExistsLink Checks the existence of the connection instance link without locking the DB.
//func (m *MockupApplicationNetworkProvider) unsafeExistsLink(connectionId string, sourceClusterId string, targetClusterId string) bool {
//	links, exists := m.connectionInstanceLinks[connectionId]
//	if exists {
//		for _, link := range links {
//			if link.SourceClusterId == sourceClusterId && link.TargetClusterId == targetClusterId {
//				return true
//			}
//		}
//	}
//	return false
//}

// AddConnectionInstance Adds a new ConnectionInstance to the system.
func (m *MockupApplicationNetworkProvider) AddConnectionInstance(toAdd entities.ConnectionInstance) derrors.Error {
	m.Lock()
	defer m.Unlock()
	compositePK := toAdd.OrganizationId + toAdd.SourceInstanceId + toAdd.TargetInstanceId + toAdd.InboundName + toAdd.OutboundName
	m.connectionInstancesPkToId[compositePK] = toAdd.ConnectionId
	if !m.unsafeExistsConnectionInstance(toAdd.ConnectionId) {
		m.connectionInstancesById[toAdd.ConnectionId] = toAdd
		return nil
	}
	return derrors.NewAlreadyExistsError(toAdd.ConnectionId)
}

// ExistsConnectionInstance Checks the existence of the connection instance using organizationId, sourceInstanceId, targetInstanceId, inboundName, and outboundName.
func (m *MockupApplicationNetworkProvider) ExistsConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) (bool, derrors.Error) {
	compositePK := organizationId + sourceInstanceId + targetInstanceId + inboundName + outboundName
	return m.ExistsConnectionInstanceById(m.connectionInstancesPkToId[compositePK])
}

// ExistsConnectionInstanceById Checks the existence of the connection instance using connectionId.
func (m *MockupApplicationNetworkProvider) ExistsConnectionInstanceById(connectionId string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsConnectionInstance(connectionId), nil
}

// GetConnectionInstance Retrieves a connection instance using organizationId, sourceInstanceId, targetInstanceId, inboundName, and outboundName.
func (m *MockupApplicationNetworkProvider) GetConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) (*entities.ConnectionInstance, derrors.Error) {
	compositePK := organizationId + sourceInstanceId + targetInstanceId + inboundName + outboundName
	return m.GetConnectionInstanceById(m.connectionInstancesPkToId[compositePK])
}

// GetConnectionInstanceById Retrieves a connection instance using connectionId.
func (m *MockupApplicationNetworkProvider) GetConnectionInstanceById(connectionId string) (*entities.ConnectionInstance, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	instance, exists := m.connectionInstancesById[connectionId]
	if exists {
		return &instance, nil
	}
	return nil, derrors.NewNotFoundError(connectionId)
}

// ListConnectionInstances Retrieves a list with all the connection instances of an organization unsing OrganizationID
func (m *MockupApplicationNetworkProvider) ListConnectionInstances(organizationId string) ([]entities.ConnectionInstance, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	ret := make([]entities.ConnectionInstance, 0)
	for _, instance := range m.connectionInstancesById {
		if instance.OrganizationId == organizationId {
			ret = append(ret, instance)
		}
	}
	return ret, nil
}

func (m *MockupApplicationNetworkProvider) RemoveConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) derrors.Error {
	compositePK := organizationId + sourceInstanceId + targetInstanceId + inboundName + outboundName
	connectionId := m.connectionInstancesPkToId[compositePK]
	if m.unsafeExistsConnectionInstance(connectionId) {
		delete(m.connectionInstancesById, connectionId)
		return nil
	}
	return derrors.NewNotFoundError(connectionId)
}

// Connection Instance Link
// ------------------------
/*
// AddConnectionInstanceLink Inserts a new connection instance link in the DB
func (m *MockupApplicationNetworkProvider) AddConnectionInstanceLink(link entities.ConnectionInstanceLink) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExistsConnectionInstance(link.ConnectionId) {
		if ! m.unsafeExistsLink(link.ConnectionId, link.SourceClusterId, link.TargetClusterId) {
			m.connectionInstanceLinks[link.ConnectionId] = append(m.connectionInstanceLinks[link.ConnectionId], link)
			return nil
		}
		return derrors.NewAlreadyExistsError("ConnectionInstanceLink").WithParams(link.ConnectionId, link.SourceClusterId, link.TargetClusterId)
	}
	return derrors.NewNotFoundError("ConnectionInstance").WithParams(link.ConnectionId)
}

// ExistsConnectionInstanceLink Checks the existence of the connection instance link
func (m *MockupApplicationNetworkProvider) ExistsConnectionInstanceLink(connectionId string, sourceClusterId string, targetClusterId string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsLink(connectionId, sourceClusterId, targetClusterId), nil
}

// GetConnectionInstanceLink Retrieves a connection instance link
func (m *MockupApplicationNetworkProvider) GetConnectionInstanceLink(connectionId string, sourceClusterId string, targetClusterId string) (*entities.ConnectionInstanceLink, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExistsConnectionInstance(connectionId) {
		if m.unsafeExistsLink(connectionId, sourceClusterId, targetClusterId) {
			var foundLink entities.ConnectionInstanceLink
			for _, link := range m.connectionInstanceLinks[connectionId] {
				if link.SourceClusterId == sourceClusterId && link.TargetClusterId == targetClusterId {
					foundLink = link
					break
				}
			}
			return &foundLink, nil
		}
		return nil, derrors.NewNotFoundError("ConnectionInstanceLink").WithParams(connectionId, sourceClusterId, targetClusterId)
	}
	return nil, derrors.NewNotFoundError("ConnectionInstance").WithParams(connectionId)
}

// ListConnectionInstanceLinks Retrieves a list with all the links from a connection instance
func (m *MockupApplicationNetworkProvider) ListConnectionInstanceLinks(connectionId string) ([]entities.ConnectionInstanceLink, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExistsConnectionInstance(connectionId) {
		return m.connectionInstanceLinks[connectionId], nil
	}
	return nil, derrors.NewNotFoundError("ConnectionInstance").WithParams(connectionId)
}

// RemoveConnectionInstanceLinks Removes all the links from a connection instance
func (m *MockupApplicationNetworkProvider) RemoveConnectionInstanceLinks(connectionId string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExistsConnectionInstance(connectionId) {
		delete(m.connectionInstanceLinks, connectionId)
		return nil
	}
	return derrors.NewNotFoundError("ConnectionInstance").WithParams(connectionId)

}
*/

func (m *MockupApplicationNetworkProvider) Clear() derrors.Error {
	m.Lock()
	defer m.Unlock()
	m.connectionInstancesById = make(map[string]entities.ConnectionInstance, 0)
	//m.connectionInstanceLinks = make(map[string][]entities.ConnectionInstanceLink, 0)
	return nil
}
