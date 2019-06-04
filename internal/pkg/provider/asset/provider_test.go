/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package asset

import (
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {

	ginkgo.BeforeEach(func() {
		provider.Clear()
	})

	ginkgo.Context("adding assets", func(){
		ginkgo.It("should be able to add full asset", func(){
			toAdd := CreateTestAsset()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			exists, err := provider.Exists(toAdd.AssetId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())
		})
		ginkgo.It("should be able to add a basic asset", func(){
			toAdd := CreateTestAsset()
			toAdd.Storage = nil
			toAdd.Hardware = nil
			toAdd.Os = nil
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			exists, err := provider.Exists(toAdd.AssetId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())
		})
	})

	ginkgo.It("should be able to update an asset", func(){
		toAdd := CreateTestAsset()
		err := provider.Add(*toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		toAdd.EicNetIp = "2.2.2.2"
		err = provider.Update(*toAdd)
		gomega.Expect(err).To(gomega.Succeed())
	})

	ginkgo.It("should be able to retrieve an asset", func(){
		toAdd := CreateTestAsset()
		err := provider.Add(*toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		exists, err := provider.Exists(toAdd.AssetId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())
		retrieved, err := provider.Get(toAdd.AssetId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(retrieved).To(gomega.Equal(toAdd))
	})

	ginkgo.It("should be able to list the assets in an organization", func(){
	    numAssets := 10
		organizationID := entities.GenerateUUID()
		for index := 0; index < numAssets; index ++{
			toAdd := CreateTestAsset()
			toAdd.OrganizationId = organizationID
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())
		}
		// Add elements to other organizations
		for index := 0; index < numAssets; index ++{
			toAdd := CreateTestAsset()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())
		}
		retrieved, err := provider.List(organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(len(retrieved)).To(gomega.Equal(numAssets))
	})

	ginkgo.It("should be able to list the assets in an organization associated with an edge controller", func(){
		numAssets := 10
		organizationID := entities.GenerateUUID()
		edgeControllerID := entities.GenerateUUID()
		for index := 0; index < numAssets; index ++{
			toAdd := CreateTestAsset()
			toAdd.OrganizationId = organizationID
			if index % 2 == 0{
				toAdd.EdgeControllerId = edgeControllerID
			}
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())
		}

		retrieved, err := provider.ListControllerAssets(edgeControllerID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(len(retrieved)).To(gomega.Equal(numAssets/2))
	})

	ginkgo.It("should be able to delete an asset", func(){
		toAdd := CreateTestAsset()
		err := provider.Add(*toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		err = provider.Remove(toAdd.AssetId)
		gomega.Expect(err).To(gomega.Succeed())
		exists, err := provider.Exists(toAdd.AssetId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeFalse())
	})

}
