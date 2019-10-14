package organization

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {

	ginkgo.BeforeEach(func() {
		provider.Clear()
	})

	// Add and organization
	ginkgo.It("Should be able to add a organization", func() {

		org := &entities.Organization{ID: "Org_0001", Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

	})

	// Get Organization
	ginkgo.It("Should be able to get a organization", func() {

		org := &entities.Organization{ID: "Org_0001", Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		org, err = provider.Get(org.ID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(org).NotTo(gomega.BeNil())

	})

	// List Organization
	ginkgo.It("Should be able to list a organization", func() {

		org := &entities.Organization{ID: "Org_0001", Name: "organization 0001", Created: 12}
		org1 := &entities.Organization{ID: "Org_0002", Name: "organization 0002", Created: 13}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.Add(*org1)
		gomega.Expect(err).To(gomega.Succeed())

		orgLst, err := provider.List()
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(orgLst).NotTo(gomega.BeNil())
		gomega.Expect(orgLst).Should(gomega.HaveLen(2))

	})

	// List Organization
	ginkgo.It("Should be able to recover a empty list of organizations", func() {

		orgLst, err := provider.List()
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(orgLst).NotTo(gomega.BeNil())
		gomega.Expect(orgLst).Should(gomega.HaveLen(0))

	})

	ginkgo.It("Should not be able to get a organization", func() {

		_, err := provider.Get("Org_0001")
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// Exists Organization
	ginkgo.It("Should be able to find a organization", func() {

		org := &entities.Organization{ID: "Org_0001", Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		exists, err := provider.Exists(org.ID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())

	})
	ginkgo.It("Should not be able to find a organization", func() {

		exists, err := provider.Exists("Org_0001")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).NotTo(gomega.BeTrue())
	})
	ginkgo.It("Should be able to find a organization by name", func() {

		org := &entities.Organization{ID: "Org_0001", Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		exists, err := provider.ExistsByName(org.Name)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())

	})
	ginkgo.It("Should not be able to find a organization by name", func() {

		exists, err := provider.ExistsByName("Org_0001")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).NotTo(gomega.BeTrue())
	})

	// --------------------------------------------------------------------------------------------------------------------

	// AddCluster
	ginkgo.It("Should be able to add a cluster in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddCluster(organizationID, "cluster001")
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to add a cluster in a organization", func() {

		organizationID := "Org_0001"

		err := provider.AddCluster(organizationID, "cluster001")
		gomega.Expect(err).NotTo(gomega.Succeed())

	})
	ginkgo.It("Should not be able to add a cluster in a organization (already exists)", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddCluster(organizationID, "cluster001")
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddCluster(organizationID, "cluster001")
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// ClusterExists
	ginkgo.It("Should be able to find a cluster in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddCluster(organizationID, "cluster001")
		gomega.Expect(err).To(gomega.Succeed())

		exists, err := provider.ClusterExists(organizationID, "cluster001")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())

	})
	ginkgo.It("Should not be able to find a cluster in a organization", func() {

		organizationID := "Org_0001"

		exists, err := provider.ClusterExists(organizationID, "cluster001")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).NotTo(gomega.BeTrue())

	})

	// ListClusters
	ginkgo.It("Should be able to get a list of the cluster in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())
		for i := 1; i <= 10; i++ {
			err = provider.AddCluster(organizationID, fmt.Sprintf("cluster00%d", i))
			gomega.Expect(err).To(gomega.Succeed())
		}

		clusters, err := provider.ListClusters(organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(clusters).NotTo(gomega.BeEmpty())

	})
	ginkgo.It("Should be able to get an empty list of the cluster in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		clusters, err := provider.ListClusters(organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(clusters).To(gomega.BeEmpty())

	})
	ginkgo.It("Should not be able to get a list of the cluster in a organization", func() {

		organizationID := "Org_0001"

		_, err := provider.ListClusters(organizationID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// DeleteCluster
	ginkgo.It("Should be able to delete a cluster in a organization", func() {

		organizationID := "Org_0001"
		clusterID := "Cluster_001"
		org := &entities.Organization{ID: organizationID, Name: "organization_0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddCluster(organizationID, clusterID)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.DeleteCluster(organizationID, clusterID)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to delete a cluster in a organization", func() {

		organizationID := "Org_0001"
		clusterID := "Cluster_001"

		err := provider.DeleteCluster(organizationID, clusterID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// --------------------------------------------------------------------------------------------------------------------

	// AddNodes
	ginkgo.It("Should be able to add a node in a organization", func() {

		organizationID := "org_XX01"
		org := &entities.Organization{ID: organizationID, Name: "organization OrgXX01", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddNode(organizationID, "node_001")
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to add a node in a organization", func() {

		organizationID := "OrgXX01"

		err := provider.AddNode(organizationID, "node")
		gomega.Expect(err).NotTo(gomega.Succeed())

	})
	ginkgo.It("Should not be able to add a node in a organization (already exists)", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddNode(organizationID, "Node001")
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddNode(organizationID, "Node001")
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// NodeExists
	ginkgo.It("Should be able to find a node in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddCluster(organizationID, "cluster001")
		gomega.Expect(err).To(gomega.Succeed())

		exists, err := provider.ClusterExists(organizationID, "cluster001")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())

	})
	ginkgo.It("Should not be able to find a node in a organization", func() {

		organizationID := "Org_0001"

		exists, err := provider.ClusterExists(organizationID, "cluster001")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).NotTo(gomega.BeTrue())

	})

	// ListNodes
	ginkgo.It("Should be able to get a list of the node in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())
		for i := 1; i <= 10; i++ {
			err = provider.AddNode(organizationID, fmt.Sprintf("node_00%d", i))
			gomega.Expect(err).To(gomega.Succeed())
		}

		clusters, err := provider.ListNodes(organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(clusters).NotTo(gomega.BeEmpty())

	})
	ginkgo.It("Should be able to get an empty list of the node in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		clusters, err := provider.ListNodes(organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(clusters).To(gomega.BeEmpty())

	})
	ginkgo.It("Should not be able to get a list of the node in a organization", func() {

		organizationID := "Org0001"

		_, err := provider.ListNodes(organizationID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// DeleteNode
	ginkgo.It("Should be able to delete a node in a organization", func() {

		organizationID := "Org_0001"
		nodeID := "Node_X01"
		org := &entities.Organization{ID: organizationID, Name: "organization_0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddNode(organizationID, nodeID)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.DeleteNode(organizationID, nodeID)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to delete a node in a organization", func() {

		organizationID := "Org_0001"
		nodeID := "node_01"

		err := provider.DeleteNode(organizationID, nodeID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// --------------------------------------------------------------------------------------------------------------------

	// AddAppDescriptors
	ginkgo.It("Should be able to add a descriptor in a organization", func() {

		organizationID := "org_XX01"
		org := &entities.Organization{ID: organizationID, Name: "organization OrgXX01", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddDescriptor(organizationID, "app_descriptor_01")
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to add a descriptor in a organization", func() {

		organizationID := "organization_id"

		err := provider.AddDescriptor(organizationID, "descriptor")
		gomega.Expect(err).NotTo(gomega.Succeed())

	})
	ginkgo.It("Should not be able to add a descriptor in a organization (already exists)", func() {

		organizationID := "organization_id"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddDescriptor(organizationID, "app_descriptor_01")
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddDescriptor(organizationID, "app_descriptor_01")
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// AppDescriptorsExists
	ginkgo.It("Should be able to find a descriptor in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddDescriptor(organizationID, "app_descriptor_01")
		gomega.Expect(err).To(gomega.Succeed())

		exists, err := provider.DescriptorExists(organizationID, "app_descriptor_01")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())

	})
	ginkgo.It("Should not be able to find a descriptor in a organization", func() {

		organizationID := "Org_0001"

		exists, err := provider.DescriptorExists(organizationID, "app_descriptor_01")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).NotTo(gomega.BeTrue())

	})

	// ListDescriptors
	ginkgo.It("Should be able to get a list of the descriptors in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())
		for i := 1; i <= 10; i++ {
			err = provider.AddDescriptor(organizationID, fmt.Sprintf("app_descriptor_%d", i))
			gomega.Expect(err).To(gomega.Succeed())
		}

		clusters, err := provider.ListDescriptors(organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(clusters).NotTo(gomega.BeEmpty())

	})
	ginkgo.It("Should be able to get an empty list of the descriptors in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		clusters, err := provider.ListDescriptors(organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(clusters).To(gomega.BeEmpty())

	})
	ginkgo.It("Should not be able to get a list of the descriptors in a organization", func() {

		organizationID := "Org0001"

		_, err := provider.ListDescriptors(organizationID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// DeleteDescriptors
	ginkgo.It("Should be able to delete a descriptor in a organization", func() {

		organizationID := "Org_0001"
		descriptorID := "app_descriptor_01"
		org := &entities.Organization{ID: organizationID, Name: "organization_0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddDescriptor(organizationID, descriptorID)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.DeleteDescriptor(organizationID, descriptorID)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to delete a descriptor in a organization", func() {

		organizationID := "Org_0001"
		descriptorID := "app_descriptor_01"

		err := provider.DeleteDescriptor(organizationID, descriptorID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// --------------------------------------------------------------------------------------------------------------------

	// AddAppInstance
	ginkgo.It("Should be able to add an instance in a organization", func() {

		organizationID := "org_XX01"
		org := &entities.Organization{ID: organizationID, Name: "organization OrgXX01", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddInstance(organizationID, "app_instance_01")
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to add an instance  in a organization", func() {

		organizationID := "organization_id"

		err := provider.AddInstance(organizationID, "instance")
		gomega.Expect(err).NotTo(gomega.Succeed())

	})
	ginkgo.It("Should not be able to add an instance  in a organization (already exists)", func() {

		organizationID := "organization_id"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddInstance(organizationID, "app_instance_01")
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddInstance(organizationID, "app_instance_01")
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// AppInstanceExists
	ginkgo.It("Should be able to find an instance  in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddInstance(organizationID, "app_instance_01")
		gomega.Expect(err).To(gomega.Succeed())

		exists, err := provider.InstanceExists(organizationID, "app_instance_01")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())

	})
	ginkgo.It("Should not be able to find an instance  in a organization", func() {

		organizationID := "Org_0001"

		exists, err := provider.InstanceExists(organizationID, "app_descriptor_01")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).NotTo(gomega.BeTrue())

	})

	// ListInstances
	ginkgo.It("Should be able to get a list of the instances in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())
		for i := 1; i <= 10; i++ {
			err = provider.AddInstance(organizationID, fmt.Sprintf("app_instance_%d", i))
			gomega.Expect(err).To(gomega.Succeed())
		}

		clusters, err := provider.ListInstances(organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(clusters).NotTo(gomega.BeEmpty())

	})
	ginkgo.It("Should be able to get an empty list of the instances in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		clusters, err := provider.ListInstances(organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(clusters).To(gomega.BeEmpty())

	})
	ginkgo.It("Should not be able to get a list of the instances in a organization", func() {

		organizationID := "Org0001"

		_, err := provider.ListInstances(organizationID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// DeleteInstance
	ginkgo.It("Should be able to delete an instance in a organization", func() {

		organizationID := "Org_0001"
		instanceID := "app_instance_01"
		org := &entities.Organization{ID: organizationID, Name: "organization_0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddInstance(organizationID, instanceID)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.DeleteInstance(organizationID, instanceID)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to delete an instance in a organization", func() {

		organizationID := "Org_0001"
		instanceID := "app_instance_01"

		err := provider.DeleteInstance(organizationID, instanceID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// --------------------------------------------------------------------------------------------------------------------

	// AddUser
	ginkgo.It("Should be able to add a new user in a organization", func() {

		organizationID := "org_XX01"
		org := &entities.Organization{ID: organizationID, Name: "organization OrgXX01", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddUser(organizationID, "email_1@daisho.group")
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to add a new user  in a organization", func() {

		organizationID := "organization_id"

		err := provider.AddInstance(organizationID, "email_1@daisho.group")
		gomega.Expect(err).NotTo(gomega.Succeed())

	})
	ginkgo.It("Should not be able to add a new user in a organization (already exists)", func() {

		organizationID := "organization_id"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddUser(organizationID, "email_1@daisho.group")
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddUser(organizationID, "email_1@daisho.group")
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// UserExists
	ginkgo.It("Should be able to find a user  in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddUser(organizationID, "email_1@daisho.group")
		gomega.Expect(err).To(gomega.Succeed())

		exists, err := provider.UserExists(organizationID, "email_1@daisho.group")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())

	})
	ginkgo.It("Should not be able to find a user  in a organization", func() {

		organizationID := "Org_0001"

		exists, err := provider.UserExists(organizationID, "email_1@daisho.group")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).NotTo(gomega.BeTrue())

	})

	// ListUsers
	ginkgo.It("Should be able to get a list of the users in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())
		for i := 1; i <= 10; i++ {
			err = provider.AddUser(organizationID, fmt.Sprintf("email_%d@daisho.group", i))
			gomega.Expect(err).To(gomega.Succeed())
		}

		clusters, err := provider.ListUsers(organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(clusters).NotTo(gomega.BeEmpty())

	})
	ginkgo.It("Should be able to get an empty list of the users in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		clusters, err := provider.ListUsers(organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(clusters).To(gomega.BeEmpty())

	})
	ginkgo.It("Should not be able to get a list of the users in a organization", func() {

		organizationID := "Org0001"

		_, err := provider.ListUsers(organizationID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// DeleteUser
	ginkgo.It("Should be able to delete a user in a organization", func() {

		organizationID := "Org_0001"
		email := "email_1@daisho.group"
		org := &entities.Organization{ID: organizationID, Name: "organization_0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddUser(organizationID, email)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.DeleteUser(organizationID, email)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to delete a user in a organization", func() {

		organizationID := "Org_0001"
		email := "email_1@daisho.group"

		err := provider.DeleteUser(organizationID, email)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// --------------------------------------------------------------------------------------------------------------------

	// AddRole
	ginkgo.It("Should be able to add a new role in a organization", func() {

		organizationID := "org_XX01"
		org := &entities.Organization{ID: organizationID, Name: "organization OrgXX01", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddRole(organizationID, "developer")
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to add a new role  in a organization", func() {

		organizationID := "organization_id"

		err := provider.AddRole(organizationID, "developer")
		gomega.Expect(err).NotTo(gomega.Succeed())

	})
	ginkgo.It("Should not be able to add a new role in a organization (already exists)", func() {

		organizationID := "organization_id"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddUser(organizationID, "mail@daisho.group")
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddUser(organizationID, "mail@daisho.group")
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// UserExists
	ginkgo.It("Should be able to find a role  in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddRole(organizationID, "developer")
		gomega.Expect(err).To(gomega.Succeed())

		exists, err := provider.RoleExists(organizationID, "developer")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())

	})
	ginkgo.It("Should not be able to find a role  in a organization", func() {

		organizationID := "Org_0001"

		exists, err := provider.RoleExists(organizationID, "developer")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).NotTo(gomega.BeTrue())

	})

	// ListUsers
	ginkgo.It("Should be able to get a list of the roles in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddRole(organizationID, "developer")
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddRole(organizationID, "root")
		gomega.Expect(err).To(gomega.Succeed())

		clusters, err := provider.ListRoles(organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(clusters).NotTo(gomega.BeEmpty())

	})
	ginkgo.It("Should be able to get an empty list of the roles in a organization", func() {

		organizationID := "Org_0001"
		org := &entities.Organization{ID: organizationID, Name: "organization 0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		clusters, err := provider.ListRoles(organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(clusters).To(gomega.BeEmpty())

	})
	ginkgo.It("Should not be able to get a list of the roles in a organization", func() {

		organizationID := "Org0001"

		_, err := provider.ListRoles(organizationID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// DeleteUser
	ginkgo.It("Should be able to delete a role in a organization", func() {

		organizationID := "Org_0001"
		roleID := "developer"
		org := &entities.Organization{ID: organizationID, Name: "organization_0001", Created: 12}

		err := provider.Add(*org)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.AddRole(organizationID, roleID)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.DeleteRole(organizationID, roleID)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to delete a role in a organization", func() {

		organizationID := "Org_0001"
		roleID := "developer"

		err := provider.DeleteRole(organizationID, roleID)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})
}
