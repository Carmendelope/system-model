/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	appProvider "github.com/nalej/system-model/internal/pkg/provider/application"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/nalej/system-model/internal/pkg/server/testhelpers"

	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"math/rand"
)

func generateRandomSpecs() * grpc_application_go.DeploySpecs {
	return &grpc_application_go.DeploySpecs{
		Cpu:int64(rand.Intn(100)),
		Memory:int64(rand.Intn(10000)),
		Replicas:int32(rand.Intn(10)),
	}
}

func generateRandomService(index int) * grpc_application_go.Service {
	credentials := &grpc_application_go.ImageCredentials{
		Username: "username",
		Password: "****",
		Email : "email@company.com",
		DockerRepository: "repo",
	}
	endpoints := make([]*grpc_application_go.Endpoint, 0)
	endpoints = append(endpoints, &grpc_application_go.Endpoint{
		Type: grpc_application_go.EndpointType_REST,
		Path: "/",
	})
	ports := make([]*grpc_application_go.Port, 0)
	ports = append(ports, &grpc_application_go.Port{
		Name : "simple-endpoint",
		InternalPort: 80,
		ExposedPort: 80,
		Endpoints: endpoints,
	})

	storage := make ([]*grpc_application_go.Storage,0)
	storage = append(storage, &grpc_application_go.Storage{
		Size: 12345,
		MountPath:"../path/",
		Type: grpc_application_go.StorageType_CLUSTER_LOCAL,
	})
	configs := make ([]*grpc_application_go.ConfigFile, 0)
	configs = append(configs, &grpc_application_go.ConfigFile{
		Name: "Config file name",
		Content: []byte{0x00, 0x01, 0x02},
		MountPath:"./path..",
	})

	return &grpc_application_go.Service{
		Name: fmt.Sprintf("service-%d", index),
		Type: grpc_application_go.ServiceType_DOCKER,
		Image: fmt.Sprintf("image:v%d", rand.Intn(10)),
		Specs: generateRandomSpecs(),
		ExposedPorts: ports,
		Credentials: credentials,
		Storage: storage,
		EnvironmentVariables: map[string]string{"env01":"env01Label", "env02":"env02Label"},
		DeployAfter: []string{"after1", "after2"},
		Labels: map[string]string {"label1":"service label 1","label2":"service label 2"},
		Configs: configs,
		RunArguments: []string{"arg1", "arg2", "arg3"},
		DeploymentSelectors:map[string]string{"clusterDeployment": "EDGE"},
	}
}

func generateServiceGroup(services []*grpc_application_go.Service) * grpc_application_go.ServiceGroup{


	return &grpc_application_go.ServiceGroup{
		Name:            "Service Group",
		Services: services,
		Policy: grpc_application_go.CollocationPolicy_SEPARATE_CLUSTERS,
		Specs: &grpc_application_go.ServiceGroupDeploymentSpecs{
			NumReplicas: 5,
			MultiClusterReplica: false,
		},
		Labels:map[string]string{"label1":"sg_label1", "label2":"sg_label2", "label3":"sg_label3"},
	}
}



func generateAddAppDescriptor(orgID string, numServices int) * grpc_application_go.AddAppDescriptorRequest {
	services := make([]*grpc_application_go.Service, 0)
	for i := 0; i < numServices; i++ {
		services = append(services, generateRandomService(i))
	}
	securityRules := make([]*grpc_application_go.SecurityRule, 0)
	for i := 0; i < (numServices ); i++ {
		securityRules = append(securityRules, &grpc_application_go.SecurityRule{
			RuleId : fmt.Sprintf("r%d", i),
			Name: fmt.Sprintf("%d -> %d", i, i+1),
			TargetServiceGroupName: fmt.Sprintf("targetServiceGroupName-%d", i),
			TargetServiceName: fmt.Sprintf("targetServiceName-%d", i),
			TargetPort: 80,
			Access: grpc_application_go.PortAccess_APP_SERVICES,
			AuthServiceGroupName: fmt.Sprintf("AuthServiceGroupName-%d", i),
			AuthServices: []string{fmt.Sprintf("s%d", i+1)},
			DeviceGroups:[]string{"device_1", "device_2"},
		})
	}
	groups := make ([]*grpc_application_go.ServiceGroup, 0)
	groups = append(groups, generateServiceGroup(services))

	return &grpc_application_go.AddAppDescriptorRequest{
		RequestId:"request_id",
		OrganizationId:orgID,
		Name: "new app",
		ConfigurationOptions: map[string]string{"conf1":"conf1", "conf2":"conf2"},
		EnvironmentVariables: map[string]string{"var1":"env1"},
		Labels: map[string]string{"label1":"eti1"},
		Rules: securityRules,
		Groups: groups,
	}
}

func generateAddAppInstance(organizationID string, appDescriptorID string) * grpc_application_go.AddAppInstanceRequest {
	return &grpc_application_go.AddAppInstanceRequest{
		OrganizationId:       organizationID,
		AppDescriptorId:      appDescriptorID,
		Name:                 fmt.Sprintf("app instance %d", rand.Int31n(100)),
	}
}

func generateUpdateAppInstance(organizationID string, appInstanceID string,
	status grpc_application_go.ApplicationStatus) * grpc_application_go.UpdateAppStatusRequest {
	return &grpc_application_go.UpdateAppStatusRequest{
		OrganizationId: organizationID,
		AppInstanceId: appInstanceID,
		Status: status,
	}
}

func generateUpdateServiceStatus(organizationID string, appInstanceID string, serviceID string,
    appDescriptorID string, status grpc_application_go.ServiceStatus) * grpc_application_go.UpdateServiceStatusRequest {
    endpoint := make([]string,0)
    endpoint = append(endpoint, "enpoint1")
    return &grpc_application_go.UpdateServiceStatusRequest{
        OrganizationId: organizationID,
        AppInstanceId: appInstanceID,
        Status: status,
		//Endpoints: endpoint,
		DeployedOnClusterId: fmt.Sprintf("Deploy on cluster - %d", rand.Int31n(100)),
    }
}


var _ = ginkgo.Describe("Applications", func(){

	const numServices = 2

	// gRPC server
	var server * grpc.Server
	// grpc test listener
	var listener * bufconn.Listener
	// client
	var client grpc_application_go.ApplicationsClient

	// Target organization.
	var targetOrganization * entities.Organization

	var targetDescriptor * grpc_application_go.AppDescriptor

	// Organization Provider
	var organizationProvider orgProvider.Provider
	var applicationProvider appProvider.Provider

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()


		// Create providers
		organizationProvider = orgProvider.NewMockupOrganizationProvider()
		applicationProvider = appProvider.NewMockupOrganizationProvider()

		manager := NewManager(organizationProvider, applicationProvider)
		handler := NewHandler(manager)
		grpc_application_go.RegisterApplicationsServer(server, handler)

		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_application_go.NewApplicationsClient(conn)
	})

	ginkgo.AfterSuite(func(){
		server.Stop()
		listener.Close()
	})

	ginkgo.BeforeEach(func(){
		ginkgo.By("cleaning the mockups", func(){
			organizationProvider.(*orgProvider.MockupOrganizationProvider).Clear()
			applicationProvider.(*appProvider.MockupApplicationProvider).Clear()

			// Initial data
			targetOrganization = testhelpers.CreateOrganization(organizationProvider)
		})
	})

	ginkgo.Context("Application descriptor", func(){
		ginkgo.Context("adding application descriptors", func(){
			ginkgo.It("should add an application descriptor", func(){
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())
				gomega.Expect(app.AppDescriptorId).ShouldNot(gomega.BeNil())
				gomega.Expect(app.Name).Should(gomega.Equal(toAdd.Name))
				gomega.Expect(len(toAdd.Groups)).Should(gomega.Equal(len(app.Groups)))
			})
			ginkgo.It("should fail on an empty request", func(){
				toAdd := &grpc_application_go.AddAppDescriptorRequest{}
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(app).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a non existing organization", func(){
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				toAdd.OrganizationId = "does not exists"
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(app).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a descriptor without services", func(){
				toAdd := generateAddAppDescriptor(targetOrganization.ID, 0)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(app).Should(gomega.BeNil())
			})
		})
		ginkgo.Context("get application descriptor", func(){
		    ginkgo.It("should get an existing app descriptor", func(){
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())
				retrieved, err := client.GetAppDescriptor(context.Background(), &grpc_application_go.AppDescriptorId{
					OrganizationId: app.OrganizationId,
					AppDescriptorId: app.AppDescriptorId,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(retrieved.Name).Should(gomega.Equal(app.Name))
		    })
		    ginkgo.It("should fail on a non existing application", func(){
				retrieved, err := client.GetAppDescriptor(context.Background(), &grpc_application_go.AppDescriptorId{
					OrganizationId: targetOrganization.ID,
					AppDescriptorId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
		    })
		    ginkgo.It("should fail on a non existing organization", func(){
				retrieved, err := client.GetAppDescriptor(context.Background(), &grpc_application_go.AppDescriptorId{
					OrganizationId: "does not exists",
					AppDescriptorId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
		    })
		})
		ginkgo.Context("listing application descriptors", func(){
			ginkgo.It("should list apps on an existing organization", func(){
			    numDescriptors := 3
			    for i := 0; i < numDescriptors; i ++ {
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
			ginkgo.It("should fail on a non existing organization", func(){
				retrieved, err := client.ListAppDescriptors(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
			ginkgo.It("should work on an organization without descriptors", func(){
				gomega.Expect(organizationProvider).ShouldNot(gomega.BeNil())
				retrieved, err := client.ListAppDescriptors(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: targetOrganization.ID,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(len(retrieved.Descriptors)).Should(gomega.Equal(0))
			})
		})

		ginkgo.Context("removing application descriptors", func(){
			ginkgo.It("should be able to remove an existing descriptor", func(){
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())

				toRemove := &grpc_application_go.AppDescriptorId{
					OrganizationId:       app.OrganizationId,
					AppDescriptorId:      app.AppDescriptorId,
				}
				success, err := client.RemoveAppDescriptor(context.Background(), toRemove)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(success).ShouldNot(gomega.BeNil())
			})
			ginkgo.It("should fail if the organization does not exists", func(){
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())
				toRemove := &grpc_application_go.AppDescriptorId{
					OrganizationId:       "unknown",
					AppDescriptorId:      app.AppDescriptorId,
				}
				success, err := client.RemoveAppDescriptor(context.Background(), toRemove)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(success).Should(gomega.BeNil())
			})
			ginkgo.It("should fail if the descriptor does not exits", func(){
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())
				toRemove := &grpc_application_go.AppDescriptorId{
					OrganizationId:       app.OrganizationId,
					AppDescriptorId:      "unknown",
				}
				success, err := client.RemoveAppDescriptor(context.Background(), toRemove)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(success).Should(gomega.BeNil())
			})
		})

	})

	ginkgo.Context("Application instance", func(){
		ginkgo.BeforeEach(func(){
			ginkgo.By("creating required descriptor", func(){
				// Initial data
				toAdd := generateAddAppDescriptor(targetOrganization.ID, numServices)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(app).ShouldNot(gomega.BeNil())
				targetDescriptor = app
			})
		})
	    ginkgo.Context("adding application instance", func(){
			ginkgo.It("should add an app instance", func(){
			    toAdd := generateAddAppInstance(targetDescriptor.OrganizationId, targetDescriptor.AppDescriptorId)
			    added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())
			    gomega.Expect(added.AppInstanceId).ShouldNot(gomega.BeEmpty())
			    gomega.Expect(added.OrganizationId).Should(gomega.Equal(targetDescriptor.OrganizationId))
			    gomega.Expect(added.AppDescriptorId).Should(gomega.Equal(targetDescriptor.AppDescriptorId))
			})
			ginkgo.It("should fail on a non existing app descriptor", func(){
				toAdd := generateAddAppInstance(targetDescriptor.OrganizationId, "does not exists")
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(added).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a non existing organization", func(){
				toAdd := generateAddAppInstance("does not exists", targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(added).Should(gomega.BeNil())
			})
	    })
	    ginkgo.Context("get application instance", func(){
			ginkgo.It("should retrieve an existing app", func(){
				toAdd := generateAddAppInstance(targetOrganization.ID, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())
				gomega.Expect(added.AppInstanceId).ShouldNot(gomega.BeEmpty())
				retrieved, err := client.GetAppInstance(context.Background(), &grpc_application_go.AppInstanceId{
					OrganizationId: added.OrganizationId,
					AppInstanceId: added.AppInstanceId,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(retrieved.Name).Should(gomega.Equal(added.Name))
			})
			ginkgo.It("should fail on a non existing instance", func(){
				retrieved, err := client.GetAppInstance(context.Background(), &grpc_application_go.AppInstanceId{
					OrganizationId: targetDescriptor.OrganizationId,
					AppInstanceId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a non existing organization", func(){
				retrieved, err := client.GetAppInstance(context.Background(), &grpc_application_go.AppInstanceId{
					OrganizationId: "does not exists",
					AppInstanceId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
	    })
	    ginkgo.Context("listing application instances", func(){
			ginkgo.It("should retrieve instances on an existing organization", func(){
				numInstances := 3
				for i := 0; i < numInstances; i ++ {
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
			ginkgo.It("should work on an organization without instances", func(){
				retrieved, err := client.ListAppInstances(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: targetOrganization.ID,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(len(retrieved.Instances)).Should(gomega.Equal(0))
			})
			ginkgo.It("should fail on a non existing organization", func(){
				retrieved, err := client.ListAppInstances(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
	    })
		ginkgo.Context("update application instance", func(){
			ginkgo.PIt("should update instance and return the new values", func(){
			})
		})

		ginkgo.Context("update service status in application instance", func(){
		    ginkgo.PIt("should update instance and return the new values", func(){
            })
        })

		ginkgo.Context("removing application instances", func(){
			ginkgo.It("should be able to remove an existing instance", func(){
				toAdd := generateAddAppInstance(targetOrganization.ID, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())
				toRemove := &grpc_application_go.AppInstanceId{
					OrganizationId:       added.OrganizationId,
					AppInstanceId:        added.AppInstanceId,
				}
				success, err := client.RemoveAppInstance(context.Background(), toRemove)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(success).ShouldNot(gomega.BeNil())
			})
			ginkgo.It("should fail if the organization does not exists", func(){
				toAdd := generateAddAppInstance(targetOrganization.ID, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())
				toRemove := &grpc_application_go.AppInstanceId{
					OrganizationId:       "unknown",
					AppInstanceId:        added.AppInstanceId,
				}
				success, err := client.RemoveAppInstance(context.Background(), toRemove)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(success).Should(gomega.BeNil())
			})
			ginkgo.It("should fail if the descriptor does not exits", func(){
				toAdd := generateAddAppInstance(targetOrganization.ID, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())
				toRemove := &grpc_application_go.AppInstanceId{
					OrganizationId:       added.OrganizationId,
					AppInstanceId:        "unknown",
				}
				success, err := client.RemoveAppInstance(context.Background(), toRemove)
				gomega.Expect(err).To(gomega.HaveOccurred())
				gomega.Expect(success).Should(gomega.BeNil())
			})
		})

		ginkgo.Context("Adding ServiceGroupInstance ", func() {
			ginkgo.It("should be able to add a service group instance", func() {
				toAdd := generateAddAppInstance(targetDescriptor.OrganizationId, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())

				sgToAdd := &grpc_application_go.AddServiceGroupInstanceRequest{
					OrganizationId:  targetDescriptor.OrganizationId,
					AppDescriptorId: targetDescriptor.AppDescriptorId,
					AppInstanceId:   added.AppInstanceId,
					ServiceGroupId:  added.Groups[0].ServiceGroupId,
				}

				sgReceived, err := client.AddServiceGroupInstance(context.Background(), sgToAdd)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(sgReceived.ServiceGroupId).Should(gomega.Equal(sgToAdd.ServiceGroupId))
			})
			ginkgo.It("should not be able to add a service group instance of a non existing group", func() {
				toAdd := generateAddAppInstance(targetDescriptor.OrganizationId, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())

				sgToAdd := &grpc_application_go.AddServiceGroupInstanceRequest{
					OrganizationId:  targetDescriptor.OrganizationId,
					AppDescriptorId: targetDescriptor.AppDescriptorId,
					AppInstanceId:   added.AppInstanceId,
					ServiceGroupId:  uuid.New().String(),
				}

				_, err = client.AddServiceGroupInstance(context.Background(), sgToAdd)
				gomega.Expect(err).NotTo(gomega.Succeed())
			})

		})

		ginkgo.Context("Adding ServiceInstance ", func() {
			ginkgo.It("should be able to add a service instance", func() {
				toAdd := generateAddAppInstance(targetDescriptor.OrganizationId, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())

				sgToAdd := &grpc_application_go.AddServiceGroupInstanceRequest{
					OrganizationId:  targetDescriptor.OrganizationId,
					AppDescriptorId: targetDescriptor.AppDescriptorId,
					AppInstanceId:   added.AppInstanceId,
					ServiceGroupId:  added.Groups[0].ServiceGroupId,
				}

				sgReceived, err := client.AddServiceGroupInstance(context.Background(), sgToAdd)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(sgReceived.ServiceGroupId).Should(gomega.Equal(sgToAdd.ServiceGroupId))

				sToAdd := &grpc_application_go.AddServiceInstanceRequest{
					OrganizationId:  targetDescriptor.OrganizationId,
					AppDescriptorId: targetDescriptor.AppDescriptorId,
					AppInstanceId:   added.AppInstanceId,
					ServiceGroupId:  sgReceived.ServiceGroupId,
					ServiceGroupInstanceId: sgReceived.ServiceGroupInstanceId,
					ServiceId: added.Groups[0].ServiceInstances[0].ServiceId,
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

				sgToAdd := &grpc_application_go.AddServiceGroupInstanceRequest{
					OrganizationId:  targetDescriptor.OrganizationId,
					AppDescriptorId: targetDescriptor.AppDescriptorId,
					AppInstanceId:   added.AppInstanceId,
					ServiceGroupId:  added.Groups[0].ServiceGroupId,
				}

				sgReceived, err := client.AddServiceGroupInstance(context.Background(), sgToAdd)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(sgReceived.ServiceGroupId).Should(gomega.Equal(sgToAdd.ServiceGroupId))

				sToAdd := &grpc_application_go.AddServiceInstanceRequest{
					OrganizationId:  targetDescriptor.OrganizationId,
					AppDescriptorId: targetDescriptor.AppDescriptorId,
					AppInstanceId:   added.AppInstanceId,
					ServiceGroupId:  sgReceived.ServiceGroupId,
					ServiceGroupInstanceId: sgReceived.ServiceGroupInstanceId,
					ServiceId: uuid.New().String(),
				}

				_, err = client.AddServiceInstance(context.Background(), sToAdd)
				gomega.Expect(err).NotTo(gomega.Succeed())

			})

		})
	})
})
