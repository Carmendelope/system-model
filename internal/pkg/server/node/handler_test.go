/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package node

import (
	"context"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-utils/pkg/test"
	clusProvider "github.com/nalej/system-model/internal/pkg/provider/cluster"
	nodeProvider "github.com/nalej/system-model/internal/pkg/provider/node"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/nalej/system-model/internal/pkg/server/testhelpers"
	"github.com/onsi/ginkgo"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/gomega"
)

func createCluster(organizationID string, orgProvider orgProvider.Provider, clusProvider clusProvider.Provider) * entities.Cluster {
	toAdd := entities.NewCluster(organizationID, "test cluster", "", "hostname", "hostname")
	err := clusProvider.Add(*toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	err = orgProvider.AddCluster(organizationID, toAdd.ClusterId)
	gomega.Expect(err).To(gomega.Succeed())
	return toAdd
}

func createAddNodeRequest(organizationID string) *grpc_infrastructure_go.AddNodeRequest {
	labels := make(map[string]string, 0)
	labels["k1"] = "v1"
	labels["k2"] = "v2"
	return &grpc_infrastructure_go.AddNodeRequest{
		RequestId:            uuid.NewV4().String(),
		OrganizationId:       organizationID,
		NodeId:               "",
		Ip:                   "127.0.0.1",
		Labels:               labels,
	}
}

var _ = ginkgo.Describe("Node service", func() {
	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client grpc_infrastructure_go.NodesClient

	// Target organization.
	var targetOrganization * entities.Organization
	var targetCluster * entities.Cluster

	// Providers
	var organizationProvider orgProvider.Provider
	var clusterProvider clusProvider.Provider
	var nProvider nodeProvider.Provider

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		// Register the service
		organizationProvider = orgProvider.NewMockupOrganizationProvider()
		clusterProvider = clusProvider.NewMockupClusterProvider()
		nProvider = nodeProvider.NewMockupNodeProvider()
		manager := NewManager(organizationProvider, clusterProvider, nProvider)
		handler := NewHandler(manager)
		grpc_infrastructure_go.RegisterNodesServer(server, handler)

		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_infrastructure_go.NewNodesClient(conn)
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
	})

	ginkgo.BeforeEach(func(){
		ginkgo.By("cleaning the mockups", func(){
			organizationProvider.(*orgProvider.MockupOrganizationProvider).Clear()
			nProvider.(*nodeProvider.MockupNodeProvider).Clear()
			clusterProvider.(*clusProvider.MockupClusterProvider).Clear()
			// Initial data
			targetOrganization = testhelpers.CreateOrganization(organizationProvider)
			targetCluster = createCluster(targetOrganization.ID, organizationProvider, clusterProvider)
		})
	})

	ginkgo.Context("With nodes", func() {
		ginkgo.It("should be able to add a node", func(){
			toAdd := createAddNodeRequest(targetOrganization.ID)
			added, err := client.AddNode(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.NodeId).ShouldNot(gomega.BeEmpty())
			gomega.Expect(added.ClusterId).Should(gomega.BeEmpty())
		})
		ginkgo.It("should fail if the request is not valid", func(){
			toAdd := createAddNodeRequest(targetOrganization.ID)
			toAdd.OrganizationId = ""
			added, err := client.AddNode(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(added).Should(gomega.BeNil())
		})
		ginkgo.FIt("should be able to update a node", func(){
			toAdd := createAddNodeRequest(targetOrganization.ID)
			added, err := client.AddNode(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			updateNodeRequest := &grpc_infrastructure_go.UpdateNodeRequest{
				OrganizationId:       added.OrganizationId,
				NodeId:               added.NodeId,
				UpdateStatus:         true,
				Status:               grpc_infrastructure_go.InfraStatus_RUNNING,
				UpdateState:          true,
				State:                grpc_infrastructure_go.NodeState_ASSIGNED,
			}
			updated, err := client.UpdateNode(context.Background(), updateNodeRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(updated.Status).Should(gomega.Equal(updateNodeRequest.Status))
			gomega.Expect(updated.State).Should(gomega.Equal(updateNodeRequest.State))
		})
		ginkgo.It("should be able to add labels to nodes", func(){
			toAdd := createAddNodeRequest(targetOrganization.ID)
			added, err := client.AddNode(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())

			newLabels := make(map[string]string, 0)
			newLabels["nk"]="nv"
			updateNodeRequest := &grpc_infrastructure_go.UpdateNodeRequest{
				OrganizationId:       added.OrganizationId,
				NodeId:               added.NodeId,
				AddLabels:         true,
				Labels:               newLabels,
			}
			updated, err := client.UpdateNode(context.Background(), updateNodeRequest)
			gomega.Expect(err).To(gomega.Succeed())
			expectedLabels := toAdd.Labels
			expectedLabels["nk"] = "nv"
			gomega.Expect(updated.Labels).Should(gomega.Equal(expectedLabels))
		})
		ginkgo.It("should be able to remove the labels from a node", func(){
			toAdd := createAddNodeRequest(targetOrganization.ID)
			added, err := client.AddNode(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())

			newLabels := make(map[string]string, 0)
			newLabels["k1"]="v1"
			updateNodeRequest := &grpc_infrastructure_go.UpdateNodeRequest{
				OrganizationId:       added.OrganizationId,
				NodeId:               added.NodeId,
				RemoveLabels:         true,
				Labels:               newLabels,
			}
			updated, err := client.UpdateNode(context.Background(), updateNodeRequest)
			gomega.Expect(err).To(gomega.Succeed())
			expectedLabels := toAdd.Labels
			delete(expectedLabels, "k1")
			gomega.Expect(updated.Labels).Should(gomega.Equal(expectedLabels))
		})
		ginkgo.It("should be able to list nodes", func(){
			toAdd := createAddNodeRequest(targetOrganization.ID)
			added, err := client.AddNode(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.NodeId).ShouldNot(gomega.BeEmpty())

			attach := &grpc_infrastructure_go.AttachNodeRequest{
				RequestId:            "req",
				OrganizationId:       targetOrganization.ID,
				ClusterId:            targetCluster.ClusterId,
				NodeId:               added.NodeId,
			}
			success, err := client.AttachNode(context.Background(), attach)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ToNot(gomega.BeNil())

			clusterID := &grpc_infrastructure_go.ClusterId{
				OrganizationId: targetOrganization.ID,
				ClusterId: targetCluster.ClusterId,
			}
			retrieved, err := client.ListNodes(context.Background(), clusterID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			gomega.Expect(len(retrieved.Nodes)).Should(gomega.Equal(1))
		})
		ginkgo.It("should not be able to list nodes on a none existing organization", func(){
			clusterID := &grpc_infrastructure_go.ClusterId{
				OrganizationId: "does not exists",
				ClusterId: targetCluster.ClusterId,
			}
			retrieved, err := client.ListNodes(context.Background(), clusterID)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(retrieved).Should(gomega.BeNil())
		})
		ginkgo.It("should not be able to list nodes on a none existing cluster", func(){
			clusterID := &grpc_infrastructure_go.ClusterId{
				OrganizationId: targetOrganization.ID,
				ClusterId: "does not exists",
			}
			retrieved, err := client.ListNodes(context.Background(), clusterID)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(retrieved).Should(gomega.BeNil())
		})
		ginkgo.It("should return an empty list on an cluster without nodes", func(){
			clusterID := &grpc_infrastructure_go.ClusterId{
				OrganizationId: targetOrganization.ID,
				ClusterId: targetCluster.ClusterId,
			}
			retrieved, err := client.ListNodes(context.Background(), clusterID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			gomega.Expect(len(retrieved.Nodes)).Should(gomega.Equal(0))
		})
		ginkgo.It("should be able to remove an existing cluster", func(){
			toAdd := createAddNodeRequest(targetOrganization.ID)
			added, err := client.AddNode(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.NodeId).ShouldNot(gomega.BeEmpty())
			// Remove nodes
			removeRequest := &grpc_infrastructure_go.RemoveNodesRequest{
				RequestId:            "removeId",
				OrganizationId:       targetOrganization.ID,
				Nodes:            []string{added.NodeId},
			}
			removed, err := client.RemoveNodes(context.Background(), removeRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(removed).ShouldNot(gomega.BeNil())
			// List nodes
			clusterID := &grpc_infrastructure_go.ClusterId{
				OrganizationId: targetOrganization.ID,
				ClusterId: targetCluster.ClusterId,
			}
			retrieved, err := client.ListNodes(context.Background(), clusterID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			gomega.Expect(len(retrieved.Nodes)).Should(gomega.Equal(0))
		})
		ginkgo.It("should not be able to remove an existing cluster", func() {
			toAdd := createAddNodeRequest(targetOrganization.ID)
			added, err := client.AddNode(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.NodeId).ShouldNot(gomega.BeEmpty())

			attach := &grpc_infrastructure_go.AttachNodeRequest{
				RequestId:            "req",
				OrganizationId:       targetOrganization.ID,
				ClusterId:            targetCluster.ClusterId,
				NodeId:               added.NodeId,
			}
			success, err := client.AttachNode(context.Background(), attach)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ToNot(gomega.BeNil())

			err = clusterProvider.DeleteNode(targetCluster.ClusterId, added.NodeId)
			gomega.Expect(err).To(gomega.Succeed())

			// Remove nodes
			removeRequest := &grpc_infrastructure_go.RemoveNodesRequest{
				RequestId:            "removeId",
				OrganizationId:       targetOrganization.ID,
				Nodes:            []string{added.NodeId},
			}
			_, err = client.RemoveNodes(context.Background(), removeRequest)
			gomega.Expect(err).NotTo(gomega.Succeed())


		})
		ginkgo.It("should not be able to remove a none existing cluster", func(){
			// Remove nodes
			removeRequest := &grpc_infrastructure_go.RemoveNodesRequest{
				RequestId:            "removeId",
				OrganizationId:       targetOrganization.ID,
				Nodes:            []string{"does not exists"},
			}
			removed, err := client.RemoveNodes(context.Background(), removeRequest)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(removed).Should(gomega.BeNil())
		})
	})

})

