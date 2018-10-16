/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-infrastructure-go"
)

// ClusterType enumeration with the types of clusters supported to manage application deployments.
type ClusterType int

const (
	KubernetesCluster  ClusterType = iota + 1
	DockerCluster
)

var ClusterTypeToGRPC = map[ClusterType]grpc_infrastructure_go.ClusterType{
	KubernetesCluster: grpc_infrastructure_go.ClusterType_KUBERNETES,
	DockerCluster: grpc_infrastructure_go.ClusterType_DOCKER_NODE,
}

var ClusterTypeFromGRPC = map[grpc_infrastructure_go.ClusterType]ClusterType{
	grpc_infrastructure_go.ClusterType_KUBERNETES: KubernetesCluster,
	grpc_infrastructure_go.ClusterType_DOCKER_NODE: DockerCluster,
}

// MultinenantSupport enumeration defining the types of multitenancy supported by the system. Notice that even
// if it is modeled as a boolean now, we leave the definition as an enumeration to support other types of multitenancy
// like restrictions to parts of an organization, or priority based options.
type MultitenantSupport int
const (
	MultitenantYes MultitenantSupport = iota + 1
	MultitenantNo
)

var MultitenantSupportToGRPC = map[MultitenantSupport]grpc_infrastructure_go.MultitenantSupport{
	MultitenantYes: grpc_infrastructure_go.MultitenantSupport_YES,
	MultitenantNo: grpc_infrastructure_go.MultitenantSupport_NO,
}

var MultitenantSupportFromGRPC = map[grpc_infrastructure_go.MultitenantSupport]MultitenantSupport{
	grpc_infrastructure_go.MultitenantSupport_YES: MultitenantYes,
	grpc_infrastructure_go.MultitenantSupport_NO: MultitenantNo,
}

// InfraStatus enumeration defining the status of an element of the infrastructure.
type InfraStatus int

const (
	// Installing state represents an infrastructure element that is being installed at the momment.
	InfraStatusInstalling InfraStatus = iota + 1
	// Running state represents an infrastucture element that has been installed as is up and running.
	InfraStatusRunning
	// Error state represents an infrastructure element that cannot be used as any of the processes failed.
	InfraStatusError
)

var InfraStatusToGRPC = map[InfraStatus]grpc_infrastructure_go.InfraStatus{
	InfraStatusInstalling: grpc_infrastructure_go.InfraStatus_INSTALLING,
	InfraStatusRunning: grpc_infrastructure_go.InfraStatus_RUNNING,
	InfraStatusError: grpc_infrastructure_go.InfraStatus_ERROR,
}

var InfraStatusFromGRPC = map[grpc_infrastructure_go.InfraStatus]InfraStatus{
	grpc_infrastructure_go.InfraStatus_INSTALLING: InfraStatusInstalling,
	grpc_infrastructure_go.InfraStatus_RUNNING: InfraStatusRunning,
	grpc_infrastructure_go.InfraStatus_ERROR: InfraStatusError,
}

// Cluster entity representing a collection of nodes that supports applicaiton orchestration. This
// abstraction is used for monitoring and orchestration purposes.
type Cluster struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// ClusterId with the cluster identifier.
	ClusterId string `json:"cluster_id,omitempty"`
	// Name of the cluster.
	Name string `json:"name,omitempty"`
	// Description of the cluster.
	Description string `json:"description,omitempty"`
	// Type of cluster.
	ClusterType ClusterType `json:"cluster_type,omitempty"`
	// Multitenant support definition.
	Multitenant MultitenantSupport `json:"multitenant,omitempty"`
	// Status of the cluster based on monitoring information.
	Status InfraStatus `json:"status,omitempty"`
	// Labels for the cluster.
	Labels map[string]string `json:"labels,omitempty"`
	// Cordon flags to signal conductor not to schedule apps in the cluster.
	Cordon               bool     `json:"cordon,omitempty"`
}

func NewCluterFromGRPC(addClusterRequest *grpc_infrastructure_go.AddClusterRequest) *Cluster {
	return nil
}

func (c * Cluster) ToGRPC() * grpc_infrastructure_go.Cluster {
	return nil
}

func ValidAddClusterRequest(addClusterRequest *grpc_infrastructure_go.AddClusterRequest) derrors.Error {
	return nil
}

func ValidRemoveClusterRequest(removeClusterRequest * grpc_infrastructure_go.RemoveClusterRequest) derrors.Error {
	return nil
}

func ValidClusterID(clusterID *grpc_infrastructure_go.ClusterId) derrors.Error {
	if clusterID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if clusterID.ClusterId == "" {
		return derrors.NewInvalidArgumentError(emptyClusterId)
	}
	return nil
}