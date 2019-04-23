/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package asset

import (
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