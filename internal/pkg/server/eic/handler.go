/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package eic

import (
	"context"
	"github.com/nalej/derrors"
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
func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}

func (h * Handler) Add(ctx context.Context, request *grpc_inventory_go.AddEdgeControllerRequest) (*grpc_inventory_go.EdgeController, error) {
	log.Debug().Str("organizationID", request.OrganizationId).
		Str("name", request.Name).Msg("add controller")
	err := entities.ValidAddEdgeControllerRequest(request)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	added, err := h.Manager.Add(request)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("edgeControllerID", added.EdgeControllerId).Msg("controller has been added")
	return added.ToGRPC(), nil
}

func (h * Handler) List(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_inventory_go.EdgeControllerList, error) {
	err := entities.ValidOrganizationID(organizationID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	controllers, err := h.Manager.List(organizationID)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	toReturn := make([]*grpc_inventory_go.EdgeController, 0, len(controllers))
	for _, c := range controllers{
		toReturn = append(toReturn, c.ToGRPC())
	}
	result := &grpc_inventory_go.EdgeControllerList{
		Controllers:               toReturn,
	}
	return result, nil
}

func (h * Handler) Remove(ctx context.Context, edgeControllerID *grpc_inventory_go.EdgeControllerId) (*grpc_common_go.Success, error) {
	err := entities.ValidEdgeControllerID(edgeControllerID)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.Remove(edgeControllerID)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

func (h * Handler) Update(ctx context.Context, request *grpc_inventory_go.UpdateEdgeControllerRequest) (*grpc_inventory_go.EdgeController, error) {
	err := entities.ValidUpdateEdgeControllerRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	updated, err := h.Manager.Update(request)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("edgeControllerID", updated.EdgeControllerId).Msg("edge controller has been updated")
	return updated.ToGRPC(), nil
}


// Get the information of an edge controller.
func (h *Handler) Get(_ context.Context, in *grpc_inventory_go.EdgeControllerId) (*grpc_inventory_go.EdgeController, error) {
	return nil, conversions.ToGRPCError(derrors.NewUnimplementedError("not implemented yet"))
}