/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package cluster

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/cluster"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
)

type Manager struct {
	OrgProvider organization.Provider
	ClusterProvider cluster.Provider
}

func NewManager(orgProvider organization.Provider, clusterProvider cluster.Provider) Manager {
	return Manager{orgProvider, clusterProvider}
}

func (m * Manager) AddCluster(addClusterRequest *grpc_infrastructure_go.AddClusterRequest) (*entities.Cluster, derrors.Error) {
	exists := m.OrgProvider.Exists(addClusterRequest.OrganizationId)
	if !exists{
		return nil, derrors.NewNotFoundError("organizationID").WithParams(addClusterRequest.OrganizationId)
	}
	toAdd := entities.NewClusterFromGRPC(addClusterRequest)
	err := m.ClusterProvider.Add(*toAdd)
	if err != nil {
		return nil, err
	}
	err = m.OrgProvider.AddCluster(toAdd.OrganizationId, toAdd.ClusterId)
	if err != nil {
		return nil, err
	}

	return toAdd, nil
}

func (m * Manager) GetCluster(clusterID *grpc_infrastructure_go.ClusterId) (*entities.Cluster, derrors.Error) {
	if ! m.OrgProvider.Exists(clusterID.OrganizationId){
		return nil, derrors.NewNotFoundError("organizationID").WithParams(clusterID.OrganizationId)
	}

	if !m.OrgProvider.ClusterExists(clusterID.OrganizationId, clusterID.ClusterId){
		return nil, derrors.NewNotFoundError("clusterID").WithParams(clusterID.OrganizationId, clusterID.ClusterId)
	}
	return m.ClusterProvider.Get(clusterID.ClusterId)
}

func (m * Manager) ListClusters(organizationID *grpc_organization_go.OrganizationId) ([] entities.Cluster, derrors.Error) {
	if !m.OrgProvider.Exists(organizationID.OrganizationId){
		return nil, derrors.NewNotFoundError("organizationID").WithParams(organizationID.OrganizationId)
	}
	clusters, err := m.OrgProvider.ListClusters(organizationID.OrganizationId)
	if err != nil {
		return nil, err
	}
	result := make([] entities.Cluster, 0)
	for _, cID := range clusters {
		toAdd, err := m.ClusterProvider.Get(cID)
		if err != nil {
			return nil, err
		}
		result = append(result, *toAdd)
	}
	return result, nil
}

func (m * Manager) RemoveCluster(removeClusterRequest *grpc_infrastructure_go.RemoveClusterRequest) derrors.Error {
	if ! m.OrgProvider.Exists(removeClusterRequest.OrganizationId){
		return derrors.NewNotFoundError("organizationID").WithParams(removeClusterRequest.OrganizationId)
	}

	if !m.OrgProvider.ClusterExists(removeClusterRequest.OrganizationId, removeClusterRequest.ClusterId){
		return derrors.NewNotFoundError("clusterID").WithParams(removeClusterRequest.OrganizationId, removeClusterRequest.ClusterId)
	}

	err := m.OrgProvider.DeleteCluster(removeClusterRequest.OrganizationId, removeClusterRequest.ClusterId)
	if err != nil {
		return err
	}
	return m.ClusterProvider.Remove(removeClusterRequest.ClusterId)
}