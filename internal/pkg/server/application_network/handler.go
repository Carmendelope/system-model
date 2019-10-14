/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package application_network

import (
	"context"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-application-network-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
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

func (h *Handler) UpdateConnection(ctx context.Context, updateConnectionRequest *grpc_application_network_go.UpdateConnectionRequest) (*grpc_common_go.Success, error) {
	if err := entities.ValidUpdateConnectionRequest(updateConnectionRequest); err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Interface("updateConnectionRequest", updateConnectionRequest).Msg("Updating connectioninstance")
	if err := h.Manager.UpdateConnectionInstance(updateConnectionRequest); err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
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

func (h *Handler) GetConnection(ctx context.Context, connectionId *grpc_application_network_go.ConnectionInstanceId) (*grpc_application_network_go.ConnectionInstance, error) {
	vErr := entities.ValidateConnectionInstanceId(connectionId)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	conn, err := h.Manager.GetConnectionInstance(connectionId)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return conn.ToGRPC(), nil
}

func (h *Handler) ListConnections(_ context.Context, orgID *grpc_organization_go.OrganizationId) (*grpc_application_network_go.ConnectionInstanceList, error) {
	log.Debug().Msg("ListConnections")
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

func (h *Handler) GetConnectionByZtNetworkId(ctx context.Context, request *grpc_application_network_go.ZTNetworkConnectionId) (*grpc_application_network_go.ConnectionInstance, error) {
	vErr := entities.ValidateZTNetworkConnectionId(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	connectionInstance, err := h.Manager.GetConnectionByZtNetworkId(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return connectionInstance.ToGRPC(), nil
}

// ListInboundConnections retrieves a list with all the connections where the appInstanceId is the target
func (h *Handler) ListInboundConnections(_ context.Context, appInstanceID *grpc_application_go.AppInstanceId) (*grpc_application_network_go.ConnectionInstanceList, error) {
	err := entities.ValidAppInstanceId(appInstanceID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	list, err := h.Manager.ListInboundConnections(appInstanceID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	connList := make([]*grpc_application_network_go.ConnectionInstance, 0)
	for _, instance := range list {
		connList = append(connList, instance.ToGRPC())
	}
	return &grpc_application_network_go.ConnectionInstanceList{
		Connections: connList,
	}, nil

}

// ListOutboundConnections retrieves a list with all the connections where the appInstanceId is the source
func (h *Handler) ListOutboundConnections(_ context.Context, appInstanceID *grpc_application_go.AppInstanceId) (*grpc_application_network_go.ConnectionInstanceList, error) {
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

// AddZTNetworkConnection adds a new zt Connection (one per an inbound and one per the inbound)
func (h *Handler) AddZTNetworkConnection(ctx context.Context, addRequest *grpc_application_network_go.ZTNetworkConnection) (*grpc_application_network_go.ZTNetworkConnection, error) {

	vErr := entities.ValidateZTNetworkConnection(addRequest)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	added, err := h.Manager.AddZTNetworkConnection(addRequest)
	if err != nil {
		return nil, err
	}
	return added.ToGRPC(), nil
}

// ListZTNetworkConnection lists the connections in one zt network (one inbound and one outbound)
func (h *Handler) ListZTNetworkConnection(ctx context.Context, ztNetworkId *grpc_application_network_go.ZTNetworkConnectionId) (*grpc_application_network_go.ZTNetworkConnectionList, error) {
	vErr := entities.ValidateZTNetworkConnectionId(ztNetworkId)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	list, err := h.Manager.ListZTNetworkConnection(ztNetworkId)
	if err != nil {
		return nil, err
	}
	ztList := make([]*grpc_application_network_go.ZTNetworkConnection, 0)
	for _, conn := range list {
		ztList = append(ztList, conn.ToGRPC())
	}

	return &grpc_application_network_go.ZTNetworkConnectionList{
		Connections: ztList,
	}, nil
}

// UpdateZTNetworkConnection updates an existing zt connection
func (h *Handler) UpdateZTNetworkConnection(ctx context.Context, updateRequest *grpc_application_network_go.UpdateZTNetworkConnectionRequest) (*grpc_common_go.Success, error) {
	vErr := entities.ValidateUpdateZTNetworkConnectionRequest(updateRequest)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	err := h.Manager.UpdateZTNetworkConnection(updateRequest)
	if err != nil {
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}

// Remove ZTNetwork removes the ztNetworkConnection (the inbound and the outbound)
func (h *Handler) RemoveZTNetworkConnection(ctx context.Context, ztNetworkId *grpc_application_network_go.ZTNetworkConnectionId) (*grpc_common_go.Success, error) {
	vErr := entities.ValidateZTNetworkConnectionId(ztNetworkId)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	err := h.Manager.RemoveZTNetworkConnection(ztNetworkId)
	if err != nil {
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}
