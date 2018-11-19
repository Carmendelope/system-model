package node

import (
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest (provider Provider) {

	var nodeOK = "nodeOK"
	var nodeKO = "nodeKO"

	labels := make(map[string]string)
	labels["label1"] = "label1"

	// AddProvider
	ginkgo.It("Should be able to add role", func(){

		node := &entities.Node{OrganizationId:"org", ClusterId:"cluster_id", NodeId: nodeOK,
		Ip:"0.0.0.0", Labels:labels, Status:entities.InfraStatusRunning, State:0}

		err := provider.Add(*node)
		gomega.Expect(err).To(gomega.Succeed())

	})

	// Update
	ginkgo.It("Should be able to update role", func(){

		labels["label2"] = "label2"
		node := &entities.Node{OrganizationId:"org", ClusterId:"clusterMODd", NodeId: nodeOK,
			Ip:"127.0.0.1", Labels:labels, Status:entities.InfraStatusInstalling, State:1}

		err := provider.Update(*node)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to update role", func(){

		node := &entities.Node{OrganizationId:"org", ClusterId:"clusterMODD", NodeId: nodeKO,
			Ip:"127.0.0.1", Labels:labels, Status:entities.InfraStatusInstalling, State:1}

		err := provider.Update(*node)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// Exists
	ginkgo.It("Should be able to find role", func(){

		exits, err := provider.Exists(nodeOK)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exits).To(gomega.BeTrue())

	})

	ginkgo.It("Should not be able to find role", func(){

		exits, err := provider.Exists(nodeKO)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exits).NotTo(gomega.BeTrue())

	})

	// Get
	ginkgo.It("Should be able to get the role", func(){

		node, err := provider.Get(nodeOK)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(node).NotTo(gomega.BeNil())
	})
	ginkgo.It("Should not be able to get the role", func(){

		node, err := provider.Get(nodeKO)
		gomega.Expect(err).NotTo(gomega.Succeed())
		gomega.Expect(node).To(gomega.BeNil())
	})

	// Remove
	ginkgo.It("Should be able to find the role", func(){

		err := provider.Remove(nodeOK)
		gomega.Expect(err).To(gomega.Succeed())
	})
	ginkgo.It("Should not be able to find the role", func(){

		err := provider.Remove(nodeKO)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})
}
