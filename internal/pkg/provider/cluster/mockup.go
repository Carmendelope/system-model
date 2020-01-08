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

package cluster

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

type MockupClusterProvider struct {
	sync.Mutex
	// Clusters indexed by cluster identifier.
	clusters map[string]entities.Cluster
	// nodes attached to a cluster
	nodes map[string][]string
}

func NewMockupClusterProvider() *MockupClusterProvider {
	return &MockupClusterProvider{
		clusters: make(map[string]entities.Cluster, 0),
		nodes:    make(map[string][]string, 0),
	}
}

func (m *MockupClusterProvider) unsafeExists(clusterID string) bool {
	_, exists := m.clusters[clusterID]
	return exists
}

func (m *MockupClusterProvider) unsafeExistsNode(clusterID string, nodeID string) bool {
	nodeList, ok := m.nodes[clusterID]
	if ok {
		for _, nID := range nodeList {
			if nID == nodeID {
				return true
			}
		}
		return false
	}
	return false
}

// Add a new cluster to the system.
func (m *MockupClusterProvider) Add(cluster entities.Cluster) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(cluster.ClusterId) {
		m.clusters[cluster.ClusterId] = cluster
		return nil
	}
	return derrors.NewAlreadyExistsError(cluster.ClusterId)
}

// Update an existing cluster in the system
func (m *MockupClusterProvider) Update(cluster entities.Cluster) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(cluster.ClusterId) {
		return derrors.NewNotFoundError(cluster.ClusterId)
	}
	m.clusters[cluster.ClusterId] = cluster
	return nil
}

// Exists checks if a cluster exists on the system.
func (m *MockupClusterProvider) Exists(clusterID string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExists(clusterID), nil
}

// Get a cluster.
func (m *MockupClusterProvider) Get(clusterID string) (*entities.Cluster, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	cluster, exists := m.clusters[clusterID]
	if exists {
		return &cluster, nil
	}
	return nil, derrors.NewNotFoundError(clusterID)
}

// Remove a cluster
func (m *MockupClusterProvider) Remove(clusterID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(clusterID) {
		return derrors.NewNotFoundError(clusterID)
	}
	delete(m.clusters, clusterID)
	return nil
}

// AddNode adds a new node ID to the cluster.
func (m *MockupClusterProvider) AddNode(clusterID string, nodeID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExists(clusterID) {
		if !m.unsafeExistsNode(clusterID, nodeID) {
			nodeList, _ := m.nodes[clusterID]
			m.nodes[clusterID] = append(nodeList, nodeID)
			return nil
		}
		return derrors.NewAlreadyExistsError("node").WithParams(clusterID, nodeID)
	}
	return derrors.NewNotFoundError("cluster").WithParams(clusterID)
}

// NodeExists checks if a node is linked to a cluster.
func (m *MockupClusterProvider) NodeExists(clusterID string, nodeID string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExistsNode(clusterID, nodeID), nil
}

// ListNodes returns a list of nodes in a cluster.
func (m *MockupClusterProvider) ListNodes(clusterID string) ([]string, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	if !m.unsafeExists(clusterID) {
		return nil, derrors.NewNotFoundError("cluster").WithParams(clusterID)
	}

	nodeList, ok := m.nodes[clusterID]
	if ok {
		return nodeList, nil
	}
	return make([]string, 0), nil
}

// DeleteNode removes a node from a cluster.
func (m *MockupClusterProvider) DeleteNode(clusterID string, nodeID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if m.unsafeExistsNode(clusterID, nodeID) {
		previous := m.nodes[clusterID]
		newList := make([]string, 0, len(previous)-1)
		for _, id := range previous {
			if id != nodeID {
				newList = append(newList, id)
			}
		}
		m.nodes[clusterID] = newList
		return nil
	}
	return derrors.NewNotFoundError("node").WithParams(clusterID, nodeID)
}

// Clear cleans the contents of the mockup.
func (m *MockupClusterProvider) Clear() derrors.Error {
	m.Lock()
	m.clusters = make(map[string]entities.Cluster, 0)
	m.nodes = make(map[string][]string, 0)
	m.Unlock()
	return nil
}
