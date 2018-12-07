package cluster

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest (provider Provider) {

	ginkgo.BeforeEach(func() {
		provider.Clear()
	})

	// AddCluster
	ginkgo.It("Should be able to add a cluster", func(){

		cluster := CreateTestCluster("ZZZ-0")

		err := provider.Add(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

	})

	// UpdateCluster
	ginkgo.It("Should be able to update the cluster", func() {

		cluster := CreateTestCluster("UUUId-0")

		err := provider.Add(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

		cluster.Description = "Modified description"
		cluster.Multitenant = entities.MultitenantSupport(1)

		err = provider.Update(*cluster)
		gomega.Expect(err).To(gomega.Succeed())
	})
	ginkgo.It("Should not be able to update the cluster", func() {

		cluster := CreateTestCluster("UUUId-0")

		err := provider.Update(*cluster)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})

	// GetCluster
	ginkgo.It("Should be able to get the cluster", func(){

		cluster := CreateTestCluster("AAA-0")

		err := provider.Add(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

		cluster, err = provider.Get(cluster.ClusterId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(cluster).NotTo(gomega.BeNil())

	})
	ginkgo.It("Should not be able to get the cluster", func(){

		clusterId := "cluster"

		cluster, err := provider.Get(clusterId)
		gomega.Expect(err).NotTo(gomega.Succeed())
		gomega.Expect(cluster).To(gomega.BeNil())

	})

	// ExistsCluster
	ginkgo.It("Should be able to find the cluster", func(){

		cluster := CreateTestCluster("AAA-0")

		err := provider.Add(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

		exists, err := provider.Exists(cluster.ClusterId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())

	})
	ginkgo.It("Should not be able to find the cluster", func(){

		clusterId := "cluster"

		cluster, err := provider.Exists(clusterId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(cluster).NotTo(gomega.BeTrue())

	})

	// DeleteCluster
	ginkgo.It("Should be able to delete the cluster", func(){

		cluster := CreateTestCluster("AAA-0")

		err := provider.Add(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.Remove(cluster.ClusterId)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to delete the cluster", func(){

		clusterId := "cluster"

		err := provider.Remove(clusterId)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// -------------------------------------------------------------------------------------------

	// AddNode
	ginkgo.It("Should be able to add the node in the cluster", func() {

		// add the cluster
		cluster := CreateTestCluster("0001")
		err := provider.Add(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

		// add the cluster-node
		err = provider.AddNode("cluster_0001", "node_0001")
		gomega.Expect(err).To(gomega.Succeed())

	})

	ginkgo.It("Should not be able to add the node in the cluster", func() {

		err := provider.AddNode("cluster_0001", "node_0001")
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// NodeExists
	ginkgo.It("Should be able to find the node of the cluster", func() {

		clusterID := "cluster_0001"
		nodeID := "node_0001"

		// add the cluster
		cluster := CreateTestCluster("0001")
		err := provider.Add(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

		// add the cluster-node
		err = provider.AddNode(clusterID, nodeID)
		gomega.Expect(err).To(gomega.Succeed())

		exists, err := provider.NodeExists(clusterID, nodeID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())

	})
	ginkgo.It("Should not be able to find the node of the cluster", func() {

		exists, err := provider.NodeExists("cluster_X", "node_0")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).NotTo(gomega.BeTrue())

	})

	// ListNodes
	ginkgo.It("Should be able to return a list of nodes", func() {

		// add the cluster
		cluster := CreateTestCluster("0001")
		err := provider.Add(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

		// add the nodes in the cluster
		clusterID := "cluster_0001"
		for i:=0; i<10; i++ {
			nodeID := fmt.Sprintf("Node_00%d", i)
			err := provider.AddNode(clusterID, nodeID)
			gomega.Expect(err).To(gomega.Succeed())
		}

		list, err := provider.ListNodes(clusterID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(list).NotTo(gomega.BeNil())
		gomega.Expect(list).NotTo(gomega.BeEmpty())


	})
	ginkgo.It("Should not be able to return a list of nodes (no cluster found)", func() {

		clusterID := "cluster_000"

		_, err := provider.ListNodes(clusterID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})
	ginkgo.It("Should not be able to return a list of nodes", func() {

		// add the cluster
		cluster := CreateTestCluster("0001")
		err := provider.Add(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

		list, err := provider.ListNodes("cluster_0001")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(list).NotTo(gomega.BeNil())
		gomega.Expect(list).To(gomega.BeEmpty())

	})

	// DeleteNode
	ginkgo.It("Should be able to delete a Node in a cluster", func() {

		cluster := CreateTestCluster("0001")

		err := provider.Add(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

		nodeID := "node0001"
		err = provider.AddNode(cluster.ClusterId, nodeID)
		err = provider.AddNode(cluster.ClusterId, "node0002")
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.DeleteNode(cluster.ClusterId, nodeID)
		gomega.Expect(err).To(gomega.Succeed())
	})
	ginkgo.It("Should not be able to delete a Node in a cluster", func() {

		clusterID := "clusterID"
		nodeID := "nodeID"

		err := provider.DeleteNode(clusterID, nodeID)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})
}