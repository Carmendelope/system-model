package application

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"math/rand"
)

var organizationId = fmt.Sprintf("organization_%d", rand.Intn(100)+1)
var appDescriptorId = fmt.Sprintf("app_descriptor_%d", rand.Intn(100)+1)
var name = "Application name"
var description = "Application description"
var confLabel = "Conf"
var confValue = "Conf_value"
var envLabel = "Env"
var envValue = "Env_value"
var labelLabel = "lab1"
var labelValue = "LABEL1"
var ruleId = "rule id 1"
var ruleName = "Rule name"
var sourceServiceId = "SourceServiceId1"
var serviceGroupId = "service_group_id1"
var serviceGroupName = "service_group name"
var serviceGroupDescription = "service_group description"
var serviceId = "service id_1"
var serviceDescription = "service description"
var serviceName = "service name"
var image = "../../image_path"


func CreateTestConfigFile () entities.ConfigFile {

	content := make([]byte, 0)
	content = append(content, 0x00, 0x01, 0x02)
	return entities.ConfigFile{
		OrganizationId: organizationId,
		AppDescriptorId:appDescriptorId,
		ConfigFileId: "Config file",
		Content: content,
		MountPath: "../../path"}
}

func CreateTestServiceInstance (appInstanceId string) entities.ServiceInstance {

	stores := make ([]entities.Storage, 0)
	stores = append(stores, entities.Storage{Size:900, MountPath:"../../mount_path", Type:entities.StorageType(1)})

	endpoints := make ([]entities.Endpoint, 0)
	endpoints = append(endpoints, entities.Endpoint {Type:entities.EndpointType(1),Path:"../../endpoint" })

	ports := make ([]entities.Port, 0)
	for i:=0; i<5; i++{
		ports = append(ports, entities.Port{Name:fmt.Sprintf("port%d", i), InternalPort:int32(i), ExposedPort:int32(i), Endpoints:endpoints})
	}

	envVariables := make(map[string]string, 0)
	envVariables["HOST"] = "HOST_VALUE"
	envVariables["PORT"] = "PORT_VALUE"

	confFile := make ([]entities.ConfigFile, 0)
	confFile = append(confFile, CreateTestConfigFile())

	labels := make (map[string]string, 0)
	for i:=0; i<4; i++{
		labels[fmt.Sprintf("label%d", i)] = fmt.Sprintf("value_%d", i)
	}

	deployAfter := make([]string, 0)
	deployAfter = append(deployAfter, "deploy after this", "and this")

	return entities.ServiceInstance{
		OrganizationId: organizationId,
		AppDescriptorId: appDescriptorId,
		AppInstanceId: appInstanceId,
		ServiceId: serviceId,
		Name: serviceName,
		Description: serviceDescription,
		Type: entities.ServiceType(1),
		Image: image,
		Credentials: &entities.ImageCredentials{
			Username: "carmen",
			Password:"*****",
			Email: "cdelope@daisho.group",
			DockerRepository: "DOCKER REPOSITORY!!!!!"},
		// DeploySpecs with the resource specs required by the service.
		Specs: &entities.DeploySpecs{
			Cpu: 1239900,
			Memory:2000,
			Replicas:2},
		Storage: stores,
		ExposedPorts: ports,
		EnvironmentVariables: envVariables,
		Configs: confFile,
		Labels: labels,
		DeployAfter: deployAfter,
		Status: entities.ServiceStatus(1),
		DeployedOnClusterId:"ClusterIDXXX",
		Endpoints:make([]string,0)}

}

func CreateTestService (appDescriptorId string) entities.Service {

	stores := make ([]entities.Storage, 0)
	stores = append(stores, entities.Storage{Size:900, MountPath:"../../mount_path", Type:entities.StorageType(1)})

	endpoints := make ([]entities.Endpoint, 0)
	endpoints = append(endpoints, entities.Endpoint {Type:entities.EndpointType(1),Path:"../../endpoint" })

	ports := make ([]entities.Port, 0)
	for i:=0; i<5; i++{
		ports = append(ports, entities.Port{Name:fmt.Sprintf("port%d", i), InternalPort:int32(i), ExposedPort:int32(i), Endpoints:endpoints})
	}

	envVariables := make(map[string]string, 0)
	envVariables["HOST"] = "HOST_VALUE"
	envVariables["PORT"] = "PORT_VALUE"

	confFile := make ([]entities.ConfigFile, 0)
	confFile = append(confFile, CreateTestConfigFile())

	labels := make (map[string]string, 0)
	for i:=0; i<4; i++{
		labels[fmt.Sprintf("label%d", i)] = fmt.Sprintf("value_%d", i)
	}

	deployAfter := make([]string, 0)
	deployAfter = append(deployAfter, "deploy after this", "and this")

	return entities.Service{
		OrganizationId: organizationId,
		AppDescriptorId: appDescriptorId,
		ServiceId: serviceId,
		Name: serviceName,
		Description: serviceDescription,
		Type: entities.ServiceType(1),
		Image: image,
		Credentials: &entities.ImageCredentials{
			Username: "carmen",
			Password:"*****",
			Email: "cdelope@daisho.group"},
		// DeploySpecs with the resource specs required by the service.
		Specs: &entities.DeploySpecs{
			Cpu: 1239900,
			Memory:2000,
			Replicas:2},
		Storage: stores,
		ExposedPorts: ports,
		EnvironmentVariables: envVariables,
		Configs: confFile,
		Labels: labels,
		DeployAfter: deployAfter}
}

func CreateTestServiceGroupInstance(appInstanceId string) entities.ServiceGroupInstance{

	servicesInstances := make ([]string, 0)
	for i:=0; i<5; i++{
		servicesInstances = append(servicesInstances, fmt.Sprintf("servicesInstances-%d",i ))
	}

	return entities.ServiceGroupInstance{
		OrganizationId: organizationId,
		AppDescriptorId: appDescriptorId,
		AppInstanceId:  appInstanceId,
		ServiceGroupId: serviceGroupId,
		Name: serviceGroupName,
		Description: serviceGroupDescription,
		ServiceInstances:servicesInstances,
		Policy: entities.CollocationPolicy(1) }
}

func CreateTestServiceGroup(appDescriptorId string) entities.ServiceGroup{

	services := make ([]string, 0)
	for i:=0; i<5; i++{
		services = append(services, fmt.Sprintf("services-%d",i ))
	}

	return entities.ServiceGroup{
		OrganizationId: organizationId,
		AppDescriptorId: appDescriptorId,
		ServiceGroupId: serviceGroupId,
		Name: serviceGroupName,
		Description: serviceGroupDescription,
		Services:services,
		Policy: entities.CollocationPolicy(1) }
}

func CreateTestRule() entities.SecurityRule {

	id := rand.Intn(10) + 1
	authServices := make ([]string, 0)
	for i:=0; i<10; i++{
		authServices = append(authServices, fmt.Sprintf("auth%d",i ))
	}
	devices := make ([]string, 0)
	for i:=0; i<6; i++{
		devices = append(devices, fmt.Sprintf("device%d",i ))
	}

	rule := entities.SecurityRule{
		OrganizationId:organizationId,
		AppDescriptorId:appDescriptorId,
		RuleId: fmt.Sprintf("ruleId_%d", id),
		Name: ruleName,
		SourceServiceId: sourceServiceId,
		SourcePort: 80,
		Access: 0,
		AuthServices: authServices,
		DeviceGroups: devices}

	return rule
}

func CreateTestApplication(id string) *entities.AppInstance {

	appInstanceId := fmt.Sprintf("App instance Id_%s", id)
	configurationOptions := make(map[string]string, 0)
	configurationOptions[confLabel] = confValue

	environmentVariables := make(map[string]string, 0)
	environmentVariables[envLabel] = envValue

	labels := make(map[string]string, 0)
	labels[labelLabel] = labelValue

	rules := make([]entities.SecurityRule, 0)
	rules = append(rules, CreateTestRule())

	groups := make ([]entities.ServiceGroupInstance, 0)
	groups = append(groups, CreateTestServiceGroupInstance(appInstanceId))

	services := make ([]entities.ServiceInstance, 0)
	services = append(services, CreateTestServiceInstance(appInstanceId))

	app := entities.AppInstance{
		OrganizationId:organizationId,
		AppDescriptorId: appDescriptorId,
		AppInstanceId: appInstanceId,
		Name: name,
		Description: description,
		ConfigurationOptions: configurationOptions,
		EnvironmentVariables:environmentVariables,
		Labels:labels,
		Rules: rules,
		Groups:groups,
		Services: services,
		Status: entities.ApplicationStatus(1),
		}

	return &app
}

func CreateTestApplicationDescriptor (id string) *entities.AppDescriptor {

	appDescriptor := fmt.Sprintf("App_descriptor_%s", id)

	tam := rand.Intn(4) + 1
	configurationOptions := make(map[string]string, 0)
	environmentVariables := make(map[string]string, 0)
	labels := make(map[string]string, 0)

	for i:= 0; i< tam; i++{
		configurationOptions[fmt.Sprintf("conf-%d", i)] = fmt.Sprintf("conf_value-%d", i)
		environmentVariables[fmt.Sprintf("env-%d", i)] = fmt.Sprintf("env_value-%d", i)
		labels[fmt.Sprintf("label-%d", i)] = fmt.Sprintf("label_value-%d", i)
	}

	rules := make([]entities.SecurityRule, 0)
	rules = append(rules, CreateTestRule())

	groups := make([]entities.ServiceGroup, 0)
	groups = append(groups, CreateTestServiceGroup(appDescriptor))

	services := make ([]entities.Service, 0)
	services = append(services, CreateTestService(appDescriptor))

	descriptor := entities.AppDescriptor{
		OrganizationId: organizationId,
		AppDescriptorId:appDescriptor,
		Name: fmt.Sprintf("%s name", appDescriptor),
		Description: fmt.Sprintf("%s description", appDescriptor),
		ConfigurationOptions:configurationOptions,
		EnvironmentVariables:environmentVariables,
		Labels:labels,
		Rules:rules,
		Groups: groups,
		Services: services}

	return &descriptor

}