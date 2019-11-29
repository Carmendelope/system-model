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
	"time"
)

func createAddEdgeControllerRequest(organizationID string) *grpc_inventory_go.AddEdgeControllerRequest {
	testController := eic.CreateTestEdgeController()

	return &grpc_inventory_go.AddEdgeControllerRequest{
		OrganizationId: organizationID,
		Name:           testController.Name,
		Labels:         testController.Labels,
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

	ginkgo.BeforeEach(func() {
		ginkgo.By("cleaning the mockups", func() {
			organizationProvider.(*orgProvider.MockupOrganizationProvider).Clear()
			controllerProvider.Clear()
			// Initial data
			targetOrganization = testhelpers.CreateOrganization(organizationProvider)
		})
	})

	ginkgo.It("should be able to add a new controller", func() {
		toAdd := createAddEdgeControllerRequest(targetOrganization.ID)
		added, err := client.Add(context.Background(), toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(added).ShouldNot(gomega.BeNil())
		gomega.Expect(added.EdgeControllerId).ShouldNot(gomega.BeEmpty())
	})

	ginkgo.It("should be able to list controllers", func() {
		numControllers := 10
		for index := 0; index < numControllers; index++ {
			toAdd := createAddEdgeControllerRequest(targetOrganization.ID)
			added, err := client.Add(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.EdgeControllerId).ShouldNot(gomega.BeEmpty())
		}
		orgID := &grpc_organization_go.OrganizationId{
			OrganizationId: targetOrganization.ID,
		}
		allControllers, err := client.List(context.Background(), orgID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(len(allControllers.Controllers)).Should(gomega.Equal(numControllers))
	})

	ginkgo.It("should be able to remove controllers", func() {
		toAdd := createAddEdgeControllerRequest(targetOrganization.ID)
		added, err := client.Add(context.Background(), toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(added).ShouldNot(gomega.BeNil())
		edgeControllerID := &grpc_inventory_go.EdgeControllerId{
			OrganizationId:   added.OrganizationId,
			EdgeControllerId: added.EdgeControllerId,
		}
		success, err := client.Remove(context.Background(), edgeControllerID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(success).ShouldNot(gomega.BeNil())
	})

	ginkgo.Context("update operations", func() {
		ginkgo.It("should be able to add new labels", func() {
			toAdd := createAddEdgeControllerRequest(targetOrganization.ID)
			added, err := client.Add(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			newLabels := make(map[string]string, 0)
			newLabels["k1"] = "v1"
			updateRequest := &grpc_inventory_go.UpdateEdgeControllerRequest{
				OrganizationId:   added.OrganizationId,
				EdgeControllerId: added.EdgeControllerId,
				AddLabels:        true,
				RemoveLabels:     false,
				Labels:           newLabels,
			}

			updated, err := client.Update(context.Background(), updateRequest)
			gomega.Expect(err).To(gomega.Succeed())
			value, exits := updated.Labels["k1"]
			gomega.Expect(exits).To(gomega.BeTrue())
			gomega.Expect(value).Should(gomega.Equal("v1"))
		})
		ginkgo.It("should be able to update last operations", func() {
			toAdd := createAddEdgeControllerRequest(targetOrganization.ID)
			added, err := client.Add(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			updateRequest := &grpc_inventory_go.UpdateEdgeControllerRequest{
				OrganizationId:      added.OrganizationId,
				EdgeControllerId:    added.EdgeControllerId,
				AddLabels:           false,
				RemoveLabels:        false,
				UpdateLastOpSummary: true,
				LastOpSummary: &grpc_inventory_go.ECOpSummary{
					OperationId: entities.GenerateUUID(),
					Timestamp:   time.Now().Unix(),
					Status:      grpc_inventory_go.OpStatus_INPROGRESS,
					Info:        "operation summary info",
				},
			}

			updated, err := client.Update(context.Background(), updateRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(updated.LastOpResult).NotTo(gomega.BeNil())
			gomega.Expect(updated.LastOpResult.Info).Should(gomega.Equal("operation summary info"))
		})
	})

})
