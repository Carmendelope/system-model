/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package cluster

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for application
type Provider interface {
	// Add a new cluster to the system.
	Add(cluster entities.Cluster) derrors.Error
	// Check if a cluster exists on the system.
	Exists(clusterID string) bool
	// Get a cluster.
	Get(clusterID string) (* entities.Cluster, derrors.Error)
	// Remove a cluster
	Remove(clusterID string) derrors.Error
}