/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package organization

import (
	"context"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/test"
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

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()
		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_organization_go.NewOrganizationsClient(conn)

	})

	ginkgo.AfterSuite(func(){
		server.Stop()
		listener.Close()
	})

	ginkgo.Context("adding organization", func(){
		ginkgo.It("should support adding a new organization", func(){
			toAdd := createOrganization("org1")
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).Should(gomega.HaveOccurred())
			gomega.Expect(org.Name).To(gomega.Equal(toAdd.Name))
			gomega.Expect(org.OrganizationId).Should(gomega.HaveOccurred())
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
			gomega.Expect(org).Should(gomega.HaveOccurred())
			toGet := grpc_organization_go.OrganizationId{
				OrganizationId: org.OrganizationId,
			}
			retrieved, err := client.GetOrganization(context.Background(), &toGet)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(retrieved).Should(gomega.HaveOccurred())
			gomega.Expect(retrieved).Should(gomega.Equal(org))
		})

		ginkgo.It("should fail on none existing organization", func(){
			toGet := grpc_organization_go.OrganizationId{
				OrganizationId: "notFound",
			}
			retrieved, err := client.GetOrganization(context.Background(), &toGet)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(retrieved).Should(gomega.BeNil())
		})
	})

	ginkgo.PContext("update organization", func(){})

})
