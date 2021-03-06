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
	"github.com/nalej/grpc-cluster-watcher-go"
	"github.com/nalej/grpc-connectivity-manager-go"
	"github.com/nalej/grpc-infrastructure-go"
)

// ClusterType enumeration with the types of clusters supported to manage application deployments.
type ClusterType int

const (
	KubernetesCluster ClusterType = iota + 1
	DockerCluster
)

var ClusterTypeToGRPC = map[ClusterType]grpc_infrastructure_go.ClusterType{
	KubernetesCluster: grpc_infrastructure_go.ClusterType_KUBERNETES,
	DockerCluster:     grpc_infrastructure_go.ClusterType_DOCKER_NODE,
}

var ClusterTypeFromGRPC = map[grpc_infrastructure_go.ClusterType]ClusterType{
	grpc_infrastructure_go.ClusterType_KUBERNETES:  KubernetesCluster,
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
	MultitenantNo:  grpc_infrastructure_go.MultitenantSupport_NO,
}

var MultitenantSupportFromGRPC = map[grpc_infrastructure_go.MultitenantSupport]MultitenantSupport{
	grpc_infrastructure_go.MultitenantSupport_YES: MultitenantYes,
	grpc_infrastructure_go.MultitenantSupport_NO:  MultitenantNo,
}

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
	InfraStatusRunning:    grpc_infrastructure_go.InfraStatus_RUNNING,
	InfraStatusError:      grpc_infrastructure_go.InfraStatus_ERROR,
}

var InfraStatusFromGRPC = map[grpc_infrastructure_go.InfraStatus]InfraStatus{
	grpc_infrastructure_go.InfraStatus_INSTALLING: InfraStatusInstalling,
	grpc_infrastructure_go.InfraStatus_RUNNING:    InfraStatusRunning,
	grpc_infrastructure_go.InfraStatus_ERROR:      InfraStatusError,
}

// ClusterStatus enumeration defining the status of an element of the infrastructure.
type ClusterStatus int

const (
	ClusterStatusUnknown ClusterStatus = iota + 1
	ClusterStatusOffline
	ClusterStatusOnline
	ClusterStatusOfflineCordon
	ClusterStatusOnlineCordon
)

var ClusterStatusToGRPC = map[ClusterStatus]grpc_connectivity_manager_go.ClusterStatus{
	ClusterStatusUnknown:       grpc_connectivity_manager_go.ClusterStatus_UNKNOWN,
	ClusterStatusOffline:       grpc_connectivity_manager_go.ClusterStatus_OFFLINE,
	ClusterStatusOnline:        grpc_connectivity_manager_go.ClusterStatus_ONLINE,
	ClusterStatusOfflineCordon: grpc_connectivity_manager_go.ClusterStatus_OFFLINE_CORDON,
	ClusterStatusOnlineCordon:  grpc_connectivity_manager_go.ClusterStatus_ONLINE_CORDON,
}

var ClusterStatusFromGRPC = map[grpc_connectivity_manager_go.ClusterStatus]ClusterStatus{
	grpc_connectivity_manager_go.ClusterStatus_UNKNOWN:        ClusterStatusUnknown,
	grpc_connectivity_manager_go.ClusterStatus_OFFLINE:        ClusterStatusOffline,
	grpc_connectivity_manager_go.ClusterStatus_ONLINE:         ClusterStatusOnline,
	grpc_connectivity_manager_go.ClusterStatus_OFFLINE_CORDON: ClusterStatusOfflineCordon,
	grpc_connectivity_manager_go.ClusterStatus_ONLINE_CORDON:  ClusterStatusOnlineCordon,
}

// Network type used by the cluster watcher
type NetworkType int

const (
	NetworkTypeCilium NetworkType = iota + 1
	NetworkTypeIstio
)

var NetworkTypeFromGRPC = map[grpc_cluster_watcher_go.NetworkType]NetworkType{
	grpc_cluster_watcher_go.NetworkType_ISTIO:  NetworkTypeIstio,
	grpc_cluster_watcher_go.NetworkType_CILIUM: NetworkTypeCilium,
}
var NetworkTypeToGRPC = map[NetworkType]grpc_cluster_watcher_go.NetworkType{
	NetworkTypeIstio:  grpc_cluster_watcher_go.NetworkType_ISTIO,
	NetworkTypeCilium: grpc_cluster_watcher_go.NetworkType_CILIUM,
}

// ClusterState defines the different states of a cluster regarding its provisioning and installation.
type ClusterState int

const (
	// Unknown indicates that the cluster has just been created in system model and no provisioning or install operation has taken place.
	Unknown ClusterState = iota + 1
	// Provisioning indicates that the cluster is now being provisioned in a cloud provider or through baremetal provisioning.
	Provisioning
	// Provisioned indicates that the cluster has been successfully provisioned and it is ready to be installed.
	Provisioned
	// InstallInProgress indicates that the Nalej platform is being installed on the target cluster.
	InstallInProgress
	// Installed indicates that the Nalej platform has been successfully deployed on the target cluster.
	Installed
	// Scaling indicates that the cluster is now being modified in terms of number of available nodes.
	Scaling
	// Failure indicates that a process related to provisioning or installation has failed.
	Failure
	// Uninstalling indicates that the cluster is now being uninstalled and the platform is being undeployed.
	Uninstalling
	// Decomisioning indicates that the hardware resources are being freed.
	Decomissioning
)

// ClusterStateToGRPC translates a ClusterState into the gRPC equivalent.
var ClusterStateToGRPC = map[ClusterState]grpc_infrastructure_go.ClusterState{
	Unknown:           grpc_infrastructure_go.ClusterState_UNKNOWN,
	Provisioning:      grpc_infrastructure_go.ClusterState_PROVISIONING,
	Provisioned:       grpc_infrastructure_go.ClusterState_PROVISIONED,
	InstallInProgress: grpc_infrastructure_go.ClusterState_INSTALL_IN_PROGRESS,
	Installed:         grpc_infrastructure_go.ClusterState_INSTALLED,
	Scaling:           grpc_infrastructure_go.ClusterState_SCALING,
	Failure:           grpc_infrastructure_go.ClusterState_FAILURE,
	Uninstalling:      grpc_infrastructure_go.ClusterState_UNINSTALLING,
	Decomissioning:    grpc_infrastructure_go.ClusterState_DECOMISIONING,
}

// ClusterStateFromGRPC translates a gRPC state into a ClusterState
var ClusterStateFromGRPC = map[grpc_infrastructure_go.ClusterState]ClusterState{
	grpc_infrastructure_go.ClusterState_UNKNOWN:             Unknown,
	grpc_infrastructure_go.ClusterState_PROVISIONING:        Provisioning,
	grpc_infrastructure_go.ClusterState_PROVISIONED:         Provisioned,
	grpc_infrastructure_go.ClusterState_INSTALL_IN_PROGRESS: InstallInProgress,
	grpc_infrastructure_go.ClusterState_INSTALLED:           Installed,
	grpc_infrastructure_go.ClusterState_SCALING:             Scaling,
	grpc_infrastructure_go.ClusterState_FAILURE:             Failure,
	grpc_infrastructure_go.ClusterState_UNINSTALLING:        Uninstalling,
	grpc_infrastructure_go.ClusterState_DECOMISIONING:       Decomissioning,
}

// DefaultMillicoresConversionFactor Default value for the millicores conversion factor
// TODO To be parametrized depending on the subjacent platform
const DefaultMillicoresConversionFactor = 1.0

// Cluster entity representing a collection of nodes that supports applicaiton orchestration. This
// abstraction is used for monitoring and orchestration purposes.
type Cluster struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// ClusterId with the cluster identifier.
	ClusterId string `json:"cluster_id,omitempty"`
	// Name of the cluster.
	Name string `json:"name,omitempty"`
	// Type of cluster.
	ClusterType ClusterType `json:"cluster_type,omitempty"`
	// Hostname of the cluster master.
	Hostname string `json:"hostname,omitempty"`
	// ControlPlaneHostname with the hostname to access K8s API.
	ControlPlaneHostname string `json:"control_plane_hostname,omitempty"`
	// Multitenant support definition.
	Multitenant MultitenantSupport `json:"multitenant,omitempty"`
	// Status of the cluster based on monitoring information.
	Status ClusterStatus `json:"status,omitempty"`
	// Labels for the cluster.
	Labels map[string]string `json:"labels,omitempty"`
	// Cordon flags to signal conductor not to schedule apps in the cluster.
	// Deprecated: will be removed TODO
	Cordon bool `json:"cordon,omitempty"`
	// Cluster watch information
	ClusterWatch ClusterWatchInfo `json:"cluster_watch,omitempty"`
	// Last alive timestamp
	LastAliveTimestamp int64 `json:"last_alive_timestamp,omitempty"`
	// MillicoresConversionFactor Contains a conversion factor for the millicores measurement that k8s exposes as it is platform dependent.
	MillicoresConversionFactor float64 `json:"millicores_conversion_factor,omitempty"`
	// State of the cluster with respect to provisioning and installation.
	State ClusterState `json:"cluster_state,omitempty"`
}

// The cluster watcher contains information to ensure the connectivity between clusters. This data
// is required by Cilium and other connectivity platforms.
type ClusterWatchInfo struct {
	// Name of the cluster
	Name string `json:"name,omitempty" cql:"name"`
	// Organization id
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// ClusterId of the cluster
	ClusterId string `json:"cluster_id,omitempty" cql:"cluster_id"`
	// IP of the cluster
	Ip string `json:"ip,omitempty" cql:"ip"`
	// Network type
	NetworkType NetworkType `json:"network_type,omitempty" cql:"network_type"`
	// Cilium data
	CiliumData CiliumCerts `json:"cilium_certs,omitempty" cql:"cilium_certs"`
	// Istio data
	IstioData IstioCerts `json:"istio_certs,omitempty" cql:"istio_certs"`
}

type CiliumCerts struct {
	// CiliumId ClusterId for the node
	CiliumId string `json:"cilium_id,omitempty" cql:"cilium_id"`
	// Cilium etcd-client-ca.crt certification authority to be used
	CiliumEtcdCaCrt string `json:"cilium_etcd_ca_crt,omitempty" cql:"cilium_etcd_ca_crt"`
	// Cilium etcd-client.crt certificate
	CiliumEtcdCrt string `json:"cilium_etcd_crt,omitempty" cql:"cilium_etcd_crt"`
	// Cilium client public key
	CiliumEtcdKey string `json:"cilium_etcd_key,omitempty" cql:"cilium_etcd_key"`
}

func (c *CiliumCerts) toGRPC() *grpc_cluster_watcher_go.ClusterWatchInfo_Cilium {
	return &grpc_cluster_watcher_go.ClusterWatchInfo_Cilium{
		Cilium: &grpc_cluster_watcher_go.CiliumCerts{
			CiliumEtcdKey:   c.CiliumEtcdKey,
			CiliumEtcdCrt:   c.CiliumEtcdCrt,
			CiliumEtcdCaCrt: c.CiliumEtcdCaCrt,
			CiliumId:        c.CiliumId,
		},
	}
}

func NewCiliumCertsFromGRPC(cilium *grpc_cluster_watcher_go.ClusterWatchInfo_Cilium) CiliumCerts {
	return CiliumCerts{
		CiliumId:        cilium.Cilium.CiliumId,
		CiliumEtcdCaCrt: cilium.Cilium.CiliumEtcdCaCrt,
		CiliumEtcdCrt:   cilium.Cilium.CiliumEtcdCrt,
		CiliumEtcdKey:   cilium.Cilium.CiliumEtcdKey,
	}
}

type IstioCerts struct {
	// Cluster name
	ClusterName string `json:"cluster_name,omitempty" cql:"cluster_name"`
	// Server name
	ServerName string `json:"server_name,omitempty" cql:"server_name"`
	// CA certificate
	CaCert string `json:"ca_cert,omitempty" cql:"ca_cert"`
	// Token
	Token string `json:"token,omitempty" cql:"cluster_token"`
}

func NewIstioCertsFromGRPC(istio *grpc_cluster_watcher_go.ClusterWatchInfo_Istio) IstioCerts {
	return IstioCerts{
		ClusterName: istio.Istio.ClusterName,
		ServerName:  istio.Istio.ServerName,
		Token:       istio.Istio.Token,
		CaCert:      istio.Istio.CaCert,
	}
}

func (c *IstioCerts) ToGRPC() *grpc_cluster_watcher_go.ClusterWatchInfo_Istio {
	return &grpc_cluster_watcher_go.ClusterWatchInfo_Istio{
		Istio: &grpc_cluster_watcher_go.IstioCerts{
			CaCert:      c.CaCert,
			Token:       c.Token,
			ServerName:  c.ServerName,
			ClusterName: c.ClusterName,
		},
	}
}

func NewClusterWatchInfo(name string, organizationId, clusterId string, ip string, networkType NetworkType,
	ciliumData CiliumCerts, istioData IstioCerts) *ClusterWatchInfo {
	return &ClusterWatchInfo{
		Name:           name,
		OrganizationId: organizationId,
		ClusterId:      clusterId,
		Ip:             ip,
		NetworkType:    networkType,
		CiliumData:     ciliumData,
		IstioData:      istioData,
	}
}

func (c *ClusterWatchInfo) ToGRPC() *grpc_cluster_watcher_go.ClusterWatchInfo {
	toReturn := &grpc_cluster_watcher_go.ClusterWatchInfo{
		Name:           c.Name,
		ClusterId:      c.ClusterId,
		Ip:             c.Ip,
		NetworkType:    NetworkTypeToGRPC[c.NetworkType],
		OrganizationId: c.OrganizationId,
	}
	switch c.NetworkType {
	case NetworkTypeCilium:
		toReturn.Security = c.CiliumData.toGRPC()
	case NetworkTypeIstio:
		toReturn.Security = c.IstioData.ToGRPC()
	}
	return toReturn
}

func ClusterWatchInfoFromGRPC(clusterWatch *grpc_cluster_watcher_go.ClusterWatchInfo) *ClusterWatchInfo {

	var cilium CiliumCerts
	var istio IstioCerts
	switch x := clusterWatch.Security.(type) {
	case *grpc_cluster_watcher_go.ClusterWatchInfo_Istio:
		istio = NewIstioCertsFromGRPC(x)
	case *grpc_cluster_watcher_go.ClusterWatchInfo_Cilium:
		cilium = NewCiliumCertsFromGRPC(x)
	}

	return NewClusterWatchInfo(clusterWatch.Name, clusterWatch.OrganizationId, clusterWatch.ClusterId,
		clusterWatch.Ip, NetworkTypeFromGRPC[clusterWatch.NetworkType], cilium, istio)
}

func NewCluster(organizationID string, name string, description string, hostname string, controlPlaneHostname string) *Cluster {
	uuid := GenerateUUID()
	return &Cluster{
		OrganizationId:       organizationID,
		ClusterId:            uuid,
		Name:                 name,
		ClusterType:          KubernetesCluster,
		Hostname:             hostname,
		ControlPlaneHostname: controlPlaneHostname,
		Multitenant:          MultitenantYes,
		Status:               ClusterStatusUnknown,
		Labels:               make(map[string]string, 0),
		Cordon:               false,
		// ClusterWatch: this is filled by external components
		// LastAliveTimestamp: this is filled by external components
	}
}

func NewClusterFromGRPC(addClusterRequest *grpc_infrastructure_go.AddClusterRequest) *Cluster {
	uuid := GenerateUUID()
	return &Cluster{
		OrganizationId:             addClusterRequest.OrganizationId,
		ClusterId:                  uuid,
		Name:                       addClusterRequest.Name,
		ClusterType:                KubernetesCluster,
		Hostname:                   addClusterRequest.Hostname,
		ControlPlaneHostname:       addClusterRequest.ControlPlaneHostname,
		Multitenant:                MultitenantYes,
		Status:                     ClusterStatusUnknown,
		Labels:                     addClusterRequest.Labels,
		Cordon:                     false,
		State:                      Provisioning,
		MillicoresConversionFactor: DefaultMillicoresConversionFactor,
		// ClusterWatch:
		// LastAliveTimestamp:
	}
}

func (c *Cluster) ToGRPC() *grpc_infrastructure_go.Cluster {
	clusterType := ClusterTypeToGRPC[c.ClusterType]
	multitenant := MultitenantSupportToGRPC[c.Multitenant]
	status := ClusterStatusToGRPC[c.Status]
	state := ClusterStateToGRPC[c.State]
	return &grpc_infrastructure_go.Cluster{
		OrganizationId:             c.OrganizationId,
		ClusterId:                  c.ClusterId,
		Name:                       c.Name,
		ClusterType:                clusterType,
		Hostname:                   c.Hostname,
		ControlPlaneHostname:       c.ControlPlaneHostname,
		Multitenant:                multitenant,
		ClusterStatus:              status,
		Labels:                     c.Labels,
		ClusterWatch:               c.ClusterWatch.ToGRPC(),
		LastAliveTimestamp:         c.LastAliveTimestamp,
		State:                      state,
		MillicoresConversionFactor: c.MillicoresConversionFactor,
	}
}

func (c *Cluster) ApplyUpdate(updateRequest grpc_infrastructure_go.UpdateClusterRequest) {
	if updateRequest.UpdateName {
		c.Name = updateRequest.Name
	}
	if updateRequest.UpdateHostname {
		c.Hostname = updateRequest.Hostname
	}
	if updateRequest.UpdateControlPlaneHostname {
		c.ControlPlaneHostname = updateRequest.ControlPlaneHostname
	}
	if updateRequest.AddLabels {
		if c.Labels == nil {
			c.Labels = make(map[string]string, 0)
		}
		for k, v := range updateRequest.Labels {
			c.Labels[k] = v
		}
	}
	if updateRequest.RemoveLabels {
		for k, _ := range updateRequest.Labels {
			delete(c.Labels, k)
		}
	}

	if updateRequest.UpdateStatus {
		c.Status = ClusterStatusFromGRPC[updateRequest.Status]
	}

	if updateRequest.UpdateClusterWatch {
		c.ClusterWatch = *ClusterWatchInfoFromGRPC(updateRequest.ClusterWatchInfo)
	}

	if updateRequest.UpdateLastClusterTimestamp {
		c.LastAliveTimestamp = updateRequest.LastClusterTimestamp
	}

	if updateRequest.UpdateClusterState {
		c.State = ClusterStateFromGRPC[updateRequest.State]
	}

	if updateRequest.UpdateMillicoresConversionFactor {
		c.MillicoresConversionFactor = updateRequest.MillicoresConversionFactor
	}
}

func ValidAddClusterRequest(addClusterRequest *grpc_infrastructure_go.AddClusterRequest) derrors.Error {
	if addClusterRequest.RequestId == "" {
		return derrors.NewInvalidArgumentError(emptyRequestId)
	}
	if addClusterRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if addClusterRequest.ClusterId != "" {
		return derrors.NewInvalidArgumentError("cluster_id must be empty, and generated by this component")
	}
	return nil
}

func ValidUpdateClusterRequest(updateClusterRequest *grpc_infrastructure_go.UpdateClusterRequest) derrors.Error {
	if updateClusterRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if updateClusterRequest.ClusterId == "" {
		return derrors.NewInvalidArgumentError(emptyClusterId)
	}
	return nil
}

func ValidRemoveClusterRequest(removeClusterRequest *grpc_infrastructure_go.RemoveClusterRequest) derrors.Error {
	if removeClusterRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if removeClusterRequest.ClusterId == "" {
		return derrors.NewInvalidArgumentError(emptyClusterId)
	}
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
