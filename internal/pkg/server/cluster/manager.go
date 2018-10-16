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
	panic("implement me")
}

func (m * Manager) GetCluster(clusterID *grpc_infrastructure_go.ClusterId) (*entities.Cluster, derrors.Error) {
	panic("implement me")
}

func (m * Manager) ListClusters(organizationID *grpc_organization_go.OrganizationId) ([] entities.Cluster, derrors.Error) {
	panic("implement me")
}

func (m * Manager) RemoveCluster(removeClusterRequest *grpc_infrastructure_go.RemoveClusterRequest) derrors.Error {
	panic("implement me")
}