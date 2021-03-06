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

package eic

import (
	"context"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
)

// Handler structure for the application requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h *Handler) Add(ctx context.Context, request *grpc_inventory_go.AddEdgeControllerRequest) (*grpc_inventory_go.EdgeController, error) {
	log.Debug().Str("organizationID", request.OrganizationId).
		Str("name", request.Name).Msg("add controller")
	err := entities.ValidAddEdgeControllerRequest(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("add edge controller request is not valid")
		return nil, conversions.ToGRPCError(err)
	}
	added, err := h.Manager.Add(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add edge controller")
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("edgeControllerID", added.EdgeControllerId).Msg("controller has been added")
	return added.ToGRPC(), nil
}

func (h *Handler) List(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_inventory_go.EdgeControllerList, error) {
	err := entities.ValidOrganizationID(organizationID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid organization identifier")
		return nil, conversions.ToGRPCError(err)
	}
	controllers, err := h.Manager.List(organizationID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot list controllers")
		return nil, conversions.ToGRPCError(err)
	}
	toReturn := make([]*grpc_inventory_go.EdgeController, 0, len(controllers))
	for _, c := range controllers {
		toReturn = append(toReturn, c.ToGRPC())
	}
	result := &grpc_inventory_go.EdgeControllerList{
		Controllers: toReturn,
	}
	return result, nil
}

func (h *Handler) Remove(ctx context.Context, edgeControllerID *grpc_inventory_go.EdgeControllerId) (*grpc_common_go.Success, error) {
	err := entities.ValidEdgeControllerID(edgeControllerID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid edge controller identifier")
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.Remove(edgeControllerID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove edge controller")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

func (h *Handler) Update(ctx context.Context, request *grpc_inventory_go.UpdateEdgeControllerRequest) (*grpc_inventory_go.EdgeController, error) {
	err := entities.ValidUpdateEdgeControllerRequest(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid update edge controller request")
		return nil, conversions.ToGRPCError(err)
	}
	updated, err := h.Manager.Update(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot update edge controller")
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("edgeControllerID", updated.EdgeControllerId).Msg("edge controller has been updated")
	return updated.ToGRPC(), nil
}

func (h *Handler) Get(ctx context.Context, edgeControllerID *grpc_inventory_go.EdgeControllerId) (*grpc_inventory_go.EdgeController, error) {
	err := entities.ValidEdgeControllerID(edgeControllerID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid edge controller identifier")
		return nil, conversions.ToGRPCError(err)
	}
	retrieved, err := h.Manager.Get(edgeControllerID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get edge controller")
		return nil, conversions.ToGRPCError(err)
	}
	return retrieved.ToGRPC(), nil
}
