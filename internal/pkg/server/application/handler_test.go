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

// TODO: Refactor this file at least into two (one for instances, another for descriptors)
package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/nalej/system-model/internal/pkg/entities"
	appProvider "github.com/nalej/system-model/internal/pkg/provider/application"
	devProvider "github.com/nalej/system-model/internal/pkg/provider/device"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/nalej/system-model/internal/pkg/server/testhelpers"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"math/rand"
	"strings"
)

func generateRandomSpecs() *grpc_application_go.DeploySpecs {
	return &grpc_application_go.DeploySpecs{
		Cpu:      int64(rand.Intn(100)),
		Memory:   int64(rand.Intn(10000)),
		Replicas: int32(rand.Intn(10)),
	}
}

func generateRandomService(index int) *grpc_application_go.Service {
	credentials := &grpc_application_go.ImageCredentials{
		Username:         "username",
		Password:         "****",
		Email:            "email@company.com",
		DockerRepository: "repo",
	}
	endpoints := make([]*grpc_application_go.Endpoint, 0)
	endpoints = append(endpoints, &grpc_application_go.Endpoint{
		Type: grpc_application_go.EndpointType_REST,
		Path: "/",
	})
	ports := make([]*grpc_application_go.Port, 0)
	ports = append(ports, &grpc_application_go.Port{
		Name:         "simple-endpoint",
		InternalPort: 80,
		ExposedPort:  80,
		Endpoints:    endpoints,
	})

	storage := make([]*grpc_application_go.Storage, 0)
	storage = append(storage, &grpc_application_go.Storage{
		Size:      12345,
		MountPath: "../path/",
		Type:      grpc_application_go.StorageType_CLUSTER_LOCAL,
	})
	configs := make([]*grpc_application_go.ConfigFile, 0)
	configs = append(configs, &grpc_application_go.ConfigFile{
		Name:      "Config file name",
		Content:   []byte{0x00, 0x01, 0x02},
		MountPath: "./path..",
	})

	return &grpc_application_go.Service{
		Name:                 fmt.Sprintf("service-%d", index),
		Type:                 grpc_application_go.ServiceType_DOCKER,
		Image:                fmt.Sprintf("image:v%d", rand.Intn(10)),
		Specs:                generateRandomSpecs(),
		ExposedPorts:         ports,
		Credentials:          credentials,
		Storage:              storage,
		EnvironmentVariables: map[string]string{"env01": "env01Label", "env02": "env02Label"},
		DeployAfter:          []string{"after1", "after2"},
		Labels:               map[string]string{"label1": "service label 1", "label2": "service label 2"},
		Configs:              configs,
		RunArguments:         []string{"arg1", "arg2", "arg3"},
	}
}

func generateServiceGroup(services []*grpc_application_go.Service) *grpc_application_go.ServiceGroup {

	return &grpc_application_go.ServiceGroup{
		Name:     "Service Group",
		Services: services,
		Policy:   grpc_application_go.CollocationPolicy_SEPARATE_CLUSTERS,
		Specs: &grpc_application_go.ServiceGroupDeploymentSpecs{
			Replicas:            5,
			MultiClusterReplica: false,
		},
		Labels: map[string]string{"label1": "sg_label1", "label2": "sg_label2", "label3": "sg_label3"},
	}
}

func generateAddAppDescriptor(orgID string, numServices int) *grpc_application_go.AddAppDescriptorRequest {
	services := make([]*grpc_application_go.Service, 0)
	for i := 0; i < numServices; i++ {
		services = append(services, generateRandomService(i))
	}
	securityRules := make([]*grpc_application_go.SecurityRule, 0)
	for i := 0; i < (numServices); i++ {
		securityRules = append(securityRules, &grpc_application_go.SecurityRule{
			RuleId:                 fmt.Sprintf("r%d", i),
			Name:                   fmt.Sprintf("%d -> %d", i, i+1),
			TargetServiceGroupName: fmt.Sprintf("targetServiceGroupName-%d", i),
			TargetServiceName:      fmt.Sprintf("targetServiceName-%d", i),
			TargetPort:             80,
			Access:                 grpc_application_go.PortAccess_APP_SERVICES,
			AuthServiceGroupName:   fmt.Sprintf("AuthServiceGroupName-%d", i),
			AuthServices:           []string{fmt.Sprintf("s%d", i+1)},
			DeviceGroupNames:       []string{"dg1", "dg2"},
		})
	}
	// update the deviceGroupsNames to be different
	if len(securityRules) > 0 {
		securityRules[0].DeviceGroupNames = []string{"dg3"}
	}

	groups := make([]*grpc_application_go.ServiceGroup, 0)
	groups = append(groups, generateServiceGroup(services))

	return &grpc_application_go.AddAppDescriptorRequest{
		RequestId:            "request_id",
		OrganizationId:       orgID,
		Name:                 "new app",
		ConfigurationOptions: map[string]string{"conf1": "conf1", "conf2": "conf2"},
		EnvironmentVariables: map[string]string{"var1": "env1"},
		Labels:               map[string]string{"label1": "eti1"},
		Rules:                securityRules,
		Groups:               groups,
	}
}

func generateParametrizedDescriptor(descriptor *grpc_application_go.AppDescriptor, instanceID string) *grpc_application_go.ParametrizedDescriptor {
	return &grpc_application_go.ParametrizedDescriptor{
		OrganizationId:       descriptor.OrganizationId,
		AppDescriptorId:      descriptor.AppDescriptorId,
		AppInstanceId:        instanceID,
		Name:                 descriptor.Name,
		ConfigurationOptions: descriptor.ConfigurationOptions,
		EnvironmentVariables: descriptor.EnvironmentVariables,
		Labels:               descriptor.Labels,
		Rules:                descriptor.Rules,
		Groups:               descriptor.Groups,
	}
}

func generateAddAppInstance(organizationID string, appDescriptorID string) *grpc_application_go.AddAppInstanceRequest {
	return &grpc_application_go.AddAppInstanceRequest{
		OrganizationId:  organizationID,
		AppDescriptorId: appDescriptorID,
		Name:            fmt.Sprintf("app instance %d", rand.Int31n(100)),
	}
}

func generateUpdateAppInstance(organizationID string, appInstanceID string,
	status grpc_application_go.ApplicationStatus) *grpc_application_go.UpdateAppStatusRequest {
	return &grpc_application_go.UpdateAppStatusRequest{
		OrganizationId: organizationID,
		AppInstanceId:  appInstanceID,
		Status:         status,
	}
}

func generateUpdateServiceStatus(organizationID string, appInstanceID string, serviceID string,
	appDescriptorID string, status grpc_application_go.ServiceStatus) *grpc_application_go.UpdateServiceStatusRequest {
	endpoint := make([]string, 0)
	endpoint = append(endpoint, "endpoint1")
	return &grpc_application_go.UpdateServiceStatusRequest{
		OrganizationId: organizationID,
		AppInstanceId:  appInstanceID,
		Status:         status,
		//Endpoints: endpoint,
		DeployedOnClusterId: fmt.Sprintf("Deploy on cluster - %d", rand.Int31n(100)),
	}
}

func InjectBadServiceName(descriptor *grpc_application_go.AddAppDescriptorRequest) {
	for g, group := range descriptor.Groups {
		for s, service := range group.Services {
			descriptor.Groups[g].Services[s].Name = fmt.Sprintf("%s #*", service.Name)
		}
	}
}

func InjectBadPortName(descriptor *grpc_application_go.AddAppDescriptorRequest) {
	for g, group := range descriptor.Groups {
		for s, service := range group.Services {
			for p, port := range service.ExposedPorts {
				descriptor.Groups[g].Services[s].ExposedPorts[p].Name = fmt.Sprintf("%s12345678912345678", port.Name)
			}
		}
	}
}

func InjectBadPortNumber(descriptor *grpc_application_go.AddAppDescriptorRequest) {
	for g, group := range descriptor.Groups {
		for s, service := range group.Services {
			for p, port := range service.ExposedPorts {
				descriptor.Groups[g].Services[s].ExposedPorts[p].ExposedPort = port.ExposedPort + 65536
				descriptor.Groups[g].Services[s].ExposedPorts[p].InternalPort = port.InternalPort + 65536
			}
		}
	}
}

func generateServiceGroupInstanceMetadata(appInstance grpc_application_go.AppInstance) *grpc_application_go.InstanceMetadata {
	return &grpc_application_go.InstanceMetadata{
		AvailableReplicas:   1,
		UnavailableReplicas: 0,
		DesiredReplicas:     1,
		AppInstanceId:       appInstance.AppInstanceId,
		InstancesId:         []string{"appMonitored001"},
		ServiceGroupId:      appInstance.Groups[0].ServiceGroupId,
		AppDescriptorId:     appInstance.AppDescriptorId,
		Info:                map[string]string{"appMonitored001": "info"},
		Type:                grpc_application_go.InstanceType_SERVICE_GROUP_INSTANCE,
		OrganizationId:      appInstance.OrganizationId,
		// MonitoredInstanceId: --> to be filled by the system model after addition
		Status: map[string]grpc_application_go.ServiceStatus{
			"service1": grpc_application_go.ServiceStatus_SERVICE_DEPLOYING,
		},
	}
}

func generateAppEndpoint(serviceName string, organizationId string) *grpc_application_go.AddAppEndpointRequest {
	appEndpoint := &grpc_application_go.AddAppEndpointRequest{
		OrganizationId:         organizationId,
		AppInstanceId:          uuid.New().String(),
		ServiceGroupInstanceId: uuid.New().String(),
		ServiceInstanceId:      uuid.New().String(),
		Protocol:               grpc_application_go.AppEndpointProtocol_HTTPS,
		ServiceName:            serviceName,
		EndpointInstance: &grpc_application_go.EndpointInstance{
			EndpointInstanceId: uuid.New().String(),
			Type:               grpc_application_go.EndpointType_IS_ALIVE,
			Port:               8080,
		},
	}

	appEndpoint.EndpointInstance.Fqdn = fmt.Sprintf("%s.%s.%s.domain", serviceName, appEndpoint.ServiceGroupInstanceId[:6],
		appEndpoint.AppInstanceId[:6])

	return appEndpoint

}

var _ = ginkgo.Describe("Applications", func() {

	const numServices = 2

	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client grpc_application_go.ApplicationsClient

	// Target organization.
	var targetOrganization *entities.Organization
	//var targetDeviceGroup * device.DeviceGroup

	var targetDescriptor *grpc_application_go.AppDescriptor
	var targetInstance *grpc_application_go.AppInstance

	// Organization Provider
	var organizationProvider orgProvider.Provider
	var applicationProvider appProvider.Provider
	var deviceProvider devProvider.Provider

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		// Create providers
		organizationProvider = orgProvider.NewMockupOrganizationProvider()
		applicationProvider = appProvider.NewMockupApplicationProvider()
		deviceProvider = devProvider.NewMockupDeviceProvider()

		manager := NewManager(organizationProvider, applicationProvider, deviceProvider, "nalej.cluster.local")
		handler := NewHandler(manager)
		grpc_application_go.RegisterApplicationsServer(server, handler)

		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_application_go.NewApplicationsClient(conn)
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
	})

	ginkgo.BeforeEach(func() {
		ginkgo.By("cleaning the mockups", func() {
			organizationProvider.(*orgProvider.MockupOrganizationProvider).Clear()
			applicationProvider.(*appProvider.MockupApplicationProvider).Clear()
			deviceProvider.(*devProvider.MockupDeviceProvider).Clear()

			// Initial data
			targetOrganization = testhelpers.CreateOrganization(organizationProvider)

			// generate deviceGroups
			testhelpers.CreateDeviceGroup(deviceProvider, targetOrganization.ID, "dg1")
			testhelpers.CreateDeviceGroup(deviceProvider, targetOrganization.ID, "dg2")
			testhelpers.CreateDeviceGroup(deviceProvider, targetOrganization.ID, "dg3")

		})
	})

	ginkgo.Context("Application descriptor", func() {
		ginkgo.Context("adding application descriptors", func() {
			ginkgo.It("should add an application descriptor", func() {
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())
				gomega.Expect(app.AppDescriptorId).ShouldNot(gomega.BeNil())
				gomega.Expect(app.Name).Should(gomega.Equal(toAdd.Name))
				gomega.Expect(len(toAdd.Groups)).Should(gomega.Equal(len(app.Groups)))
			})
			ginkgo.It("should fail on an empty request", func() {
				toAdd := &grpc_application_go.AddAppDescriptorRequest{}
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(app).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a non existing organization", func() {
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				toAdd.OrganizationId = "does not exists"
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(app).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a descriptor without services", func() {
				toAdd := generateAddAppDescriptor(targetOrganization.ID, 0)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(app).Should(gomega.BeNil())
			})

			// AddDescriptor with BadServiceName
			ginkgo.It("Should fail to add a descriptor with bad service name", func() {

				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				InjectBadServiceName(toAdd)
				_, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).NotTo(gomega.Succeed())
			})
			// AddDescriptor with Bad portname
			ginkgo.It("Should fail to add a descriptor with bad port name", func() {

				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				InjectBadPortName(toAdd)
				_, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).NotTo(gomega.Succeed())
			})
			// AddDescriptor with Bad portname
			ginkgo.It("Should fail to add a descriptor with bad port number", func() {

				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				InjectBadPortNumber(toAdd)
				_, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).NotTo(gomega.Succeed())
			})

		})
		ginkgo.Context("get application descriptor", func() {
			ginkgo.It("should get an existing app descriptor", func() {
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())
				retrieved, err := client.GetAppDescriptor(context.Background(), &grpc_application_go.AppDescriptorId{
					OrganizationId:  app.OrganizationId,
					AppDescriptorId: app.AppDescriptorId,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(retrieved.Name).Should(gomega.Equal(app.Name))
			})
			ginkgo.It("should fail on a non existing application", func() {
				retrieved, err := client.GetAppDescriptor(context.Background(), &grpc_application_go.AppDescriptorId{
					OrganizationId:  targetOrganization.ID,
					AppDescriptorId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a non existing organization", func() {
				retrieved, err := client.GetAppDescriptor(context.Background(), &grpc_application_go.AppDescriptorId{
					OrganizationId:  "does not exists",
					AppDescriptorId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
		})
		ginkgo.Context("listing application descriptors", func() {
			ginkgo.It("should list apps on an existing organization", func() {
				numDescriptors := 3
				for i := 0; i < numDescriptors; i++ {
					toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
					app, err := client.AddAppDescriptor(context.Background(), toAdd)
					gomega.Expect(err).Should(gomega.Succeed())
					gomega.Expect(app).ShouldNot(gomega.BeNil())
				}
				retrieved, err := client.ListAppDescriptors(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: targetOrganization.ID,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(len(retrieved.Descriptors)).Should(gomega.Equal(numDescriptors))
			})
			ginkgo.It("should fail on a non existing organization", func() {
				retrieved, err := client.ListAppDescriptors(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
			ginkgo.It("should work on an organization without descriptors", func() {
				gomega.Expect(organizationProvider).ShouldNot(gomega.BeNil())
				retrieved, err := client.ListAppDescriptors(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: targetOrganization.ID,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(len(retrieved.Descriptors)).Should(gomega.Equal(0))
			})
		})
		ginkgo.Context("removing application descriptors", func() {
			ginkgo.It("should be able to remove an existing descriptor", func() {
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())

				toRemove := &grpc_application_go.AppDescriptorId{
					OrganizationId:  app.OrganizationId,
					AppDescriptorId: app.AppDescriptorId,
				}
				success, err := client.RemoveAppDescriptor(context.Background(), toRemove)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(success).ShouldNot(gomega.BeNil())
			})
			ginkgo.It("should fail if the organization does not exists", func() {
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())
				toRemove := &grpc_application_go.AppDescriptorId{
					OrganizationId:  "unknown",
					AppDescriptorId: app.AppDescriptorId,
				}
				success, err := client.RemoveAppDescriptor(context.Background(), toRemove)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(success).Should(gomega.BeNil())
			})
			ginkgo.It("should fail if the descriptor does not exits", func() {
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())
				toRemove := &grpc_application_go.AppDescriptorId{
					OrganizationId:  app.OrganizationId,
					AppDescriptorId: "unknown",
				}
				success, err := client.RemoveAppDescriptor(context.Background(), toRemove)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(success).Should(gomega.BeNil())
			})
		})
		ginkgo.Context("Getting the params of a descriptor", func() {
			ginkgo.It("should be able to retrieve the parameters of a descriptor", func() {
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)

				toAdd.Parameters = []*grpc_application_go.AppParameter{
					{
						Name:         "Param1",
						Description:  "Param1 Descriptor",
						Path:         "Path1",
						Type:         grpc_application_go.ParamDataType_STRING,
						DefaultValue: "default",
						Category:     grpc_application_go.ParamCategory_ADVANCED,
					},
					{
						Name:         "Param2",
						Description:  "Param2 Descriptor",
						Path:         "Path2",
						Type:         grpc_application_go.ParamDataType_BOOLEAN,
						DefaultValue: "true",
						Category:     grpc_application_go.ParamCategory_BASIC,
					},
				}

				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())

				params, err := client.GetDescriptorAppParameters(context.Background(), &grpc_application_go.AppDescriptorId{
					OrganizationId:  app.OrganizationId,
					AppDescriptorId: app.AppDescriptorId,
				})
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(params).NotTo(gomega.BeNil())
				gomega.Expect(len(params.Parameters)).Should(gomega.Equal(2))

			})
		})
	})

	ginkgo.Context("Application instance", func() {
		ginkgo.BeforeEach(func() {
			ginkgo.By("creating required descriptor", func() {
				// Initial data
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())
				targetDescriptor = app
			})
		})
		ginkgo.Context("adding application instance", func() {
			ginkgo.It("should add an app instance", func() {
				toAdd := generateAddAppInstance(targetDescriptor.OrganizationId, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())
				gomega.Expect(added.AppInstanceId).ShouldNot(gomega.BeEmpty())
				gomega.Expect(added.OrganizationId).Should(gomega.Equal(targetDescriptor.OrganizationId))
				gomega.Expect(added.AppDescriptorId).Should(gomega.Equal(targetDescriptor.AppDescriptorId))
			})
			ginkgo.It("should fail on a non existing app descriptor", func() {
				toAdd := generateAddAppInstance(targetDescriptor.OrganizationId, "does not exists")
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(added).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a non existing organization", func() {
				toAdd := generateAddAppInstance("does not exists", targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(added).Should(gomega.BeNil())
			})
		})
		ginkgo.Context("get application instance", func() {
			ginkgo.It("should retrieve an existing app", func() {
				toAdd := generateAddAppInstance(targetOrganization.ID, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())
				gomega.Expect(added.AppInstanceId).ShouldNot(gomega.BeEmpty())
				retrieved, err := client.GetAppInstance(context.Background(), &grpc_application_go.AppInstanceId{
					OrganizationId: added.OrganizationId,
					AppInstanceId:  added.AppInstanceId,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(retrieved.Name).Should(gomega.Equal(added.Name))
			})
			ginkgo.It("should fail on a non existing instance", func() {
				retrieved, err := client.GetAppInstance(context.Background(), &grpc_application_go.AppInstanceId{
					OrganizationId: targetDescriptor.OrganizationId,
					AppInstanceId:  "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a non existing organization", func() {
				retrieved, err := client.GetAppInstance(context.Background(), &grpc_application_go.AppInstanceId{
					OrganizationId: "does not exists",
					AppInstanceId:  "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
		})
		ginkgo.Context("listing application instances", func() {
			ginkgo.It("should retrieve instances on an existing organization", func() {
				numInstances := 3
				for i := 0; i < numInstances; i++ {
					toAdd := generateAddAppInstance(targetOrganization.ID, targetDescriptor.AppDescriptorId)
					added, err := client.AddAppInstance(context.Background(), toAdd)
					gomega.Expect(err).Should(gomega.Succeed())
					gomega.Expect(added).ShouldNot(gomega.BeNil())
				}
				retrieved, err := client.ListAppInstances(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: targetOrganization.ID,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(len(retrieved.Instances)).Should(gomega.Equal(numInstances))
			})
			ginkgo.It("should work on an organization without instances", func() {
				retrieved, err := client.ListAppInstances(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: targetOrganization.ID,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(len(retrieved.Instances)).Should(gomega.Equal(0))
			})
			ginkgo.It("should fail on a non existing organization", func() {
				retrieved, err := client.ListAppInstances(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
		})
		ginkgo.Context("update application instance", func() {
			ginkgo.PIt("should update instance and return the new values", func() {
			})
		})

		ginkgo.Context("update service status in application instance", func() {
			ginkgo.It("should update instance and return the new values with the global Fqdn", func() {
				toAdd := generateAddAppInstance(targetOrganization.ID, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())

				parametrized := generateParametrizedDescriptor(targetDescriptor, added.AppInstanceId)

				newDesc, err := client.AddParametrizedDescriptor(context.Background(), parametrized)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(newDesc).ShouldNot(gomega.BeNil())

				// add serviceGroupInstance
				list, err := client.AddServiceGroupInstances(context.Background(), &grpc_application_go.AddServiceGroupInstancesRequest{
					OrganizationId:  targetOrganization.ID,
					AppDescriptorId: targetDescriptor.AppDescriptorId,
					AppInstanceId:   added.AppInstanceId,
					ServiceGroupId:  targetDescriptor.Groups[0].ServiceGroupId,
					NumInstances:    1,
				})
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(list).NotTo(gomega.BeNil())

				// update status
				toUpdate := &grpc_application_go.UpdateServiceStatusRequest{
					OrganizationId:         targetOrganization.ID,
					AppInstanceId:          added.AppInstanceId,
					ServiceGroupInstanceId: list.ServiceGroupInstances[0].ServiceGroupInstanceId,
					ServiceInstanceId:      list.ServiceGroupInstances[0].ServiceInstances[0].ServiceInstanceId,
					Status:                 grpc_application_go.ServiceStatus_SERVICE_RUNNING,
					Endpoints: []*grpc_application_go.EndpointInstance{{
						EndpointInstanceId: uuid.New().String(),
						Type:               grpc_application_go.EndpointType_IS_ALIVE,
						Fqdn: fmt.Sprintf("%s.%s.%s.appcluster.nalej.com",
							list.ServiceGroupInstances[0].ServiceInstances[0].Name, list.ServiceGroupInstances[0].ServiceGroupInstanceId, list.ServiceGroupInstances[0].ServiceInstances[0].ServiceInstanceId),
					},
					},
				}
				success, err := client.UpdateServiceStatus(context.Background(), toUpdate)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(success).NotTo(gomega.BeNil())

				// add AppEndpoint
				success, err = client.AddAppEndpoint(context.Background(), &grpc_application_go.AddAppEndpointRequest{
					OrganizationId:         targetOrganization.ID,
					AppInstanceId:          added.AppInstanceId,
					ServiceGroupInstanceId: list.ServiceGroupInstances[0].ServiceGroupInstanceId,
					ServiceInstanceId:      list.ServiceGroupInstances[0].ServiceInstances[0].ServiceInstanceId,
					EndpointInstance:       toUpdate.Endpoints[0],
					ServiceName:            "service_name",
				})
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(success).NotTo(gomega.BeNil())

				// get instance
				instance, err := client.GetAppInstance(context.Background(), &grpc_application_go.AppInstanceId{
					OrganizationId: targetOrganization.ID,
					AppInstanceId:  added.AppInstanceId,
				})
				gomega.Expect(instance).NotTo(gomega.BeNil())
				gomega.Expect(err).To(gomega.BeNil())
			})
			ginkgo.It("should update instance and return the new values with the global Fqdn (IP address)", func() {
				toAdd := generateAddAppInstance(targetOrganization.ID, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())

				parametrized := generateParametrizedDescriptor(targetDescriptor, added.AppInstanceId)

				newDesc, err := client.AddParametrizedDescriptor(context.Background(), parametrized)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(newDesc).ShouldNot(gomega.BeNil())

				// add serviceGroupInstance
				list, err := client.AddServiceGroupInstances(context.Background(), &grpc_application_go.AddServiceGroupInstancesRequest{
					OrganizationId:  targetOrganization.ID,
					AppDescriptorId: targetDescriptor.AppDescriptorId,
					AppInstanceId:   added.AppInstanceId,
					ServiceGroupId:  targetDescriptor.Groups[0].ServiceGroupId,
					NumInstances:    1,
				})
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(list).NotTo(gomega.BeNil())

				// update status
				toUpdate := &grpc_application_go.UpdateServiceStatusRequest{
					OrganizationId:         targetOrganization.ID,
					AppInstanceId:          added.AppInstanceId,
					ServiceGroupInstanceId: list.ServiceGroupInstances[0].ServiceGroupInstanceId,
					ServiceInstanceId:      list.ServiceGroupInstances[0].ServiceInstances[0].ServiceInstanceId,
					Status:                 grpc_application_go.ServiceStatus_SERVICE_RUNNING,
					Endpoints: []*grpc_application_go.EndpointInstance{{
						EndpointInstanceId: uuid.New().String(),
						Type:               grpc_application_go.EndpointType_IS_ALIVE,
						Fqdn:               "100.100.200.1:900",
					},
					},
				}
				success, err := client.UpdateServiceStatus(context.Background(), toUpdate)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(success).NotTo(gomega.BeNil())

				// add AppEndpoint
				success, err = client.AddAppEndpoint(context.Background(), &grpc_application_go.AddAppEndpointRequest{
					OrganizationId:         targetOrganization.ID,
					AppInstanceId:          added.AppInstanceId,
					ServiceGroupInstanceId: list.ServiceGroupInstances[0].ServiceGroupInstanceId,
					ServiceInstanceId:      list.ServiceGroupInstances[0].ServiceInstances[0].ServiceInstanceId,
					EndpointInstance:       toUpdate.Endpoints[0],
					ServiceName:            "service_name",
				})
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(success).NotTo(gomega.BeNil())

				// get instance
				instance, err := client.GetAppInstance(context.Background(), &grpc_application_go.AppInstanceId{
					OrganizationId: targetOrganization.ID,
					AppInstanceId:  added.AppInstanceId,
				})
				gomega.Expect(instance).NotTo(gomega.BeNil())
				gomega.Expect(err).To(gomega.BeNil())
			})
			ginkgo.It("should update instance and list the new values with the global Fqdn", func() {
				toAdd := generateAddAppInstance(targetOrganization.ID, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())

				parametrized := generateParametrizedDescriptor(targetDescriptor, added.AppInstanceId)

				newDesc, err := client.AddParametrizedDescriptor(context.Background(), parametrized)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(newDesc).ShouldNot(gomega.BeNil())

				// add serviceGroupInstance
				list, err := client.AddServiceGroupInstances(context.Background(), &grpc_application_go.AddServiceGroupInstancesRequest{
					OrganizationId:  targetOrganization.ID,
					AppDescriptorId: targetDescriptor.AppDescriptorId,
					AppInstanceId:   added.AppInstanceId,
					ServiceGroupId:  targetDescriptor.Groups[0].ServiceGroupId,
					NumInstances:    1,
				})
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(list).NotTo(gomega.BeNil())

				// update status
				toUpdate := &grpc_application_go.UpdateServiceStatusRequest{
					OrganizationId:         targetOrganization.ID,
					AppInstanceId:          added.AppInstanceId,
					ServiceGroupInstanceId: list.ServiceGroupInstances[0].ServiceGroupInstanceId,
					ServiceInstanceId:      list.ServiceGroupInstances[0].ServiceInstances[0].ServiceInstanceId,
					Status:                 grpc_application_go.ServiceStatus_SERVICE_RUNNING,
					Endpoints: []*grpc_application_go.EndpointInstance{{
						EndpointInstanceId: uuid.New().String(),
						Type:               grpc_application_go.EndpointType_IS_ALIVE,
						Fqdn: fmt.Sprintf("%s.%s.%s.appcluster.nalej.com",
							list.ServiceGroupInstances[0].ServiceInstances[0].Name, list.ServiceGroupInstances[0].ServiceGroupInstanceId, list.ServiceGroupInstances[0].ServiceInstances[0].ServiceInstanceId),
					},
					},
				}
				success, err := client.UpdateServiceStatus(context.Background(), toUpdate)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(success).NotTo(gomega.BeNil())

				// add AppEndpoint
				success, err = client.AddAppEndpoint(context.Background(), &grpc_application_go.AddAppEndpointRequest{
					OrganizationId:         targetOrganization.ID,
					AppInstanceId:          added.AppInstanceId,
					ServiceGroupInstanceId: list.ServiceGroupInstances[0].ServiceGroupInstanceId,
					ServiceInstanceId:      list.ServiceGroupInstances[0].ServiceInstances[0].ServiceInstanceId,
					EndpointInstance:       toUpdate.Endpoints[0],
					ServiceName:            "service_name",
				})
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(success).NotTo(gomega.BeNil())

				// get instance
				instance, err := client.ListAppInstances(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: targetOrganization.ID,
				})
				gomega.Expect(instance).NotTo(gomega.BeNil())
				gomega.Expect(err).To(gomega.BeNil())
			})
		})

		ginkgo.Context("removing application instances", func() {
			ginkgo.It("should be able to remove an existing instance", func() {
				toAdd := generateAddAppInstance(targetOrganization.ID, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())
				toRemove := &grpc_application_go.AppInstanceId{
					OrganizationId: added.OrganizationId,
					AppInstanceId:  added.AppInstanceId,
				}
				success, err := client.RemoveAppInstance(context.Background(), toRemove)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(success).ShouldNot(gomega.BeNil())
			})
			ginkgo.It("should fail if the organization does not exists", func() {
				toAdd := generateAddAppInstance(targetOrganization.ID, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())
				toRemove := &grpc_application_go.AppInstanceId{
					OrganizationId: "unknown",
					AppInstanceId:  added.AppInstanceId,
				}
				success, err := client.RemoveAppInstance(context.Background(), toRemove)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(success).Should(gomega.BeNil())
			})
			ginkgo.It("should fail if the descriptor does not exits", func() {
				toAdd := generateAddAppInstance(targetOrganization.ID, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())
				toRemove := &grpc_application_go.AppInstanceId{
					OrganizationId: added.OrganizationId,
					AppInstanceId:  "unknown",
				}
				success, err := client.RemoveAppInstance(context.Background(), toRemove)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(success).Should(gomega.BeNil())
			})
			ginkgo.It("should be able to remove an existing instance and its params", func() {

				toAdd := generateAddAppInstance(targetOrganization.ID, targetDescriptor.AppDescriptorId)
				list := []*grpc_application_go.InstanceParameter{{ParameterName: "name1", Value: "value1"}, {ParameterName: "name2", Value: "value2"}}
				toAdd.Parameters = &grpc_application_go.InstanceParameterList{Parameters: list}
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())
				toRemove := &grpc_application_go.AppInstanceId{
					OrganizationId: added.OrganizationId,
					AppInstanceId:  added.AppInstanceId,
				}
				success, err := client.RemoveAppInstance(context.Background(), toRemove)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(success).ShouldNot(gomega.BeNil())

				// check the instance paramas have been removed
				params, err := applicationProvider.GetInstanceParameters(added.AppInstanceId)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(params).ShouldNot(gomega.BeNil())
				gomega.Expect(len(params)).Should(gomega.Equal(0))

			})

		})

		ginkgo.Context("Adding ServiceGroupInstance ", func() {
			ginkgo.It("should be able to add a service group instance", func() {
				toAdd := generateAddAppInstance(targetDescriptor.OrganizationId, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())

				parametrized := generateParametrizedDescriptor(targetDescriptor, added.AppInstanceId)

				newDesc, err := client.AddParametrizedDescriptor(context.Background(), parametrized)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(newDesc).ShouldNot(gomega.BeNil())

				sgToAdd := &grpc_application_go.AddServiceGroupInstancesRequest{
					OrganizationId:  targetDescriptor.OrganizationId,
					AppDescriptorId: targetDescriptor.AppDescriptorId,
					AppInstanceId:   added.AppInstanceId,
					ServiceGroupId:  targetDescriptor.Groups[0].ServiceGroupId,
					NumInstances:    1,
				}

				sgReceived, err := client.AddServiceGroupInstances(context.Background(), sgToAdd)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(sgReceived.ServiceGroupInstances[0].ServiceGroupId).Should(gomega.Equal(sgToAdd.ServiceGroupId))
			})
			ginkgo.It("should not be able to add a service group instance of a non existing group", func() {
				toAdd := generateAddAppInstance(targetDescriptor.OrganizationId, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())

				sgToAdd := &grpc_application_go.AddServiceGroupInstancesRequest{
					OrganizationId:  targetDescriptor.OrganizationId,
					AppDescriptorId: targetDescriptor.AppDescriptorId,
					AppInstanceId:   added.AppInstanceId,
					ServiceGroupId:  uuid.New().String(),
					NumInstances:    1,
				}

				_, err = client.AddServiceGroupInstances(context.Background(), sgToAdd)
				gomega.Expect(err).NotTo(gomega.Succeed())
			})

		})
		ginkgo.Context("Getting the params of an instance", func() {
			ginkgo.It("should be able to retrieve the parameters of an instance", func() {

				toAdd := generateAddAppInstance(targetDescriptor.OrganizationId, targetDescriptor.AppDescriptorId)
				list := []*grpc_application_go.InstanceParameter{{ParameterName: "name1", Value: "value1"}, {ParameterName: "name2", Value: "value2"}}
				toAdd.Parameters = &grpc_application_go.InstanceParameterList{Parameters: list}
				app, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())

				params, err := client.GetInstanceParameters(context.Background(), &grpc_application_go.AppInstanceId{
					OrganizationId: app.OrganizationId,
					AppInstanceId:  app.AppInstanceId,
				})
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(params).NotTo(gomega.BeNil())
				gomega.Expect(len(params.Parameters)).Should(gomega.Equal(2))

			})
			ginkgo.It("should be able to retrieve an empty list if the instance has no parameters", func() {

				toAdd := generateAddAppInstance(targetDescriptor.OrganizationId, targetDescriptor.AppDescriptorId)

				app, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())

				params, err := client.GetInstanceParameters(context.Background(), &grpc_application_go.AppInstanceId{
					OrganizationId: app.OrganizationId,
					AppInstanceId:  app.AppInstanceId,
				})
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(params).NotTo(gomega.BeNil())
				gomega.Expect(len(params.Parameters)).Should(gomega.Equal(0))

			})
		})

		// Service instances are
		ginkgo.Context("Adding ServiceInstance ", func() {
			ginkgo.It("should be able to add a service instance", func() {
				toAdd := generateAddAppInstance(targetDescriptor.OrganizationId, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())

				parametrized := generateParametrizedDescriptor(targetDescriptor, added.AppInstanceId)

				newDesc, err := client.AddParametrizedDescriptor(context.Background(), parametrized)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(newDesc).ShouldNot(gomega.BeNil())

				sgToAdd := &grpc_application_go.AddServiceGroupInstancesRequest{
					OrganizationId:  targetDescriptor.OrganizationId,
					AppDescriptorId: targetDescriptor.AppDescriptorId,
					AppInstanceId:   added.AppInstanceId,
					ServiceGroupId:  targetDescriptor.Groups[0].ServiceGroupId,
					NumInstances:    1,
				}

				sgReceived, err := client.AddServiceGroupInstances(context.Background(), sgToAdd)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(sgReceived.ServiceGroupInstances[0].ServiceGroupId).Should(gomega.Equal(sgToAdd.ServiceGroupId))

				sToAdd := &grpc_application_go.AddServiceInstanceRequest{
					OrganizationId:         targetDescriptor.OrganizationId,
					AppDescriptorId:        targetDescriptor.AppDescriptorId,
					AppInstanceId:          added.AppInstanceId,
					ServiceGroupId:         sgReceived.ServiceGroupInstances[0].ServiceGroupId,
					ServiceGroupInstanceId: sgReceived.ServiceGroupInstances[0].ServiceGroupInstanceId,
					ServiceId:              sgReceived.ServiceGroupInstances[0].ServiceInstances[0].ServiceId,
				}

				serviceInstance, err := client.AddServiceInstance(context.Background(), sToAdd)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(serviceInstance.ServiceId).Should(gomega.Equal(sToAdd.ServiceId))
				gomega.Expect(serviceInstance.ServiceInstanceId).NotTo(gomega.BeNil())

			})
			ginkgo.It("should not be able to add a service instance (service instance no exists)", func() {
				toAdd := generateAddAppInstance(targetDescriptor.OrganizationId, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())

				parametrized := generateParametrizedDescriptor(targetDescriptor, added.AppInstanceId)

				newDesc, err := client.AddParametrizedDescriptor(context.Background(), parametrized)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(newDesc).ShouldNot(gomega.BeNil())

				sgToAdd := &grpc_application_go.AddServiceGroupInstancesRequest{
					OrganizationId:  targetDescriptor.OrganizationId,
					AppDescriptorId: targetDescriptor.AppDescriptorId,
					AppInstanceId:   added.AppInstanceId,
					ServiceGroupId:  targetDescriptor.Groups[0].ServiceGroupId,
					NumInstances:    1,
				}

				sgReceived, err := client.AddServiceGroupInstances(context.Background(), sgToAdd)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(sgReceived.ServiceGroupInstances[0].ServiceGroupId).Should(gomega.Equal(sgToAdd.ServiceGroupId))

				sToAdd := &grpc_application_go.AddServiceInstanceRequest{
					OrganizationId:         targetDescriptor.OrganizationId,
					AppDescriptorId:        targetDescriptor.AppDescriptorId,
					AppInstanceId:          added.AppInstanceId,
					ServiceGroupId:         sgReceived.ServiceGroupInstances[0].ServiceGroupId,
					ServiceGroupInstanceId: sgReceived.ServiceGroupInstances[0].ServiceGroupInstanceId,
					ServiceId:              uuid.New().String(),
				}

				_, err = client.AddServiceInstance(context.Background(), sToAdd)
				gomega.Expect(err).NotTo(gomega.Succeed())

			})

		})

	})

	ginkgo.Context("App Endpoint", func() {
		ginkgo.It("should be able to add an app endpoint", func() {
			endPoint := generateAppEndpoint("serviceName", uuid.New().String())
			success, err := client.AddAppEndpoint(context.Background(), endPoint)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ShouldNot(gomega.BeNil())
		})
		ginkgo.It("should be able to get app endpoints list", func() {
			organizationID := uuid.New().String()

			endPoint := generateAppEndpoint("serviceName", organizationID)
			success, err := client.AddAppEndpoint(context.Background(), endPoint)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ShouldNot(gomega.BeNil())

			fqdnSplit := strings.Split(endPoint.EndpointInstance.Fqdn, ".")
			globalFqdn := fmt.Sprintf("%s-%d.%s.%s.%s.globaldomain.com", fqdnSplit[0], endPoint.EndpointInstance.Port,
				fqdnSplit[1], fqdnSplit[2], organizationID[:8])

			list, err := client.GetAppEndpoints(context.Background(), &grpc_application_go.GetAppEndPointRequest{
				Fqdn: globalFqdn,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(list).NotTo(gomega.BeNil())
			gomega.Expect(len(list.AppEndpoints)).Should(gomega.Equal(1))

		})

		ginkgo.It("Should not be able to app endpoints list (several organizations)", func() {
			endPoint1 := &grpc_application_go.AddAppEndpointRequest{
				OrganizationId:         "xxxxxxxxx1",
				AppInstanceId:          "aaaaaaaaa1",
				ServiceGroupInstanceId: "ggggggggg1",
				ServiceInstanceId:      "sssssssss1",
				ServiceName:            "service1",
				Protocol:               grpc_application_go.AppEndpointProtocol_HTTPS,
				EndpointInstance: &grpc_application_go.EndpointInstance{
					EndpointInstanceId: "1",
					Type:               grpc_application_go.EndpointType_IS_ALIVE,
					Fqdn:               "service.gggggg.aaaaaa.domain",
					Port:               80,
				},
			}
			success, err := client.AddAppEndpoint(context.Background(), endPoint1)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ShouldNot(gomega.BeNil())

			endPoint2 := &grpc_application_go.AddAppEndpointRequest{
				OrganizationId:         "xxxxxxxxx2",
				AppInstanceId:          "aaaaaaaaa2",
				ServiceGroupInstanceId: "ggggggggg2",
				ServiceInstanceId:      "sssssssss2",
				ServiceName:            "service2",
				Protocol:               grpc_application_go.AppEndpointProtocol_HTTPS,
				EndpointInstance: &grpc_application_go.EndpointInstance{
					EndpointInstanceId: "1",
					Type:               grpc_application_go.EndpointType_IS_ALIVE,
					Fqdn:               "service2.gggggg.aaaaaa.domain",
					Port:               80,
				},
			}
			success, err = client.AddAppEndpoint(context.Background(), endPoint2)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ShouldNot(gomega.BeNil())

			_, err = client.GetAppEndpoints(context.Background(), &grpc_application_go.GetAppEndPointRequest{
				Fqdn: "service.gggggg.aaaaaa.xxxxxxxx.globaldomain",
			})
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		ginkgo.It("should be able to remove an app endpoint", func() {
			endPoint := generateAppEndpoint("serviceName", uuid.New().String())
			success, err := client.AddAppEndpoint(context.Background(), endPoint)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ShouldNot(gomega.BeNil())

			success, err = client.RemoveAppEndpoints(context.Background(), &grpc_application_go.RemoveAppEndpointRequest{
				OrganizationId: endPoint.OrganizationId,
				AppInstanceId:  endPoint.AppInstanceId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ShouldNot(gomega.BeNil())
		})
	})

	ginkgo.Context("app zt network", func() {

		ginkgo.It("should be able to add a new zt network", func() {

			appNetwork := entities.AppZtNetwork{
				OrganizationId: "org001",
				AppInstanceId:  "app001",
				ZtNetworkId:    "ztnetwork001",
			}

			request := grpc_application_go.AddAppZtNetworkRequest{
				OrganizationId: appNetwork.OrganizationId,
				AppInstanceId:  appNetwork.AppInstanceId,
				NetworkId:      appNetwork.ZtNetworkId,
			}

			success, err := client.AddAppZtNetwork(context.Background(), &request)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ShouldNot(gomega.BeNil())
		})
		ginkgo.It("should be able to remove an existing zt network", func() {
			// first add a network
			appNetwork := entities.AppZtNetwork{
				OrganizationId: "org001",
				AppInstanceId:  "appToRemove",
				ZtNetworkId:    "ztnetwork001",
			}

			request := grpc_application_go.AddAppZtNetworkRequest{
				OrganizationId: appNetwork.OrganizationId,
				AppInstanceId:  appNetwork.AppInstanceId,
				NetworkId:      appNetwork.ZtNetworkId,
			}

			success, err := client.AddAppZtNetwork(context.Background(), &request)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ShouldNot(gomega.BeNil())

			// then remove it
			removeRequest := grpc_application_go.RemoveAppZtNetworkRequest{
				OrganizationId: appNetwork.OrganizationId,
				AppInstanceId:  appNetwork.AppInstanceId,
			}
			success, err = client.RemoveAppZtNetwork(context.Background(), &removeRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ShouldNot(gomega.BeNil())
		})

		ginkgo.It("should be able to add a ztMember", func() {
			ztMember, err := client.AddAuthorizedZtNetworkMember(context.Background(), &grpc_application_go.AddAuthorizedZtNetworkMemberRequest{
				OrganizationId:               uuid.New().String(),
				NetworkId:                    uuid.New().String(),
				MemberId:                     uuid.New().String(),
				AppInstanceId:                uuid.New().String(),
				ServiceGroupInstanceId:       uuid.New().String(),
				ServiceApplicationInstanceId: uuid.New().String(),
				IsProxy:                      false,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(ztMember).ShouldNot(gomega.BeNil())

		})
		ginkgo.It("should be able to add two ztMembers", func() {

			toAdd := &grpc_application_go.AddAuthorizedZtNetworkMemberRequest{
				OrganizationId:               uuid.New().String(),
				NetworkId:                    uuid.New().String(),
				MemberId:                     uuid.New().String(),
				AppInstanceId:                uuid.New().String(),
				ServiceGroupInstanceId:       uuid.New().String(),
				ServiceApplicationInstanceId: uuid.New().String(),
				IsProxy:                      false,
			}

			ztMember, err := client.AddAuthorizedZtNetworkMember(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(ztMember).ShouldNot(gomega.BeNil())

			toAdd.IsProxy = true
			toAdd.MemberId = uuid.New().String()
			ztMember, err = client.AddAuthorizedZtNetworkMember(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(ztMember).ShouldNot(gomega.BeNil())

		})
		ginkgo.It("should be able to list the members of a zt-network", func() {
			toAdd := &grpc_application_go.AddAuthorizedZtNetworkMemberRequest{
				OrganizationId:               uuid.New().String(),
				NetworkId:                    uuid.New().String(),
				MemberId:                     uuid.New().String(),
				AppInstanceId:                uuid.New().String(),
				ServiceGroupInstanceId:       uuid.New().String(),
				ServiceApplicationInstanceId: uuid.New().String(),
				IsProxy:                      false,
			}

			ztMember, err := client.AddAuthorizedZtNetworkMember(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(ztMember).ShouldNot(gomega.BeNil())

			toAdd.IsProxy = true
			toAdd.MemberId = uuid.New().String()
			ztMember, err = client.AddAuthorizedZtNetworkMember(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(ztMember).ShouldNot(gomega.BeNil())

			// new record (same networkId and different service_group_instance_id and service_application_instance_id

			toAdd.ServiceApplicationInstanceId = uuid.New().String()
			toAdd.ServiceGroupInstanceId = uuid.New().String()
			ztMember, err = client.AddAuthorizedZtNetworkMember(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(ztMember).ShouldNot(gomega.BeNil())

			list, err := client.ListAuthorizedZTNetworkMembers(context.Background(), &grpc_application_go.ListAuthorizedZtNetworkMemberRequest{
				OrganizationId: toAdd.OrganizationId,
				AppInstanceId:  toAdd.AppInstanceId,
				ZtNetworkId:    toAdd.NetworkId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(list).ShouldNot(gomega.BeNil())
			gomega.Expect(len(list.Members)).Should(gomega.Equal(3))

		})
	})

	ginkgo.Context("Parametrized Descriptor", func() {
		ginkgo.BeforeEach(func() {
			ginkgo.By("creating required descriptor", func() {
				// Initial data
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())
				targetDescriptor = app

				instToAdd := generateAddAppInstance(targetDescriptor.OrganizationId, targetDescriptor.AppDescriptorId)
				inst, err := client.AddAppInstance(context.Background(), instToAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(inst).ShouldNot(gomega.BeNil())
				targetInstance = inst
			})
		})

		ginkgo.It("should be able to add parametrized descriptor", func() {
			parameterDescriptor := generateParametrizedDescriptor(targetDescriptor, targetInstance.AppInstanceId)

			newDesc, err := client.AddParametrizedDescriptor(context.Background(), parameterDescriptor)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(newDesc).ShouldNot(gomega.BeNil())
		})
		ginkgo.It("should not be able to add parametrized descriptor of a non-existent organization", func() {
			parameterDescriptor := generateParametrizedDescriptor(targetDescriptor, targetInstance.AppInstanceId)

			parameterDescriptor.OrganizationId = uuid.New().String()

			newDesc, err := client.AddParametrizedDescriptor(context.Background(), parameterDescriptor)
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(newDesc).Should(gomega.BeNil())
		})
		ginkgo.It("should not be able to add parametrized descriptor of a non-existent descriptor", func() {
			parameterDescriptor := generateParametrizedDescriptor(targetDescriptor, targetInstance.AppInstanceId)

			parameterDescriptor.AppDescriptorId = uuid.New().String()

			newDesc, err := client.AddParametrizedDescriptor(context.Background(), parameterDescriptor)
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(newDesc).Should(gomega.BeNil())
		})
		ginkgo.It("should not be able to add parametrized descriptor of a non-existent instance", func() {
			parameterDescriptor := generateParametrizedDescriptor(targetDescriptor, targetInstance.AppInstanceId)

			parameterDescriptor.AppInstanceId = uuid.New().String()

			newDesc, err := client.AddParametrizedDescriptor(context.Background(), parameterDescriptor)
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(newDesc).Should(gomega.BeNil())
		})

		ginkgo.It("should be able to get parametrized descriptor", func() {
			parameterDescriptor := generateParametrizedDescriptor(targetDescriptor, targetInstance.AppInstanceId)

			newDesc, err := client.AddParametrizedDescriptor(context.Background(), parameterDescriptor)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(newDesc).ShouldNot(gomega.BeNil())

			retrieved, err := client.GetParametrizedDescriptor(context.Background(), &grpc_application_go.AppInstanceId{
				OrganizationId: parameterDescriptor.OrganizationId,
				AppInstanceId:  parameterDescriptor.AppInstanceId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).NotTo(gomega.BeNil())
			gomega.Expect(retrieved.Name).Should(gomega.Equal(retrieved.Name))

		})
		ginkgo.It("should not be able to get a non-existent parametrized descriptor", func() {
			retrieved, err := client.GetParametrizedDescriptor(context.Background(), &grpc_application_go.AppInstanceId{
				OrganizationId: targetDescriptor.OrganizationId,
				AppInstanceId:  uuid.New().String(),
			})
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(retrieved).To(gomega.BeNil())

		})

	})

})
