/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package node

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

type MockupNodeProvider struct {
	sync.Mutex
	// Nodes indexed by node identifier.
	nodes map[string]entities.Node
}

func NewMockupNodeProvider() * MockupNodeProvider {
	return &MockupNodeProvider{
		nodes: make(map[string]entities.Node, 0),
	}
}

func (m * MockupNodeProvider) unsafeExists(nodeID string) bool {
	_, exists := m.nodes[nodeID]
	return exists
}

// Add a new node to the system.
func (m * MockupNodeProvider) Add(node entities.Node) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(node.NodeId){
		m.nodes[node.NodeId] = node
		return nil
	}
	return derrors.NewAlreadyExistsError(node.NodeId)
}

// Update an existing node in the system
func (m * MockupNodeProvider) Update(node entities.Node) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(node.NodeId){
		return derrors.NewNotFoundError(node.NodeId)
	}
	m.nodes[node.NodeId] = node
	return nil
}

// Exists checks if a node exists on the system.
func (m * MockupNodeProvider) Exists(nodeID string) (bool,derrors.Error) {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExists(nodeID), nil
}

// Get a node.
func (m * MockupNodeProvider) Get(nodeID string) (*entities.Node, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	node, exists := m.nodes[nodeID]
	if exists {
		return &node, nil
	}
	return nil, derrors.NewNotFoundError(nodeID)
}

// Remove a node
func (m * MockupNodeProvider) Remove(nodeID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(nodeID){
		return derrors.NewNotFoundError(nodeID)
	}
	delete(m.nodes, nodeID)
	return nil
}

// Clear cleans the contents of the mockup.
func (m * MockupNodeProvider) Clear() derrors.Error{
	m.Lock()
	m.nodes = make(map[string]entities.Node, 0)
	m.Unlock()
	return nil
}
