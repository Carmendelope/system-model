/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
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
}

func NewMockupClusterProvider() * MockupClusterProvider {
	return &MockupClusterProvider{
		clusters: make(map[string]entities.Cluster, 0),
	}
}

func (m * MockupClusterProvider) unsafeExists(clusterID string) bool {
	_, exists := m.clusters[clusterID]
	return exists
}


func (m * MockupClusterProvider) Add(cluster entities.Cluster) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(cluster.ClusterId){
		m.clusters[cluster.ClusterId] = cluster
		return nil
	}
	return derrors.NewAlreadyExistsError(cluster.ClusterId)
}

func (m * MockupClusterProvider) Exists(clusterID string) bool {
	m.Lock()
	defer m.Unlock()
	return m.unsafeExists(clusterID)
}

func (m * MockupClusterProvider) Get(clusterID string) (*entities.Cluster, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	cluster, exists := m.clusters[clusterID]
	if exists {
		return &cluster, nil
	}
	return nil, derrors.NewNotFoundError(clusterID)
}

// Remove a cluster
func (m * MockupClusterProvider) Remove(clusterID string) derrors.Error {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExists(clusterID){
		return derrors.NewNotFoundError(clusterID)
	}
	delete(m.clusters, clusterID)
	return nil
}

func (m * MockupClusterProvider) Clear() {
	m.Lock()
	m.clusters = make(map[string]entities.Cluster, 0)
	m.Unlock()
}
