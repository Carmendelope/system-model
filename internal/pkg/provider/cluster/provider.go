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
	// Update an existing cluster in the system
	Update(cluster entities.Cluster) derrors.Error
	// Exists checks if a cluster exists on the system.
	Exists(clusterID string) (bool, derrors.Error)
	// Get a cluster.
	Get(clusterID string) (*entities.Cluster, derrors.Error)
	// Remove a cluster
	Remove(clusterID string) derrors.Error

	// AddNode adds a new node ID to the cluster.
	AddNode(clusterID string, nodeID string) derrors.Error
	// NodeExists checks if a node is linked to a cluster.
	NodeExists(clusterID string, nodeID string) (bool, derrors.Error)
	// ListNodes returns a list of nodes in a cluster.
	ListNodes(clusterID string) ([]string, derrors.Error)
	// DeleteNode removes a node from a cluster.
	DeleteNode(clusterID string, nodeID string) derrors.Error
	// clear the cluster information
	Clear() derrors.Error
}
