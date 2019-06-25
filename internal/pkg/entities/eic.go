/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-inventory-go"
	"time"
)

const DefaultLocation = "undefined"

type InventoryLocation struct {
	Geolocation string `json:"geolocation,omitempty" cql:"geolocation"`
	Geohash     string `json:"geohash,omitempty" cql:"geohash"`
}

func (location *InventoryLocation) ToGRPC() *grpc_inventory_go.InventoryLocation {
	if location == nil {
		return nil
	}
	return &grpc_inventory_go.InventoryLocation{
		Geolocation: location.Geolocation,
		Geohash: location.Geohash,
	}
}

func NewLocationFromGRPC(location *grpc_inventory_go.InventoryLocation)*InventoryLocation {
	if location == nil {
		return nil
	}
	return &InventoryLocation{
		Geolocation: location.Geolocation,
		Geohash: location.Geohash,
	}
}

// EdgeController entity.
type EdgeController struct {
	// OrganizationId with the organization identifier.
	OrganizationId		string `json:"organization_id,omitempty" cql:"organization_id"`
	// EdgeControllerId with the EIC identifier.
	EdgeControllerId 	string `json:"edge_controller_id,omitempty" cql:"edge_controller_id"`
	// Show flag to determine if this asset should be shown on the UI. This flag is internally used
	// for the async uninstall/removal of the asset.
	Show 				bool `json:"show,omitempty" cql:"show"`
	// Created time
	Created 			int64 `json:"created,omitempty" cql:"created"`
	// Name of the EIC.
	Name 				string `json:"name,omitempty" cql:"name"`
	// Labels defined by the user.
	Labels				map[string]string `json:"labels,omitempty" cql:"labels"`
	// LastAliveTimestamp contains the last alive message received
	LastAliveTimestamp 	int64    `json:"last_alive_timestamp,omitempty" cql:"last_alive_timestamp"`
	// location with the EC location
	Location            *InventoryLocation `json:"location,omitempty" cql:"location"`
	Os 					*OperatingSystemInfo    `json:"os,omitempty" cql:"os"`
	Hardware 			*HardwareInfo           `json:"hardware,omitempty" cql:"hardware"`
	Storage 			[]*StorageHardwareInfo  `json:"storage,omitempty" cql:"storage"`

}

func NewEdgeControllerFromGRPC(eic * grpc_inventory_go.AddEdgeControllerRequest) * EdgeController{
	if eic == nil{
		return nil
	}
	if eic.Geolocation == "" {
		eic.Geolocation = DefaultLocation
	}
	var os *OperatingSystemInfo
	var hardware 		*HardwareInfo
	var storage 		[]*StorageHardwareInfo
	storage = make ([]*StorageHardwareInfo, 0)

	if eic.AssetInfo != nil {
		os = NewOperatingSystemInfoFromGRPC(eic.AssetInfo.Os)
		hardware = NewHardwareInfoFromGRPC(eic.AssetInfo.Hardware)
		for _, sto := range eic.AssetInfo.Storage {
			storage = append(storage, NewStorageHardwareInfoFromGRPC(sto))
		}
	}

	return &EdgeController{
		OrganizationId:		eic.OrganizationId,
		EdgeControllerId: 	GenerateUUID(),
		Show:             	true,
		Created:          	time.Now().Unix(),
		Name:             	eic.Name,
		Labels:           	eic.Labels,
		Location:         	&InventoryLocation{
			Geolocation:eic.Geolocation,
		},
		Os:         		os,
		Hardware: 			hardware,
		Storage: 			storage,
	}
}

func (ec * EdgeController) ToGRPC() *grpc_inventory_go.EdgeController{
	if ec == nil{
		return nil
	}
	storage := make ([]*grpc_inventory_go.StorageHardwareInfo, 0)
	for _, sto := range ec.Storage {
		storage = append(storage, sto.ToGRPC())
	}
	return &grpc_inventory_go.EdgeController{
		OrganizationId:       ec.OrganizationId,
		EdgeControllerId:     ec.EdgeControllerId,
		Show:                 ec.Show,
		Created:              ec.Created,
		Name:                 ec.Name,
		Labels:               ec.Labels,
		LastAliveTimestamp:   ec.LastAliveTimestamp,
		Location: 			  ec.Location.ToGRPC(),
		AssetInfo: 			  &grpc_inventory_go.AssetInfo{
			Os:       ec.Os.ToGRPC(),
			Hardware: ec.Hardware.ToGRPC(),
			Storage:  storage,
		},
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
	if request.UpdateLastAlive {
		ec.LastAliveTimestamp = request.LastAliveTimestamp
	}
	if request.UpdateGeolocation {
		if ec.Location == nil {
			ec.Location = &InventoryLocation{
				Geolocation: request.Geolocation,
			}
		}else{
			ec.Location.Geolocation = request.Geolocation
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