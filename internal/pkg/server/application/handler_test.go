/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package application

import (
	"context"
	"fmt"
	"github.com/nalej/grpc-organization-go"

	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"

	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/satori/go.uuid"
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

func generateRandomService() * grpc_application_go.Service {
	return &grpc_application_go.Service{
		Name: uuid.NewV1().String(),
		Description : uuid.NewV1().String(),
		Image: fmt.Sprintf("image:v%d", rand.Intn(10)),
		Specs: generateRandomSpecs(),
		Type: grpc_application_go.ServiceType_DOCKER,
	}
}

func generateAddAppDescriptor(orgID string, numServices int) * grpc_application_go.AddAppDescriptorRequest {
	services := make([]*grpc_application_go.Service, 0)
	for i := 0; i < numServices; i++ {
		services = append(services, generateRandomService())
	}
	return &grpc_application_go.AddAppDescriptorRequest{
		RequestId:"request_id",
		OrganizationId:orgID,
		Name: "new app",
		Description:"description",
		Services: services,
	}
}

func createOrganization(orgProvider orgProvider.Provider) * grpc_organization_go.Organization {

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

	var targetOrganization grpc_organization_go.Organization

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()
		test.LaunchServer(server, listener)

		// Register the service
		//appProvider := nil
		manager := NewManager(nil)
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
		})
		ginkgo.Context("get application descriptor", func(){
		    ginkgo.It("should get an existing app descriptor", func(){

		    })
		    ginkgo.It("should fail on a non existing application", func(){

		    })
		    ginkgo.It("should fail on a non existing organization", func(){

		    })
		})
		ginkgo.Context("listing application descriptors", func(){
			ginkgo.It("should list apps on an existing organization", func(){
			    
			})
			ginkgo.It("should fail on a non existing organization", func(){
			    
			})
			ginkgo.It("should work on an organization without descriptors", func(){

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
