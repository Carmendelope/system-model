/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-go"
)

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
	// SourceServiceId defining the service onto which the security rule is defined.
	SourceServiceId string `json:"source_service_id,omitempty" cql:"source_service_id"`
	// SourcePort defining the port that is affected by the current rule.
	SourcePort int32 `json:"source_port,omitempty" cql:"source_port"`
	// Access level to that port defining who can access the port.
	Access PortAccess `json:"access,omitempty" cql:"access"`
	// AuthServices defining a list of services that can access the port.
	AuthServices []string `json:"auth_services,omitempty" cql:"auth_services"`
	// DeviceGroups defining a list of device groups that can access the port.
	DeviceGroups []string `json:"device_groups,omitempty" cql:"device_groups"`
}

func NewSecurityRuleFromGRPC(appDescriptorID string, rule *grpc_application_go.SecurityRule) *SecurityRule {
	if rule == nil {
		return nil
	}
	uuid := GenerateUUID()
	access := PortAccessFromGRPC[rule.Access]
	return &SecurityRule{
		OrganizationId:  rule.OrganizationId,
		AppDescriptorId: appDescriptorID,
		RuleId:          uuid,
		Name:            rule.Name,
		SourceServiceId: rule.SourceServiceId,
		SourcePort:      rule.SourcePort,
		Access:          access,
		AuthServices:    rule.AuthServices,
		DeviceGroups:    rule.DeviceGroups,
	}
}

func (sr *SecurityRule) ToGRPC() *grpc_application_go.SecurityRule {
	access, _ := PortAccessToGRPC[sr.Access]
	return &grpc_application_go.SecurityRule{
		OrganizationId:  sr.OrganizationId,
		AppDescriptorId: sr.AppDescriptorId,
		RuleId:          sr.RuleId,
		Name:            sr.Name,
		SourceServiceId: sr.SourceServiceId, SourcePort: sr.SourcePort,
		Access:       access,
		AuthServices: sr.AuthServices,
		DeviceGroups: sr.DeviceGroups,
	}
}

// ServiceGroupDeploymentSpecs -- //
type ServiceGroupDeploymentSpecs struct {
	// How many times this service group must be replicated
	NumReplicas int32 `json:"num_replicas,omitempty" cql:"num_replicas"`
	// Indicate if this service group must be replicated in every cluster
	MultiClusterReplica  bool     `json:"multi_cluster_replica,omitempty" cql:"multi_cluster_replica"`
}

func NewServiceGroupDeploymentSpecsFromGRPC(specs * grpc_application_go.ServiceGroupDeploymentSpecs) * ServiceGroupDeploymentSpecs{
	if specs == nil {
		return nil
	}
	return &ServiceGroupDeploymentSpecs{
		NumReplicas:         specs.NumReplicas,
		MultiClusterReplica: specs.MultiClusterReplica,
	}
}

func (sp * ServiceGroupDeploymentSpecs) ToGRPC() *grpc_application_go.ServiceGroupDeploymentSpecs  {
	if sp == nil {
		return nil
	}
	return &grpc_application_go.ServiceGroupDeploymentSpecs{
		NumReplicas:          sp.NumReplicas,
		MultiClusterReplica:  sp.MultiClusterReplica,
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
	// Description of the service group.
	Description string `json:"description,omitempty" cql:"description"`
	// Services defining a list of service identifiers that belong to the group.
	Services []string `json:"services,omitempty" cql:"services"`
	// Policy indicating the deployment collocation policy.
	Policy CollocationPolicy `json:"policy,omitempty" cql:"policy"`
	// Particular deployment specs for this service
	Specs  *ServiceGroupDeploymentSpecs `json:"specs,omitempty" cql:"specs"`
}

func NewServiceGroupFromGRPC(appDescriptorID string, group * grpc_application_go.ServiceGroup) * ServiceGroup {
	if group == nil {
		return nil
	}
	policy, _ := CollocationPolicyFromGRPC[group.Policy]
	return &ServiceGroup{
		OrganizationId: group.OrganizationId,
		AppDescriptorId: appDescriptorID,
		ServiceGroupId: GenerateUUID(),
		Name : group.Name,
		Description: group.Description,
		Services: group.Services,
		Policy: policy,
		Specs: NewServiceGroupDeploymentSpecsFromGRPC(group.Specs),
	}
}

func (sg *ServiceGroup) ToGRPC() *grpc_application_go.ServiceGroup {
	policy, _ := CollocationPolicyToGRPC[sg.Policy]
	return &grpc_application_go.ServiceGroup{
		OrganizationId:  sg.OrganizationId,
		AppDescriptorId: sg.AppDescriptorId,
		ServiceGroupId:  sg.ServiceGroupId,
		Name:            sg.Name,
		Description:     sg.Description,
		Services:        sg.Services,
		Policy:          policy,
		Specs:           sg.Specs.ToGRPC(),
	}
}

//func (s * Service) ToServiceInstance(appInstanceID string) * ServiceInstance {
func (sg * ServiceGroup) GetService(serviceId string, descriptor * AppDescriptor) * Service {

	for _, service := range descriptor.Services {
		if service.ServiceId == serviceId {
			return &service
		}
	}

	return nil
}
// -- ServiceGroup -> ServiceGroupInstance -- //
func (sg * ServiceGroup) ToServiceGroupInstance(appInstID string, descriptor * AppDescriptor) * ServiceGroupInstance {

	// get Services from descriptor
	serviceInstances := make([]ServiceInstance, 0)
	for _,serviceId := range sg.Services {
		service := sg.GetService(serviceId, descriptor)
		if service != nil {
			serviceInstances = append(serviceInstances, *service.ToServiceInstance(appInstID))
		}
	}

	return &ServiceGroupInstance{
		OrganizationId:   sg.OrganizationId,
		AppDescriptorId:  sg.AppDescriptorId,
		AppInstanceId:    appInstID,
		ServiceGroupId:   sg.ServiceGroupId,
		Name:             sg.Name,
		Description:      sg.Description,
		ServiceInstances: serviceInstances,
		Policy:           sg.Policy,
		Specs:			  sg.Specs,
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
	// Name of the service group.
	Name string `json:"name,omitempty" cql:"name"`
	// Description of the service group.
	Description string `json:"description,omitempty" cql:"description"`
	// ServicesInstances defining a list of service identifiers that belong to the group.
	ServiceInstances []ServiceInstance `json:"service_instances,omitempty" cql:"service_instances"`
	// Policy indicating the deployment collocation policy.
	Policy               CollocationPolicy `json:"policy,omitempty" cql:"policy"`
	// The status for this service group instance will be the worst status of its services
	Status ServiceStatus `json:"status,omitempty" cql:"status"`
	// Metadata for this service group
	Metadata * InstanceMetadata `json:"metadata,omitempty" cql:"metadata"`
	// Particular deployment specs for this service
	Specs                *ServiceGroupDeploymentSpecs `json:"specs,omitempty" cql:"specs""`
}

func (sgi *ServiceGroupInstance) ToGRPC() *grpc_application_go.ServiceGroupInstance {

	services := make ([]*grpc_application_go.ServiceInstance, 0)
	for _, instance := range sgi.ServiceInstances {
		services = append(services, instance.ToGRPC() )
	}

	policy, _ := CollocationPolicyToGRPC[sgi.Policy]
	return &grpc_application_go.ServiceGroupInstance{
		OrganizationId:       sgi.OrganizationId,
		AppDescriptorId:      sgi.AppDescriptorId,
		AppInstanceId:        sgi.AppInstanceId,
		ServiceGroupId:       sgi.ServiceGroupId,
		Name:                 sgi.Name,
		Description:          sgi.Description,
		ServiceInstances:     services,
		Policy:               policy,
		Status:				  ServiceStatusToGRPC[sgi.Status],
		Metadata:			  sgi.Metadata.ToGRPC(),
		Specs: 				  sgi.Specs.ToGRPC(),
	}
}

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
	return &grpc_application_go.DeploySpecs{
		Cpu:      ds.Cpu,
		Memory:   ds.Memory,
		Replicas: ds.Replicas,
	}
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
)

var EndpointTypeToGRPC = map[EndpointType]grpc_application_go.EndpointType{
	IsAlive:    grpc_application_go.EndpointType_IS_ALIVE,
	Rest:       grpc_application_go.EndpointType_REST,
	Web:        grpc_application_go.EndpointType_WEB,
	Prometheus: grpc_application_go.EndpointType_PROMETHEUS,
}

var EndpointTypeFromGRPC = map[grpc_application_go.EndpointType]EndpointType{
	grpc_application_go.EndpointType_IS_ALIVE:   IsAlive,
	grpc_application_go.EndpointType_REST:       Rest,
	grpc_application_go.EndpointType_WEB:        Web,
	grpc_application_go.EndpointType_PROMETHEUS: Prometheus,
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
	// Content of the configuration file.
	Content []byte `json:"content,omitempty" cql:"content"`
	// MountPath of the configuration file in the service instance.
	MountPath string `json:"mount_path,omitempty" cql:"mount_path"`
}

func NewConfigFileFromGRPC(appDescriptorID string, config * grpc_application_go.ConfigFile) * ConfigFile {
	if config == nil {
		return nil
	}
	return &ConfigFile{
		OrganizationId:  config.OrganizationId,
		AppDescriptorId: appDescriptorID,
		ConfigFileId:    GenerateUUID(),
		Content:         config.Content,
		MountPath:       config.MountPath,
	}
}

func (cf *ConfigFile) ToGRPC() *grpc_application_go.ConfigFile {
	return &grpc_application_go.ConfigFile{
		OrganizationId:  cf.OrganizationId,
		AppDescriptorId: cf.AppDescriptorId,
		ConfigFileId:    cf.ConfigFileId,
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
	// ServiceId with the service identifier.
	ServiceId string `json:"service_id,omitempty" cql:"service_id"`
	// Name of the service.
	Name string `json:"name,omitempty" cql:"name"`
	// Description of the service.
	Description string `json:"description,omitempty" cql:"description"`
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

func NewServiceFromGRPC(appDescriptorID string, service *grpc_application_go.Service) * Service {
	if service == nil{
		return nil
	}

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
		configs = append(configs, *NewConfigFileFromGRPC(appDescriptorID, cf))
	}

	serviceType, _ := ServiceTypeFromGRPC[service.Type]
	return &Service{
		OrganizationId:       service.OrganizationId,
		AppDescriptorId:      appDescriptorID,
		ServiceId:            service.ServiceId,
		Name:                 service.Name,
		Description:          service.Description,
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
		ServiceId:            s.ServiceId,
		Name:                 s.Name,
		Description:          s.Description,
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
func (s * Service) ToServiceInstance(appInstanceID string) * ServiceInstance {

	return &ServiceInstance{
		OrganizationId:       s.OrganizationId,
		AppDescriptorId:      s.AppDescriptorId,
		AppInstanceId:        appInstanceID,
		ServiceId:            s.ServiceId,
		Name:                 s.Name,
		Description:          s.Description,
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
	// ServiceId with the service identifier.
	ServiceId string `json:"service_id,omitempty" cql:"service_id"`
	// Name of the service.
	Name string `json:"name,omitempty" cql:"name"`
	// Description of the service.
	Description string `json:"description,omitempty" cql:"description"`
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
	Endpoints []string `json:"endpoints,omitempty" cql:"endpoints"`
	// DeployedOnClusterId specifies which is the cluster where the service is running.
	DeployedOnClusterId  string  `json:"deployed_on_cluster_id,omitempty" cql:"deployed_on_cluster_id"`
	// RunArguments containts a list of arguments
	RunArguments [] string `json:"run_arguments" cql:"run_arguments"`
	Info                 string   `json:"info,omitempty" cql:"info"`

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
	return &grpc_application_go.ServiceInstance{
		OrganizationId:       si.OrganizationId,
		AppDescriptorId:      si.AppDescriptorId,
		AppInstanceId:        si.AppInstanceId,
		ServiceId:           si.ServiceId,
		Name:                 si.Name,
		Type:                 serviceType,
		Image:                si.Image,
		Credentials:          si.Credentials.ToGRPC(),
		Specs:                si.Specs.ToGRPC(),
		Storage:              storage,
		ExposedPorts:         exposedPorts,
		EnvironmentVariables: si.EnvironmentVariables,
		Configs:              configs,
		Labels:               si.Labels,
		DeployAfter:          si.DeployAfter,
		Status:               serviceStatus,
		Endpoints:            si.Endpoints,
		DeployedOnClusterId:  si.DeployedOnClusterId,
		RunArguments: 		  si.RunArguments,

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
	// Services of the application.
	Services []Service `json:"services,omitempty" cql:"services"`
}

func NewAppDescriptor(organizationID string, appDescriptorID string, name string,
	configOptions map[string]string, envVars map[string]string,
	labels map[string]string,
	rules []SecurityRule, groups []ServiceGroup, services []Service) *AppDescriptor {
	return &AppDescriptor{
		organizationID, appDescriptorID,
		name,
		configOptions,
		envVars,
		labels,
		rules,
		groups,
		services,
	}
}

func NewAppDescriptorFromGRPC(addRequest * grpc_application_go.AddAppDescriptorRequest) * AppDescriptor {

	if addRequest == nil {
		return nil
	}

	uuid := GenerateUUID()

	rules := make([]SecurityRule, 0)
	for _, r := range addRequest.Rules {
		rules = append(rules, *NewSecurityRuleFromGRPC(uuid, r))
	}
	groups := make([]ServiceGroup, 0)
	for _, sg := range addRequest.Groups{
		groups = append(groups, *NewServiceGroupFromGRPC(uuid, sg))
	}
	services := make([]Service, 0)
	for _, s := range addRequest.Services{
		services = append(services, *NewServiceFromGRPC(uuid, s))
	}

	return NewAppDescriptor(
		addRequest.OrganizationId,
		uuid,
		addRequest.Name,
		addRequest.ConfigurationOptions,
		addRequest.EnvironmentVariables,
		addRequest.Labels,
		rules, groups, services)
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
	services := make([]*grpc_application_go.Service, 0)
	for _, s := range d.Services {
		services = append(services, s.ToGRPC())
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
		Services:             services,
	}
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

	if len(toAdd.Services) == 0 {
		return derrors.NewInvalidArgumentError("expecting at least one service")
	}

	for _, s := range toAdd.Services {
		err := ValidAddService(s)
		if err != nil{
			return err
		}
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

}

func NewAppInstanceFromGRPC(addRequest * grpc_application_go.AddAppInstanceRequest, descriptor * AppDescriptor) * AppInstance {
	uuid := GenerateUUID()

	groups := make([]ServiceGroupInstance, 0)
	for _, g := range descriptor.Groups {
		groups = append(groups, *g.ToServiceGroupInstance(uuid, descriptor))
	}


	return &AppInstance{
		OrganizationId:       addRequest.OrganizationId,
		AppDescriptorId:      addRequest.AppDescriptorId,
		AppInstanceId:        uuid,
		Name:                 addRequest.Name,
		ConfigurationOptions: descriptor.ConfigurationOptions,
		EnvironmentVariables: descriptor.EnvironmentVariables,
		Labels:               descriptor.Labels,
		Rules:                descriptor.Rules,
		Groups:               groups,
		Status: Queued,
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
	}
}

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
	if updateRequest.OrganizationId == "" || updateRequest.AppInstanceId == "" || updateRequest.ServiceId == "" {
			return derrors.NewInvalidArgumentError("expecting orgazation_id, app_instance_id, app_descriptor_id " +
				"and service_id")
	}
	return nil
}