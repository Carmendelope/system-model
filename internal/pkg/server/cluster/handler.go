/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package cluster

import (
	"context"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
)
// Handler structure for the cluster requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}

// AddCluster adds a new cluster to the system.
func (h * Handler) AddCluster(ctx context.Context, addClusterRequest *grpc_infrastructure_go.AddClusterRequest) (*grpc_infrastructure_go.Cluster, error) {
	log.Debug().Str("organizationID", addClusterRequest.OrganizationId).
		Str("name", addClusterRequest.Name).
		Str("hostname", addClusterRequest.Hostname).Msg("add cluster")
	err := entities.ValidAddClusterRequest(addClusterRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	cluster, err := h.Manager.AddCluster(addClusterRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("clusterID", cluster.ClusterId).Msg("cluster has been added")
	return cluster.ToGRPC(), nil
}

// UpdateCluster updates the information of a cluster.
func (h * Handler) UpdateCluster(ctx context.Context, updateClusterRequest *grpc_infrastructure_go.UpdateClusterRequest) (*grpc_infrastructure_go.Cluster, error){
	log.Debug().Str("organizationID", updateClusterRequest.OrganizationId).
		Str("clusterID", updateClusterRequest.ClusterId).Msg("update cluster")
	err := entities.ValidUpdateClusterRequest(updateClusterRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	cluster, err := h.Manager.UpdateCluster(updateClusterRequest)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("clusterID", cluster.ClusterId).Msg("cluster has been updated")
	return cluster.ToGRPC(), nil
}

// GetCluster retrieves the cluster information.
func (h * Handler) GetCluster(ctx context.Context, clusterID *grpc_infrastructure_go.ClusterId) (*grpc_infrastructure_go.Cluster, error) {
	err := entities.ValidClusterID(clusterID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	cluster, err := h.Manager.GetCluster(clusterID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return cluster.ToGRPC(), nil
}

// ListClusters obtains a list of the clusters in the organization.
func (h * Handler) ListClusters(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_infrastructure_go.ClusterList, error) {
	err := entities.ValidOrganizationID(organizationID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	clusters, err := h.Manager.ListClusters(organizationID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	toReturn := make([]*grpc_infrastructure_go.Cluster, 0)
	for _, c := range clusters {
		toReturn = append(toReturn, c.ToGRPC())
	}
	result := &grpc_infrastructure_go.ClusterList{
		Clusters:          toReturn,
	}
	return result, nil
}

// RemoveCluster removes a cluster from an organization. Notice that removing a cluster implies draining the cluster
// of running applications.
func (h * Handler) RemoveCluster(ctx context.Context, removeClusterRequest *grpc_infrastructure_go.RemoveClusterRequest) (*grpc_common_go.Success, error) {
	err := entities.ValidRemoveClusterRequest(removeClusterRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.RemoveCluster(removeClusterRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}


// Cordon a cluster. The cluster will not accept any new application deployment request.
func (h * Handler) CordonCluster(ctx context.Context, clusterID *grpc_infrastructure_go.ClusterId) (*grpc_common_go.Success, error) {
	err := entities.ValidClusterID(clusterID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.CordonCluster(clusterID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}
// Uncordon a cluster. The cordon flag will be disabled for this cluster.
func (h * Handler) UncordonCluster(ctx context.Context, clusterID *grpc_infrastructure_go.ClusterId) (*grpc_common_go.Success, error) {
	err := entities.ValidClusterID(clusterID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.UncordonCluster(clusterID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

