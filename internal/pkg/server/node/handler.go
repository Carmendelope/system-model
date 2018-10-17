/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package node

import (
	"context"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
)

type Handler struct {
	Manager Manager
}

func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}

func (h *Handler) AddNode(ctx context.Context, addNodeRequest *grpc_infrastructure_go.AddNodeRequest) (*grpc_infrastructure_go.Node, error) {
	err := entities.ValidAddNodeRequest(addNodeRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	added, err := h.Manager.AddNode(addNodeRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return added.ToGRPC(), nil
}

func (h *Handler) AttachNode(ctx context.Context, attachNodeRequest *grpc_infrastructure_go.AttachNodeRequest) (*grpc_common_go.Success, error) {
	err := entities.ValidAttachNodeRequest(attachNodeRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.AttachNode(attachNodeRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

func (h *Handler) ListNodes(ctx context.Context, clusterID *grpc_infrastructure_go.ClusterId) (*grpc_infrastructure_go.NodeList, error) {
	err := entities.ValidClusterID(clusterID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	nodes, err := h.Manager.ListNodes(clusterID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	toReturn := make([]*grpc_infrastructure_go.Node, 0)
	for _, n := range nodes {
		toReturn = append(toReturn, n.ToGRPC())
	}
	result := &grpc_infrastructure_go.NodeList{
		Nodes:          toReturn,
	}
	return result, nil
}

func (h *Handler) RemoveNodes(ctx context.Context, removeNodesRequest *grpc_infrastructure_go.RemoveNodesRequest) (*grpc_common_go.Success, error) {
	err := entities.ValidRemoveNodesRequest(removeNodesRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.RemoveNodes(removeNodesRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}
