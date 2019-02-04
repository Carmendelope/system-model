package application

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nalej/system-model/internal/pkg/entities"
	"math/rand"
)

//var organizationId = fmt.Sprintf("organization_%d", rand.Intn(100)+1)
//var appDescriptorId = fmt.Sprintf("app_descriptor_%d", rand.Intn(100)+1)
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


func CreateTestConfigFile (organizationID string, appDescriptorID string) entities.ConfigFile {

	content := make([]byte, 0)
	content = append(content, 0x00, 0x01, 0x02)
	return entities.ConfigFile{
		OrganizationId: organizationID,
		AppDescriptorId:appDescriptorID,
		ConfigFileId: "Config file",
		Content: content,
		MountPath: "../../path"}
}

func CreateTestServiceInstance (organizationID string, appDescriptorID string, appInstanceId string) entities.ServiceInstance {

	stores := make ([]entities.Storage, 0)
	stores = append(stores, entities.Storage{Size:900, MountPath:"../../mount_path", Type:entities.StorageType(1)})

	endpoints := make ([]entities.Endpoint, 0)
	endpoints = append(endpoints, entities.Endpoint {Type:entities.EndpointType(1),Path:"../../endpoint" })

	ports := make ([]entities.Port, 0)
	for i:=0; i<1; i++{
		ports = append(ports, entities.Port{Name:fmt.Sprintf("port%d", i), InternalPort:int32(i), ExposedPort:int32(i), Endpoints:endpoints})
	}

	envVariables := make(map[string]string, 0)
	envVariables["HOST"] = "HOST_VALUE"
	envVariables["PORT"] = "PORT_VALUE"

	confFile := make ([]entities.ConfigFile, 0)
	confFile = append(confFile, CreateTestConfigFile(organizationID, appDescriptorID))

	labels := make (map[string]string, 0)
	for i:=0; i<4; i++{
		labels[fmt.Sprintf("label%d", i)] = fmt.Sprintf("value_%d", i)
	}

	deployAfter := make([]string, 0)
	deployAfter = append(deployAfter, "deploy after this", "and this")

	return entities.ServiceInstance{
		OrganizationId: organizationID,
		AppDescriptorId: appDescriptorID,
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
		Endpoints:[]string{"endpoint1", "endpoint2", "endpoint3"},
		RunArguments: []string{"arg1, agr2"},
		Info: "info",
	}

}

func CreateTestService (organizationID string, appDescriptorId string) entities.Service {

	stores := make ([]entities.Storage, 0)
	stores = append(stores, entities.Storage{Size:900, MountPath:"../../mount_path", Type:entities.StorageType(1)})

	endpoints := make ([]entities.Endpoint, 0)
	endpoints = append(endpoints, entities.Endpoint {Type:entities.EndpointType(1),Path:"../../endpoint" })

	ports := make ([]entities.Port, 0)
	for i:=0; i<5; i++{
		ports = append(ports, entities.Port{Name:fmt.Sprintf("port%d", i), InternalPort:int32(i), ExposedPort:int32(i), Endpoints:endpoints})
	}

	envVariables := map[string]string{"HOST":"HOST_VALUE", "PORT":"PORT_VALUE"}

	confFile := make ([]entities.ConfigFile, 0)
	confFile = append(confFile, CreateTestConfigFile(organizationID, appDescriptorId))

	labels := make (map[string]string, 0)
	for i:=0; i<4; i++{
		labels[fmt.Sprintf("label_%d", i)] = fmt.Sprintf("value_%d", i)
	}

	deployAfter := []string{"deploy after this", "and this"}
	runArguments := [] string{"arg1", "arg2", "arg3", "arg4"}

	return entities.Service{
		OrganizationId: organizationID,
		AppDescriptorId: appDescriptorId,
		ServiceId: serviceId,
		Name: serviceName,
		Description: serviceDescription,
		Type: entities.ServiceType(1),
		Image: image,
		Credentials: &entities.ImageCredentials{
			Username: "carmen",
			Password:"*****",
			Email: "cdelope@daisho.group",
			DockerRepository:"DOCKER"},
		// DeploySpecs with the resource specs required by the service.
		Specs: &entities.DeploySpecs{
			Cpu: 1239900,
			Memory:2000,
			Replicas:2,
		},
		Storage: stores,
		ExposedPorts: ports,
		EnvironmentVariables: envVariables,
		Configs: confFile,
		Labels: labels,
		DeployAfter: deployAfter,
		RunArguments:runArguments,
	}
}

func CreateTestServiceGroupInstance(organizationID string, appDescriptorID string, appInstanceId string) entities.ServiceGroupInstance{

	servicesInstances := make ([]entities.ServiceInstance, 0)
	for i:=0; i<1; i++{
		servicesInstances = append(servicesInstances, CreateTestServiceInstance(organizationID, appDescriptorID, appInstanceId))
	}

	return entities.ServiceGroupInstance{
		OrganizationId: organizationID,
		AppDescriptorId: appDescriptorID,
		AppInstanceId:  appInstanceId,
		ServiceGroupId: serviceGroupId,
		Name: serviceGroupName,
		Description: serviceGroupDescription,
		Policy: entities.CollocationPolicy(1),
		ServiceInstances: servicesInstances,
		Status: entities.ServiceScheduled,
		Metadata: &entities.InstanceMetadata{
			OrganizationId: organizationID,
			AppDescriptorId: appDescriptorID,
			AppInstanceId:  appInstanceId,
			MonitoredInstanceId: uuid.New().String(),
			Type: entities.ServiceGroupInstanceType,
			InstancesId: []string{"inst1", "inst2", "inst3"},
			DesiredReplicas: 3,
			AvailableReplicas:3,
			UnavailableReplicas: 0,
			Status: map[string]entities.ServiceStatus {"status1": entities.ServiceError, "status2": entities.ServiceDeploying},
			Info:map[string]string{"sgInfo1": "info1"},
		},
		Specs: &entities.ServiceGroupDeploymentSpecs{
			NumReplicas: 3,
			MultiClusterReplica: true,
		},
	}
}

func CreateTestServiceGroup(organizationID string, appDescriptorId string) entities.ServiceGroup{

	services := make ([]string, 0)
	for i:=0; i<1; i++{
		services = append(services, fmt.Sprintf("services-%d",i ))
	}

	specs := &entities.ServiceGroupDeploymentSpecs {
		NumReplicas: 5,
		MultiClusterReplica: false,
	}

	return entities.ServiceGroup{
		OrganizationId: organizationID,
		AppDescriptorId: appDescriptorId,
		ServiceGroupId: serviceGroupId,
		Name: serviceGroupName,
		Description: serviceGroupDescription,
		Services:services,
		Policy: entities.CollocationPolicy(1),
		Specs: specs,
	}
}

func CreateTestRule(organizationID string, appDescriptorID string) entities.SecurityRule {

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
		OrganizationId:organizationID,
		AppDescriptorId:appDescriptorID,
		RuleId: fmt.Sprintf("ruleId_%d", id),
		Name: ruleName,
		SourceServiceId: sourceServiceId,
		SourcePort: 80,
		Access: 0,
		AuthServices: authServices,
		DeviceGroups: devices}

	return rule
}

func CreateTestApplication(organizationID string, appDescriptorID string) *entities.AppInstance {

	id:= uuid.New().String()

	configurationOptions := map[string]string {"conf1":"value1", "conf2":"value2","conf3":"value3"}
	environmentVariables := map[string]string {"env01":"value1", "env02":"value2"}
	labels := map[string]string{"label1":"value1", "label2":"value2","label3":"value3", "label4":"value4"}

	rules := make([]entities.SecurityRule, 0)
	rules = append(rules, CreateTestRule(organizationID, appDescriptorID))

	groups := make ([]entities.ServiceGroupInstance, 0)
	groups = append(groups, CreateTestServiceGroupInstance(organizationID, appDescriptorID, id))

	services := make ([]entities.ServiceInstance, 0)
	services = append(services, CreateTestServiceInstance(organizationID, appDescriptorID, id))

	app := entities.AppInstance{
		OrganizationId:organizationID,
		AppDescriptorId: appDescriptorID,
		AppInstanceId: id,
		Name: "App instance Name",
		ConfigurationOptions: configurationOptions,
		EnvironmentVariables:environmentVariables,
		Labels:labels,
		Rules: rules,
		Groups:groups,
		Status: entities.ApplicationStatus(1),
		}

	return &app
}

func CreateTestApplicationDescriptor (organizationID string) *entities.AppDescriptor {

	id := uuid.New().String()

	tam := rand.Intn(4) + 1
	configurationOptions := make(map[string]string, 0)
	environmentVariables := make(map[string]string, 0)
	labels := make(map[string]string, 0)

	for i:= 0; i< tam; i++{
		configurationOptions[fmt.Sprintf("conf_%d", i)] = fmt.Sprintf("conf_value_%d", i)
		environmentVariables[fmt.Sprintf("env_%d", i)] = fmt.Sprintf("env_value_%d", i)
		labels[fmt.Sprintf("label-%d", i)] = fmt.Sprintf("label_value-%d", i)
	}

	rules := make([]entities.SecurityRule, 0)
	rules = append(rules, CreateTestRule(organizationID, id))

	groups := make([]entities.ServiceGroup, 0)
	groups = append(groups, CreateTestServiceGroup(organizationID,id))

	services := make ([]entities.Service, 0)
	services = append(services, CreateTestService(organizationID, id))

	descriptor := entities.AppDescriptor{
		OrganizationId: organizationID,
		AppDescriptorId: id,
		Name: "Descriptor-Name",
		ConfigurationOptions:configurationOptions,
		EnvironmentVariables:environmentVariables,
		Labels:labels,
		Rules:rules,
		Groups: groups,
		Services: services}

	return &descriptor

}