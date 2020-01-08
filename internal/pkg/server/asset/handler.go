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

package asset

import (
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

// Handler structure for the cluster requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

// Add a new asset to the system.
func (h *Handler) Add(ctx context.Context, addRequest *grpc_inventory_go.AddAssetRequest) (*grpc_inventory_go.Asset, error) {
	log.Debug().Str("organizationID", addRequest.OrganizationId).
		Str("agentID", addRequest.AgentId).Msg("add asset")
	err := entities.ValidAddAssetRequest(addRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid add asset request")
		return nil, conversions.ToGRPCError(err)
	}
	added, err := h.Manager.Add(addRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add asset")
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("assetID", added.AssetId).Msg("asset has been added")
	return added.ToGRPC(), nil
}

func (h *Handler) Get(ctx context.Context, assetID *grpc_inventory_go.AssetId) (*grpc_inventory_go.Asset, error) {
	err := entities.ValidAssetID(assetID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid asset identifier")
		return nil, conversions.ToGRPCError(err)
	}
	asset, err := h.Manager.Get(assetID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get asset")
		return nil, conversions.ToGRPCError(err)
	}
	return asset.ToGRPC(), nil
}

// List the assets of an organization.
func (h *Handler) List(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_inventory_go.AssetList, error) {
	err := entities.ValidOrganizationID(organizationID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid organization identifier")
		return nil, conversions.ToGRPCError(err)
	}
	assets, err := h.Manager.List(organizationID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot list assets")
		return nil, conversions.ToGRPCError(err)
	}
	toReturn := make([]*grpc_inventory_go.Asset, 0, len(assets))
	for _, a := range assets {
		toReturn = append(toReturn, a.ToGRPC())
	}
	result := &grpc_inventory_go.AssetList{
		Assets: toReturn,
	}
	return result, nil
}

// Remove a given assets from an organization.
func (h *Handler) Remove(ctx context.Context, assetID *grpc_inventory_go.AssetId) (*grpc_common_go.Success, error) {
	err := entities.ValidAssetID(assetID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid asset identifier")
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.Remove(assetID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove asset")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// Update the information of an asset.
func (h *Handler) Update(ctx context.Context, updateRequest *grpc_inventory_go.UpdateAssetRequest) (*grpc_inventory_go.Asset, error) {
	err := entities.ValidUpdateAssetRequest(updateRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid update asset request")
		return nil, conversions.ToGRPCError(err)
	}
	updated, err := h.Manager.Update(updateRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot update asset")
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("assetID", updated.AssetId).Msg("asset has been updated")
	return updated.ToGRPC(), nil
}

func (h *Handler) ListControllerAssets(ctx context.Context, edgeControllerId *grpc_inventory_go.EdgeControllerId) (*grpc_inventory_go.AssetList, error) {
	err := entities.ValidEdgeControllerID(edgeControllerId)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid edge controller identifier")
		return nil, conversions.ToGRPCError(err)
	}
	assets, err := h.Manager.ListControllerAssets(edgeControllerId)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot list controller assets")
		return nil, conversions.ToGRPCError(err)
	}
	toReturn := make([]*grpc_inventory_go.Asset, 0, len(assets))
	for _, a := range assets {
		toReturn = append(toReturn, a.ToGRPC())
	}
	result := &grpc_inventory_go.AssetList{
		Assets: toReturn,
	}
	return result, nil
}
