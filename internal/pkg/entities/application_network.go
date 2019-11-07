/*
 * Copyright 2019 Nalej
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
 *
 */

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-network-go"
)

type ConnectionStatus int32

const (
	ConnectionStatusWaiting ConnectionStatus = iota + 1
	ConnectionStatusEstablished
	ConnectionStatusTerminated
	ConnectionStatusFailed
)

var ConnectionStatusToGRPC = map[ConnectionStatus]grpc_application_network_go.ConnectionStatus{
	ConnectionStatusWaiting:     grpc_application_network_go.ConnectionStatus_WAITING,
	ConnectionStatusEstablished: grpc_application_network_go.ConnectionStatus_ESTABLISHED,
	ConnectionStatusTerminated:  grpc_application_network_go.ConnectionStatus_TERMINATED,
	ConnectionStatusFailed:      grpc_application_network_go.ConnectionStatus_FAILED,
}

var ConnectionStatusFromGRPC = map[grpc_application_network_go.ConnectionStatus]ConnectionStatus{
	grpc_application_network_go.ConnectionStatus_WAITING:     ConnectionStatusWaiting,
	grpc_application_network_go.ConnectionStatus_ESTABLISHED: ConnectionStatusEstablished,
	grpc_application_network_go.ConnectionStatus_TERMINATED:  ConnectionStatusTerminated,
	grpc_application_network_go.ConnectionStatus_FAILED:      ConnectionStatusFailed,
}

// ConnectionInstance model with the info of a connection between two application instances
type ConnectionInstance struct {
	// OrganizationId with the organization identifier
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// ConnectionId with the connection identifier
	ConnectionId string `json:"connection_id,omitempty" cql:"connection_id"`
	// SourceInstanceId with the instance identifier of the connection source
	SourceInstanceId string `json:"source_instance_id,omitempty" cql:"source_instance_id"`
	// SourceInstanceName with the instance name of the connection source
	SourceInstanceName string `json:"source_instance_name,omitempty" cql:"source_instance_name"`
	// TargetInstanceId with the instance identifier of the connection target
	TargetInstanceId string `json:"target_instance_id,omitempty" cql:"target_instance_id"`
	// TargetInstanceName with the instance name of the connection target
	TargetInstanceName string `json:"target_instance_name,omitempty" cql:"target_instance_name"`
	// InboundName with the name of the inbound network interface
	InboundName string `json:"inbound_name,omitempty" cql:"inbound_name"`
	// OutboundName with the name of the outbound network interface
	OutboundName string `json:"outbound_name,omitempty" cql:"outbound_name"`
	// OutboundRequired with the flag `required` of the outbound network interface
	OutboundRequired bool `json:"outbound_required,omitempty" cql:"outbound_required"`
	// Status with the status of the connection instance
	Status ConnectionStatus `json:"status,omitempty" cql:"status"`
	// IpRange with the IP range of the connection
	IpRange string `json:"ip_range,omitempty" cql:"ip_range"`
	// ZtNetworkId with the ID of the ZT network created for the connection
	ZtNetworkId string `json:"zt_network_id,omitempty" cql:"zt_network_id"`
}

// NewConnectionInstanceFromGRPC Creates a new entities.ConnectionInstance using an grpc_application_network_go.AddConnectionRequest, source and target names, and outbound required flag.
func NewConnectionInstanceFromGRPC(request grpc_application_network_go.AddConnectionRequest, sourceInstanceName string, targetInstanceName string, outboundRequired bool) *ConnectionInstance {
	return &ConnectionInstance{
		OrganizationId:     request.GetOrganizationId(),
		ConnectionId:       GenerateUUID(),
		SourceInstanceId:   request.GetSourceInstanceId(),
		SourceInstanceName: sourceInstanceName,
		TargetInstanceId:   request.GetTargetInstanceId(),
		TargetInstanceName: targetInstanceName,
		InboundName:        request.GetInboundName(),
		OutboundName:       request.GetOutboundName(),
		OutboundRequired:   outboundRequired,
		Status:             ConnectionStatusWaiting,
		IpRange:            request.IpRange,
		ZtNetworkId:        request.ZtNetworkId,
	}
}

// ToGRPC Converts a entities.ConnectionInstance to a grpc_application_network_go.ConnectionInstance and returns its pointer.
func (c *ConnectionInstance) ToGRPC() *grpc_application_network_go.ConnectionInstance {
	if c == nil {
		return nil
	}
	return &grpc_application_network_go.ConnectionInstance{
		OrganizationId:     c.OrganizationId,
		ConnectionId:       c.ConnectionId,
		SourceInstanceId:   c.SourceInstanceId,
		SourceInstanceName: c.SourceInstanceName,
		TargetInstanceId:   c.TargetInstanceId,
		TargetInstanceName: c.TargetInstanceName,
		InboundName:        c.InboundName,
		OutboundName:       c.OutboundName,
		OutboundRequired:   c.OutboundRequired,
		Status:             ConnectionStatusToGRPC[c.Status],
		IpRange:            c.IpRange,
		ZtNetworkId:        c.ZtNetworkId,
	}
}

func (c *ConnectionInstance) ApplyUpdate(updateConnectionRequest *grpc_application_network_go.UpdateConnectionRequest) {
	if updateConnectionRequest.UpdateStatus {
		c.Status = ConnectionStatusFromGRPC[updateConnectionRequest.Status]
	}
	if updateConnectionRequest.UpdateIpRange {
		c.IpRange = updateConnectionRequest.IpRange
	}
	if updateConnectionRequest.UpdateZtNetworkId {
		c.ZtNetworkId = updateConnectionRequest.ZtNetworkId
	}
}

func ValidAddConnectionRequest(request *grpc_application_network_go.AddConnectionRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("expecting an OrganizationId")
	}
	if request.InboundName == "" {
		return derrors.NewInvalidArgumentError("expecting an InboundName")
	}
	if request.SourceInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting an SourceInstanceId")
	}
	if request.OutboundName == "" {
		return derrors.NewInvalidArgumentError("expecting an OutboundName")
	}
	if request.TargetInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting an TargetInstanceId")
	}
	return nil
}

func ValidUpdateConnectionRequest(request *grpc_application_network_go.UpdateConnectionRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("expecting an OrganizationId")
	}
	if request.InboundName == "" {
		return derrors.NewInvalidArgumentError("expecting an InboundName")
	}
	if request.SourceInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting a SourceInstanceId")
	}
	if request.OutboundName == "" {
		return derrors.NewInvalidArgumentError("expecting an OutboundName")
	}
	if request.TargetInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting a TargetInstanceId")
	}
	if request.UpdateIpRange && request.IpRange == "" {
		return derrors.NewInvalidArgumentError("expecting a RangeIp")
	}
	return nil
}

func ValidRemoveConnectionRequest(request *grpc_application_network_go.RemoveConnectionRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("expecting an OrganizationId")
	}
	if request.InboundName == "" {
		return derrors.NewInvalidArgumentError("expecting an InboundName")
	}
	if request.SourceInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting an SourceInstanceId")
	}
	if request.OutboundName == "" {
		return derrors.NewInvalidArgumentError("expecting an OutboundName")
	}
	if request.TargetInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting an TargetInstanceId")
	}
	return nil
}

func ValidateConnectionInstanceId(request *grpc_application_network_go.ConnectionInstanceId) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("expecting an OrganizationId")
	}
	if request.InboundName == "" {
		return derrors.NewInvalidArgumentError("expecting an InboundName")
	}
	if request.SourceInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting an SourceInstanceId")
	}
	if request.OutboundName == "" {
		return derrors.NewInvalidArgumentError("expecting an OutboundName")
	}
	if request.TargetInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting an TargetInstanceId")
	}
	return nil
}

// ConnectionInstanceLink model with the info of a connection between two fragments on each cluster
type ConnectionInstanceLink struct {
	// OrganizationId with the organization identifier
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// ConnectionId with the connection identifier
	ConnectionId string `json:"connection_id,omitempty" cql:"connection_id"`
	// SourceInstanceId with the instance identifier of the connection source
	SourceInstanceId string `json:"source_instance_id,omitempty" cql:"source_instance_id"`
	// SourceClusterId with the cluster identifier where the source fragment is deployed
	SourceClusterId string `json:"source_cluster_id,omitempty" cql:"source_cluster_id"`
	// TargetInstanceId with the instance identifier of the connection target
	TargetInstanceId string `json:"target_instance_id,omitempty" cql:"target_instance_id"`
	// TargetClusterId with the cluster identifier where the target fragment is deployed
	TargetClusterId string `json:"target_cluster_id,omitempty" cql:"target_cluster_id"`
	// InboundName with the name of the inbound network interface
	InboundName string `json:"inbound_name,omitempty" cql:"inbound_name"`
	// OutboundName with the name of the outbound network interface
	OutboundName string `json:"outbound_name,omitempty" cql:"outbound_name"`
	// Status with the status of the connection instance
	Status ConnectionStatus `json:"status,omitempty" cql:"status"`
}

// toGRPC
func (c *ConnectionInstanceLink) toGRPC() *grpc_application_network_go.ConnectionInstanceLink {
	if c == nil {
		return nil
	}
	return &grpc_application_network_go.ConnectionInstanceLink{
		OrganizationId:   c.OrganizationId,
		ConnectionId:     c.ConnectionId,
		SourceInstanceId: c.SourceInstanceId,
		SourceClusterId:  c.SourceClusterId,
		TargetInstanceId: c.TargetInstanceId,
		TargetClusterId:  c.TargetClusterId,
		InboundName:      c.InboundName,
		OutboundName:     c.OutboundName,
		Status:           ConnectionStatusToGRPC[c.Status],
	}
}

type ConnectionSide int

const (
	ConnectionSideInbound ConnectionSide = iota + 1
	ConnectionSideOutbound
)

var ConnectionSideToGRPC = map[ConnectionSide]grpc_application_network_go.ConnectionSide{
	ConnectionSideInbound:  grpc_application_network_go.ConnectionSide_SIDE_INBOUND,
	ConnectionSideOutbound: grpc_application_network_go.ConnectionSide_SIDE_OUTBOUND,
}

var ConnectionSideFromGRPC = map[grpc_application_network_go.ConnectionSide]ConnectionSide{
	grpc_application_network_go.ConnectionSide_SIDE_INBOUND:  ConnectionSideInbound,
	grpc_application_network_go.ConnectionSide_SIDE_OUTBOUND: ConnectionSideOutbound,
}

type ZTNetworkConnection struct {
	// Organization ID
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// ZtNetworkId zero tier network identifier
	ZtNetworkId string `json:"zt_network_id,omitempty" cql:"zt_network_id"`
	// Application instance identifier
	AppInstanceId string `json:"app_instance_id,omitempty" cql:"app_instance_id"`
	// ServiceId with the service identifier
	ServiceId string ` json:"service_id,omitempty" cql:"service_id"`
	// ZtMember zero tier member identifier
	ZtMember string `json:"zt_member,omitempty" cql:"zt_member"`
	// ZtIp
	ZtIp string `json:"zt_ip,omitempty" cql:"zt_ip"`
	// ClusterId with the cluster identifier
	ClusterId string `json:"cluster_id,omitempty" cql:"cluster_id"`
	// Side to indicate if the instance is the inbound or the outbound in the connection
	Side ConnectionSide `json:"side,omitempty" cql:"side"`
}

func NewZTNetworkConnectionFromGRPC(zt *grpc_application_network_go.ZTNetworkConnection) *ZTNetworkConnection {
	if zt == nil {
		return nil
	}
	return &ZTNetworkConnection{
		OrganizationId: zt.OrganizationId,
		ZtNetworkId:    zt.ZtNetworkId,
		AppInstanceId:  zt.AppInstanceId,
		ServiceId:      zt.ServiceId,
		ZtMember:       zt.ZtMember,
		ZtIp:           zt.ZtIp,
		ClusterId:      zt.ClusterId,
		Side:           ConnectionSideFromGRPC[zt.Side],
	}
}
func (zt *ZTNetworkConnection) ToGRPC() *grpc_application_network_go.ZTNetworkConnection {
	return &grpc_application_network_go.ZTNetworkConnection{
		OrganizationId: zt.OrganizationId,
		ZtNetworkId:    zt.ZtNetworkId,
		AppInstanceId:  zt.AppInstanceId,
		ServiceId:      zt.ServiceId,
		ZtMember:       zt.ZtMember,
		ZtIp:           zt.ZtIp,
		ClusterId:      zt.ClusterId,
		Side:           ConnectionSideToGRPC[zt.Side],
	}
}

func (zt *ZTNetworkConnection) ApplyUpdate(updateRequest *grpc_application_network_go.UpdateZTNetworkConnectionRequest) {
	if updateRequest.UpdateZtIp {
		zt.ZtIp = updateRequest.ZtIp
	}
	if updateRequest.UpdateZtMember {
		zt.ZtMember = updateRequest.ZtMember
	}
}

func ValidateZTNetworkConnection(request *grpc_application_network_go.ZTNetworkConnection) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("expecting an OrganizationId")
	}
	if request.ZtNetworkId == "" {
		return derrors.NewInvalidArgumentError("expecting a ZtNetworkId")
	}
	if request.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting a AppInstanceId")
	}
	if request.ServiceId == "" {
		return derrors.NewInvalidArgumentError("expecting a ServiceId")
	}
	if request.ClusterId == "" {
		return derrors.NewInvalidArgumentError("expecting a ClusterId")
	}
	return nil
}

func ValidateZTNetworkConnectionId(request *grpc_application_network_go.ZTNetworkConnectionId) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("expecting an OrganizationId")
	}
	if request.ZtNetworkId == "" {
		return derrors.NewInvalidArgumentError("expecting a ZtNetworkId")
	}
	if request.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting a AppInstanceId")
	}
	if request.ServiceId == "" {
		return derrors.NewInvalidArgumentError("expecting a ServiceId")
	}
	if request.ClusterId == "" {
		return derrors.NewInvalidArgumentError("expecting a ClusterId")
	}
	return nil
}

func ValidateUpdateZTNetworkConnectionRequest(request *grpc_application_network_go.UpdateZTNetworkConnectionRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("expecting an OrganizationId")
	}
	if request.ZtNetworkId == "" {
		return derrors.NewInvalidArgumentError("expecting a ZtNetworkId")
	}
	if request.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting a AppInstanceId")
	}
	if request.ServiceId == "" {
		return derrors.NewInvalidArgumentError("expecting a ServiceId")
	}
	return nil
}

func ValidateZTNetworkId(ztNetworkId *grpc_application_network_go.ZTNetworkId) derrors.Error {
	if ztNetworkId.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("expecting an OrganizationId")
	}
	if ztNetworkId.ZtNetworkId == "" {
		return derrors.NewInvalidArgumentError("expecting a ZtNetworkId")
	}
	return nil
}
