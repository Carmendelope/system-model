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
)

type Handler struct {
	Manager Manager
}

func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}

func (h * Handler) AddCluster(ctx context.Context, addClusterRequest *grpc_infrastructure_go.AddClusterRequest) (*grpc_infrastructure_go.Cluster, error) {
	err := entities.ValidAddClusterRequest(addClusterRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	cluster, err := h.Manager.AddCluster(addClusterRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return cluster.ToGRPC(), nil
}

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


