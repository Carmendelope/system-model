package node

import (
	"os"

	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {

	labels := make(map[string]string)
	labels["label1"] = "label1"

	ginkgo.BeforeEach(func() {
		var clearProvider = os.Getenv("IT_CLEAR_PROVIDER")
		if clearProvider == "true" {
			provider.Clear()
		}
	})

	// Add
	ginkgo.It("Should be able to add node", func() {

		// TODO: This entity creation should be in a helper method rather than here
		node := &entities.Node{OrganizationId: "org", ClusterId: "cluster_id", NodeId: "node",
			Ip: "0.0.0.0", Labels: labels, Status: entities.InfraStatusRunning, State: 0}

		err := provider.Add(*node)
		gomega.Expect(err).To(gomega.Succeed())

		_ = provider.Remove(node.NodeId)
	})

	// Update
	ginkgo.It("Should be able to update node", func() {

		// add a node
		node := &entities.Node{OrganizationId: "org", ClusterId: "cluster_id", NodeId: "node",
			Ip: "0.0.0.0", Labels: labels, Status: entities.InfraStatusRunning, State: 0}

		err := provider.Add(*node)
		gomega.Expect(err).To(gomega.Succeed())

		// uodate it
		labels["label2"] = "label2"
		node.OrganizationId = "org_MOD"
		node.Labels = labels

		err = provider.Update(*node)
		gomega.Expect(err).To(gomega.Succeed())

		_ = provider.Remove(node.NodeId)
	})
	ginkgo.It("Should not be able to update node", func() {

		node := &entities.Node{OrganizationId: "org", ClusterId: "clusterMODD", NodeId: "node",
			Ip: "127.0.0.1", Labels: labels, Status: entities.InfraStatusInstalling, State: 1}

		err := provider.Update(*node)
		gomega.Expect(err).NotTo(gomega.Succeed())

		_ = provider.Remove(node.NodeId)
	})

	// Exists
	ginkgo.It("Should be able to find node", func() {

		// add a node
		node := &entities.Node{OrganizationId: "org", ClusterId: "cluster_id", NodeId: "node",
			Ip: "0.0.0.0", Labels: labels, Status: entities.InfraStatusRunning, State: 0}

		err := provider.Add(*node)
		gomega.Expect(err).To(gomega.Succeed())

		// ask if it exists
		exits, err := provider.Exists(node.NodeId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exits).To(gomega.BeTrue())

		_ = provider.Remove(node.NodeId)
	})

	ginkgo.It("Should not be able to find node", func() {

		exits, err := provider.Exists("node")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exits).NotTo(gomega.BeTrue())

	})

	// Get
	ginkgo.It("Should be able to get the node", func() {

		// add a node
		node := &entities.Node{OrganizationId: "org", ClusterId: "cluster_id", NodeId: "node",
			Ip: "0.0.0.0", Labels: labels, Status: entities.InfraStatusRunning, State: 0}

		err := provider.Add(*node)
		gomega.Expect(err).To(gomega.Succeed())

		// ask for it
		node, err = provider.Get(node.NodeId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(node).NotTo(gomega.BeNil())

		_ = provider.Remove(node.NodeId)
	})
	ginkgo.It("Should not be able to get the node", func() {

		node, err := provider.Get("node")
		gomega.Expect(err).NotTo(gomega.Succeed())
		gomega.Expect(node).To(gomega.BeNil())
	})

	// Remove
	ginkgo.It("Should be able to remove the node", func() {

		// add a node
		node := &entities.Node{OrganizationId: "org", ClusterId: "cluster_id", NodeId: "node",
			Ip: "0.0.0.0", Labels: labels, Status: entities.InfraStatusRunning, State: 0}

		err := provider.Add(*node)
		gomega.Expect(err).To(gomega.Succeed())

		// remove it
		err = provider.Remove(node.NodeId)
		gomega.Expect(err).To(gomega.Succeed())
	})
	ginkgo.It("Should not be able to remove the node", func() {

		err := provider.Remove("node")
		gomega.Expect(err).NotTo(gomega.Succeed())
	})
}
