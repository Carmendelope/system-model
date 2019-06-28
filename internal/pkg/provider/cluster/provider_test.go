package cluster

import (
	"fmt"
	"os"

	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {

	ginkgo.BeforeEach(func() {
		var clearProvider = os.Getenv("IT_CLEAR_PROVIDER")
		if clearProvider == "true" {
			provider.Clear()
		}
	})

	// AddCluster
	ginkgo.It("Should be able to add a cluster", func() {

		cluster := CreateTestCluster("ZZZ-0")

		err := provider.Add(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

		_ = provider.Remove(cluster.ClusterId)

	})

	// UpdateCluster
	ginkgo.It("Should be able to update the cluster", func() {

		cluster := CreateTestCluster("UUUId-0")

		err := provider.Add(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

		cluster.Multitenant = entities.MultitenantSupport(1)

		err = provider.Update(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

		_ = provider.Remove(cluster.ClusterId)
	})
	ginkgo.It("Should not be able to update the cluster", func() {

		cluster := CreateTestCluster("UUUId-0")

		err := provider.Update(*cluster)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})

	// GetCluster
	ginkgo.It("Should be able to get the cluster", func() {

		cluster := CreateTestCluster("AAA-0")

		err := provider.Add(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

		cluster, err = provider.Get(cluster.ClusterId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(cluster).NotTo(gomega.BeNil())

		//NP-396 check that ControlPaneHostname exists.
		gomega.Expect(cluster.ControlPlaneHostname).NotTo(gomega.BeNil())
		gomega.Expect(cluster.ControlPlaneHostname).Should(gomega.Equal("cp_host_AAA-0"))

		_ = provider.Remove(cluster.ClusterId)

	})
	ginkgo.It("Should not be able to get the cluster", func() {

		clusterID := "cluster"

		cluster, err := provider.Get(clusterID)
		gomega.Expect(err).NotTo(gomega.Succeed())
		gomega.Expect(cluster).To(gomega.BeNil())

	})

	// ExistsCluster
	ginkgo.It("Should be able to find the cluster", func() {

		cluster := CreateTestCluster("AAA-0")

		err := provider.Add(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

		exists, err := provider.Exists(cluster.ClusterId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())

		_ = provider.Remove(cluster.ClusterId)

	})
	ginkgo.It("Should not be able to find the cluster", func() {

		clusterID := "cluster"

		cluster, err := provider.Exists(clusterID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(cluster).NotTo(gomega.BeTrue())

	})

	// DeleteCluster
	ginkgo.It("Should be able to delete the cluster", func() {

		cluster := CreateTestCluster("AAA-0")

		err := provider.Add(*cluster)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.Remove(cluster.ClusterId)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to delete the cluster", func() {

		clusterID := "cluster"

		err := provider.Remove(clusterID)
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
		clusterID := "cluster_0001"
		nodeID := "node_0001"
		err = provider.AddNode(clusterID, nodeID)
		gomega.Expect(err).To(gomega.Succeed())

		_ = provider.DeleteNode(clusterID, nodeID)
		_ = provider.Remove(cluster.ClusterId)

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
		for i := 0; i < 10; i++ {
			nodeID := fmt.Sprintf("Node_00%d", i)
			err := provider.AddNode(clusterID, nodeID)
			gomega.Expect(err).To(gomega.Succeed())
		}

		list, err := provider.ListNodes(clusterID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(list).NotTo(gomega.BeNil())
		gomega.Expect(list).NotTo(gomega.BeEmpty())

		for i := 0; i < 10; i++ {
			nodeID := fmt.Sprintf("Node_00%d", i)
			_ = provider.DeleteNode(clusterID, nodeID)
		}
		_ = provider.Remove(cluster.ClusterId)

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

		_ = provider.Remove(cluster.ClusterId)
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

		_ = provider.DeleteNode(cluster.ClusterId, "node0002")
		_ = provider.Remove(cluster.ClusterId)
	})
	ginkgo.It("Should not be able to delete a Node in a cluster", func() {

		clusterID := "clusterID"
		nodeID := "nodeID"

		err := provider.DeleteNode(clusterID, nodeID)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})
}
