/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package entities

import "github.com/nalej/grpc-application-go"

type PortAccess int

const (
	AllAppServices PortAccess = iota + 1
	AppServices
	Public
	DeviceGroup
)

var PortAccessToGRPC = map[PortAccess]grpc_application_go.PortAccess{
	AllAppServices : grpc_application_go.PortAccess_ALL_APP_SERVICES,
	AppServices : grpc_application_go.PortAccess_APP_SERVICES,
	Public : grpc_application_go.PortAccess_PUBLIC,
	DeviceGroup : grpc_application_go.PortAccess_DEVICE_GROUP,
}

type CollocationPolicy int

const (
	SameCluster CollocationPolicy = iota +1
	SeparateClusters
)

var CollocationPolicyToGRPC = map[CollocationPolicy] grpc_application_go.CollocationPolicy {
	SameCluster : grpc_application_go.CollocationPolicy_SAME_CLUSTER,
	SeparateClusters : grpc_application_go.CollocationPolicy_SEPARATE_CLUSTERS,
}

type SecurityRule struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// RuleId with the security rule identifier.
	RuleId string `json:"rule_id,omitempty"`
	// Name of the security rule.
	Name string `json:"name,omitempty"`
	// SourceServiceId defining the service onto which the security rule is defined.
	SourceServiceId string `json:"source_service_id,omitempty"`
	// SourcePort defining the port that is affected by the current rule.
	SourcePort int32 `json:"source_port,omitempty"`
	// Access level to that port defining who can access the port.
	Access PortAccess `json:"access,omitempty"`
	// AuthServices defining a list of services that can access the port.
	AuthServices []string `json:"auth_services,omitempty"`
	// DeviceGroups defining a list of device groups that can access the port.
	DeviceGroups []string `json:"device_groups,omitempty"`

}

func (sr * SecurityRule) ToGRPC() * grpc_application_go.SecurityRule {
	access, _ := PortAccessToGRPC[sr.Access]
	return &grpc_application_go.SecurityRule{
		OrganizationId: sr.OrganizationId,
		AppDescriptorId:sr.AppDescriptorId,
		RuleId:sr.RuleId,
		Name:sr.Name,
		SourceServiceId:sr.SourceServiceId, SourcePort:sr.SourcePort,
		Access:access,
		AuthServices:sr.AuthServices,
		DeviceGroups:sr.DeviceGroups,
	}
}

type ServiceGroup struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// ServiceGroupId with the group identifier.
	ServiceGroupId string `json:"service_group_id,omitempty"`
	// Name of the service group.
	Name string `protobuf:"json:"name,omitempty"`
	// Description of the service group.
	Description string `json:"description,omitempty"`
	// Services defining a list of service identifiers that belong to the group.
	Services []string `json:"services,omitempty"`
	// Policy indicating the deployment collocation policy.
	Policy CollocationPolicy `json:"policy,omitempty"`

}

func (sg * ServiceGroup) ToGRPC() *grpc_application_go.ServiceGroup {
	policy, _ := CollocationPolicyToGRPC[sg.Policy]
	return &grpc_application_go.ServiceGroup{
		OrganizationId: sg.OrganizationId,
		AppDescriptorId: sg.AppDescriptorId,
		ServiceGroupId: sg.ServiceGroupId,
		Name: sg.Name,
		Description: sg.Description,
		Services: sg.Services,
		Policy: policy,
	}
}

type ServiceType int32

const (
	DockerService ServiceType = iota + 1
)

var ServiceTypeToGRPC = map[ServiceType] grpc_application_go.ServiceType {
	DockerService : grpc_application_go.ServiceType_DOCKER,
}

type ImageCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
}

func (ic * ImageCredentials) ToGRPC() *grpc_application_go.ImageCredentials {
	return &grpc_application_go.ImageCredentials{
		Username: ic.Username,
		Password: ic.Password,
		Email: ic.Email,
	}
}

type DeploySpecs struct {
	Cpu      int64 `json:"cpu,omitempty"`
	Memory   int64 `json:"memory,omitempty"`
	Replicas int32 `json:"replicas,omitempty"`
}

func (ds * DeploySpecs) ToGRPC() *grpc_application_go.DeploySpecs {
	return &grpc_application_go.DeploySpecs{
		Cpu: ds.Cpu,
		Memory: ds.Memory,
		Replicas: ds.Replicas,
	}
}

type StorageType int32

const (
	Ephemeral        StorageType = iota + 1
	ClusterLocal
	ClusterReplica
	CloudPersistent
)

var StorageTypeToGRPC = map[StorageType] grpc_application_go.StorageType {
	Ephemeral : grpc_application_go.StorageType_EPHEMERAL,
	ClusterLocal : grpc_application_go.StorageType_CLUSTER_LOCAL,
	ClusterReplica : grpc_application_go.StorageType_CLUSTER_REPLICA,
	CloudPersistent : grpc_application_go.StorageType_CLOUD_PERSISTENT,
}

type Storage struct {
	Size      int64       `json:"size,omitempty"`
	MountPath string      `json:"mount_path,omitempty"`
	Type      StorageType `json:"type,omitempty"`
}

func (s * Storage) ToGRPC() *grpc_application_go.Storage {
	convertedType, _ := StorageTypeToGRPC[s.Type]
	return &grpc_application_go.Storage{
		Size: s.Size,
		MountPath: s.MountPath,
		Type: convertedType,
	}
}

type EndpointType int

const(
	IsAlive   EndpointType = iota + 1
	Rest
	Web
	Prometheus
)

var EndpointTypeToGRPC = map[EndpointType] grpc_application_go.EndpointType {
	IsAlive : grpc_application_go.EndpointType_IS_ALIVE,
	Rest : grpc_application_go.EndpointType_REST,
	Web : grpc_application_go.EndpointType_WEB,
	Prometheus : grpc_application_go.EndpointType_PROMETHEUS,
}

type Endpoint struct {
	Type EndpointType `json:"type,omitempty"`
	Path string       `json:"path,omitempty"`
}

func (e * Endpoint) ToGRPC() *grpc_application_go.Endpoint {
	convertedType, _ := EndpointTypeToGRPC[e.Type]
	return &grpc_application_go.Endpoint{
		Type: convertedType,
		Path: e.Path,
	}
}


type Port struct {
	Name         string      `json:"name,omitempty"`
	InternalPort int32       `json:"internal_port,omitempty"`
	ExposedPort  int32       `json:"exposed_port,omitempty"`
	Endpoints    []Endpoint `json:"endpoints,omitempty"`
}

func (p *Port) ToGRPC() *grpc_application_go.Port  {
	endpoints := make([]*grpc_application_go.Endpoint, 0)

	for _, ep := range p.Endpoints {
		endpoints = append(endpoints, ep.ToGRPC())
	}

	return &grpc_application_go.Port{
		Name: p.Name,
		InternalPort: p.InternalPort,
		ExposedPort: p.ExposedPort,
		Endpoints: endpoints,
	}
}

type ConfigFile struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// ConfigFileId with the config file identifier.
	ConfigFileId string `json:"config_file_id,omitempty"`
	// Content of the configuration file.
	Content []byte `json:"content,omitempty"`
	// MountPath of the configuration file in the service instance.
	MountPath string `json:"mount_path,omitempty"`
}

func (cf * ConfigFile) ToGRPC() *grpc_application_go.ConfigFile {
	return &grpc_application_go.ConfigFile{
		OrganizationId:       cf.OrganizationId,
		AppDescriptorId:      cf.AppDescriptorId,
		ConfigFileId:         cf.ConfigFileId,
		Content:              cf.Content,
		MountPath:            cf.MountPath,
	}
}

type Service struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// ServiceId with the service identifier.
	ServiceId string `json:"service_id,omitempty"`
	// Name of the service.
	Name string `json:"name,omitempty"`
	// Description of the service.
	Description string `json:"description,omitempty"`
	// ServiceType represents the underlying technology of the service to be launched.
	Type ServiceType `json:"type,omitempty"`
	// Image contains the URL/name of the image to be executed.
	Image string `json:"image,omitempty"`
	// ImageCredentials with the data required to access the repository the image is available at.
	Credentials *ImageCredentials `json:"credentials,omitempty"`
	// DeploySpecs with the resource specs required by the service.
	Specs DeploySpecs `json:"specs,omitempty"`
	// Storage restrictions
	Storage []Storage `json:"storage,omitempty"`
	// ExposedPorts contains the list of ports exposed by the current service.
	ExposedPorts []Port `json:"exposed_ports,omitempty"`
	// EnvironmentVariables defines a key-value map of environment variables and values that will be passed to all
	// running services.
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty"`
	// Configs contains the configuration files required by the service.
	Configs []ConfigFile `json:"configs,omitempty"`
	// Labels with the user defined labels.
	Labels map[string]string `json:"labels,omitempty"`
	// DeployAfter contains the list of services that must be running before launching a service.
	DeployAfter          []string `json:"deploy_after,omitempty"`

}

func (s * Service) ToGRPC() *grpc_application_go.Service {
	serviceType, _ := ServiceTypeToGRPC[s.Type]
	storage := make([]*grpc_application_go.Storage, 0)
	for _, s := range s.Storage{
		storage = append(storage, s.ToGRPC())
	}
	exposedPorts := make([]*grpc_application_go.Port, 0)
	for _, ep := range s.ExposedPorts{
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
	}
}

type AppDescriptor struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty"`
	// AppDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty"`
	// Name of the application.
	Name string `json:"name,omitempty"`
	// Description of the application.
	Description string `json:"description,omitempty"`
	// ConfigurationOptions defines a key-value map of configuration options.
	ConfigurationOptions map[string]string `json:"configuration_options,omitempty"`
	// EnvironmentVariables defines a key-value map of environment variables and values that will be passed to all
	// running services.
	EnvironmentVariables map[string]string `json:"environment_variables,omitempty"`
	// Labels defined by the user.
	Labels map[string]string `json:"labels,omitempty"`
	// Rules that define the connectivity between the elements of an application.
	Rules []SecurityRule `json:"rules,omitempty"`
	// Groups with the Service collocation strategies.
	Groups []ServiceGroup `json:"groups,omitempty"`
	// Services of the application.
	Services             []Service `json:"services,omitempty"`
}

func NewAppDescriptor(organizationID string, name string, description string,
	configOptions map[string]string, envVars map[string]string,
	labels map[string]string,
	rules []SecurityRule, groups []ServiceGroup, services []Service) * AppDescriptor {
	uuid := GenerateUUID(AppDescPrefix)
	return &AppDescriptor{
		organizationID, uuid,
		name, description,
		configOptions,
		envVars,
		labels,
		rules,
		groups,
		services,
	}
}

func (d * AppDescriptor) ToGRPC() *grpc_application_go.AppDescriptor {

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
		Description:          d.Description,
		ConfigurationOptions: d.ConfigurationOptions,
		EnvironmentVariables: d.EnvironmentVariables,
		Labels:               d.Labels,
		Rules:                rules,
		Groups:               groups,
		Services:             services,
	}
}

type AppInstance struct {

}

func (i * AppInstance) ToGRPC() * grpc_application_go.AppInstance {
	return nil
}