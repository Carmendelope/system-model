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

func generateRandomService(index int) * grpc_application_go.Service {
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
		ServiceId: fmt.Sprintf("s%d", index),
		Name: fmt.Sprintf("Service %d", index),
		Description : fmt.Sprintf("Descriptin s%d", index),
		Image: fmt.Sprintf("image:v%d", rand.Intn(10)),
		Specs: generateRandomSpecs(),
		Type: grpc_application_go.ServiceType_DOCKER,
	}
}

func generateAddAppDescriptor(orgID string, numServices int) * grpc_application_go.AddAppDescriptorRequest {
	services := make([]*grpc_application_go.Service, 0)
	for i := 0; i < numServices; i++ {
		services = append(services, generateRandomService(i))
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

func createOrganization(orgProvider orgProvider.Provider) * entities.Organization {
	toAdd := entities.NewOrganization("test org")
	err := orgProvider.Add(*toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	return toAdd
}

var _ = ginkgo.Describe("Applications", func(){

	const orgID = "existingOrg"
	const numServices = 10

	// gRPC server
	var server * grpc.Server
	// grpc test listener
	var listener * bufconn.Listener
	// client
	var client grpc_application_go.ApplicationsClient

	// Target organization.
	var targetOrganization * entities.Organization

	// Organization Provider
	var organizationProvider orgProvider.Provider
	var applicationProvider appProvider.Provider

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()
		test.LaunchServer(server, listener)

		// Create providers
		organizationProvider = orgProvider.NewMockupOrganizationProvider()

		// Initial data
		targetOrganization = createOrganization(organizationProvider)

		manager := NewManager(organizationProvider, applicationProvider)
		handler := NewHandler(manager)
		grpc_application_go.RegisterApplicationsServer(server, handler)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_application_go.NewApplicationsClient(conn)
	})

	ginkgo.AfterSuite(func(){
		server.Stop()
		listener.Close()
	})

	ginkgo.Context("Application descriptor", func(){
		ginkgo.Context("adding application descriptors", func(){
			ginkgo.It("should add an application descriptor", func(){
				toAdd := generateAddAppDescriptor(orgID, numServices)
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
				toAdd := generateAddAppDescriptor(orgID, numServices)
				toAdd.OrganizationId = "does not exists"
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(app).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a descriptor without services", func(){
				toAdd := generateAddAppDescriptor(orgID, 0)
				app, err := client.AddAppDescriptor(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(app).Should(gomega.BeNil())
			})
		})
		ginkgo.Context("get application descriptor", func(){
		    ginkgo.It("should get an existing app descriptor", func(){
				toAdd := generateAddAppDescriptor(orgID, numServices)
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
					OrganizationId: orgID,
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
			ginkgo.BeforeEach("clear mockup provider", func(){

			})
			ginkgo.It("should list apps on an existing organization", func(){
			    numDescriptors := 3
			    for i := 0; i < numDescriptors; i ++ {
					toAdd := generateAddAppDescriptor(orgID, numServices)
					app, err := client.AddAppDescriptor(context.Background(), toAdd)
					gomega.Expect(err).Should(gomega.Succeed())
					gomega.Expect(app).ShouldNot(gomega.BeNil())
				}
			    retrieved, err := client.GetAppDescriptors(context.Background(), &grpc_organization_go.OrganizationId{
			    	OrganizationId: orgID,
				})
			    gomega.Expect(err).Should(gomega.Succeed())
			    gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			    gomega.Expect(len(retrieved.Descriptors)).ShouldNot(gomega.Equal(numDescriptors))
			})
			ginkgo.It("should fail on a non existing organization", func(){
				retrieved, err := client.GetAppDescriptors(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
			ginkgo.It("should work on an organization without descriptors", func(){
				retrieved, err := client.GetAppDescriptors(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: orgID,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(len(retrieved.Descriptors)).ShouldNot(gomega.Equal(0))
			})
		})

	})

	ginkgo.PContext("Application instance", func(){
	    ginkgo.Context("adding application instance", func(){
			ginkgo.PIt("should add an app instance", func(){
			    
			})
			ginkgo.PIt("should fail on a non existing app descriptor", func(){
			    
			})
			ginkgo.PIt("should fail on a non existing organization", func(){
			    
			})
	    })
	    ginkgo.Context("get application instance", func(){
			ginkgo.PIt("should retrieve an existing app", func(){

			})
			ginkgo.PIt("should fail on a non existing instance", func(){

			})
			ginkgo.PIt("should fail on a non existing organization", func(){

			})
	    })
	    ginkgo.Context("listing application instances", func(){
			ginkgo.PIt("should retrieve instances on an existing organization", func(){

			})
			ginkgo.PIt("should work on an organization without instances", func(){

			})
			ginkgo.PIt("should fail on a non existing organization", func(){

			})
	    })
	})
})
