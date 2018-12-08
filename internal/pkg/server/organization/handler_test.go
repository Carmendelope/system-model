/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package organization

import (
	"context"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func createOrganization(name string) * grpc_organization_go.AddOrganizationRequest {
	return &grpc_organization_go.AddOrganizationRequest{
		Name : name,
	}
}

var _ = ginkgo.Describe("Organization service", func(){
	// gRPC server
	var server * grpc.Server
	// grpc test listener
	var listener * bufconn.Listener
	// client
	var client grpc_organization_go.OrganizationsClient

	var orgProvider organization.Provider

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()
		test.LaunchServer(server, listener)

		// Register the service
		orgProvider = organization.NewMockupOrganizationProvider()
		manager := NewManager(orgProvider)
		handler := NewHandler(manager)
		grpc_organization_go.RegisterOrganizationsServer(server, handler)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_organization_go.NewOrganizationsClient(conn)
	})

	ginkgo.AfterSuite(func(){
		server.Stop()
		listener.Close()
	})

	ginkgo.BeforeEach(func() {
		orgProvider.Clear()
	})

	ginkgo.Context("adding organization", func(){
		ginkgo.It("should support adding a new organization", func(){
			toAdd := createOrganization("org1")
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())
			gomega.Expect(org.Name).To(gomega.Equal(toAdd.Name))
			gomega.Expect(org.OrganizationId).ShouldNot(gomega.BeNil())
		})

		ginkgo.It("should fail if the organization name is not specified", func(){
		    toAdd := &grpc_organization_go.AddOrganizationRequest{}
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(org).Should(gomega.BeNil())
		})
	})

	ginkgo.Context("retrieve organization", func(){
		ginkgo.It("should work on existing organization", func(){
			toAdd := createOrganization("org2")
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())
			toGet := grpc_organization_go.OrganizationId{
				OrganizationId: org.OrganizationId,
			}
			retrieved, err := client.GetOrganization(context.Background(), &toGet)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			gomega.Expect(retrieved).Should(gomega.Equal(org))
		})

		ginkgo.It("should recover a list of organizations", func(){
			toAdd := createOrganization("org2")
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())


			toAdd = createOrganization("org3")
			org, err = client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			retrieved, err := client.ListOrganizations(context.Background(), new(grpc_common_go.Empty))
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			gomega.Expect(retrieved.Organizations).ShouldNot(gomega.BeNil())
			gomega.Expect(retrieved.Organizations).Should(gomega.HaveLen(2))
		})

		ginkgo.It("should recover a list empty", func(){
			retrieved, err := client.ListOrganizations(context.Background(), new(grpc_common_go.Empty))
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			gomega.Expect(retrieved.Organizations).Should(gomega.HaveLen(0))
		})

		ginkgo.It("should fail on none existing organization", func(){
			toGet := grpc_organization_go.OrganizationId{
				OrganizationId: "notFound",
			}
			retrieved, err := client.GetOrganization(context.Background(), &toGet)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(retrieved).Should(gomega.BeNil())
		})

		ginkgo.It("should fail on empty request", func(){
			toGet := grpc_organization_go.OrganizationId{}
			retrieved, err := client.GetOrganization(context.Background(), &toGet)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(retrieved).Should(gomega.BeNil())
		})
	})

	ginkgo.PContext("update organization", func(){
		ginkgo.PIt("should support updating an existing organization", func(){

		})

		ginkgo.PIt("should fail on non-existing organization", func(){

		})
	})

	ginkgo.PContext("remove organization", func(){
		ginkgo.PIt("should support removing an existing organization", func(){

		})
		ginkgo.PIt("should fail on non-existing organization", func(){

		})
	})

})
