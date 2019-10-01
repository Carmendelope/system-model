/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
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
		IpRange:            "",
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
	}
}

func (c *ConnectionInstance) ApplyUpdate(updateConnectionRequest *grpc_application_network_go.UpdateConnectionRequest) {
	if updateConnectionRequest.UpdateStatus {
		c.Status = ConnectionStatusFromGRPC[updateConnectionRequest.Status]
	}
	if updateConnectionRequest.UpdateIpRange {
		c.IpRange = updateConnectionRequest.IpRange
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
