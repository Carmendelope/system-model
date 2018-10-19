/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package application

import (
	"context"
	"fmt"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/system-model/internal/pkg/entities"

	appProvider "github.com/nalej/system-model/internal/pkg/provider/application"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"

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

func generateRandomService(orgID string, index int) * grpc_application_go.Service {
	endpoints := make([]*grpc_application_go.Endpoint, 0)
	endpoints = append(endpoints, &grpc_application_go.Endpoint{
		Type: grpc_application_go.EndpointType_REST,
		Path: "/",
	})
	ports := make([]*grpc_application_go.Port, 0)
	ports = append(ports, &grpc_application_go.Port{
		Name : "simple endpoint",
		InternalPort: 80,
		ExposedPort: 80,
		Endpoints: endpoints,
	})
	return &grpc_application_go.Service{
		OrganizationId: orgID,
		ServiceId: fmt.Sprintf("s%d", index),
		Name: fmt.Sprintf("Service %d", index),
		Description : fmt.Sprintf("Description s%d", index),
		Type: grpc_application_go.ServiceType_DOCKER,
		Image: fmt.Sprintf("image:v%d", rand.Intn(10)),
		Specs: generateRandomSpecs(),
		ExposedPorts: ports,
	}
}

func generateAddAppDescriptor(orgID string, numServices int) * grpc_application_go.AddAppDescriptorRequest {
	services := make([]*grpc_application_go.Service, 0)
	for i := 0; i < numServices; i++ {
		services = append(services, generateRandomService(orgID, i))
	}
	securityRules := make([]*grpc_application_go.SecurityRule, 0)
	for i := 0; i < (numServices - 1); i++ {
		securityRules = append(securityRules, &grpc_application_go.SecurityRule{
			OrganizationId: orgID,
			RuleId : fmt.Sprintf("r%d", i),
			Name: fmt.Sprintf("%d -> %d", i, i+1),
			SourceServiceId: fmt.Sprintf("s%d", i),
			SourcePort: 80,
			Access: grpc_application_go.PortAccess_APP_SERVICES,
			AuthServices: []string{fmt.Sprintf("s%d", i+1)},
		})
	}
	envVars := make(map[string]string, 0)
	envVars["VAR1"] = "VALUE1"
	return &grpc_application_go.AddAppDescriptorRequest{
		RequestId:"request_id",
		OrganizationId:orgID,
		Name: "new app",
		Description:"description",
		EnvironmentVariables: envVars,
		Rules: securityRules,
		Services: services,
	}
}

func generateAddAppInstance(organizationID string, appDescriptorID string) * grpc_application_go.AddAppInstanceRequest {
	return &grpc_application_go.AddAppInstanceRequest{
		OrganizationId:       organizationID,
		AppDescriptorId:      appDescriptorID,
		Name:                 fmt.Sprintf("app instance %d", rand.Int31n(100)),
		Description:          "app instance description",
	}
}

func createOrganization(orgProvider orgProvider.Provider) * entities.Organization {
	toAdd := entities.NewOrganization("test org")
	err := orgProvider.Add(*toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	return toAdd
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
    return &grpc_application_go.UpdateServiceStatusRequest{
        OrganizationId: organizationID,
        AppInstanceId: appInstanceID,
        ServiceId: serviceID,
        Status: status,
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
			targetOrganization = createOrganization(organizationProvider)
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
				gomega.Expect(len(toAdd.Services)).Should(gomega.Equal(len(app.Services)))
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
			    retrieved, err := client.GetAppDescriptors(context.Background(), &grpc_organization_go.OrganizationId{
			    	OrganizationId: targetOrganization.ID,
				})
			    gomega.Expect(err).Should(gomega.Succeed())
			    gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			    gomega.Expect(len(retrieved.Descriptors)).Should(gomega.Equal(numDescriptors))
			})
			ginkgo.It("should fail on a non existing organization", func(){
				retrieved, err := client.GetAppDescriptors(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
			ginkgo.It("should work on an organization without descriptors", func(){
				gomega.Expect(organizationProvider).ShouldNot(gomega.BeNil())
				retrieved, err := client.GetAppDescriptors(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: targetOrganization.ID,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(len(retrieved.Descriptors)).Should(gomega.Equal(0))
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
				retrieved, err := client.GetAppInstances(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: targetOrganization.ID,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(len(retrieved.Instances)).Should(gomega.Equal(numInstances))
			})
			ginkgo.It("should work on an organization without instances", func(){
				retrieved, err := client.GetAppInstances(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: targetOrganization.ID,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(len(retrieved.Instances)).Should(gomega.Equal(0))
			})
			ginkgo.It("should fail on a non existing organization", func(){
				retrieved, err := client.GetAppInstances(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
	    })
		ginkgo.Context("update application instance", func(){
			ginkgo.It("should update instance and return the new values", func(){
				toAdd := generateAddAppInstance(targetOrganization.ID, targetDescriptor.AppDescriptorId)
				added, err := client.AddAppInstance(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(added).ShouldNot(gomega.BeNil())
				gomega.Expect(added.AppInstanceId).ShouldNot(gomega.BeEmpty())
				// update
				req := generateUpdateAppInstance(targetOrganization.ID, added.AppInstanceId,
					grpc_application_go.ApplicationStatus_RUNNING)
				_, err = client.UpdateAppStatus(context.Background(), req)
				gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
				// recover and check the changes
				recovered, err := client.GetAppInstance(context.Background(),
					&grpc_application_go.AppInstanceId{OrganizationId:req.OrganizationId,AppInstanceId:req.AppInstanceId})
				gomega.Expect(recovered.Status).To(gomega.Equal(req.Status))
				gomega.Expect(recovered.AppInstanceId).To(gomega.Equal(req.AppInstanceId))

			})
		})

		ginkgo.Context("update service status in application instance", func(){
		    ginkgo.It("should update intance and return the new values", func(){
                toAdd := generateAddAppInstance(targetOrganization.ID, targetDescriptor.AppDescriptorId)
                added, err := client.AddAppInstance(context.Background(), toAdd)
                gomega.Expect(err).Should(gomega.Succeed())
                gomega.Expect(added).ShouldNot(gomega.BeNil())
                gomega.Expect(added.AppInstanceId).ShouldNot(gomega.BeEmpty())
                // update it
                req := generateUpdateServiceStatus(added.OrganizationId, added.AppInstanceId,
                     added.Services[0].ServiceId, added.AppDescriptorId, grpc_application_go.ServiceStatus_SERVICE_RUNNING)
                _, err = client.UpdateServiceStatus(context.Background(), req)
                // recover changes
                recovered, err := client.GetAppInstance(context.Background(),
                    &grpc_application_go.AppInstanceId{OrganizationId:added.OrganizationId,AppInstanceId:added.AppInstanceId})
                gomega.Expect(err).Should(gomega.BeNil())
                gomega.Expect(recovered.AppInstanceId).To(gomega.Equal(added.AppInstanceId))
                gomega.Expect(recovered.Services[0].Status).To(gomega.Equal(grpc_application_go.ServiceStatus_SERVICE_RUNNING))
            })
        })
	})
})
