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
}

func NewMockupClusterProvider() * MockupClusterProvider {
	return &MockupClusterProvider{
	}
}

func (m * MockupClusterProvider) Add(org entities.Cluster) derrors.Error {
	panic("implement me")
}

func (m * MockupClusterProvider) Exists(clusterID string) bool {
	panic("implement me")
}

func (m * MockupClusterProvider) Get(clusterID string) (*entities.Cluster, derrors.Error) {
	panic("implement me")
}

func (m * MockupClusterProvider) Clear() {
	m.Lock()
	m.Unlock()
}
