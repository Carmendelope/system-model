/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-inventory-go"
	"time"
)

// EdgeController entity.
type EdgeController struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// EdgeControllerId with the EIC identifier.
	EdgeControllerId string `json:"edge_controller_id,omitempty"`
	// Show flag to determine if this asset should be shown on the UI. This flag is internally used
	// for the async uninstall/removal of the asset.
	Show bool `json:"show,omitempty"`
	// Created time
	Created int64 `json:"created,omitempty"`
	// Name of the EIC.
	Name string `json:"name,omitempty"`
	// Labels defined by the user.
	Labels               map[string]string `json:"labels,omitempty"`
}

func NewEdgeControllerFromGRPC(eic * grpc_inventory_go.AddEdgeControllerRequest) * EdgeController{
	if eic == nil{
		return nil
	}
	return &EdgeController{
		OrganizationId:   eic.OrganizationId,
		EdgeControllerId: GenerateUUID(),
		Show:             true,
		Created:          time.Now().Unix(),
		Name:             eic.Name,
		Labels:           eic.Labels,
	}
}

func (ec * EdgeController) ToGRPC() *grpc_inventory_go.EdgeController{
	if ec == nil{
		return nil
	}
	return &grpc_inventory_go.EdgeController{
		OrganizationId:       ec.OrganizationId,
		EdgeControllerId:     ec.EdgeControllerId,
		Show:                 ec.Show,
		Created:              ec.Created,
		Name:                 ec.Name,
		Labels:               ec.Labels,
	}
}

func (ec * EdgeController) ApplyUpdate(request * grpc_inventory_go.UpdateEdgeControllerRequest){
	if request.AddLabels {
		if ec.Labels == nil {
			ec.Labels = make(map[string]string, 0)
		}
		for k, v := range request.Labels {
			ec.Labels[k] = v
		}
	}
	if request.RemoveLabels {
		for k, _ := range request.Labels {
			delete(ec.Labels, k)
		}
	}
}

func ValidAddEdgeControllerRequest(request *grpc_inventory_go.AddEdgeControllerRequest) derrors.Error{
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	return nil
}

func ValidEdgeControllerID(edgeControllerID *grpc_inventory_go.EdgeControllerId) derrors.Error{
	if edgeControllerID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if edgeControllerID.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError(emptyEdgeControllerId)
	}
	return nil
}

func ValidUpdateEdgeControllerRequest(request * grpc_inventory_go.UpdateEdgeControllerRequest) derrors.Error{
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError(emptyEdgeControllerId)
	}
	return nil
}