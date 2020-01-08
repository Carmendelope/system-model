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
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/cluster"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/rs/zerolog/log"
)

// Manager structure with the required providers for cluster operations.
type Manager struct {
	OrgProvider     organization.Provider
	ClusterProvider cluster.Provider
}

// NewManager creates a Manager using a set of providers.
func NewManager(orgProvider organization.Provider, clusterProvider cluster.Provider) Manager {
	return Manager{orgProvider, clusterProvider}
}

// AddCluster adds a new cluster to the system.
func (m *Manager) AddCluster(addClusterRequest *grpc_infrastructure_go.AddClusterRequest) (*entities.Cluster, derrors.Error) {
	exists, err := m.OrgProvider.Exists(addClusterRequest.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(addClusterRequest.OrganizationId)
	}
	toAdd := entities.NewClusterFromGRPC(addClusterRequest)
	err = m.ClusterProvider.Add(*toAdd)
	if err != nil {
		return nil, err
	}
	err = m.OrgProvider.AddCluster(toAdd.OrganizationId, toAdd.ClusterId)
	if err != nil {
		return nil, err
	}

	return toAdd, nil
}

func (m *Manager) UpdateCluster(updateRequest *grpc_infrastructure_go.UpdateClusterRequest) (*entities.Cluster, derrors.Error) {
	exists, err := m.OrgProvider.Exists(updateRequest.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(updateRequest.OrganizationId)
	}

	exists, err = m.OrgProvider.ClusterExists(updateRequest.OrganizationId, updateRequest.ClusterId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("clusterID").WithParams(updateRequest.OrganizationId, updateRequest.ClusterId)
	}
	old, err := m.ClusterProvider.Get(updateRequest.ClusterId)
	if err != nil {
		return nil, err
	}
	old.ApplyUpdate(*updateRequest)
	err = m.ClusterProvider.Update(*old)
	if err != nil {
		return nil, err
	}
	return old, nil
}

// GetCluster retrieves the cluster information.
func (m *Manager) GetCluster(clusterID *grpc_infrastructure_go.ClusterId) (*entities.Cluster, derrors.Error) {
	exists, err := m.OrgProvider.Exists(clusterID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(clusterID.OrganizationId)
	}

	exists, err = m.OrgProvider.ClusterExists(clusterID.OrganizationId, clusterID.ClusterId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("clusterID").WithParams(clusterID.OrganizationId, clusterID.ClusterId)
	}
	return m.ClusterProvider.Get(clusterID.ClusterId)
}

// ListClusters obtains a list of the clusters in the organization.
func (m *Manager) ListClusters(organizationID *grpc_organization_go.OrganizationId) ([]entities.Cluster, derrors.Error) {
	exists, err := m.OrgProvider.Exists(organizationID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(organizationID.OrganizationId)
	}
	clusters, err := m.OrgProvider.ListClusters(organizationID.OrganizationId)
	if err != nil {
		return nil, err
	}
	result := make([]entities.Cluster, 0)
	for _, cID := range clusters {
		toAdd, err := m.ClusterProvider.Get(cID)
		if err != nil {
			return nil, err
		}
		result = append(result, *toAdd)
	}
	return result, nil
}

// RemoveCluster removes a cluster from an organization. Notice that removing a cluster implies draining the cluster
// of running applications.
func (m *Manager) RemoveCluster(removeClusterRequest *grpc_infrastructure_go.RemoveClusterRequest) derrors.Error {
	exists, err := m.OrgProvider.Exists(removeClusterRequest.OrganizationId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organizationID").WithParams(removeClusterRequest.OrganizationId)
	}

	exists, err = m.OrgProvider.ClusterExists(removeClusterRequest.OrganizationId, removeClusterRequest.ClusterId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("clusterID").WithParams(removeClusterRequest.OrganizationId, removeClusterRequest.ClusterId)
	}

	err = m.OrgProvider.DeleteCluster(removeClusterRequest.OrganizationId, removeClusterRequest.ClusterId)
	if err != nil {
		return err
	}
	err = m.ClusterProvider.Remove(removeClusterRequest.ClusterId)
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("Error removing cluster. Rollback!")
		rollbackError := m.OrgProvider.AddCluster(removeClusterRequest.OrganizationId, removeClusterRequest.ClusterId)
		if rollbackError != nil {
			log.Error().Str("trace", conversions.ToDerror(rollbackError).DebugReport()).
				Str("removeClusterRequest.OrganizationId", removeClusterRequest.OrganizationId).
				Str("removeClusterRequest.ClusterId", removeClusterRequest.ClusterId).
				Msg("error in Rollback")
		}
	}
	return err
}

func (m *Manager) CordonCluster(clusterID *grpc_infrastructure_go.ClusterId) derrors.Error {
	exists, err := m.OrgProvider.Exists(clusterID.OrganizationId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organizationID").WithParams(clusterID.OrganizationId)
	}

	exists, err = m.OrgProvider.ClusterExists(clusterID.OrganizationId, clusterID.ClusterId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("clusterID").WithParams(clusterID.OrganizationId, clusterID.ClusterId)
	}

	old, err := m.ClusterProvider.Get(clusterID.ClusterId)
	if err != nil {
		return err
	}
	if old.Status == entities.ClusterStatusUnknown {
		return derrors.NewFailedPreconditionError("cannot cordon cluster with unknown state")
	}
	// Check the the cordon status
	if old.Status == entities.ClusterStatusOfflineCordon || old.Status == entities.ClusterStatusOnlineCordon {
		return nil
	}

	newStatus := entities.ClusterStatusUnknown
	switch old.Status {
	case entities.ClusterStatusOnline:
		{
			newStatus = entities.ClusterStatusOnlineCordon
			break
		}
	case entities.ClusterStatusOffline:
		{
			newStatus = entities.ClusterStatusOfflineCordon
			break
		}
	}

	// this is going to be cordoned
	old.Status = newStatus
	err = m.ClusterProvider.Update(*old)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UncordonCluster(clusterID *grpc_infrastructure_go.ClusterId) derrors.Error {
	exists, err := m.OrgProvider.Exists(clusterID.OrganizationId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organizationID").WithParams(clusterID.OrganizationId)
	}

	exists, err = m.OrgProvider.ClusterExists(clusterID.OrganizationId, clusterID.ClusterId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("clusterID").WithParams(clusterID.OrganizationId, clusterID.ClusterId)
	}

	old, err := m.ClusterProvider.Get(clusterID.ClusterId)
	if err != nil {
		return err
	}
	if old.Status == entities.ClusterStatusUnknown {
		return derrors.NewFailedPreconditionError("cannot uncordon cluster with unknown state")
	}
	// Check the the cordon status for already uncordoned
	if old.Status == entities.ClusterStatusOffline || old.Status == entities.ClusterStatusOnline {
		return nil
	}

	newStatus := entities.ClusterStatusUnknown
	switch old.Status {
	case entities.ClusterStatusOnlineCordon:
		{
			newStatus = entities.ClusterStatusOnline
			break
		}
	case entities.ClusterStatusOfflineCordon:
		{
			newStatus = entities.ClusterStatusOffline
			break
		}
	}
	old.Status = newStatus
	err = m.ClusterProvider.Update(*old)
	if err != nil {
		return err
	}

	return nil
}
