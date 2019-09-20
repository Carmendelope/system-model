/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package application_network

import (
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
	connectionInstances map[string]entities.ConnectionInstance
	// connectionInstanceLinks derived from a connectionInstance.
	connectionInstanceLinks map[string][]entities.ConnectionInstanceLink
}

// NewMockupApplicationNetworkProvider Create a new mockup provider for the application network domain.
func NewMockupApplicationNetworkProvider() *MockupApplicationNetworkProvider {
	return &MockupApplicationNetworkProvider{
		connectionInstances:     make(map[string]entities.ConnectionInstance, 0),
		connectionInstanceLinks: make(map[string][]entities.ConnectionInstanceLink, 0),
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
	if !m.unsafeExistsConnectionInstance(compositePK) {
		m.connectionInstances[compositePK] = toAdd
		return nil
	}
	return derrors.NewAlreadyExistsError(toAdd.ConnectionId)
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

// GetConnectionInstanceById Retrieves a connection instance using connectionId.
func (m *MockupApplicationNetworkProvider) GetConnectionInstanceById(connectionId string) (*entities.ConnectionInstance, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	instance, exists := m.connectionInstances[connectionId]
	if exists {
		return &instance, nil
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
			ret = append(ret, instance)
		}
	}
	return ret, nil
}

func (m *MockupApplicationNetworkProvider) RemoveConnectionInstance(organizationId string, sourceInstanceId string, targetInstanceId string, inboundName string, outboundName string) derrors.Error {
	compositePK := getCompositePK(organizationId, sourceInstanceId, targetInstanceId, inboundName, outboundName)
	if m.unsafeExistsConnectionInstance(compositePK) {
		delete(m.connectionInstances, compositePK)
		return nil
	}
	return derrors.NewNotFoundError("connectionInstance").WithParams(compositePK)
}

// ListInboundConnections retrieve all the connections where instance is the target
func (m *MockupApplicationNetworkProvider) ListInboundConnections(organizationId string, appInstanceId string)([]entities.ConnectionInstance, derrors.Error){
	m.Lock()
	defer m.Unlock()
	ret := make([]entities.ConnectionInstance, 0)
	for _, instance := range m.connectionInstances {
		if instance.OrganizationId == organizationId && instance.TargetInstanceId == appInstanceId{
			ret = append(ret, instance)
		}
	}
	return ret, nil
}
// ListOutboundConnections retrieve all the connections where instance is the source
func (m *MockupApplicationNetworkProvider) ListOutboundConnections(organizationId string, appInstanceId string)([]entities.ConnectionInstance, derrors.Error){
	m.Lock()
	defer m.Unlock()
	ret := make([]entities.ConnectionInstance, 0)
	for _, instance := range m.connectionInstances {
		if instance.OrganizationId == organizationId && instance.SourceInstanceId == appInstanceId{
			ret = append(ret, instance)
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

func (m *MockupApplicationNetworkProvider) Clear() derrors.Error {
	m.Lock()
	defer m.Unlock()
	m.connectionInstances = make(map[string]entities.ConnectionInstance, 0)
	m.connectionInstanceLinks = make(map[string][]entities.ConnectionInstanceLink, 0)
	return nil
}
