/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package cluster

import (
	"context"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/test"
	clusProvider "github.com/nalej/system-model/internal/pkg/provider/cluster"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/onsi/ginkgo"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/gomega"
)

func createOrganization(orgProvider orgProvider.Provider) * entities.Organization {
	toAdd := entities.NewOrganization("test org")
	err := orgProvider.Add(*toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	return toAdd
}

func createAddClusterRequest(organizationID string) *grpc_infrastructure_go.AddClusterRequest {
	labels := make(map[string]string, 0)
	labels["k1"] = "v1"
	labels["k2"] = "v2"
	return &grpc_infrastructure_go.AddClusterRequest{
		RequestId:            uuid.NewV4().String(),
		OrganizationId:       organizationID,
		Name:                 "name",
		Description:          "description",
		Hostname: "hostname",
		Labels:               labels,
	}
}

var _ = ginkgo.Describe("Cluster service", func() {
	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client grpc_infrastructure_go.ClustersClient

	// Target organization.
	var targetOrganization * entities.Organization

	// Providers
	var organizationProvider orgProvider.Provider
	var clusterProvider clusProvider.Provider

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()
		test.LaunchServer(server, listener)

		// Register the service
		organizationProvider = orgProvider.NewMockupOrganizationProvider()
		clusterProvider = clusProvider.NewMockupClusterProvider()
		manager := NewManager(organizationProvider, clusterProvider)
		handler := NewHandler(manager)
		grpc_infrastructure_go.RegisterClustersServer(server, handler)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_infrastructure_go.NewClustersClient(conn)
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
	})

	ginkgo.BeforeEach(func(){
		ginkgo.By("cleaning the mockups", func(){
			organizationProvider.(*orgProvider.MockupOrganizationProvider).Clear()
			clusterProvider.(*clusProvider.MockupClusterProvider).Clear()
			// Initial data
			targetOrganization = createOrganization(organizationProvider)
		})
	})

	ginkgo.Context("With clusters", func() {
		ginkgo.It("should be able to add a cluster", func(){
			toAdd := createAddClusterRequest(targetOrganization.ID)
			added, err := client.AddCluster(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.ClusterId).ShouldNot(gomega.BeEmpty())
		})
		ginkgo.It("should fail if the request is not valid", func(){
			toAdd := createAddClusterRequest(targetOrganization.ID)
			toAdd.OrganizationId = ""
			added, err := client.AddCluster(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(added).Should(gomega.BeNil())
		})
		ginkgo.It("should be able to get an existing cluster", func(){
			toAdd := createAddClusterRequest(targetOrganization.ID)
			added, err := client.AddCluster(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.ClusterId).ShouldNot(gomega.BeEmpty())

			clusterID := &grpc_infrastructure_go.ClusterId{
				OrganizationId:       added.OrganizationId,
				ClusterId:            added.ClusterId,
			}
			retrieved, err := client.GetCluster(context.Background(), clusterID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			gomega.Expect(retrieved).Should(gomega.Equal(added))
		})
		ginkgo.It("should be able to update a cluster", func(){
			toAdd := createAddClusterRequest(targetOrganization.ID)
			added, err := client.AddCluster(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.ClusterId).ShouldNot(gomega.BeEmpty())

			newLabels := make(map[string]string, 0)
			newLabels["nk"]="nv"
			updateClusterReq := &grpc_infrastructure_go.UpdateClusterRequest{
				OrganizationId:       targetOrganization.ID,
				ClusterId:            added.ClusterId,
				UpdateName:           true,
				Name:                 "newName",
				UpdateDescription:    true,
				Description:          "newDescription",
				UpdateHostname:       true,
				Hostname:             "newHostname",
				UpdateLabels:         true,
				Labels:               newLabels,
				UpdateStatus:         true,
				Status:               grpc_infrastructure_go.InfraStatus_RUNNING,
			}
			updated, err := client.UpdateCluster(context.Background(), updateClusterReq)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(updated.Name).Should(gomega.Equal(updateClusterReq.Name))
			gomega.Expect(updated.Description).Should(gomega.Equal(updateClusterReq.Description))
			gomega.Expect(updated.Hostname).Should(gomega.Equal(updateClusterReq.Hostname))
			gomega.Expect(updated.Labels).Should(gomega.Equal(updateClusterReq.Labels))
			gomega.Expect(updated.Status).Should(gomega.Equal(updateClusterReq.Status))
		})
		ginkgo.It("should be able to list clusters", func(){
			toAdd := createAddClusterRequest(targetOrganization.ID)
			added, err := client.AddCluster(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.ClusterId).ShouldNot(gomega.BeEmpty())

			organizationID := &grpc_organization_go.OrganizationId{
				OrganizationId: targetOrganization.ID,
			}
			retrieved, err := client.ListClusters(context.Background(), organizationID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			gomega.Expect(len(retrieved.Clusters)).Should(gomega.Equal(1))
		})
		ginkgo.It("should not be able to list clusters on a none existing organization", func(){
			organizationID := &grpc_organization_go.OrganizationId{
				OrganizationId: "does not exists",
			}
			retrieved, err := client.ListClusters(context.Background(), organizationID)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(retrieved).Should(gomega.BeNil())
		})
		ginkgo.It("should return an empty list on an organization without clusters", func(){
			organizationID := &grpc_organization_go.OrganizationId{
				OrganizationId: targetOrganization.ID,
			}
			retrieved, err := client.ListClusters(context.Background(), organizationID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			gomega.Expect(len(retrieved.Clusters)).Should(gomega.Equal(0))
		})
		ginkgo.It("should be able to remove an existing cluster", func(){
			toAdd := createAddClusterRequest(targetOrganization.ID)
			added, err := client.AddCluster(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.ClusterId).ShouldNot(gomega.BeEmpty())
			// Remove cluster
			removeRequest := &grpc_infrastructure_go.RemoveClusterRequest{
				RequestId:            "removeId",
				OrganizationId:       targetOrganization.ID,
				ClusterId:            added.ClusterId,
			}
			removed, err := client.RemoveCluster(context.Background(), removeRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(removed).ShouldNot(gomega.BeNil())
			// List clusters
			organizationID := &grpc_organization_go.OrganizationId{
				OrganizationId: targetOrganization.ID,
			}
			retrieved, err := client.ListClusters(context.Background(), organizationID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			gomega.Expect(len(retrieved.Clusters)).Should(gomega.Equal(0))
		})
		ginkgo.It("should not be able to remove a none existing cluster", func(){
			// Remove cluster
			removeRequest := &grpc_infrastructure_go.RemoveClusterRequest{
				RequestId:            "removeId",
				OrganizationId:       targetOrganization.ID,
				ClusterId:            "does not exists",
			}
			removed, err := client.RemoveCluster(context.Background(), removeRequest)
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(removed).Should(gomega.BeNil())
		})
	})

	ginkgo.PContext("With nodes", func() {

	})

})

