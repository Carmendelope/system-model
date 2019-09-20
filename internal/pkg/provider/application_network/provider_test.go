/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package application_network

import (
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {
	ginkgo.AfterEach(func() {
		_ = provider.Clear()
	})

	ginkgo.It("should be able to add a ConnectionInstance (also check existence methods)", func() {
		toAdd := entities.ConnectionInstance{
			OrganizationId:     entities.GenerateUUID(),
			ConnectionId:       entities.GenerateUUID(),
			SourceInstanceId:   entities.GenerateUUID(),
			SourceInstanceName: entities.GenerateUUID(),
			TargetInstanceId:   entities.GenerateUUID(),
			TargetInstanceName: entities.GenerateUUID(),
			InboundName:        entities.GenerateUUID(),
			OutboundName:       entities.GenerateUUID(),
			OutboundRequired:   false,
		}
		err := provider.AddConnectionInstance(toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		exists, err := provider.ExistsConnectionInstance(
			toAdd.OrganizationId,
			toAdd.SourceInstanceId,
			toAdd.TargetInstanceId,
			toAdd.InboundName,
			toAdd.OutboundName,
		)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())
	})

	ginkgo.It("should be able to retrieve a previously inserted ConnectionInstance using the composite PK", func() {
		toAdd := entities.ConnectionInstance{
			OrganizationId:     entities.GenerateUUID(),
			ConnectionId:       entities.GenerateUUID(),
			SourceInstanceId:   entities.GenerateUUID(),
			SourceInstanceName: entities.GenerateUUID(),
			TargetInstanceId:   entities.GenerateUUID(),
			TargetInstanceName: entities.GenerateUUID(),
			InboundName:        entities.GenerateUUID(),
			OutboundName:       entities.GenerateUUID(),
			OutboundRequired:   false,
		}
		err := provider.AddConnectionInstance(toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		connectionInstance, err := provider.GetConnectionInstance(
			toAdd.OrganizationId,
			toAdd.SourceInstanceId,
			toAdd.TargetInstanceId,
			toAdd.InboundName,
			toAdd.OutboundName,
		)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(*connectionInstance).To(gomega.Equal(toAdd))
	})

	ginkgo.It("should be able to retrieve a list of inserted ConnectionInstances", func() {
		organizationId := entities.GenerateUUID()
		toAdd := []entities.ConnectionInstance{
			{
				OrganizationId:     organizationId,
				ConnectionId:       entities.GenerateUUID(),
				SourceInstanceId:   entities.GenerateUUID(),
				SourceInstanceName: entities.GenerateUUID(),
				TargetInstanceId:   entities.GenerateUUID(),
				TargetInstanceName: entities.GenerateUUID(),
				InboundName:        entities.GenerateUUID(),
				OutboundName:       entities.GenerateUUID(),
				OutboundRequired:   false,
			},
			{
				OrganizationId:     organizationId,
				ConnectionId:       entities.GenerateUUID(),
				SourceInstanceId:   entities.GenerateUUID(),
				SourceInstanceName: entities.GenerateUUID(),
				TargetInstanceId:   entities.GenerateUUID(),
				TargetInstanceName: entities.GenerateUUID(),
				InboundName:        entities.GenerateUUID(),
				OutboundName:       entities.GenerateUUID(),
				OutboundRequired:   false,
			},
		}
		for _, instance := range toAdd {
			err := provider.AddConnectionInstance(instance)
			gomega.Expect(err).To(gomega.Succeed())
		}
		connectionInstances, err := provider.ListConnectionInstances(organizationId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(connectionInstances).To(gomega.ConsistOf(toAdd))
	})

	ginkgo.It("should be able to remove a ConnectionInstance from the DB using the composite PK", func() {
		organizationId := entities.GenerateUUID()
		toAdd := []entities.ConnectionInstance{
			{
				OrganizationId:     organizationId,
				ConnectionId:       entities.GenerateUUID(),
				SourceInstanceId:   entities.GenerateUUID(),
				SourceInstanceName: entities.GenerateUUID(),
				TargetInstanceId:   entities.GenerateUUID(),
				TargetInstanceName: entities.GenerateUUID(),
				InboundName:        entities.GenerateUUID(),
				OutboundName:       entities.GenerateUUID(),
				OutboundRequired:   false,
			},
			{
				OrganizationId:     organizationId,
				ConnectionId:       entities.GenerateUUID(),
				SourceInstanceId:   entities.GenerateUUID(),
				SourceInstanceName: entities.GenerateUUID(),
				TargetInstanceId:   entities.GenerateUUID(),
				TargetInstanceName: entities.GenerateUUID(),
				InboundName:        entities.GenerateUUID(),
				OutboundName:       entities.GenerateUUID(),
				OutboundRequired:   false,
			},
		}
		for _, instance := range toAdd {
			err := provider.AddConnectionInstance(instance)
			gomega.Expect(err).To(gomega.Succeed())
		}
		err := provider.RemoveConnectionInstance(
			toAdd[0].OrganizationId,
			toAdd[0].SourceInstanceId,
			toAdd[0].TargetInstanceId,
			toAdd[0].InboundName,
			toAdd[0].OutboundName,
		)
		gomega.Expect(err).To(gomega.Succeed())
		connectionInstances, err := provider.ListConnectionInstances(organizationId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(connectionInstances).To(gomega.ConsistOf(toAdd[1:]))
	})

	ginkgo.It("should be able to list ConnectionInstances where the instance is the target", func() {
		organizationId := entities.GenerateUUID()
		targetInstanceId := entities.GenerateUUID()
		numConnections := 5
		for i := 0; i < numConnections; i++ {
			toAdd := entities.ConnectionInstance{
				OrganizationId:     organizationId,
				ConnectionId:       entities.GenerateUUID(),
				SourceInstanceId:   entities.GenerateUUID(),
				SourceInstanceName: entities.GenerateUUID(),
				TargetInstanceId:   targetInstanceId,
				TargetInstanceName: targetInstanceId,
				InboundName:        entities.GenerateUUID(),
				OutboundName:       entities.GenerateUUID(),
				OutboundRequired:   false,
			}
			err := provider.AddConnectionInstance(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
		}
		list, err := provider.ListInboundConnections(organizationId, targetInstanceId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(list).NotTo(gomega.BeNil())
		gomega.Expect(len(list)).Should(gomega.Equal(numConnections))
	})
	ginkgo.It("should be able to retrieve an empty list ConnectionInstance when there are no connections where the instance is the target", func() {
		organizationId := entities.GenerateUUID()
		targetInstanceId := entities.GenerateUUID()

		list, err := provider.ListInboundConnections(organizationId, targetInstanceId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(list).NotTo(gomega.BeNil())
		gomega.Expect(len(list)).Should(gomega.Equal(0))
	})
	ginkgo.It("should be able to list ConnectionInstances where the instance is the source", func() {
		organizationId := entities.GenerateUUID()
		sourceInstanceId := entities.GenerateUUID()
		numConnections := 5
		for i := 0; i < numConnections; i++ {
			toAdd := entities.ConnectionInstance{
				OrganizationId:     organizationId,
				ConnectionId:       entities.GenerateUUID(),
				SourceInstanceId:   sourceInstanceId,
				SourceInstanceName: "source instance",
				TargetInstanceId:   entities.GenerateUUID(),
				TargetInstanceName: entities.GenerateUUID(),
				InboundName:        entities.GenerateUUID(),
				OutboundName:       entities.GenerateUUID(),
				OutboundRequired:   false,
			}
			err := provider.AddConnectionInstance(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
		}
		list, err := provider.ListOutboundConnections(organizationId, sourceInstanceId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(list).NotTo(gomega.BeNil())
		gomega.Expect(len(list)).Should(gomega.Equal(numConnections))

	})
	ginkgo.It("should be able to retrieve an empty list ConnectionInstance when there are no connections where the instance is the source", func() {
		organizationId := entities.GenerateUUID()
		sourceInstanceId := entities.GenerateUUID()
		list, err := provider.ListOutboundConnections(organizationId, sourceInstanceId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(list).NotTo(gomega.BeNil())
		gomega.Expect(len(list)).Should(gomega.Equal(0))
	})

	// Connection Instance Link
	// ------------------------
	ginkgo.It("should be able to add a ConnectionInstanceLink", func() {
		instance := entities.ConnectionInstance{
			OrganizationId:     entities.GenerateUUID(),
			ConnectionId:       entities.GenerateUUID(),
			SourceInstanceId:   entities.GenerateUUID(),
			SourceInstanceName: entities.GenerateUUID(),
			TargetInstanceId:   entities.GenerateUUID(),
			TargetInstanceName: entities.GenerateUUID(),
			InboundName:        entities.GenerateUUID(),
			OutboundName:       entities.GenerateUUID(),
			OutboundRequired:   false,
		}
		err := provider.AddConnectionInstance(instance)
		gomega.Expect(err).To(gomega.Succeed())

		toAdd := entities.ConnectionInstanceLink{
			OrganizationId:   instance.OrganizationId,
			ConnectionId:     instance.ConnectionId,
			SourceInstanceId: instance.SourceInstanceId,
			SourceClusterId:  entities.GenerateUUID(),
			TargetInstanceId: instance.TargetInstanceId,
			TargetClusterId:  entities.GenerateUUID(),
			InboundName:      instance.InboundName,
			OutboundName:     instance.OutboundName,
		}
		err = provider.AddConnectionInstanceLink(toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		exists, err := provider.ExistsConnectionInstanceLink(
			toAdd.OrganizationId,
			toAdd.SourceInstanceId,
			toAdd.TargetInstanceId,
			toAdd.SourceClusterId,
			toAdd.TargetClusterId,
			toAdd.InboundName,
			toAdd.OutboundName,
		)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())
	})

	ginkgo.It("should be able to retrieve a previously inserted ConnectionInstanceLink", func() {
		instance := entities.ConnectionInstance{
			OrganizationId:     entities.GenerateUUID(),
			ConnectionId:       entities.GenerateUUID(),
			SourceInstanceId:   entities.GenerateUUID(),
			SourceInstanceName: entities.GenerateUUID(),
			TargetInstanceId:   entities.GenerateUUID(),
			TargetInstanceName: entities.GenerateUUID(),
			InboundName:        entities.GenerateUUID(),
			OutboundName:       entities.GenerateUUID(),
			OutboundRequired:   false,
		}
		_ = provider.AddConnectionInstance(instance)

		toAdd := entities.ConnectionInstanceLink{
			OrganizationId:   instance.OrganizationId,
			ConnectionId:     instance.ConnectionId,
			SourceInstanceId: instance.SourceInstanceId,
			SourceClusterId:  entities.GenerateUUID(),
			TargetInstanceId: instance.TargetInstanceId,
			TargetClusterId:  entities.GenerateUUID(),
			InboundName:      instance.InboundName,
			OutboundName:     instance.OutboundName,
		}
		err := provider.AddConnectionInstanceLink(toAdd)
		link, err := provider.GetConnectionInstanceLink(
			toAdd.OrganizationId,
			toAdd.SourceInstanceId,
			toAdd.TargetInstanceId,
			toAdd.SourceClusterId,
			toAdd.TargetClusterId,
			toAdd.InboundName,
			toAdd.OutboundName,
		)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(*link).To(gomega.Equal(toAdd))
	})

	ginkgo.It("should be able to list all the ConnectionInstanceLinks associated to a ConnectionInstance", func() {
		instance := entities.ConnectionInstance{
			OrganizationId:     entities.GenerateUUID(),
			ConnectionId:       entities.GenerateUUID(),
			SourceInstanceId:   entities.GenerateUUID(),
			SourceInstanceName: entities.GenerateUUID(),
			TargetInstanceId:   entities.GenerateUUID(),
			TargetInstanceName: entities.GenerateUUID(),
			InboundName:        entities.GenerateUUID(),
			OutboundName:       entities.GenerateUUID(),
			OutboundRequired:   false,
		}
		_ = provider.AddConnectionInstance(instance)

		toAdd := []entities.ConnectionInstanceLink{
			{
				OrganizationId:   instance.OrganizationId,
				ConnectionId:     instance.ConnectionId,
				SourceInstanceId: instance.SourceInstanceId,
				SourceClusterId:  entities.GenerateUUID(),
				TargetInstanceId: instance.TargetInstanceId,
				TargetClusterId:  entities.GenerateUUID(),
				InboundName:      instance.InboundName,
				OutboundName:     instance.OutboundName,
			},
			{
				OrganizationId:   instance.OrganizationId,
				ConnectionId:     instance.ConnectionId,
				SourceInstanceId: instance.SourceInstanceId,
				SourceClusterId:  entities.GenerateUUID(),
				TargetInstanceId: instance.TargetInstanceId,
				TargetClusterId:  entities.GenerateUUID(),
				InboundName:      instance.InboundName,
				OutboundName:     instance.OutboundName,
			},
		}
		for _, link := range toAdd {
			err := provider.AddConnectionInstanceLink(link)
			gomega.Expect(err).To(gomega.Succeed())
		}
		links, err := provider.ListConnectionInstanceLinks(instance.OrganizationId, instance.SourceInstanceId, instance.TargetInstanceId, instance.InboundName, instance.OutboundName)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(links).To(gomega.ConsistOf(toAdd))
	})

	ginkgo.It("should be able to remove all the ConnectionInstanceLinks from a connectionInstance", func() {
		instance := entities.ConnectionInstance{
			OrganizationId:     entities.GenerateUUID(),
			ConnectionId:       entities.GenerateUUID(),
			SourceInstanceId:   entities.GenerateUUID(),
			SourceInstanceName: entities.GenerateUUID(),
			TargetInstanceId:   entities.GenerateUUID(),
			TargetInstanceName: entities.GenerateUUID(),
			InboundName:        entities.GenerateUUID(),
			OutboundName:       entities.GenerateUUID(),
			OutboundRequired:   false,
		}
		err := provider.AddConnectionInstance(instance)
		gomega.Expect(err).To(gomega.Succeed())

		toAdd := []entities.ConnectionInstanceLink{
			{
				OrganizationId:   instance.OrganizationId,
				ConnectionId:     instance.ConnectionId,
				SourceInstanceId: instance.SourceInstanceId,
				SourceClusterId:  entities.GenerateUUID(),
				TargetInstanceId: instance.TargetInstanceId,
				TargetClusterId:  entities.GenerateUUID(),
				InboundName:      instance.InboundName,
				OutboundName:     instance.OutboundName,
			},
			{
				OrganizationId:   instance.OrganizationId,
				ConnectionId:     instance.ConnectionId,
				SourceInstanceId: instance.SourceInstanceId,
				SourceClusterId:  entities.GenerateUUID(),
				TargetInstanceId: instance.TargetInstanceId,
				TargetClusterId:  entities.GenerateUUID(),
				InboundName:      instance.InboundName,
				OutboundName:     instance.OutboundName,
			},
		}
		for _, link := range toAdd {
			err = provider.AddConnectionInstanceLink(link)
			gomega.Expect(err).To(gomega.Succeed())
		}
		err = provider.RemoveConnectionInstanceLinks(instance.OrganizationId, instance.SourceInstanceId, instance.TargetInstanceId, instance.InboundName, instance.OutboundName)
		gomega.Expect(err).To(gomega.Succeed())
		links, err := provider.ListConnectionInstanceLinks(instance.OrganizationId, instance.SourceInstanceId, instance.TargetInstanceId, instance.InboundName, instance.OutboundName)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(links).To(gomega.BeEmpty())
	})
}
