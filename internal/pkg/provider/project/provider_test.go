/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package project

import (
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {

	ginkgo.AfterEach(func() {
		provider.Clear()
	})

	ginkgo.Context("adding project", func() {
		ginkgo.It("should be able to add a project", func() {
			toAdd := CreateProject()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("should not be able to add a project twice", func() {
			toAdd := CreateProject()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.Add(*toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})

	ginkgo.Context("getting project", func() {
		ginkgo.It("should be able to get a project", func(){
			toAdd := CreateProject()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			retrieve, err := provider.Get(toAdd.OwnerAccountId, toAdd.ProjectId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieve).NotTo(gomega.BeNil())
			gomega.Expect(retrieve).Should(gomega.Equal(toAdd))
		})
		ginkgo.It("should not be able to get a non existing project", func(){
			_, err := provider.Get(entities.GenerateUUID(), entities.GenerateUUID())
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})

	ginkgo.Context("removing project", func() {
		ginkgo.It("should be able to remove a project", func(){
			toAdd := CreateProject()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.Remove(toAdd.OwnerAccountId, toAdd.ProjectId)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("should not be able to remove a non existing project", func(){
			err := provider.Remove(entities.GenerateUUID(), entities.GenerateUUID())
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})
	ginkgo.Context("updating project", func() {
		ginkgo.It("should be able to update a project", func(){
			toAdd := CreateProject()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			// update Account
			toAdd.Name = "updated name"
			toAdd.State = entities.ProjectState_Deactivated
			toAdd.StateInfo = "deactivated for test"

			err = provider.Update(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			// check the update works
			retrieve, err := provider.Get(toAdd.OwnerAccountId, toAdd.ProjectId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieve).NotTo(gomega.BeNil())
			gomega.Expect(retrieve).Should(gomega.Equal(toAdd))

		})
		ginkgo.It("should not be able to update a non existing project", func(){
			toAdd := CreateProject()

			err := provider.Update(*toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})
	ginkgo.Context("checking if exists a project", func() {
		ginkgo.It("should be able to check a project exists", func(){
			toAdd := CreateProject()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			exists, err := provider.Exists(toAdd.OwnerAccountId, toAdd.ProjectId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())
		})
		ginkgo.It("should be able to check a project does not exist", func(){
			exists, err := provider.Exists(entities.GenerateUUID(), entities.GenerateUUID())
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())
		})
	})
	ginkgo.Context("listing projects of an account", func() {
		ginkgo.It("should be able to list existing projects", func(){
			toAdd := CreateProject()
			numProjects := 10
			for i:= 0; i<numProjects; i++ {
				toAdd.ProjectId = entities.GenerateUUID()
				err := provider.Add(*toAdd)
				gomega.Expect(err).To(gomega.Succeed())
			}

			projects, err := provider.ListAccountProjects(toAdd.OwnerAccountId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(projects).NotTo(gomega.BeNil())
			gomega.Expect(len(projects)).Should(gomega.Equal(numProjects))
		})
		ginkgo.It("should be able to list existing projects", func(){

			projects, err := provider.ListAccountProjects(entities.GenerateUUID())
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(projects).NotTo(gomega.BeNil())
			gomega.Expect(len(projects)).Should(gomega.BeZero())
		})

	})
}