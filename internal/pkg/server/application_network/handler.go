/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package application_network

import (
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-application-network-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type Handler struct {
	Manager Manager
}

func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h *Handler) AddConnection(ctx context.Context, addConnectionRequest *grpc_application_network_go.AddConnectionRequest) (*grpc_application_network_go.ConnectionInstance, error) {
	if err := entities.ValidAddConnectionRequest(addConnectionRequest); err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Interface("addConnectionRequest", addConnectionRequest).Msg("Adding connection instance")
	added, err := h.Manager.AddConnectionInstance(addConnectionRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return added.ToGRPC(), nil
}

func (h *Handler) RemoveConnection(ctx context.Context, removeConnectionRequest *grpc_application_network_go.RemoveConnectionRequest) (*grpc_common_go.Success, error) {
	if err := entities.ValidRemoveConnectionRequest(removeConnectionRequest); err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Interface("removeConnectionRequest", removeConnectionRequest).Msg("Removing connection instance")
	if err := h.Manager.RemoveConnectionInstance(removeConnectionRequest); err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

func (h *Handler) ListConnections(_ context.Context, orgID *grpc_organization_go.OrganizationId) (*grpc_application_network_go.ConnectionInstanceList, error) {
	if err := entities.ValidOrganizationID(orgID); err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	list, err := h.Manager.ListConnectionInstances(orgID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	grpcArray := make([]*grpc_application_network_go.ConnectionInstance, 0, len(list))
	for _, instance := range list {
		grpcArray = append(grpcArray, instance.ToGRPC())
	}
	result := &grpc_application_network_go.ConnectionInstanceList{
		Connections: grpcArray,
	}
	return result, nil
}

// ListInboundConnections retrieves a list with all the connections where the appInstanceId is the target
func (h *Handler) ListInboundConnections(_ context.Context, appInstanceID *grpc_application_go.AppInstanceId) (*grpc_application_network_go.ConnectionInstanceList, error){
	err := entities.ValidAppInstanceId(appInstanceID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	list, err := h.Manager.ListInboundConnections(appInstanceID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	connList := make([]*grpc_application_network_go.ConnectionInstance, 0, len(list))
	for _, instance := range list {
		connList = append(connList, instance.ToGRPC())
	}
	return &grpc_application_network_go.ConnectionInstanceList{
		Connections: connList,
	}, nil

}

// ListOutboundConnections retrieves a list with all the connections where the appInstanceId is the source
func (h *Handler) ListOutboundConnections(_ context.Context, appInstanceID *grpc_application_go.AppInstanceId) (*grpc_application_network_go.ConnectionInstanceList, error){
	err := entities.ValidAppInstanceId(appInstanceID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	list, err := h.Manager.ListOutboundConnections(appInstanceID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	connList := make([]*grpc_application_network_go.ConnectionInstance, 0, len(list))
	for _, instance := range list {
		connList = append(connList, instance.ToGRPC())
	}
	return &grpc_application_network_go.ConnectionInstanceList{
		Connections: connList,
	}, nil
}
