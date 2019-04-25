/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package eic

import (
	"context"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/eic"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/nalej/system-model/internal/pkg/server/testhelpers"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func createAddEdgeControllerRequest(organizationID string) *grpc_inventory_go.AddEdgeControllerRequest{
	testController := eic.CreateTestEdgeController()

	return &grpc_inventory_go.AddEdgeControllerRequest{
		OrganizationId:       organizationID,
		Name:                 testController.Name,
		Labels:               testController.Labels,
	}

}

var _ = ginkgo.Describe("Asset service", func() {
	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client grpc_inventory_go.ControllersClient

	// Target organization.
	var targetOrganization *entities.Organization

	// Providers
	var organizationProvider orgProvider.Provider
	var controllerProvider eic.Provider

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()
		test.LaunchServer(server, listener)

		// Register the service
		organizationProvider = orgProvider.NewMockupOrganizationProvider()
		controllerProvider = eic.NewMockupEICProvider()
		manager := NewManager(controllerProvider, organizationProvider)
		handler := NewHandler(manager)
		grpc_inventory_go.RegisterControllersServer(server, handler)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_inventory_go.NewControllersClient(conn)
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
	})

	ginkgo.BeforeEach(func(){
		ginkgo.By("cleaning the mockups", func(){
			organizationProvider.(*orgProvider.MockupOrganizationProvider).Clear()
			controllerProvider.Clear()
			// Initial data
			targetOrganization = testhelpers.CreateOrganization(organizationProvider)
		})
	})

	ginkgo.It("should be able to add a new controller", func(){
		toAdd := createAddEdgeControllerRequest(targetOrganization.ID)
		added, err := client.Add(context.Background(), toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(added).ShouldNot(gomega.BeNil())
		gomega.Expect(added.EdgeControllerId).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should be able to list controllers", func(){
		numControllers := 10
		for index := 0; index < numControllers; index++{
			toAdd := createAddEdgeControllerRequest(targetOrganization.ID)
			added, err := client.Add(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.EdgeControllerId).ShouldNot(gomega.BeEmpty())
		}
		orgID := &grpc_organization_go.OrganizationId{
			OrganizationId:       targetOrganization.ID,
		}
		allControllers, err := client.List(context.Background(), orgID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(len(allControllers.Controllers)).Should(gomega.Equal(numControllers))
	})

	ginkgo.It("should be able to remove controllers", func(){
		toAdd := createAddEdgeControllerRequest(targetOrganization.ID)
		added, err := client.Add(context.Background(), toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(added).ShouldNot(gomega.BeNil())
		edgeControllerID := &grpc_inventory_go.EdgeControllerId{
			OrganizationId:       added.OrganizationId,
			EdgeControllerId:              added.EdgeControllerId,
		}
		success, err := client.Remove(context.Background(), edgeControllerID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(success).ShouldNot(gomega.BeNil())
	})

	ginkgo.Context("update operations", func(){
		ginkgo.It("should be able to add new labels", func(){
			toAdd := createAddEdgeControllerRequest(targetOrganization.ID)
			added, err := client.Add(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			newLabels := make(map[string]string, 0)
			newLabels["k1"]="v1"
			updateRequest := &grpc_inventory_go.UpdateEdgeControllerRequest{
				OrganizationId:       added.OrganizationId,
				EdgeControllerId:              added.EdgeControllerId,
				AddLabels:            true,
				RemoveLabels:         false,
				Labels:               newLabels,
			}

			updated, err := client.Update(context.Background(), updateRequest)
			gomega.Expect(err).To(gomega.Succeed())
			value, exits := updated.Labels["k1"]
			gomega.Expect(exits).To(gomega.BeTrue())
			gomega.Expect(value).Should(gomega.Equal("v1"))
		})
	})

})