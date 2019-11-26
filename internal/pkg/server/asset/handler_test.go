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

package asset

import (
	"context"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/asset"
	assetProvider "github.com/nalej/system-model/internal/pkg/provider/asset"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/nalej/system-model/internal/pkg/server/testhelpers"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func createAddAssetRequest(organizationID string) *grpc_inventory_go.AddAssetRequest {
	testAsset := asset.CreateTestAsset()

	storage := make([]*grpc_inventory_go.StorageHardwareInfo, 0)
	for _, sto := range testAsset.Storage {
		storage = append(storage, sto.ToGRPC())
	}

	return &grpc_inventory_go.AddAssetRequest{
		OrganizationId:   organizationID,
		EdgeControllerId: testAsset.EdgeControllerId,
		AgentId:          testAsset.AgentId,
		Labels:           testAsset.Labels,
		Os:               testAsset.Os.ToGRPC(),
		Hardware:         testAsset.Hardware.ToGRPC(),
		Storage:          storage,
	}
}

var _ = ginkgo.Describe("Asset service", func() {
	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client grpc_inventory_go.AssetsClient

	// Target organization.
	var targetOrganization *entities.Organization

	// Providers
	var organizationProvider orgProvider.Provider
	var aProvider assetProvider.Provider

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		// Register the service
		organizationProvider = orgProvider.NewMockupOrganizationProvider()
		aProvider = assetProvider.NewMockupAssetProvider()
		manager := NewManager(organizationProvider, aProvider)
		handler := NewHandler(manager)
		grpc_inventory_go.RegisterAssetsServer(server, handler)

		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_inventory_go.NewAssetsClient(conn)
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
	})

	ginkgo.BeforeEach(func() {
		ginkgo.By("cleaning the mockups", func() {
			organizationProvider.(*orgProvider.MockupOrganizationProvider).Clear()
			aProvider.Clear()
			// Initial data
			targetOrganization = testhelpers.CreateOrganization(organizationProvider)
		})
	})

	ginkgo.It("should be able to add a new asset", func() {
		toAdd := createAddAssetRequest(targetOrganization.ID)
		added, err := client.Add(context.Background(), toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(added).ShouldNot(gomega.BeNil())
		gomega.Expect(added.AssetId).ShouldNot(gomega.BeEmpty())

		assetID := &grpc_inventory_go.AssetId{
			OrganizationId: added.OrganizationId,
			AssetId:        added.AssetId,
		}
		retrieved, err := client.Get(context.Background(), assetID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(retrieved.EdgeControllerId).Should(gomega.Equal(toAdd.EdgeControllerId))
	})

	ginkgo.It("should be able to list assets", func() {
		numAssets := 10
		for index := 0; index < numAssets; index++ {
			toAdd := createAddAssetRequest(targetOrganization.ID)
			added, err := client.Add(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.AssetId).ShouldNot(gomega.BeEmpty())
		}
		orgID := &grpc_organization_go.OrganizationId{
			OrganizationId: targetOrganization.ID,
		}
		allAssets, err := client.List(context.Background(), orgID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(len(allAssets.Assets)).Should(gomega.Equal(numAssets))
	})

	ginkgo.It("should be able to remove assets", func() {
		toAdd := createAddAssetRequest(targetOrganization.ID)
		added, err := client.Add(context.Background(), toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(added).ShouldNot(gomega.BeNil())
		assetID := &grpc_inventory_go.AssetId{
			OrganizationId: added.OrganizationId,
			AssetId:        added.AssetId,
		}
		success, err := client.Remove(context.Background(), assetID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(success).ShouldNot(gomega.BeNil())
	})

	ginkgo.Context("update operations", func() {
		ginkgo.It("should be able to add new labels", func() {
			toAdd := createAddAssetRequest(targetOrganization.ID)
			added, err := client.Add(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			newLabels := make(map[string]string, 0)
			newLabels["k1"] = "v1"
			updateRequest := &grpc_inventory_go.UpdateAssetRequest{
				OrganizationId: added.OrganizationId,
				AssetId:        added.AssetId,
				AddLabels:      true,
				RemoveLabels:   false,
				Labels:         newLabels,
			}

			updated, err := client.Update(context.Background(), updateRequest)
			gomega.Expect(err).To(gomega.Succeed())
			value, exits := updated.Labels["k1"]
			gomega.Expect(exits).To(gomega.BeTrue())
			gomega.Expect(value).Should(gomega.Equal("v1"))
		})
	})

})
