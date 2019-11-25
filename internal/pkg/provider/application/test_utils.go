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
 */

package application

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nalej/system-model/internal/pkg/entities"
	"math/rand"
)

func CreateTestConfigFile(organizationID string, appDescriptorID string) entities.ConfigFile {

	content := make([]byte, 0)
	content = append(content, 0x00, 0x01, 0x02)
	return entities.ConfigFile{
		OrganizationId:  organizationID,
		AppDescriptorId: appDescriptorID,
		ConfigFileId:    "Config file",
		Content:         content,
		MountPath:       "../../path",
		Name:            "configFileName",
	}
}

func CreateTestInstanceMetadata(organizationID string, AppDescriptorID string, AppInstanceID string) entities.InstanceMetadata {

	return entities.InstanceMetadata{
		OrganizationId:      organizationID,
		AppDescriptorId:     AppDescriptorID,
		AppInstanceId:       AppInstanceID,
		MonitoredInstanceId: uuid.New().String(),
		Type:                entities.ServiceGroupInstanceType,
		InstancesId:         []string{"instance1", "instance2", "instance3"},
		DesiredReplicas:     2,
		AvailableReplicas:   2,
		UnavailableReplicas: 0,
		Status:              map[string]entities.ServiceStatus{"instance1": entities.ServiceScheduled, "instance2": entities.ServiceError},
		Info:                map[string]string{"info1": "value1(InstanceMetadata)", "info2": "value2(InstanceMetadata)"},
	}
}

func CreateTestServiceInstance(organizationID string, appDescriptorID string, appInstanceId string,
	serviceGroupId string, ServiceGroupInstanceId string) entities.ServiceInstance {

	id := uuid.New().String()
	stores := make([]entities.Storage, 0)
	stores = append(stores, entities.Storage{Size: 900, MountPath: "../../mount_path", Type: entities.StorageType(1)})

	endpoints := make([]entities.EndpointInstance, 0)
	endpoints = []entities.EndpointInstance{
		{
			EndpointInstanceId: uuid.New().String(),
			Fqdn:               "../Fqdn_path",
			Type:               entities.IsAlive,
		},
	}

	ports := make([]entities.Port, 0)
	ports = append(ports, entities.Port{
		Name:         "port name",
		InternalPort: 80,
		ExposedPort:  80,
		Endpoints: []entities.Endpoint{
			{
				Path: "../enpoint_path",
				Type: entities.IsAlive,
			},
		}})

	confFile := make([]entities.ConfigFile, 0)
	confFile = append(confFile, CreateTestConfigFile(organizationID, appDescriptorID))

	labels := make(map[string]string, 0)
	for i := 0; i < 4; i++ {
		labels[fmt.Sprintf("label%d", i)] = fmt.Sprintf("value_%d", i)
	}

	deployAfter := make([]string, 0)
	deployAfter = append(deployAfter, "deploy after this", "and this")

	return entities.ServiceInstance{
		OrganizationId:         organizationID,
		AppDescriptorId:        appDescriptorID,
		AppInstanceId:          appInstanceId,
		ServiceGroupId:         serviceGroupId,
		ServiceGroupInstanceId: ServiceGroupInstanceId,
		ServiceId:              id,
		ServiceInstanceId:      uuid.New().String(),
		Name:                   "Service Instance Name",
		Type:                   entities.DockerService,
		Image:                  "../image.txt",
		Credentials: &entities.ImageCredentials{
			Username:         "carmen",
			Password:         "*****",
			Email:            "cdelope@daisho.group",
			DockerRepository: "Docker Repository"},
		Specs: &entities.DeploySpecs{
			Cpu:      1239900,
			Memory:   2000,
			Replicas: 2},
		Storage:              stores,
		ExposedPorts:         ports,
		EnvironmentVariables: map[string]string{"env1": "env1(serviceInstance)", "env2": "env2(ServiceInstance)"},
		Configs:              confFile,
		Labels:               map[string]string{"label1": "label1(serviceInstance)"},
		DeployAfter:          []string{"this", "and this"},
		Status:               entities.ServiceScheduled,
		DeployedOnClusterId:  "Cluster id",
		Endpoints:            endpoints,
		RunArguments:         []string{"arg1, agr2"},
		Info:                 "info",
	}

}

func CreateTestService(organizationID string, appDescriptorId string, serviceGroupId string) entities.Service {

	serviceID := uuid.New().String()

	endpoints := []entities.Endpoint{
		{
			Type: entities.EndpointType(1),
			Path: "../../endpoint",
		},
	}

	ports := make([]entities.Port, 0)
	// port number should be greater than zero
	for i := 0; i < 5; i++ {
		ports = append(ports, entities.Port{Name: fmt.Sprintf("port%d", i), InternalPort: int32(i + 1), ExposedPort: int32(i + 1), Endpoints: endpoints})
	}

	confFile := make([]entities.ConfigFile, 0)
	confFile = append(confFile, CreateTestConfigFile(organizationID, appDescriptorId))

	return entities.Service{
		OrganizationId:  organizationID,
		AppDescriptorId: appDescriptorId,
		ServiceGroupId:  serviceGroupId,
		ServiceId:       serviceID,
		Name:            "service-name",
		Type:            entities.DockerService,
		Image:           "../image.txt",
		Credentials: &entities.ImageCredentials{
			Username:         "carmen",
			Password:         "*****",
			Email:            "cdelope@daisho.group",
			DockerRepository: "Docker_repo"},
		Specs: &entities.DeploySpecs{
			Cpu:      1239900,
			Memory:   2000,
			Replicas: 2,
		},
		Storage: []entities.Storage{
			{
				Size:      900,
				MountPath: "../../mount_path",
				Type:      entities.StorageType(1),
			},
		},
		ExposedPorts:         ports,
		EnvironmentVariables: map[string]string{"HOST": "HOST_VALUE", "PORT": "PORT_VALUE"},
		Configs:              confFile,
		Labels:               map[string]string{"eti1": "label1(Service)", "eti2": "label2(Service)"},
		DeployAfter:          []string{"deploy after this", "and this"},
		RunArguments:         []string{"arg1", "arg2", "arg3", "arg4"},
	}
}

func CreateTestServiceGroupInstance(organizationID string, appDescriptorID string, appInstanceId string) entities.ServiceGroupInstance {

	id := uuid.New().String()
	id2 := uuid.New().String()
	servicesInstances := make([]entities.ServiceInstance, 0)
	for i := 0; i < 1; i++ {
		servicesInstances = append(servicesInstances, CreateTestServiceInstance(organizationID, appDescriptorID, appInstanceId, id, id2))
	}

	return entities.ServiceGroupInstance{
		OrganizationId:         organizationID,
		AppDescriptorId:        appDescriptorID,
		AppInstanceId:          appInstanceId,
		ServiceGroupId:         id,
		ServiceGroupInstanceId: id2,
		Name:                   "Service group Instance name",
		Policy:                 entities.SameCluster,
		ServiceInstances:       servicesInstances,
		Status:                 entities.ServiceScheduled,
		Metadata: &entities.InstanceMetadata{
			OrganizationId:      organizationID,
			AppDescriptorId:     appDescriptorID,
			AppInstanceId:       appInstanceId,
			MonitoredInstanceId: uuid.New().String(),
			Type:                entities.ServiceGroupInstanceType,
			InstancesId:         []string{"inst1", "inst2", "inst3"},
			DesiredReplicas:     3,
			AvailableReplicas:   3,
			UnavailableReplicas: 0,
			Status:              map[string]entities.ServiceStatus{"status1": entities.ServiceError, "status2": entities.ServiceDeploying},
			Info:                map[string]string{"sgInfo1": "info1"},
		},
		Specs: &entities.ServiceGroupDeploymentSpecs{
			Replicas:            3,
			MultiClusterReplica: true,
		},
		Labels: map[string]string{"label1": "label1(servicegroupinstance)"},
	}
}

func CreateTestServiceGroup(organizationID string, appDescriptorId string) entities.ServiceGroup {

	serviceGroupID := uuid.New().String()

	services := make([]entities.Service, 0)
	for i := 0; i < 1; i++ {
		services = append(services, CreateTestService(organizationID, appDescriptorId, serviceGroupID))
	}

	specs := &entities.ServiceGroupDeploymentSpecs{
		Replicas:            5,
		MultiClusterReplica: false,
		DeploymentSelectors: map[string]string{"deploy1": "select1", "deploy2": "select2"},
	}

	return entities.ServiceGroup{
		OrganizationId:  organizationID,
		AppDescriptorId: appDescriptorId,
		ServiceGroupId:  serviceGroupID,
		Name:            "Service Group Test",
		Services:        services,
		Policy:          entities.CollocationPolicy(1),
		Specs:           specs,
		Labels:          map[string]string{"label1": "value1(ServiceGroup)"},
	}
}

func CreateTestRule(organizationID string, appDescriptorID string) entities.SecurityRule {
	id := uuid.New().String()

	rule := entities.SecurityRule{
		OrganizationId:         organizationID,
		AppDescriptorId:        appDescriptorID,
		RuleId:                 id,
		Name:                   "Rule name",
		TargetServiceGroupName: "target service group name",
		TargetServiceName:      "target service name",
		TargetPort:             80,
		Access:                 entities.AllAppServices,
		AuthServiceGroupName:   "auth service group name",
		AuthServices:           []string{"authService1", "authService2"},
		DeviceGroupNames:       []string{"deviceGroup1", "deviceGroup2"},
		DeviceGroupIds:         []string{"device001", "device002"},
		InboundNetInterface:    "inbound1",
		OutboundNetInterface:   "outbound1",
	}

	return rule
}

func CreateTestApplication(organizationID string, appDescriptorID string) *entities.AppInstance {

	id := uuid.New().String()

	rules := make([]entities.SecurityRule, 0)
	rules = append(rules, CreateTestRule(organizationID, appDescriptorID))

	groups := make([]entities.ServiceGroupInstance, 0)
	groups = append(groups, CreateTestServiceGroupInstance(organizationID, appDescriptorID, id))

	metadata := make([]entities.InstanceMetadata, 0)
	metadata = append(metadata, CreateTestInstanceMetadata(organizationID, appDescriptorID, id))

	app := entities.AppInstance{
		OrganizationId:        organizationID,
		AppDescriptorId:       appDescriptorID,
		AppInstanceId:         id,
		Name:                  "App instance Name",
		ConfigurationOptions:  map[string]string{"conf1": "value1(appInstance)", "conf2": "value2(appInstance)", "conf3": "value3(appInstance)"},
		EnvironmentVariables:  map[string]string{"env01": "value1(appInstance)", "env02": "value2(appInstance)"},
		Labels:                map[string]string{"label1": "value1(appInstance)", "label2": "value2(appInstance)", "label3": "value3(appInstance)", "label4": "value4(appInstance)"},
		Rules:                 rules,
		Groups:                groups,
		Status:                entities.Queued,
		Metadata:              metadata,
		InboundNetInterfaces:  []entities.InboundNetworkInterface{{Name: "inbound1"}, {Name: "inbound2"}},
		OutboundNetInterfaces: []entities.OutboundNetworkInterface{{Name: "inbound1", Required: true}, {Name: "inbound2", Required: false}},
	}

	return &app
}

func CreateParametrizedDescriptor(organizationID string) *entities.ParametrizedDescriptor {
	rules := make([]entities.SecurityRule, 0)
	rules = append(rules, CreateTestRule(organizationID, uuid.New().String()))

	groups := make([]entities.ServiceGroup, 0)
	groups = append(groups, CreateTestServiceGroup(organizationID, uuid.New().String()))

	descriptor := entities.ParametrizedDescriptor{
		OrganizationId:       organizationID,
		AppDescriptorId:      uuid.New().String(),
		AppInstanceId:        uuid.New().String(),
		Name:                 "App descriptor Test",
		ConfigurationOptions: map[string]string{"conf1": "value1", "conf2": "value2"},
		EnvironmentVariables: map[string]string{"env1": "value1", "env2": "value2"},
		Labels:               map[string]string{"label1": "value1", "label2": "value2", "label3": "value3"},
		Rules:                rules,
		Groups:               groups,
	}

	return &descriptor
}

func CreateTestApplicationDescriptor(organizationID string) *entities.AppDescriptor {

	id := uuid.New().String()

	rules := make([]entities.SecurityRule, 0)
	rules = append(rules, CreateTestRule(organizationID, id))

	groups := make([]entities.ServiceGroup, 0)
	groups = append(groups, CreateTestServiceGroup(organizationID, id))

	descriptor := entities.AppDescriptor{
		OrganizationId:       organizationID,
		AppDescriptorId:      id,
		Name:                 "App descriptor Test",
		ConfigurationOptions: map[string]string{"conf1": "value1", "conf2": "value2"},
		EnvironmentVariables: map[string]string{"env1": "value1", "env2": "value2"},
		Labels:               map[string]string{"label1": "value1", "label2": "value2", "label3": "value3"},
		Rules:                rules,
		Groups:               groups,
		Parameters: []entities.Parameter{{
			Name:         "Param1",
			Description:  "param1 descriptor",
			Path:         "xpath",
			Type:         entities.String,
			DefaultValue: "default",
			Category:     entities.Advanced,
			Required:     true},
		},
		InboundNetInterfaces:  []entities.InboundNetworkInterface{{Name: "inbound1"}, {Name: "inbound2"}},
		OutboundNetInterfaces: []entities.OutboundNetworkInterface{{Name: "inbound1", Required: true}, {Name: "inbound2", Required: false}},
	}

	return &descriptor

}

func CreateApplicationDescriptorWithParameters(organizationID string) *entities.AppDescriptor {
	id := uuid.New().String()

	rules := make([]entities.SecurityRule, 0)
	rules = append(rules, CreateTestRule(organizationID, id))

	groups := make([]entities.ServiceGroup, 0)
	groups = append(groups, CreateTestServiceGroup(organizationID, id))

	descriptor := entities.AppDescriptor{
		OrganizationId:       organizationID,
		AppDescriptorId:      id,
		Name:                 "Test-Descriptor",
		ConfigurationOptions: map[string]string{"conf1": "value1", "conf2": "value2"},
		EnvironmentVariables: map[string]string{"env1": "value1", "env2": "value2"},
		Labels:               map[string]string{"label1": "value1", "label2": "value2", "label3": "value3"},
		Rules:                rules,
		Groups:               groups,
		Parameters: []entities.Parameter{
			{
				Name:         "param_name1",
				Description:  "param_name1 description",
				Path:         "path1",
				Type:         entities.Boolean,
				DefaultValue: "true",
				Category:     entities.Basic,
				Required:     true,
			}, {
				Name:         "param_name2",
				Description:  "param_name2 description",
				Path:         "path2",
				Type:         entities.Enum,
				DefaultValue: "ENUM1",
				Category:     entities.Basic,
				EnumValues:   []string{"ENUM1, ENUM2"},
				Required:     false,
			},
		},
	}

	return &descriptor
}

func CreateAppEndPoint() *entities.AppEndpoint {
	return &entities.AppEndpoint{
		OrganizationId:         uuid.New().String(),
		AppInstanceId:          uuid.New().String(),
		ServiceGroupInstanceId: uuid.New().String(),
		ServiceInstanceId:      uuid.New().String(),
		Port:                   8080,
		Protocol:               entities.HTTP,
		EndpointInstanceId:     uuid.New().String(),
		Type:                   entities.IsAlive,
		Fqdn:                   "fqdn.domain.es",
		GlobalFqdn:             fmt.Sprintf("%d.globaldomain.es", rand.Int()),
	}
}
