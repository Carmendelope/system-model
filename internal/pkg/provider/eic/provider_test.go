/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package eic

import (
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {

	ginkgo.BeforeEach(func() {
		provider.Clear()
	})

	ginkgo.It("should be able to add a controller", func(){
		toAdd := CreateTestEdgeController()
		err := provider.Add(*toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		exists, err := provider.Exists(toAdd.EdgeControllerId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())
	})

	ginkgo.It("should be able to update an EIC", func(){
		toAdd := CreateTestEdgeController()
		err := provider.Add(*toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		toAdd.Name = "newName"
		err = provider.Update(*toAdd)
		gomega.Expect(err).To(gomega.Succeed())
	})

	ginkgo.It("should be able to retrieve an EIC", func(){
		toAdd := CreateTestEdgeController()
		err := provider.Add(*toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		exists, err := provider.Exists(toAdd.EdgeControllerId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())
		retrieved, err := provider.Get(toAdd.EdgeControllerId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(retrieved).To(gomega.Equal(toAdd))
	})

	ginkgo.It("should be able to list the EIC of an organization", func(){
	    numEIC := 10
	    organizationID := entities.GenerateUUID()
	    for index := 0; index < numEIC; index ++{
	    	toAdd := CreateTestEdgeController()
	    	toAdd.OrganizationId = organizationID
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())
		}
	    // Add elements to other organizations
		for index := 0; index < numEIC; index ++{
			toAdd := CreateTestEdgeController()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())
		}
		retrieved, err := provider.List(organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(len(retrieved)).To(gomega.Equal(numEIC))
	})

	ginkgo.It("should be able to remove an EIC", func(){
		toAdd := CreateTestEdgeController()
		err := provider.Add(*toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		err = provider.Remove(toAdd.EdgeControllerId)
		gomega.Expect(err).To(gomega.Succeed())
		exists, err := provider.Exists(toAdd.EdgeControllerId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeFalse())
	})
	ginkgo.It("should be able to update an EIC", func(){
		toAdd := CreateTestEdgeController()
		err := provider.Add(*toAdd)
		gomega.Expect(err).To(gomega.Succeed())

		toAdd.Name = "new Name"
		err = provider.Update(*toAdd)
		gomega.Expect(err).To(gomega.Succeed())

		ec, err := provider.Get(toAdd.EdgeControllerId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(ec).NotTo(gomega.BeNil())
		gomega.Expect(ec.Name).Should(gomega.Equal("new Name"))


	})

}