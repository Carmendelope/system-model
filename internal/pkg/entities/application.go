/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package entities

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-go"
	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/util/validation"
	"strings"
)

// DefaultEndPointInstance is used when the endpoint recived from GRPC has no endpoint
var DefaultEndpointInstance = &grpc_application_go.EndpointInstance{
	EndpointInstanceId: "",
	Type: grpc_application_go.EndpointType_IS_ALIVE,
	Fqdn: "",
	Port: 0,
}

// regular expresion for IP:port address
var IPAddressRegExp = string("(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])(.(25[0-5]|2[0-4][0-9]|[01]?[0-9]?[0-9])){3}(:(6553[0-5]|655[0-2][0-9]|65[0-4][0-9]{2}|6[0-4][0-9]{3}|[0-5]?([0-9]){0,3}[0-9]))?")

// characters of service_group_instance_id and service_instance_id to create gloabl_fqdn
const InstPrefixLength = 6
// characters of organization_id to create gloabl_fqdn
const OrgPrefixLength = 8


// Enumerate with the type of instances we can deploy in the system.
type InstanceType int32

const (
	ServiceInstanceType InstanceType = iota + 1
	ServiceGroupInstanceType
)

var InstanceTypeToGRPC = map[InstanceType]grpc_application_go.InstanceType{
	ServiceInstanceType : grpc_application_go.InstanceType_SERVICE_INSTANCE,
	ServiceGroupInstanceType : grpc_application_go.InstanceType_SERVICE_GROUP_INSTANCE,
}

var InstanceTypeFromGRPC = map[grpc_application_go.InstanceType]InstanceType {
	grpc_application_go.InstanceType_SERVICE_INSTANCE:ServiceInstanceType,
	grpc_application_go.InstanceType_SERVICE_GROUP_INSTANCE:ServiceGroupInstanceType,
}

type PortAccess int

const (
	AllAppServices PortAccess = iota + 1
	AppServices
	Public
	DeviceGroup
)

var PortAccessToGRPC = map[PortAccess]grpc_application_go.PortAccess{
	AllAppServices: grpc_application_go.PortAccess_ALL_APP_SERVICES,
	AppServices:    grpc_application_go.PortAccess_APP_SERVICES,
	Public:         grpc_application_go.PortAccess_PUBLIC,
	DeviceGroup:    grpc_application_go.PortAccess_DEVICE_GROUP,
}

var PortAccessFromGRPC = map[grpc_application_go.PortAccess]PortAccess{
	grpc_application_go.PortAccess_ALL_APP_SERVICES: AllAppServices,
	grpc_application_go.PortAccess_APP_SERVICES:     AppServices,
	grpc_application_go.PortAccess_PUBLIC:           Public,
	grpc_application_go.PortAccess_DEVICE_GROUP:     DeviceGroup,
}

type CollocationPolicy int

const (
	SameCluster CollocationPolicy = iota + 1
	SeparateClusters
)

var CollocationPolicyToGRPC = map[CollocationPolicy]grpc_application_go.CollocationPolicy{
	SameCluster:      grpc_application_go.CollocationPolicy_SAME_CLUSTER,
	SeparateClusters: grpc_application_go.CollocationPolicy_SEPARATE_CLUSTERS,
}

var CollocationPolicyFromGRPC = map[grpc_application_go.CollocationPolicy]CollocationPolicy{
	grpc_application_go.CollocationPolicy_SAME_CLUSTER:      SameCluster,
	grpc_application_go.CollocationPolicy_SEPARATE_CLUSTERS: SeparateClusters,
}

// -- SecurityRule -- //
type SecurityRule struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty" cql:"app_descriptor_id"`
	// RuleId with the security rule identifier.
	RuleId string `json:"rule_id,omitempty" cql:"rule_id"`
	// Name of the security rule.
	Name string `json:"name,omitempty" cql:"name"`
	// TargetServiceGroupName defining the name of the service to be accessed.
	TargetServiceGroupName string `json:"target_service_group_name,omitempty" cql:"target_service_group_name"`
	// TargetServiceName name of the service belonging to be source group mentioned above to be accessed.
	TargetServiceName string `json:"target_service_name,omitempty" cql:"target_service_name"`
	// TargetPort defining the port that is affected by the current rule.
	TargetPort int32 `json:"target_port,omitempty" cql:"target_port"`
	// Access level to that port defining who can access the port.
	Access PortAccess `json:"access,omitempty" cql:"access"`
	// Name of the service group
	AuthServiceGroupName string `json:"auth_service_group_name,omitempty" cql:"auth_service_group_name"`
	// AuthServices defining a list of services that can access the port.
	AuthServices []string `json:"auth_services,omitempty" cql:"auth_services"`
	// DeviceGroupIds defining a list of device groups that can access the port.
	DeviceGroupIds []string `json:"device_groups,omitempty" cql:"device_group_ids"`
	// DeviceGroupIds defining a list of device groups that can access the port.
	DeviceGroupNames []string `json:"device_group_names,omitempty" cql:"device_group_names"`
}

// NewSecurityRuleFromGRPC converts a grpc_application_go.SecurityRule into SecurityRule
// deviceGroupIds is a map of deviceGroupIds indexed by deviceGroupNames (it contains ALL the devices in the appDescriptor)
func NewSecurityRuleFromGRPC(organizationID string, appDescriptorID string, rule *grpc_application_go.SecurityRule, deviceGroupIds map[string]string) (*SecurityRule, derrors.Error) {
	if rule == nil {
		return nil, nil
	}

	ids := make ([]string, 0)
	if rule != nil {
		for _, name := range rule.DeviceGroupNames {
			deviceGroupId, exists := deviceGroupIds[name]
			if ! exists {
				log.Error().Str("deviceName", name).Msg("Device id not found")
				return nil, derrors.NewNotFoundError("device group id").WithParams(name)
			} else {
				ids = append(ids, deviceGroupId)
			}
		}
	}else{
		log.Debug().Msg("rule empty")
	}

	uuid := GenerateUUID()
	access := PortAccessFromGRPC[rule.Access]
	return &SecurityRule{
		OrganizationId:  		organizationID,
		AppDescriptorId: 		appDescriptorID,
		RuleId:          		uuid,
		Name:            		rule.Name,
		TargetServiceGroupName:	rule.TargetServiceGroupName,
		TargetServiceName: 		rule.TargetServiceName,
		TargetPort: 			rule.TargetPort,
		Access:          		access,
		AuthServiceGroupName: 	rule.AuthServiceGroupName,
		AuthServices:    		rule.AuthServices,
		DeviceGroupNames:  		rule.DeviceGroupNames,
		DeviceGroupIds:         ids,
	}, nil
}


// NewSecurityRuleFromGRPC converts a grpc_application_go.SecurityRule into SecurityRule
// TODO revisit if it is necessary to have the other version of this function running
func NewSecurityRuleFromInstantiatedGRPC(rule *grpc_application_go.SecurityRule) (*SecurityRule, derrors.Error) {
	if rule == nil {
		return nil, nil
	}

	access := PortAccessFromGRPC[rule.Access]
	return &SecurityRule{
		OrganizationId:  		rule.OrganizationId,
		AppDescriptorId: 		rule.AppDescriptorId,
		RuleId:          		rule.RuleId,
		Name:            		rule.Name,
		TargetServiceGroupName:	rule.TargetServiceGroupName,
		TargetServiceName: 		rule.TargetServiceName,
		TargetPort: 			rule.TargetPort,
		Access:          		access,
		AuthServiceGroupName: 	rule.AuthServiceGroupName,
		AuthServices:    		rule.AuthServices,
		DeviceGroupNames:  		rule.DeviceGroupNames,
		DeviceGroupIds:         rule.DeviceGroupIds,
	}, nil
}

func (sr *SecurityRule) ToGRPC() *grpc_application_go.SecurityRule {
	access, _ := PortAccessToGRPC[sr.Access]
	return &grpc_application_go.SecurityRule{
		OrganizationId:  		sr.OrganizationId,
		AppDescriptorId: 		sr.AppDescriptorId,
		RuleId:          		sr.RuleId,
		Name:            		sr.Name,
		TargetServiceGroupName:	sr.TargetServiceGroupName,
		TargetServiceName: 		sr.TargetServiceName,
		TargetPort: 			sr.TargetPort,
		Access:       			access,
		AuthServiceGroupName: 	sr.AuthServiceGroupName,
		AuthServices: 			sr.AuthServices,
		DeviceGroupNames:		sr.DeviceGroupNames,
		DeviceGroupIds:  		sr.DeviceGroupIds,
	}
}

// ServiceGroupDeploymentSpecs -- //
type ServiceGroupDeploymentSpecs struct {
	// How many times this service group must be replicated
	Replicas int32 `json:"replicas,omitempty" cql:"replicas"`
	// Indicate if this service group must be replicated in every cluster
	MultiClusterReplica  bool     `json:"multi_cluster_replica,omitempty" cql:"multi_cluster_replica"`
	// DeploymentSelectors defines a key-value map of deployment selectors
	DeploymentSelectors  map[string]string `json:"deployment_selectors,omitempty" cql:"deployment_selectors"`
}

func NewServiceGroupDeploymentSpecsFromGRPC(specs * grpc_application_go.ServiceGroupDeploymentSpecs) * ServiceGroupDeploymentSpecs{
	if specs == nil {
		return nil
	}
	return &ServiceGroupDeploymentSpecs{
		Replicas:            specs.Replicas,
		MultiClusterReplica: specs.MultiClusterReplica,
		DeploymentSelectors: specs.DeploymentSelectors,
	}
}

func (sp * ServiceGroupDeploymentSpecs) ToGRPC() *grpc_application_go.ServiceGroupDeploymentSpecs  {
	if sp == nil {
		return nil
	}
	return &grpc_application_go.ServiceGroupDeploymentSpecs{
		Replicas:          sp.Replicas,
		MultiClusterReplica:  sp.MultiClusterReplica,
		DeploymentSelectors:  sp.DeploymentSelectors,
	}
}

// -- ServiceGroup -- //
type ServiceGroup struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty" cql:"app_descriptor_id"`
	// ServiceGroupId with the group identifier.
	ServiceGroupId string `json:"service_group_id,omitempty" cql:"service_group_id"`
	// Name of the service group.
	Name string `json:"name,omitempty" cql:"name"`
	// Services defining a list of service identifiers that belong to the group.
	Services []Service `json:"services,omitempty" cql:"services"`
	// Policy indicating the deployment collocation policy.
	Policy CollocationPolicy `json:"policy,omitempty" cql:"policy"`
	// Particular deployment specs for this service
	Specs *ServiceGroupDeploymentSpecs `json:"specs,omitempty" cql:"specs"`
	// Labels defined by the user.
	Labels map[string]string `json:"labels,omitempty" cql:"labels"`
}

func NewServiceGroupFromGRPC(organizationID string, appDescriptorID string, group * grpc_application_go.ServiceGroup) * ServiceGroup {
	if group == nil {
		return nil
	}
	id := GenerateUUID()

	services := make ([]Service, 0)
	for _, service := range group.Services {
		services = append(services, *NewServiceFromGRPC(organizationID, appDescriptorID, id, service))
	}

	policy, _ := CollocationPolicyFromGRPC[group.Policy]
	return &ServiceGroup{
		OrganizationId:		organizationID,
		AppDescriptorId:	appDescriptorID,
		ServiceGroupId: 	id,
		Name : 				group.Name,
		Services: 			services,
		Policy: 			policy,
		Specs: 				NewServiceGroupDeploymentSpecsFromGRPC(group.Specs),
		Labels: 			group.Labels,
	}
}

func (sg *ServiceGroup) ToGRPC() *grpc_application_go.ServiceGroup {

	services := make([]*grpc_application_go.Service, 0)
	for _, service := range sg.Services {
		services = append(services, service.ToGRPC() )
	}

	policy, _ := CollocationPolicyToGRPC[sg.Policy]
	return &grpc_application_go.ServiceGroup{
		OrganizationId:  	sg.OrganizationId,
		AppDescriptorId:	sg.AppDescriptorId,
		ServiceGroupId:  	sg.ServiceGroupId,
		Name:            	sg.Name,
		Services:        	services,
		Policy:          	policy,
		Specs:           	sg.Specs.ToGRPC(),
		Labels: 		 	sg.Labels,
	}
}

// -- InstanceMetadata -- //
type InstanceMetadata struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty" cql:"app_descriptor_id"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty" cql:"app_instance_id"`
	// ServiceGroupId with the service group id this entity belongs to.
	ServiceGroupId string `json:"service_group_id,omitempty" cql:"service_group_id"`
	// Identifier of the monitored entity
	MonitoredInstanceId string `json:"monitored_instance_id,omitempty" cql:"monitored_instance_id"`
	// Type of instance this metadata refers to
	Type InstanceType `json:"type,omitempty" cql:"type"`
	// List of instances supervised byu this metadata structure
	InstancesId []string `json:"instances_id,omitempty" cql:"instance_id"`
	// Number of desired replicas specified in the descriptor
	DesiredReplicas int32 `json:"desired_replicas,omitempty" cql:"desired_replicas"`
	// Number of available replicas for this instance
	AvailableReplicas int32 `json:"available_replicas,omitempty" cql:"available_replicas"`
	// Number of unavaiable replicas for this descriptor
	UnavailableReplicas int32 `json:"unavailable_replicas,omitempty" cql:"unavailable_replicas"`
	// Status of every item monitored by this metadata entry
	Status map[string]ServiceStatus `json:"status,omitempty" cql:"status"`
	// Relevant information for every monitored instance
	Info map[string]string `json:"info,omitempty" cql:"info"`
}
func (md *InstanceMetadata) ToGRPC() *grpc_application_go.InstanceMetadata {
	if md == nil {
		return nil
	}
	status := make (map[string]grpc_application_go.ServiceStatus, 0)
	for key, value := range md.Status{
		status[key] = ServiceStatusToGRPC[value]
	}

	return &grpc_application_go.InstanceMetadata{
		OrganizationId: md.OrganizationId,
		AppDescriptorId: md.AppDescriptorId,
		AppInstanceId: md.AppInstanceId,
		ServiceGroupId: md.ServiceGroupId,
		MonitoredInstanceId: md.MonitoredInstanceId,
		Type: InstanceTypeToGRPC[md.Type],
		InstancesId: md.InstancesId,
		DesiredReplicas: md.DesiredReplicas,
		AvailableReplicas: md.AvailableReplicas,
		UnavailableReplicas: md.UnavailableReplicas,
		Status: status,
		Info: md.Info,
	}
}

func NewMetadataFromGRPC (metadata * grpc_application_go.InstanceMetadata) * InstanceMetadata {
	if metadata == nil {
		return nil
	}

	status := make (map[string]ServiceStatus, 0)
	for key, value := range metadata.Status{
		status[key] = ServiceStatusFromGRPC[value]
	}

	return &InstanceMetadata{
		OrganizationId: metadata.OrganizationId,
		AppDescriptorId: metadata.AppDescriptorId,
		AppInstanceId: metadata.AppInstanceId,
		ServiceGroupId: metadata.ServiceGroupId,
		MonitoredInstanceId: metadata.MonitoredInstanceId,
		Type: InstanceTypeFromGRPC[metadata.Type],
		InstancesId: metadata.InstancesId,
		DesiredReplicas: metadata.DesiredReplicas,
		AvailableReplicas: metadata.AvailableReplicas,
		UnavailableReplicas: metadata.UnavailableReplicas,
		Status: status,
		Info: metadata.Info,
	}
}

// -- ServiceGroupInstance -- //
type ServiceGroupInstance struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty" cql:"app_descriptor_id"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty" cql:"app_instance_id"`
	// ServiceGroupId with the group identifier.
	ServiceGroupId string `json:"service_group_id,omitempty" cql:"service_group_id"`
	// Unique identifier for this instance
	ServiceGroupInstanceId string `json:"service_group_instance_id,omitempty" cql:"service_group_instance_id" `
	// Name of the service group.
	Name string `json:"name,omitempty" cql:"name"`
	// ServicesInstances defining a list of service identifiers that belong to the group.
	ServiceInstances []ServiceInstance `json:"service_instances,omitempty" cql:"service_instances"`
	// Policy indicating the deployment collocation policy.
	Policy CollocationPolicy `json:"policy,omitempty" cql:"policy"`
	// The status for this service group instance will be the worst status of its services
	Status ServiceStatus `json:"status,omitempty" cql:"status"`
	// Metadata for this service group
	Metadata * InstanceMetadata `json:"metadata,omitempty" cql:"metadata"`
	// Particular deployment specs for this service
	Specs * ServiceGroupDeploymentSpecs `json:"specs,omitempty" cql:"specs"`
	// Labels defined by the user.
	Labels map[string]string `json:"labels,omitempty" cql:"labels"`
	// GlobalFqdn
	GlobalFqdn  []string `json:"global_fqdn,omitempty"`

}

func (sgi *ServiceGroupInstance) ToGRPC() *grpc_application_go.ServiceGroupInstance {

	services := make ([]*grpc_application_go.ServiceInstance, 0)
	for _, instance := range sgi.ServiceInstances {
		services = append(services, instance.ToGRPC() )
	}

	policy, _ := CollocationPolicyToGRPC[sgi.Policy]
	return &grpc_application_go.ServiceGroupInstance{
		OrganizationId:     sgi.OrganizationId,
		AppDescriptorId:    sgi.AppDescriptorId,
		AppInstanceId:      sgi.AppInstanceId,
		ServiceGroupId:     sgi.ServiceGroupId,
		ServiceGroupInstanceId: sgi.ServiceGroupInstanceId,
		Name:               sgi.Name,
		ServiceInstances:	services,
		Policy:             policy,
		Status:				ServiceStatusToGRPC[sgi.Status],
		Metadata:			sgi.Metadata.ToGRPC(),
		Specs: 				sgi.Specs.ToGRPC(),
		Labels: 			sgi.Labels,
		GlobalFqdn:         sgi.GlobalFqdn,
	}
}

// Use this function to fill the metadata object for an initial status.
func (sgi *ServiceGroupInstance) FillMetadata(totalReplicas int) {
	// fill the list of ids for the monitored instances
	instancesId := make([]string, len(sgi.ServiceInstances))
	statuses := make(map[string]ServiceStatus,len(sgi.ServiceInstances))
	info := make(map[string]string,len(sgi.ServiceInstances))
	for i, s := range sgi.ServiceInstances {
		instancesId[i] = s.ServiceInstanceId
		statuses[s.ServiceInstanceId] = ServiceScheduled
		info[s.ServiceInstanceId] = ""
	}

	metadata := &InstanceMetadata{
		AppInstanceId: sgi.AppInstanceId,
		ServiceGroupId: sgi.ServiceGroupId,
		AppDescriptorId: sgi.AppDescriptorId,
		OrganizationId: sgi.OrganizationId,
		MonitoredInstanceId: sgi.ServiceGroupInstanceId,
		Type: ServiceGroupInstanceType,
		InstancesId: instancesId,
		Status: statuses,
		DesiredReplicas: int32(totalReplicas),
		AvailableReplicas: 0,
		UnavailableReplicas: 0,
		Info: info,
	}
	sgi.Metadata = metadata
}

// ----

type ServiceType int32

const (
	DockerService ServiceType = iota + 1
)

var ServiceTypeToGRPC = map[ServiceType]grpc_application_go.ServiceType{
	DockerService: grpc_application_go.ServiceType_DOCKER,
}

var ServiceTypeFromGRPC = map[grpc_application_go.ServiceType]ServiceType{
	grpc_application_go.ServiceType_DOCKER: DockerService,
}

// -- ImageCredentials -- //
type ImageCredentials struct {
	Username string `json:"username,omitempty" cql:"username"`
	Password string `json:"password,omitempty" cql:"password"`
	Email    string `json:"email,omitempty" cql:"email"`
	DockerRepository  string   `json:"docker_repository,omitempty" cql:"docker_repository"`
}

func NewImageCredentialsFromGRPC(credentials * grpc_application_go.ImageCredentials) *ImageCredentials {
	if credentials == nil {
		return nil
	}
	return &ImageCredentials{
		Username: credentials.Username,
		Password: credentials.Password,
		Email:    credentials.Email,
		DockerRepository: credentials.DockerRepository,
	}
}

func (ic *ImageCredentials) ToGRPC() *grpc_application_go.ImageCredentials {
	if ic == nil {
		return nil
	}
	return &grpc_application_go.ImageCredentials{
		Username: ic.Username,
		Password: ic.Password,
		Email:    ic.Email,
		DockerRepository: ic.DockerRepository,
	}
}

// -- DeploySpecs -- //
type DeploySpecs struct {
	Cpu      int64 `json:"cpu,omitempty" cql:"cpu"`
	Memory   int64 `json:"memory,omitempty" cql:"memory"`
	Replicas int32 `json:"replicas,omitempty" cql:"replicas"`
}

func NewDeploySpecsFromGRPC(specs * grpc_application_go.DeploySpecs) * DeploySpecs {
	if specs == nil {
		return nil
	}
	return &DeploySpecs{
		Cpu:      specs.Cpu,
		Memory:   specs.Memory,
		Replicas: specs.Replicas,
	}
}

func (ds *DeploySpecs) ToGRPC() *grpc_application_go.DeploySpecs {

	spec := &grpc_application_go.DeploySpecs{
		Replicas: 1,
	}

	if ds != nil {
		spec.Cpu = ds.Cpu
		spec.Memory = ds.Memory
		spec.Replicas = ds.Replicas
	}

	return spec
}

type StorageType int32

const (
	Ephemeral StorageType = iota + 1
	ClusterLocal
	ClusterReplica
	CloudPersistent
)

var StorageTypeToGRPC = map[StorageType]grpc_application_go.StorageType{
	Ephemeral:       grpc_application_go.StorageType_EPHEMERAL,
	ClusterLocal:    grpc_application_go.StorageType_CLUSTER_LOCAL,
	ClusterReplica:  grpc_application_go.StorageType_CLUSTER_REPLICA,
	CloudPersistent: grpc_application_go.StorageType_CLOUD_PERSISTENT,
}

var StorageTypeFromGRPC = map[grpc_application_go.StorageType]StorageType{
	grpc_application_go.StorageType_EPHEMERAL:        Ephemeral,
	grpc_application_go.StorageType_CLUSTER_LOCAL:    ClusterLocal,
	grpc_application_go.StorageType_CLUSTER_REPLICA:  ClusterReplica,
	grpc_application_go.StorageType_CLOUD_PERSISTENT: CloudPersistent,
}

// -- Storage -- //
type Storage struct {
	Size      int64       `json:"size,omitempty" cql:"size"`
	MountPath string      `json:"mount_path,omitempty" cql:"mount_path"`
	Type      StorageType `json:"type,omitempty" cql:"type"`
}

func NewStorageFromGRPC(storage * grpc_application_go.Storage) * Storage{
	if storage == nil {
		return nil
	}
	storageType, _ := StorageTypeFromGRPC[storage.Type]
	return &Storage{
		Size:      storage.Size,
		MountPath: storage.MountPath,
		Type:      storageType,
	}
}

func (s *Storage) ToGRPC() *grpc_application_go.Storage {
	convertedType, _ := StorageTypeToGRPC[s.Type]
	return &grpc_application_go.Storage{
		Size:      s.Size,
		MountPath: s.MountPath,
		Type:      convertedType,
	}
}

type EndpointType int

const (
	IsAlive EndpointType = iota + 1
	Rest
	Web
	Prometheus
	Ingestion
)

var EndpointTypeToGRPC = map[EndpointType]grpc_application_go.EndpointType{
	IsAlive:    grpc_application_go.EndpointType_IS_ALIVE,
	Rest:       grpc_application_go.EndpointType_REST,
	Web:        grpc_application_go.EndpointType_WEB,
	Prometheus: grpc_application_go.EndpointType_PROMETHEUS,
	Ingestion:	grpc_application_go.EndpointType_INGESTION,
}

var EndpointTypeFromGRPC = map[grpc_application_go.EndpointType]EndpointType{
	grpc_application_go.EndpointType_IS_ALIVE:   IsAlive,
	grpc_application_go.EndpointType_REST:       Rest,
	grpc_application_go.EndpointType_WEB:        Web,
	grpc_application_go.EndpointType_PROMETHEUS: Prometheus,
	grpc_application_go.EndpointType_INGESTION:  Ingestion,
}

// -- Endpoint -- //
type Endpoint struct {
	Type EndpointType `json:"type,omitempty" cql:"type"`
	Path string       `json:"path,omitempty" cql:"path"`
}

func NewEndpointFromGRPC( endpoint * grpc_application_go.Endpoint) * Endpoint {
	if endpoint == nil {
		return nil
	}
	endpointType, _ := EndpointTypeFromGRPC[endpoint.Type]
	return &Endpoint{
		Type: endpointType,
		Path: endpoint.Path,
	}
}

func (e *Endpoint) ToGRPC() *grpc_application_go.Endpoint {
	convertedType, _ := EndpointTypeToGRPC[e.Type]
	return &grpc_application_go.Endpoint{
		Type: convertedType,
		Path: e.Path,
	}
}

// -- Port -- //
type Port struct {
	Name         string     `json:"name,omitempty" cql:"name"`
	InternalPort int32      `json:"internal_port,omitempty" cql:"internal_port"`
	ExposedPort  int32      `json:"exposed_port,omitempty" cql:"exposed_port"`
	Endpoints    []Endpoint `json:"endpoints,omitempty" cql:"endpoint"`
}

func NewPortFromGRPC(port *grpc_application_go.Port) * Port {
	if port == nil {
		return nil
	}
	endpoints := make([]Endpoint, 0)
	for _, e := range port.Endpoints{
		endpoints = append(endpoints, *NewEndpointFromGRPC(e))
	}
	return &Port{
		Name:         port.Name,
		InternalPort: port.InternalPort,
		ExposedPort:  port.ExposedPort,
		Endpoints:    endpoints,
	}
}

func (p *Port) ToGRPC() *grpc_application_go.Port {
	endpoints := make([]*grpc_application_go.Endpoint, 0)

	for _, ep := range p.Endpoints {
		endpoints = append(endpoints, ep.ToGRPC())
	}

	return &grpc_application_go.Port{
		Name:         p.Name,
		InternalPort: p.InternalPort,
		ExposedPort:  p.ExposedPort,
		Endpoints:    endpoints,
	}
}

// -- ConfigFile -- //
type ConfigFile struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty" cql:"app_descriptor_id"`
	// ConfigFileId with the config file identifier.
	ConfigFileId string `json:"config_file_id,omitempty" cql:"config_file_id"`
	Name string `json:"name" cql:"name"`
	// Content of the configuration file.
	Content []byte `json:"content,omitempty" cql:"content"`
	// MountPath of the configuration file in the service instance.
	MountPath string `json:"mount_path,omitempty" cql:"mount_path"`
}

func NewConfigFileFromGRPC(organizationID string, appDescriptorID string, config * grpc_application_go.ConfigFile) * ConfigFile {
	if config == nil {
		return nil
	}
	return &ConfigFile{
		OrganizationId:  organizationID,
		AppDescriptorId: appDescriptorID,
		ConfigFileId:    GenerateUUID(),
		Name: 			 config.Name,
		Content:         config.Content,
		MountPath:       config.MountPath,
	}
}

func (cf *ConfigFile) ToGRPC() *grpc_application_go.ConfigFile {
	return &grpc_application_go.ConfigFile{
		OrganizationId:  cf.OrganizationId,
		AppDescriptorId: cf.AppDescriptorId,
		ConfigFileId:    cf.ConfigFileId,
		Name:    		 cf.Name,
		Content:         cf.Content,
		MountPath:       cf.MountPath,
	}
}


// -- Service -- //
type Service struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty" cql:"app_descriptor_id"`
	// Service group id this service belongs to.
	ServiceGroupId string `json:"service_group_id,omitempty" cql:"service_group_id"`
	// ServiceId with the service identifier.
	ServiceId string `json:"service_id,omitempty" cql:"service_id"`
	// Name of the service.
	Name string `json:"name,omitempty" cql:"name"`
	// ServiceType represents the underlying technology of the service to be launched.
	Type ServiceType `json:"type,omitempty" cql:"type"`
	// Image contains the URL/name of the image to be executed.
	Image string `json:"image,omitempty" cql:"image"`
	// ImageCredentials with the data required to access the repository the image is available at.
	Credentials * ImageCredentials `json:"credentials,omitempty" cql:"credentials"`
	// DeploySpecs with the resource specs required by the service.
	Specs * DeploySpecs `json:"specs,omitempty" cql:"specs"`
	// Storage restrictions
	Storage []Storage `json:"storage,omitempty" cql:"storage"`
	// ExposedPorts contains the list of ports exposed by the current service.
	ExposedPorts []Port `json:"exposed_ports,omitempty" cql:"exposed_ports"`
	// EnvironmentVariables defines a key-value map of environment variables and values that will be passed to all
	// running services.
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty" cql:"environment_variables"`
	// Configs contains the configuration files required by the service.
	Configs []ConfigFile `json:"configs,omitempty" cql:"configs"`
	// Labels with the user defined labels.
	Labels map[string]string `json:"labels,omitempty" cql:"labels"`
	// DeployAfter contains the list of services that must be running before launching a service.
	DeployAfter []string `json:"deploy_after,omitempty" cql:"deploy_after"`
	// RunArguments contains the list of arguments
	RunArguments [] string `json:"run_arguments" cql:"run_arguments"`
}

func NewServiceFromGRPC(organizationID string, appDescriptorID string, serviceGroupId string, service *grpc_application_go.Service) * Service {
	if service == nil{
		return nil
	}

	id := GenerateUUID()

	storage := make([]Storage, 0)
	for _, s := range service.Storage {
		storage = append(storage, *NewStorageFromGRPC(s))
	}
	ports := make([]Port, 0)
	for _, p := range service.ExposedPorts {
		ports = append(ports, *NewPortFromGRPC(p))
	}
	configs := make([]ConfigFile, 0)
	for _, cf := range service.Configs {
		configs = append(configs, *NewConfigFileFromGRPC(organizationID, appDescriptorID, cf))
	}

	serviceType, _ := ServiceTypeFromGRPC[service.Type]
	return &Service{
		OrganizationId:       organizationID,
		AppDescriptorId:      appDescriptorID,
		ServiceGroupId:       serviceGroupId,
		ServiceId:            id,
		Name:                 service.Name,
		Type:                 serviceType,
		Image:                service.Image,
		Credentials:          NewImageCredentialsFromGRPC(service.Credentials),
		Specs:                NewDeploySpecsFromGRPC(service.Specs),
		Storage:              storage,
		ExposedPorts:         ports,
		EnvironmentVariables: service.EnvironmentVariables,
		Configs:              configs,
		Labels:               service.Labels,
		DeployAfter:          service.DeployAfter,
		RunArguments: 		  service.RunArguments,
	}
}

func (s *Service) ToGRPC() *grpc_application_go.Service {
	serviceType, _ := ServiceTypeToGRPC[s.Type]
	storage := make([]*grpc_application_go.Storage, 0)
	for _, s := range s.Storage {
		storage = append(storage, s.ToGRPC())
	}
	exposedPorts := make([]*grpc_application_go.Port, 0)
	for _, ep := range s.ExposedPorts {
		exposedPorts = append(exposedPorts, ep.ToGRPC())
	}
	configs := make([]*grpc_application_go.ConfigFile, 0)
	for _, c := range s.Configs {
		configs = append(configs, c.ToGRPC())
	}
	return &grpc_application_go.Service{
		OrganizationId:       s.OrganizationId,
		AppDescriptorId:      s.AppDescriptorId,
		ServiceGroupId:       s.ServiceGroupId,
		ServiceId:            s.ServiceId,
		Name:                 s.Name,
		Type:                 serviceType,
		Image:                s.Image,
		Credentials:          s.Credentials.ToGRPC(),
		Specs:                s.Specs.ToGRPC(),
		Storage:              storage,
		ExposedPorts:         exposedPorts,
		EnvironmentVariables: s.EnvironmentVariables,
		Configs:              configs,
		Labels:               s.Labels,
		DeployAfter:          s.DeployAfter,
		RunArguments:         s.RunArguments,
	}
}

// -- Service -> ServiceInstance -- //
func (s * Service) ToServiceInstance(appInstanceID string, serviceGroupInstanceID string) * ServiceInstance {

	return &ServiceInstance{
		OrganizationId:       s.OrganizationId,
		AppDescriptorId:      s.AppDescriptorId,
		AppInstanceId:        appInstanceID,
		ServiceGroupId:       s.ServiceGroupId,
		ServiceGroupInstanceId: serviceGroupInstanceID,
		ServiceId:            s.ServiceId,
		ServiceInstanceId:    uuid.New().String(),
		Name:                 s.Name,
		Type:                 s.Type,
		Image:                s.Image,
		Credentials:          s.Credentials,
		Specs:                s.Specs,
		Storage:              s.Storage,
		ExposedPorts:         s.ExposedPorts,
		EnvironmentVariables: s.EnvironmentVariables,
		Configs:              s.Configs,
		Labels:               s.Labels,
		DeployAfter:          s.DeployAfter,
		Status:               ServiceWaiting,
		RunArguments:         s.RunArguments,
	}
}

// EndpointInstance represents a working endpoint exposing data to the outside world. The main difference between
// and endpoint and its instance is the availability of FQDN exposing all the information.
type EndpointInstance struct {
	// EndpointInstanceId unique id for this endpoint
	EndpointInstanceId string `json:"endpoint_instance_id,omitempty" cql:"endpoint_instance_id"`
	// Type of endpoint
	Type EndpointType `json:"type,omitempty" cql:"type"`
	// FQDN to be accessed by any client
	Fqdn string   `json:"fqdn,omitempty" cql:"fqdn"`
	// Port port in the endpoint
	Port                 int32    `json:"port,omitempty" cql:"port"`
}

func (ep * EndpointInstance) ToGRPC () *grpc_application_go.EndpointInstance {
	convertedType, _ := EndpointTypeToGRPC[ep.Type]
	return & grpc_application_go.EndpointInstance{
		EndpointInstanceId: ep.EndpointInstanceId,
		Type : 				convertedType,
		Fqdn: 				ep.Fqdn,
		Port:               ep.Port,
	}
}

func EndpointInstanceFromGRPC(endpoint *grpc_application_go.EndpointInstance) EndpointInstance{
	return EndpointInstance{
		EndpointInstanceId: endpoint.EndpointInstanceId,
		Fqdn: endpoint.Fqdn,
		Type: EndpointTypeFromGRPC[endpoint.Type],
		Port: endpoint.Port,
	}
}

type ServiceStatus int

const (
	ServiceScheduled ServiceStatus = iota + 1
	ServiceWaiting
	ServiceDeploying
	ServiceRunning
	ServiceError
)

var ServiceStatusToGRPC = map[ServiceStatus]grpc_application_go.ServiceStatus{
	ServiceScheduled:    grpc_application_go.ServiceStatus_SERVICE_SCHEDULED,
	ServiceWaiting: grpc_application_go.ServiceStatus_SERVICE_WAITING,
	ServiceDeploying:       grpc_application_go.ServiceStatus_SERVICE_DEPLOYING,
	ServiceRunning:        grpc_application_go.ServiceStatus_SERVICE_RUNNING,
	ServiceError: grpc_application_go.ServiceStatus_SERVICE_ERROR,
}

var ServiceStatusFromGRPC = map[grpc_application_go.ServiceStatus]ServiceStatus{
	grpc_application_go.ServiceStatus_SERVICE_SCHEDULED : ServiceScheduled,
	grpc_application_go.ServiceStatus_SERVICE_WAITING : ServiceWaiting,
	grpc_application_go.ServiceStatus_SERVICE_DEPLOYING : ServiceDeploying,
	grpc_application_go.ServiceStatus_SERVICE_RUNNING : ServiceRunning,
	grpc_application_go.ServiceStatus_SERVICE_ERROR : ServiceError,
}
// -- ServiceInstance -- //
type ServiceInstance struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty" cql:"app_descriptor_id"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty" cql:"app_instance_id"`
	// ServiceGroupId with the group identifier.
	ServiceGroupId string `json:"service_group_id,omitempty" cql:"service_group_id"`
	// ServiceGroupInstanceId with the service group instance identifier.
	ServiceGroupInstanceId string `json:"service_group_instance_id,omitempty" cql:"service_group_instance_id"`
	// ServiceId with the service identifier.
	ServiceId string `json:"service_id,omitempty" cql:"service_id"`
	// Unique identifier for this instance
	ServiceInstanceId string `json:"service_instance_id,omitempty" cql:"service_instance_id"`
	// Name of the service.
	Name string `json:"name,omitempty" cql:"name"`
	// ServiceType represents the underlying technology of the service to be launched.
	Type ServiceType `json:"type,omitempty" cql:"type"`
	// Image contains the URL/name of the image to be executed.
	Image string `json:"image,omitempty" cql:"image"`
	// ImageCredentials with the data required to access the repository the image is available at.
	Credentials * ImageCredentials `json:"credentials,omitempty" cql:"credentials"`
	// DeploySpecs with the resource specs required by the service.
	Specs * DeploySpecs `json:"specs,omitempty" cql:"specs"`
	// Storage restrictions
	Storage []Storage `json:"storage,omitempty" cql:"storage"`
	// ExposedPorts contains the list of ports exposed by the current service.
	ExposedPorts []Port `json:"exposed_ports,omitempty" cql:"exposed_ports"`
	// EnvironmentVariables defines a key-value map of environment variables and values that will be passed to all
	// running services.
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty" cql:"environment_variables"`
	// Configs contains the configuration files required by the service.
	Configs []ConfigFile `json:"configs,omitempty" cql:"configs"`
	// Labels with the user defined labels.
	Labels map[string]string `json:"labels,omitempty" cql:"labels"`
	// DeployAfter contains the list of services that must be running before launching a service.
	DeployAfter []string `json:"deploy_after,omitempty" cql:"deploy_after"`
	// Status of the deployed service
	Status ServiceStatus `json:"status,omitempty" cql:"status"`
	// Endpoints exposed to the users by the service.
	Endpoints []EndpointInstance `json:"endpoints,omitempty" cql:"endpoints"`
	// DeployedOnClusterId specifies which is the cluster where the service is running.
	DeployedOnClusterId  string  `json:"deployed_on_cluster_id,omitempty" cql:"deployed_on_cluster_id"`
	// RunArguments containts a list of arguments
	RunArguments [] string `json:"run_arguments" cql:"run_arguments"`
	// Relevant information about this instance
	Info string   `json:"info,omitempty" cql:"info"`
}

func NewServiceInstanceFromGRPC(serviceInstance *grpc_application_go.ServiceInstance) ServiceInstance {

	endpoints := make([]EndpointInstance,len(serviceInstance.Endpoints))
	for i, e := range serviceInstance.Endpoints {
		endpoints[i] = EndpointInstance{Type: EndpointTypeFromGRPC[e.Type], EndpointInstanceId: e.EndpointInstanceId,
			Fqdn: e.Fqdn, Port: e.Port}
	}

	storage := make([]Storage, len(serviceInstance.Storage))
	for i, s := range serviceInstance.Storage {
		storage[i] = *NewStorageFromGRPC(s)
	}

	exposedPorts := make([]Port, len(serviceInstance.ExposedPorts))
	for i, p := range serviceInstance.ExposedPorts {
		exposedPorts[i] = *NewPortFromGRPC(p)
	}

	configs := make([]ConfigFile, len(serviceInstance.Configs))
	for i, c := range serviceInstance.Configs {
		configs[i] = ConfigFile{AppDescriptorId: c.AppDescriptorId, OrganizationId: c.OrganizationId,
			Name: c.Name, ConfigFileId: c.ConfigFileId, Content: c.Content, MountPath: c.MountPath}
	}

	return ServiceInstance{
		Name: serviceInstance.Name,
		Labels: serviceInstance.Labels,
		Status: ServiceStatusFromGRPC[serviceInstance.Status],
		ServiceGroupInstanceId: serviceInstance.ServiceGroupInstanceId,
		AppInstanceId: serviceInstance.AppInstanceId,
		ServiceGroupId: serviceInstance.ServiceGroupId,
		OrganizationId: serviceInstance.OrganizationId,
		AppDescriptorId: serviceInstance.AppDescriptorId,
		Info: serviceInstance.Info,
		EnvironmentVariables: serviceInstance.EnvironmentVariables,
		ServiceInstanceId: serviceInstance.ServiceInstanceId,
		Type: ServiceTypeFromGRPC[serviceInstance.Type],
		ServiceId: serviceInstance.ServiceId,
		Image: serviceInstance.Image,
		DeployedOnClusterId: serviceInstance.DeployedOnClusterId,
		Specs: NewDeploySpecsFromGRPC(serviceInstance.Specs),
		Endpoints: endpoints,
		Credentials: NewImageCredentialsFromGRPC(serviceInstance.Credentials),
		Storage: storage,
		ExposedPorts: exposedPorts,
		Configs: configs,
		DeployAfter: serviceInstance.DeployAfter,
		RunArguments: serviceInstance.RunArguments,
	}
}

func (si *ServiceInstance) ToGRPC() *grpc_application_go.ServiceInstance {
	serviceType, _ := ServiceTypeToGRPC[si.Type]
	serviceStatus, _ := ServiceStatusToGRPC[si.Status]
	storage := make([]*grpc_application_go.Storage, 0)
	for _, s := range si.Storage {
		storage = append(storage, s.ToGRPC())
	}
	exposedPorts := make([]*grpc_application_go.Port, 0)
	for _, ep := range si.ExposedPorts {
		exposedPorts = append(exposedPorts, ep.ToGRPC())
	}
	configs := make([]*grpc_application_go.ConfigFile, 0)
	for _, c := range si.Configs {
		configs = append(configs, c.ToGRPC())
	}
	endpoints := make ([]*grpc_application_go.EndpointInstance, 0)
	for _,ep := range si.Endpoints {
		endpoints = append(endpoints, ep.ToGRPC())
	}
	return &grpc_application_go.ServiceInstance{
		OrganizationId:       	si.OrganizationId,
		AppDescriptorId:      	si.AppDescriptorId,
		AppInstanceId:        	si.AppInstanceId,
		ServiceGroupId: 	  	si.ServiceGroupId,
		ServiceGroupInstanceId:	si.ServiceGroupInstanceId,
		ServiceId:           	si.ServiceId,
		ServiceInstanceId: 		si.ServiceInstanceId,
		Name:                 	si.Name,
		Type:                 	serviceType,
		Image:                	si.Image,
		Credentials:          	si.Credentials.ToGRPC(),
		Specs:                	si.Specs.ToGRPC(),
		Storage:              	storage,
		ExposedPorts:         	exposedPorts,
		EnvironmentVariables: 	si.EnvironmentVariables,
		Configs:              	configs,
		Labels:               	si.Labels,
		DeployAfter:          	si.DeployAfter,
		Status:               	serviceStatus,
		Endpoints:            	endpoints,
		DeployedOnClusterId:  	si.DeployedOnClusterId,
		RunArguments: 		  	si.RunArguments,
		Info:            		si.Info,
	}

}

// -- AppDecriptor -- //
type AppDescriptor struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty" cql:"app_descriptor_id"`
	// Name of the application.
	Name string `json:"name,omitempty" cql:"name"`
	// ConfigurationOptions defines a key-value map of configuration options.
	ConfigurationOptions map[string]string `json:"configuration_options,omitempty" cql:"configuration_options"`
	// EnvironmentVariables defines a key-value map of environment variables and values that will be passed to all
	// running services.
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty" cql:"environment_variables"`
	// Labels defined by the user.
	Labels map[string]string `json:"labels,omitempty" cql:"labels"`
	// Rules that define the connectivity between the elements of an application.
	Rules []SecurityRule `json:"rules,omitempty" cql:"rules"`
	// Groups with the Service collocation strategies.
	Groups []ServiceGroup `json:"groups,omitempty" cql:"groups"`
}

func NewAppDescriptor(organizationID string, appDescriptorID string, name string,
	configOptions map[string]string, envVars map[string]string,
	labels map[string]string,
	rules []SecurityRule, groups []ServiceGroup) *AppDescriptor {
	return &AppDescriptor{
		OrganizationId: 		organizationID,
	 	AppDescriptorId: 		appDescriptorID,
		Name: 					name,
		ConfigurationOptions:	configOptions,
		EnvironmentVariables:	envVars,
		Labels:					labels,
		Rules:					rules,
		Groups:					groups,
		}
}

func NewAppDescriptorFromGRPC(addRequest * grpc_application_go.AddAppDescriptorRequest, deviceGroupIds map[string]string) (*AppDescriptor, derrors.Error) {

	if addRequest == nil {
		return nil, nil
	}

	uuid := GenerateUUID()

	rules := make([]SecurityRule, 0)
	if addRequest.Rules != nil {
		for _, r := range addRequest.Rules {
			rule, err := NewSecurityRuleFromGRPC(addRequest.OrganizationId, uuid, r, deviceGroupIds)
			if err != nil {
				return nil, err
			}
			rules = append(rules, *rule)
		}
	}
	groups := make([]ServiceGroup, 0)
	for _, sg := range addRequest.Groups{
		groups = append(groups, *NewServiceGroupFromGRPC(addRequest.OrganizationId, uuid, sg))
	}
	return NewAppDescriptor(
		addRequest.OrganizationId,
		uuid,
		addRequest.Name,
		addRequest.ConfigurationOptions,
		addRequest.EnvironmentVariables,
		addRequest.Labels,
		rules, groups), nil
}

func (d *AppDescriptor) ToGRPC() *grpc_application_go.AppDescriptor {

	rules := make([]*grpc_application_go.SecurityRule, 0)
	for _, r := range d.Rules {
		rules = append(rules, r.ToGRPC())
	}
	groups := make([]*grpc_application_go.ServiceGroup, 0)
	for _, g := range d.Groups {
		groups = append(groups, g.ToGRPC())
	}
	return &grpc_application_go.AppDescriptor{
		OrganizationId:       d.OrganizationId,
		AppDescriptorId:      d.AppDescriptorId,
		Name:                 d.Name,
		ConfigurationOptions: d.ConfigurationOptions,
		EnvironmentVariables: d.EnvironmentVariables,
		Labels:               d.Labels,
		Rules:                rules,
		Groups:               groups,
	}
}

// -------------

type AppEndpointProtocol int

const (
	HTTP AppEndpointProtocol = iota + 1
	HTTPS
)

var AppEndpointProtocolToGRPC = map[AppEndpointProtocol]grpc_application_go.AppEndpointProtocol{
	HTTP:    grpc_application_go.AppEndpointProtocol_HTTP,
	HTTPS:   grpc_application_go.AppEndpointProtocol_HTTPS,
}

var AppEndpointProtocolFromGRPC = map[grpc_application_go.AppEndpointProtocol]AppEndpointProtocol{
	grpc_application_go.AppEndpointProtocol_HTTP:   HTTP,
	grpc_application_go.AppEndpointProtocol_HTTPS:  HTTPS,
}

type AppEndpoint struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty" cql:"app_instance_id"`
	// ServiceGroupInstanceId the identifier of the group instance.
	ServiceGroupInstanceId string `json:"service_group_instance_id,omitempty" cql:"service_group_instance_id"`
	// ServiceInstanceId the identifier of the service instance.
	ServiceInstanceId string `json:"service_instance_id,omitempty" cql:"service_instance_id"`
	// Port port in the endpoint
	Port int32 `json:"port,omitempty" cql:"port"`
	// protocol (http, https)
	Protocol AppEndpointProtocol `json:"protocol,omitempty" cql:"protocol"`
	// EndpointInstanceId unique id for this endpoint
	EndpointInstanceId string `json:"endpoint_instance_id,omitempty" cql:"endpoint_instance_id"`
	// Type of endpoint
	Type EndpointType `json:"type,omitempty" cql:"type"`
	// FQDN to be accessed by any client
	Fqdn string   `json:"fqdn,omitempty" cql:"fqdn"`
	// GlobalFqdn
	GlobalFqdn string `json:"global_fqdn,omitempty" cql:"global_fqdn"`
}

func (ep * AppEndpoint) ToGRPC () *grpc_application_go.AppEndpoint {
	convertedType, _ := EndpointTypeToGRPC[ep.Type]
	convertedProtocol, _ := AppEndpointProtocolToGRPC[ep.Protocol]
	return & grpc_application_go.AppEndpoint{
		OrganizationId: ep.OrganizationId,
		AppInstanceId: ep.AppInstanceId,
		ServiceGroupInstanceId: ep.ServiceGroupInstanceId,
		ServiceInstanceId: ep.ServiceInstanceId,
		Protocol: convertedProtocol,
		EndpointInstance: &grpc_application_go.EndpointInstance{
			EndpointInstanceId: ep.EndpointInstanceId,
			Type:convertedType,
			Fqdn: ep.Fqdn,
			Port: ep.Port,
		},
	}
}


// getNamePrefixes returns prefix to fill the globalFQDN
// 1) "service-name"-"port"
// 2) service_group_instanceID (6 characters)
// 3) appInstance (6 characters)
// 4) organizationID (8 characters)
func getNamePrefixes(ep *grpc_application_go.AddAppEndpointRequest) (string, string, string, string){
	serviceName := ep.ServiceName

	if ep.EndpointInstance != nil && ep.EndpointInstance.Port != 80 && ep.EndpointInstance.Port != 0 {
		serviceName = fmt.Sprintf("%s-%d", ep.ServiceName, ep.EndpointInstance.Port)
	}
	serviceGroupInstPrefix := ep.ServiceGroupInstanceId
	if len(serviceGroupInstPrefix) > InstPrefixLength {
		serviceGroupInstPrefix = serviceGroupInstPrefix[0:InstPrefixLength]
	}
	appInstPrefix := ep.AppInstanceId
	if len(appInstPrefix) > InstPrefixLength {
		appInstPrefix = appInstPrefix[0:InstPrefixLength]
	}
	orgPrefix := ep.OrganizationId
	if len(orgPrefix) > OrgPrefixLength {
		orgPrefix = orgPrefix[0:OrgPrefixLength]
	}
	return serviceName, serviceGroupInstPrefix, appInstPrefix, orgPrefix
}

// createGlobalFqdn returns the globalFqdn for a endpoinFqnd given
func createGlobalFqdn(endpoint *grpc_application_go.AddAppEndpointRequest) string {

	// Option1 - Fqdn: serv.A.B.domain
	// where:
	// A: service_group_id
	// B: app_instance_id

	// Option2 - Fqdn: IP:port

	// We need to store:
	// Global Fqdn: serv.A.B.C.domain
	// where
	// A: service_group_id (6 characters)
	// B: app_instance_id (6 characters)
	// C: organization_id (8 characters)
	// the domain is not stored

	serviceName, serviceGroupId, appInstanceId, organizationId := getNamePrefixes(endpoint)

	return fmt.Sprintf("%s.%s.%s.%s", serviceName, serviceGroupId, appInstanceId, organizationId)

}

func NewAppEndpointFromGRPC(endpoint *grpc_application_go.AddAppEndpointRequest) (* AppEndpoint, derrors.Error){

	if endpoint.EndpointInstance == nil {
		endpoint.EndpointInstance = DefaultEndpointInstance
	}

	return &AppEndpoint{
		OrganizationId: endpoint.OrganizationId,
		AppInstanceId: endpoint.AppInstanceId,
		ServiceGroupInstanceId:endpoint.ServiceGroupInstanceId,
		ServiceInstanceId: endpoint.ServiceInstanceId,
		Protocol: AppEndpointProtocolFromGRPC[endpoint.Protocol],
		Port: endpoint.EndpointInstance.Port,
		EndpointInstanceId:endpoint.EndpointInstance.EndpointInstanceId,
		Type:  EndpointTypeFromGRPC[endpoint.EndpointInstance.Type],
		Fqdn: endpoint.EndpointInstance.Fqdn,
		GlobalFqdn:createGlobalFqdn(endpoint),
	}, nil
}

// -------------

func (d * AppDescriptor) ApplyUpdate(request grpc_application_go.UpdateAppDescriptorRequest){
	if request.AddLabels {
		for k, v := range request.Labels {
			d.Labels[k] = v
		}
	}
	if request.RemoveLabels {
		for k, _ := range request.Labels {
			delete(d.Labels, k)
		}
	}
}

func ValidGroup(group * grpc_application_go.ServiceGroup) derrors.Error {
	if group.Name == "" {
		return derrors.NewInvalidArgumentError("expecting name")
	}
	if len(group.Services) == 0 {
		return derrors.NewInvalidArgumentError("expecting at least one service")
	}
	return nil
}


func ValidAddService(service * grpc_application_go.Service) derrors.Error {
	if service.OrganizationId == "" || service.ServiceId == "" {
		return derrors.NewInvalidArgumentError("expecting organization_id, service_id")
	}
	return nil
}

func ValidAddAppDescriptorRequest(toAdd * grpc_application_go.AddAppDescriptorRequest) derrors.Error {
	if toAdd.OrganizationId == "" || toAdd.Name == "" || toAdd.RequestId == "" {
		return derrors.NewInvalidArgumentError("expecting organization_id, name, and request_id")
	}

	if len(toAdd.Groups) == 0 {
		return derrors.NewInvalidArgumentError("expecting at least one service group")
	}

	for _, g := range toAdd.Groups {
		err := ValidGroup (g)
		if err != nil{
			return err
		}
	}

	return nil
}

func ValidUpdateAppDescriptorRequest(request * grpc_application_go.UpdateAppDescriptorRequest) derrors.Error{
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if request.AppDescriptorId == ""{
		return derrors.NewInvalidArgumentError(emptyAppDescriptorId)
	}
	return nil
}


type ApplicationStatus int

const (
	Queued ApplicationStatus = iota +1
	Planning
	Scheduled
	Deploying
	Running
	Incomplete
	PlanningError
	DeploymentError
	Error
)

var AppStatusToGRPC = map[ApplicationStatus]grpc_application_go.ApplicationStatus{
	Queued: grpc_application_go.ApplicationStatus_QUEUED,
	Planning: grpc_application_go.ApplicationStatus_PLANNING,
	Scheduled: grpc_application_go.ApplicationStatus_SCHEDULED,
	Deploying: grpc_application_go.ApplicationStatus_DEPLOYING,
	Running: grpc_application_go.ApplicationStatus_RUNNING,
	Incomplete: grpc_application_go.ApplicationStatus_INCOMPLETE,
	PlanningError: grpc_application_go.ApplicationStatus_PLANNING_ERROR,
	DeploymentError: grpc_application_go.ApplicationStatus_DEPLOYMENT_ERROR,
	Error: grpc_application_go.ApplicationStatus_ERROR,
}

var AppStatusFromGRPC = map[grpc_application_go.ApplicationStatus]ApplicationStatus{
	grpc_application_go.ApplicationStatus_QUEUED:Queued,
	grpc_application_go.ApplicationStatus_PLANNING:Planning,
	grpc_application_go.ApplicationStatus_SCHEDULED:Scheduled,
	grpc_application_go.ApplicationStatus_DEPLOYING:Deploying,
	grpc_application_go.ApplicationStatus_RUNNING:Running,
	grpc_application_go.ApplicationStatus_INCOMPLETE:Incomplete,
	grpc_application_go.ApplicationStatus_PLANNING_ERROR:PlanningError,
	grpc_application_go.ApplicationStatus_DEPLOYMENT_ERROR:DeploymentError,
	grpc_application_go.ApplicationStatus_ERROR:Error,
}

// -- AppInstance -- //
type AppInstance struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty" cql: "app_descriptor_id"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty" cql:"app_instance_id"`
	// Name of the application.
	Name string `json:"name,omitempty" cql:"name"`
	// ConfigurationOptions defines a key-value map of configuration options.
	ConfigurationOptions map[string]string `json:"configuration_options,omitempty" cql:"configuration_options"`
	// EnvironmentVariables defines a key-value map of environment variables and values that will be passed to all
	// running services.
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty" cql:"environment_variables"`
	// Labels defined by the user.
	Labels map[string]string `json:"labels,omitempty" cql:"labels"`
	// Rules that define the connectivity between the elements of an application.
	Rules []SecurityRule `json:"rules,omitempty" cql:"rules"`
	// Groups with the Service collocation strategies.
	Groups []ServiceGroupInstance `json:"groups,omitempty" cql:"groups"`
	// Status of the deployed instance.
	Status  ApplicationStatus `json:"status,omitempty" cql:"status"`
	// Metadata descriptor for the instances triggered by this app
	Metadata []InstanceMetadata `json:"metadata,omitempty" cql:"metadata"`
	// Textual information for this application instance
	Info string `json:"info,omitempty" cql:"info"`

}
func (sg * ServiceGroup) ToServiceGroupInstance(appInstanceID string) *ServiceGroupInstance {
	serviceGroupInstanceID := uuid.New().String()
	services := make ([]ServiceInstance, 0)
	for _, s := range sg.Services {
		services = append(services, *s.ToServiceInstance(appInstanceID, serviceGroupInstanceID))
	}
	return &ServiceGroupInstance{
		OrganizationId: 		sg.OrganizationId,
		AppDescriptorId: 		sg.AppDescriptorId,
		AppInstanceId: 			appInstanceID,
		ServiceGroupId: 		sg.ServiceGroupId,
		ServiceGroupInstanceId: serviceGroupInstanceID,
		Name: 					sg.Name,
		ServiceInstances:		services,
		Policy: 				sg.Policy,
		Status: 				ServiceScheduled,
		Specs: 					sg.Specs,
		Labels: 				sg.Labels,
	}
}

func (sg * ServiceGroup) ToEmptyServiceGroupInstance(appInstanceID string) *ServiceGroupInstance {

	serviceGroupInstanceID := uuid.New().String()

	instances := make([]ServiceInstance,len(sg.Services))
	// New service instances for this service group
	for i, g := range sg.Services {
		instance := g.ToServiceInstance(appInstanceID, serviceGroupInstanceID)
		instances[i] = *instance
	}

	return &ServiceGroupInstance{
		OrganizationId: 		sg.OrganizationId,
		AppDescriptorId: 		sg.AppDescriptorId,
		AppInstanceId: 			appInstanceID,
		ServiceGroupId: 		sg.ServiceGroupId,
		ServiceGroupInstanceId: serviceGroupInstanceID,
		Name: 					sg.Name,
		ServiceInstances:		instances,
		Policy: 				sg.Policy,
		Status: 				ServiceScheduled,
		Specs: 					sg.Specs,
		Labels: 				sg.Labels,
	}
}

// This function returns a local object from an incoming grpc service group instance
func NewServiceGroupInstanceFromGRPC(group *grpc_application_go.ServiceGroupInstance) *ServiceGroupInstance {

	serviceInstances := make([]ServiceInstance,0)
	for _, serv := range group.ServiceInstances {
		serviceInstances = append(serviceInstances, NewServiceInstanceFromGRPC(serv))
	}

	return &ServiceGroupInstance{
		AppDescriptorId: group.AppDescriptorId,
		OrganizationId: group.OrganizationId,
		AppInstanceId: group.AppInstanceId,
		Name: group.Name,
		Metadata: NewMetadataFromGRPC(group.Metadata),
		Labels: group.Labels,
		Status: ServiceStatusFromGRPC[group.Status],
		ServiceGroupId: group.ServiceGroupId,
		ServiceGroupInstanceId: group.ServiceGroupInstanceId,
		Policy: CollocationPolicyFromGRPC[group.Policy],
		ServiceInstances: serviceInstances,
		GlobalFqdn: group.GlobalFqdn,
		Specs: NewServiceGroupDeploymentSpecsFromGRPC(group.Specs),
	}
}


func NewAppInstanceFromAddInstanceRequestGRPC(addRequest * grpc_application_go.AddAppInstanceRequest, descriptor * AppDescriptor) * AppInstance {
	uuid := GenerateUUID()

	return &AppInstance{
		OrganizationId:       addRequest.OrganizationId,
		AppDescriptorId:      addRequest.AppDescriptorId,
		AppInstanceId:        uuid,
		Name:                 addRequest.Name,
		ConfigurationOptions: descriptor.ConfigurationOptions,
		EnvironmentVariables: descriptor.EnvironmentVariables,
		Labels:               descriptor.Labels,
		Rules:                descriptor.Rules,
		// ServiceGroupInstances are added using the addservicegroupinstances function
		//Groups:               groups,
		Status: Queued,
		Info:                 "",
	}
}

func NewAppInstanceFromGRPC(appInstance * grpc_application_go.AppInstance) * AppInstance {
	groups := make([]ServiceGroupInstance, 0)
	for _, g := range appInstance.Groups {
		groups = append(groups, *NewServiceGroupInstanceFromGRPC(g))
	}
	metadata := make([]InstanceMetadata,0)
	for _, m := range appInstance.Metadata {
		metadata = append(metadata, *NewMetadataFromGRPC(m))
	}
	rules := make([]SecurityRule,0)
	for _, r := range appInstance.Rules {
		newR, err := NewSecurityRuleFromInstantiatedGRPC(r)
		if err != nil {
			return nil
		}
		rules = append(rules, *newR)
	}

	return &AppInstance{
		OrganizationId: appInstance.OrganizationId,
		AppInstanceId: appInstance.AppInstanceId,
		Groups: groups,
		Info: appInstance.Info,
		Status: AppStatusFromGRPC[appInstance.Status],
		Metadata: metadata,
		Name: appInstance.Name,
		Labels: appInstance.Labels,
		AppDescriptorId: appInstance.AppInstanceId,
		EnvironmentVariables: appInstance.EnvironmentVariables,
		Rules: rules,
		ConfigurationOptions: appInstance.ConfigurationOptions,
	}
}

func (i *AppInstance) ToGRPC() *grpc_application_go.AppInstance {
	rules := make([]*grpc_application_go.SecurityRule, 0)
	for _, r := range i.Rules {
		rules = append(rules, r.ToGRPC())
	}
	groups := make([]*grpc_application_go.ServiceGroupInstance, 0)
	for _, g := range i.Groups {
		groups = append(groups, g.ToGRPC())
	}
	metadata := make ([]*grpc_application_go.InstanceMetadata, 0)
	for _, md := range i.Metadata {
		metadata = append(metadata, md.ToGRPC())
	}

	status, _ := AppStatusToGRPC[i.Status]

	return &grpc_application_go.AppInstance{
		OrganizationId:       i.OrganizationId,
		AppDescriptorId:      i.AppDescriptorId,
		AppInstanceId:        i.AppInstanceId,
		Name:                 i.Name,
		ConfigurationOptions: i.ConfigurationOptions,
		EnvironmentVariables: i.EnvironmentVariables,
		Labels:               i.Labels,
		Rules:                rules,
		Groups:               groups,
		Status:               status,
		Metadata:             metadata,
		Info:                 i.Info,
	}
}

// AppZtNetwork
type AppZtNetwork struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty" cql: "app_instance_id"`
	// ZtNetworkId zero-tier network identifier.
	ZtNetworkId string `json:"zt_network_id,omitempty" cql:"zt_network_id"`
}

func NewAppZtNetworkFromGRPC(req *grpc_application_go.AppZtNetwork ) *AppZtNetwork {
	return &AppZtNetwork{
		OrganizationId: req.OrganizationId,
		AppInstanceId: req.AppInstanceId,
		ZtNetworkId: req.NetworkId,
	}
}

func(a *AppZtNetwork) ToGRPC() *grpc_application_go.AppZtNetwork {
	return &grpc_application_go.AppZtNetwork{
		NetworkId: a.ZtNetworkId,
		AppInstanceId: a.AppInstanceId,
		OrganizationId: a.OrganizationId,
	}
}

// Validation functions

func ValidAddAppInstanceRequest(toAdd * grpc_application_go.AddAppInstanceRequest) derrors.Error {
	if toAdd.OrganizationId == "" || toAdd.Name == "" || toAdd.AppDescriptorId == "" {
		return derrors.NewInvalidArgumentError("expecting organization_id, name, and descriptor_id")
	}
	return nil
}

func ValidUpdateAppStatusRequest(updateRequest *grpc_application_go.UpdateAppStatusRequest) derrors.Error {
    if updateRequest.OrganizationId == "" || updateRequest.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting organization_id and app_instance_id")
	}
	return nil
}

func ValidUpdateServiceStatusRequest (updateRequest *grpc_application_go.UpdateServiceStatusRequest) derrors.Error {
	if updateRequest.OrganizationId == "" || updateRequest.AppInstanceId == ""  ||
		updateRequest.ServiceGroupInstanceId == "" || updateRequest.ServiceInstanceId == "" {
			return derrors.NewInvalidArgumentError("expecting organization_id, app_instance_id, app_service_instance_id " +
				"and service_instance_id")
	}
	return nil
}

func ValidAddServiceGroupInstanceRequest (request *grpc_application_go.AddServiceGroupInstancesRequest) derrors.Error {
	if request.OrganizationId == "" || request.AppDescriptorId == "" ||
		request.AppInstanceId == "" || request.ServiceGroupId == ""  || request.NumInstances <= 0 {
		return derrors.NewInvalidArgumentError("expecting organization_id, app_descriptor_id, app_instance_id, " +
			"service_group_id, metadata, numInstances greater than zero")
	}

	return nil
}

func ValidAddServiceInstanceRequest(request *grpc_application_go.AddServiceInstanceRequest) derrors.Error {
	if request.OrganizationId == "" || request.AppDescriptorId == "" ||
		request.AppInstanceId == "" || request.ServiceGroupId == "" ||
		request.ServiceGroupInstanceId == "" || request.ServiceId == "" {
		return derrors.NewInvalidArgumentError("expecting organization_id, app_descriptor_id, app_instance_id, " +
			"service_group_id, service_group_instance_id, service_id")
	}
	return nil
}

// ValidateDescriptor checks validity of object names, ports meeting Kubernetes specs.
func  ValidateDescriptor(descriptor AppDescriptor) derrors.Error {
	// for each group
	for _, group := range descriptor.Groups {
		for _, service := range group.Services {
			// Validate service name
			kerr := validation.IsDNS1123Label(service.Name)
			if len(kerr) > 0 {
				return derrors.NewInvalidArgumentError("Service Name").WithParams(service.Name).WithParams(kerr)
			}
			// validate Exposed Port Name and Number
			for _, port := range service.ExposedPorts {
				kerr = validation.IsValidPortName(port.Name)
				if len(kerr) > 0 {
					return derrors.NewInvalidArgumentError("Port Name").WithParams(port.Name).WithParams(kerr)
				}
				kerr = validation.IsValidPortNum(int(port.ExposedPort))
				if len(kerr) > 0 {
					return derrors.NewInvalidArgumentError("Exposed Port").WithParams(port.ExposedPort).WithParams(kerr)
				}
				kerr = validation.IsValidPortNum(int(port.InternalPort))
				if len(kerr) > 0 {
					return derrors.NewInvalidArgumentError("Internal Port").WithParams(port.InternalPort).WithParams(kerr)
				}
			}
		}
	}
	return nil
}

func ValidGetServiceGroupInstanceMetadataRequest(request *grpc_application_go.GetServiceGroupInstanceMetadataRequest) derrors.Error {
	if request.OrganizationId == "" || request.AppInstanceId == "" || request.ServiceGroupInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting organization_id, app_instance_id, " +
			"service_group_instance_id")
	}
	return nil
}

func ValidateRemoveServiceGroupInstancesRequest(request *grpc_application_go.RemoveServiceGroupInstancesRequest) derrors.Error {
	if request.OrganizationId == "" || request.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting organization_id, app_instance_id")
	}
	return nil
}

func ValidUpdateInstanceMetadata(request *grpc_application_go.InstanceMetadata) derrors.Error {
	if request.OrganizationId == "" || request.AppInstanceId == "" || request.ServiceGroupId == "" ||
		request.AppDescriptorId == "" || request.MonitoredInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting organization_id, app_instance_id, " +
			"service_group_instance_id, app_descriptor_id, monitored_instance_id")
	}
	return nil
}

func ValidAddAppEndpointRequest(request *grpc_application_go.AddAppEndpointRequest) derrors.Error {
	if request.AppInstanceId == "" || request.OrganizationId == "" || request.ServiceGroupInstanceId == "" ||
		request.ServiceInstanceId == "" || request.EndpointInstance.Fqdn == "" || request.ServiceName == "" {
			return derrors.NewInvalidArgumentError("expecting organization_id, app_instance_id, " +
				"service_group_instance_id, service_instance_id, service_name, fqdn")
	}

	if request.EndpointInstance == nil || request.EndpointInstance.Fqdn == "" {
		return  derrors.NewInvalidArgumentError("expecting fqdn")
	}
	fqdnSplit := strings.Split(request.EndpointInstance.Fqdn, ".")
	if len(fqdnSplit) < 4 {
		return derrors.NewInvalidArgumentError("fqdn has incorrect format").WithParams(request.EndpointInstance.Fqdn)
	}

	if len(request.OrganizationId) < 8 {
		return derrors.NewInvalidArgumentError("OrganizationId is too short").WithParams(request.OrganizationId)
	}
	return nil
}

func ValidGetAppEndPointRequest(request *grpc_application_go.GetAppEndPointRequest) derrors.Error{
	if request.Fqdn == "" {
		return derrors.NewInvalidArgumentError("expecting fqdn")
	}
	split := strings.Split(request.Fqdn, ".")
	if len(split) < 5 {
		return derrors.NewInvalidArgumentError("fqdn has incorrect format").WithParams(request.Fqdn)
	}
	return nil
}

func ValidRemoveEndpointRequest(request * grpc_application_go.RemoveEndpointRequest)  derrors.Error{
	if request.AppInstanceId == "" || request.OrganizationId == ""  {
		return derrors.NewInvalidArgumentError("expecting organization_id, app_instance_id")
	}
	return nil
}

func ValidAddAppZtNetworkRequest(request * grpc_application_go.AddAppZtNetworkRequest) derrors.Error {
	if request.OrganizationId == "" || request.AppInstanceId == "" || request.NetworkId == "" {
		return derrors.NewInvalidArgumentError("expecting organization_id, app_instance_id, network_id")
	}
	return nil
}

func ValidRemoveAppZtNetworkRequest(request * grpc_application_go.RemoveAppZtNetworkRequest) derrors.Error {
	if request.OrganizationId == "" || request.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting organization_id, app_instance_id")
	}
	return nil
}

func ValidGetAppZtNetworkRequest(request * grpc_application_go.GetAppZtNetworkRequest) derrors.Error {
	if request.OrganizationId == "" || request.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError("expecting organization_id, app_instance_id")
	}
	return nil
}
