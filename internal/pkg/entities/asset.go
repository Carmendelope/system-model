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

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-inventory-go"
	"time"
)

// AgentOpSummary contains the result of an asset operation
// this is a provisional result!
type AgentOpSummary struct {
	// OperationId with the operation identifier.
	OperationId string `json:"operation_id,omitempty" cql:"operation_id"`
	// Timestamp of the response.
	Timestamp int64 `json:"timestamp,omitempty" cql:"timestamp"`
	// Status indicates if the operation was successfull
	Status OpStatus `json:"status,omitempty" cql:"status"`
	// Info with additional information for an operation.
	Info string `json:"info,omitempty" cql:"info"`
}

func (a *AgentOpSummary) ToGRPC() *grpc_inventory_go.AgentOpSummary {
	if a == nil {
		return nil
	}
	return &grpc_inventory_go.AgentOpSummary{
		OperationId: a.OperationId,
		Timestamp:   a.Timestamp,
		Status:      OpStatusToGRPC[a.Status],
		Info:        a.Info,
	}
}

func NewAgentOpSummaryFromGRPC(op *grpc_inventory_go.AgentOpSummary) *AgentOpSummary {
	return &AgentOpSummary{
		OperationId: op.OperationId,
		Timestamp:   op.Timestamp,
		Status:      OpStatusFromGRPC[op.Status],
		Info:        op.Info,
	}
}

// Asset represents an element in the network from which we register some type of information. Example of
// assets could be workstations, nodes in a cluster, or other type of hardware.
type Asset struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// EdgeControllerId with the EIC identifier
	EdgeControllerId string `json:"edge_controller_id,omitempty"`
	// AssetId with the asset identifier.
	AssetId string `json:"asset_id,omitempty"`
	// AgentId with the agent identifier that is monitoring this asset if any.
	AgentId string `json:"agent_id,omitempty"`
	// Show flag to determine if this asset should be shown on the UI. This flag is internally used
	// for the async uninstall/removal of the asset.
	Show bool `json:"show,omitempty"`
	// Created time
	Created int64 `json:"created,omitempty"`
	// Labels defined by the user.
	Labels map[string]string `json:"labels,omitempty"`
	// OS contains Operating System information.
	Os *OperatingSystemInfo `json:"os,omitempty" cql:"os"`
	// Hardware information.
	Hardware *HardwareInfo `json:"hardware,omitempty" cql:"hardware"`
	// Storage information.
	Storage []StorageHardwareInfo `json:"storage,omitempty" cql:"storage"`
	// EicNetIp contains the current IP address that connects the asset to the EIC.
	EicNetIp string `json:"eic_net_ip,omitempty" cql:"eic_net_ip"`
	// AgentOpSummary contains the result of the last operation fr this asset
	LastOpResult *AgentOpSummary `json:"last_op_result,omitempty" cql: "eic_net_ip"`
	// LastAliveTimestamp contains the last alive message received
	LastAliveTimestamp int64 `json:"last_alive_timestamp,omitempty" cql:"last_alive_timestamp"`
	// Location contains the location of the asset
	Location *InventoryLocation `json:"location,omitempty"`
}

func NewAssetFromGRPC(addRequest *grpc_inventory_go.AddAssetRequest) *Asset {

	storage := make([]StorageHardwareInfo, 0)
	for _, sto := range addRequest.Storage {
		storage = append(storage, *NewStorageHardwareInfoFromGRPC(sto))
	}

	location := &InventoryLocation{
		Geolocation: DefaultLocation,
		Geohash:     "",
	}

	if addRequest.Location != nil && addRequest.Location.Geolocation != "" {
		location.Geolocation = addRequest.Location.Geolocation
	}

	return &Asset{
		OrganizationId:   addRequest.OrganizationId,
		EdgeControllerId: addRequest.EdgeControllerId,
		AssetId:          GenerateUUID(),
		AgentId:          addRequest.AgentId,
		Show:             true,
		Created:          time.Now().Unix(),
		Labels:           addRequest.Labels,
		Os:               NewOperatingSystemInfoFromGRPC(addRequest.Os),
		Hardware:         NewHardwareInfoFromGRPC(addRequest.Hardware),
		Storage:          storage,
		Location:         location,
	}
}

func (a *Asset) ToGRPC() *grpc_inventory_go.Asset {

	storage := make([]*grpc_inventory_go.StorageHardwareInfo, 0)
	for _, sto := range a.Storage {
		storage = append(storage, sto.ToGRPC())
	}

	return &grpc_inventory_go.Asset{
		OrganizationId:     a.OrganizationId,
		EdgeControllerId:   a.EdgeControllerId,
		AssetId:            a.AssetId,
		AgentId:            a.AgentId,
		Show:               a.Show,
		Created:            a.Created,
		Labels:             a.Labels,
		Os:                 a.Os.ToGRPC(),
		Hardware:           a.Hardware.ToGRPC(),
		Storage:            storage,
		EicNetIp:           a.EicNetIp,
		LastAliveTimestamp: a.LastAliveTimestamp,
		LastOpResult:       a.LastOpResult.ToGRPC(),
		Location:           a.Location.ToGRPC(),
	}
}

func (a *Asset) ApplyUpdate(request *grpc_inventory_go.UpdateAssetRequest) {
	if request.AddLabels {
		if a.Labels == nil {
			a.Labels = make(map[string]string, 0)
		}
		for k, v := range request.Labels {
			a.Labels[k] = v
		}
	}
	if request.RemoveLabels {
		for k, _ := range request.Labels {
			delete(a.Labels, k)
		}
	}
	if request.UpdateLastAlive {
		a.LastAliveTimestamp = request.LastAliveTimestamp
	}
	if request.UpdateLastOpSummary {
		a.LastOpResult = NewAgentOpSummaryFromGRPC(request.LastOpSummary)
	}
	if request.UpdateIp {
		a.EicNetIp = request.EicNetIp
	}
	if request.UpdateLocation {
		if a.Location == nil {
			a.Location = &InventoryLocation{
				Geolocation: request.Location.Geolocation,
				Geohash:     request.Location.Geohash,
			}
		} else {
			a.Location.Geolocation = request.Location.Geolocation
			a.Location.Geohash = request.Location.Geohash
		}
	}
}

func ValidAddAssetRequest(addRequest *grpc_inventory_go.AddAssetRequest) derrors.Error {
	if addRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if addRequest.EdgeControllerId == "" {
		return derrors.NewInvalidArgumentError(emptyEdgeControllerId)
	}
	return nil
}

func ValidAssetID(assetID *grpc_inventory_go.AssetId) derrors.Error {
	if assetID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if assetID.AssetId == "" {
		return derrors.NewInvalidArgumentError(emptyAssetId)
	}

	return nil
}

func ValidUpdateAssetRequest(request *grpc_inventory_go.UpdateAssetRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.AssetId == "" {
		return derrors.NewInvalidArgumentError(emptyAssetId)
	}
	return nil
}
