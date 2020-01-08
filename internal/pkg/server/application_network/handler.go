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
		log.Error().Str("trace", err.DebugReport()).Msg("invalid add connection request")
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Interface("addConnectionRequest", addConnectionRequest).Msg("Adding connection instance")
	added, err := h.Manager.AddConnectionInstance(addConnectionRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add connection instance")
		return nil, conversions.ToGRPCError(err)
	}
	return added.ToGRPC(), nil
}

func (h *Handler) UpdateConnection(ctx context.Context, updateConnectionRequest *grpc_application_network_go.UpdateConnectionRequest) (*grpc_common_go.Success, error) {
	if err := entities.ValidUpdateConnectionRequest(updateConnectionRequest); err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid update connection request")
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Interface("updateConnectionRequest", updateConnectionRequest).Msg("Updating connection instance")
	if err := h.Manager.UpdateConnectionInstance(updateConnectionRequest); err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot update connection instance")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

func (h *Handler) RemoveConnection(ctx context.Context, removeConnectionRequest *grpc_application_network_go.RemoveConnectionRequest) (*grpc_common_go.Success, error) {
	if err := entities.ValidRemoveConnectionRequest(removeConnectionRequest); err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid remove connection request")
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Interface("removeConnectionRequest", removeConnectionRequest).Msg("Removing connection instance")
	if err := h.Manager.RemoveConnectionInstance(removeConnectionRequest); err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove connection instance")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

func (h *Handler) ExistsConnection(ctx context.Context, connectionId *grpc_application_network_go.ConnectionInstanceId) (*grpc_common_go.Exists, error) {
	vErr := entities.ValidateConnectionInstanceId(connectionId)
	if vErr != nil {
		log.Error().Str("trace", vErr.DebugReport()).Msg("invalid connection instance identifier")
		return nil, conversions.ToGRPCError(vErr)
	}
	exists, err := h.Manager.ExistsConnectionInstance(connectionId)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot determine if the connection instance exists")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Exists{Exists: exists}, nil
}

func (h *Handler) GetConnection(ctx context.Context, connectionId *grpc_application_network_go.ConnectionInstanceId) (*grpc_application_network_go.ConnectionInstance, error) {
	vErr := entities.ValidateConnectionInstanceId(connectionId)
	if vErr != nil {
		log.Error().Str("trace", vErr.DebugReport()).Msg("invalid connection instance identifier")
		return nil, conversions.ToGRPCError(vErr)
	}
	conn, err := h.Manager.GetConnectionInstance(connectionId)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get connection instance")
		return nil, conversions.ToGRPCError(err)
	}
	return conn.ToGRPC(), nil
}

func (h *Handler) ListConnections(_ context.Context, orgID *grpc_organization_go.OrganizationId) (*grpc_application_network_go.ConnectionInstanceList, error) {
	log.Debug().Msg("ListConnections")
	if err := entities.ValidOrganizationID(orgID); err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid organization identifier")
		return nil, conversions.ToGRPCError(err)
	}
	list, err := h.Manager.ListConnectionInstances(orgID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot list connection instances")
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

func (h *Handler) GetConnectionByZtNetworkId(ctx context.Context, request *grpc_application_network_go.ZTNetworkId) (*grpc_application_network_go.ConnectionInstance, error) {
	vErr := entities.ValidateZTNetworkId(request)
	if vErr != nil {
		log.Error().Str("trace", vErr.DebugReport()).Msg("invalid ZT network identifier")
		return nil, conversions.ToGRPCError(vErr)
	}
	connectionInstance, err := h.Manager.GetConnectionByZtNetworkId(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get connection by ZT network identifier")
		return nil, conversions.ToGRPCError(err)
	}
	return connectionInstance.ToGRPC(), nil
}

// ListInboundConnections retrieves a list with all the connections where the appInstanceId is the target
func (h *Handler) ListInboundConnections(_ context.Context, appInstanceID *grpc_application_go.AppInstanceId) (*grpc_application_network_go.ConnectionInstanceList, error) {
	err := entities.ValidAppInstanceId(appInstanceID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid application instance identifier")
		return nil, conversions.ToGRPCError(err)
	}
	list, err := h.Manager.ListInboundConnections(appInstanceID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot list inbound connections")
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
		log.Error().Str("trace", err.DebugReport()).Msg("invalid application instance identifier")
		return nil, conversions.ToGRPCError(err)
	}
	list, err := h.Manager.ListOutboundConnections(appInstanceID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot list outbound connections")
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
		log.Error().Str("trace", vErr.DebugReport()).Msg("invalid ZT network connection")
		return nil, conversions.ToGRPCError(vErr)
	}
	added, err := h.Manager.AddZTNetworkConnection(addRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add ZT network connection")
		return nil, err
	}
	return added.ToGRPC(), nil
}

// ListZTNetworkConnection lists the connections in one zt network (one inbound and one outbound)
func (h *Handler) ListZTNetworkConnection(ctx context.Context, ztNetworkId *grpc_application_network_go.ZTNetworkId) (*grpc_application_network_go.ZTNetworkConnectionList, error) {
	vErr := entities.ValidateZTNetworkId(ztNetworkId)
	if vErr != nil {
		log.Error().Str("trace", vErr.DebugReport()).Msg("invalid ZT network identifier")
		return nil, conversions.ToGRPCError(vErr)
	}
	list, err := h.Manager.ListZTNetworkConnection(ztNetworkId)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot list ZT network connection")
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
		log.Error().Str("trace", vErr.DebugReport()).Msg("invalid update ZT network connection request")
		return nil, conversions.ToGRPCError(vErr)
	}
	err := h.Manager.UpdateZTNetworkConnection(updateRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot update ZT network connection")
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}

// Remove ZTNetwork removes the ztNetworkConnection (the inbound and the outbound)
func (h *Handler) RemoveZTNetworkConnection(ctx context.Context, connection *grpc_application_network_go.ZTNetworkConnectionId) (*grpc_common_go.Success, error) {
	vErr := entities.ValidateZTNetworkConnectionId(connection)
	if vErr != nil {
		log.Error().Str("trace", vErr.DebugReport()).Msg("invalid ZT network connection identifier")
		return nil, conversions.ToGRPCError(vErr)
	}
	err := h.Manager.RemoveZTNetworkConnection(connection)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove ZT network connection")
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}

func (h *Handler) RemoveZTNetworkConnectionByNetworkId(_ context.Context, ztNetworkId *grpc_application_network_go.ZTNetworkId) (*grpc_common_go.Success, error) {
	vErr := entities.ValidateZTNetworkId(ztNetworkId)
	if vErr != nil {
		log.Error().Str("trace", vErr.DebugReport()).Msg("invalid ZT network identifier")
		return nil, conversions.ToGRPCError(vErr)
	}
	err := h.Manager.RemoveZTNetworkConnectionByNetworkId(ztNetworkId)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove ZT network connection by network identifier")
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}
