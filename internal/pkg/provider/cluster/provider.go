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
	Add(org entities.Cluster) derrors.Error
	// Check if an organization exists on the system.
	Exists(clusterID string) bool
	// Get an organization.
	Get(clusterID string) (* entities.Cluster, derrors.Error)
}