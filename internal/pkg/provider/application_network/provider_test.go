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
			InboundName:        "",
			OutboundName:       "",
			OutboundRequired:   false,
		}
		err := provider.AddConnectionInstance(toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		exists, err := provider.ExistsConnectionInstanceById(toAdd.ConnectionId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())
		exists, err = provider.ExistsConnectionInstance(
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
			InboundName:        "",
			OutboundName:       "",
			OutboundRequired:   false,
		}
		_ = provider.AddConnectionInstance(toAdd)
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

	ginkgo.It("should be able to retrieve a previously inserted ConnectionInstance using connectionId", func() {
		toAdd := entities.ConnectionInstance{
			OrganizationId:     entities.GenerateUUID(),
			ConnectionId:       entities.GenerateUUID(),
			SourceInstanceId:   entities.GenerateUUID(),
			SourceInstanceName: entities.GenerateUUID(),
			TargetInstanceId:   entities.GenerateUUID(),
			TargetInstanceName: entities.GenerateUUID(),
			InboundName:        "",
			OutboundName:       "",
			OutboundRequired:   false,
		}
		_ = provider.AddConnectionInstance(toAdd)
		connectionInstance, err := provider.GetConnectionInstanceById(toAdd.ConnectionId)
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
				InboundName:        "",
				OutboundName:       "",
				OutboundRequired:   false,
			},
			{
				OrganizationId:     organizationId,
				ConnectionId:       entities.GenerateUUID(),
				SourceInstanceId:   entities.GenerateUUID(),
				SourceInstanceName: entities.GenerateUUID(),
				TargetInstanceId:   entities.GenerateUUID(),
				TargetInstanceName: entities.GenerateUUID(),
				InboundName:        "",
				OutboundName:       "",
				OutboundRequired:   false,
			},
		}
		for _, instance := range toAdd {
			_ = provider.AddConnectionInstance(instance)
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
				InboundName:        "",
				OutboundName:       "",
				OutboundRequired:   false,
			},
			{
				OrganizationId:     organizationId,
				ConnectionId:       entities.GenerateUUID(),
				SourceInstanceId:   entities.GenerateUUID(),
				SourceInstanceName: entities.GenerateUUID(),
				TargetInstanceId:   entities.GenerateUUID(),
				TargetInstanceName: entities.GenerateUUID(),
				InboundName:        "",
				OutboundName:       "",
				OutboundRequired:   false,
			},
		}
		for _, instance := range toAdd {
			_ = provider.AddConnectionInstance(instance)
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

	// Connection Instance Link
	// ------------------------
	/*
		ginkgo.It("should be able to add a ConnectionInstanceLink", func() {
			instance := entities.ConnectionInstance{
				OrganizationId:     entities.GenerateUUID(),
				ConnectionId:       entities.GenerateUUID(),
				SourceInstanceId:   entities.GenerateUUID(),
				SourceInstanceName: entities.GenerateUUID(),
				TargetInstanceId:   entities.GenerateUUID(),
				TargetInstanceName: entities.GenerateUUID(),
				InboundName:        "",
				OutboundName:       "",
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
				InboundName:      entities.GenerateUUID(),
				OutboundName:     entities.GenerateUUID(),
			}
			err := provider.AddConnectionInstanceLink(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			exists, err := provider.ExistsConnectionInstanceLink(toAdd.ConnectionId, toAdd.SourceClusterId, toAdd.TargetClusterId)
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
				InboundName:        "",
				OutboundName:       "",
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
				InboundName:      entities.GenerateUUID(),
				OutboundName:     entities.GenerateUUID(),
			}
			err := provider.AddConnectionInstanceLink(toAdd)
			link, err := provider.GetConnectionInstanceLink(toAdd.ConnectionId, toAdd.SourceClusterId, toAdd.TargetClusterId)
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
				InboundName:        "",
				OutboundName:       "",
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
					InboundName:      entities.GenerateUUID(),
					OutboundName:     entities.GenerateUUID(),
				},
				{
					OrganizationId:   instance.OrganizationId,
					ConnectionId:     instance.ConnectionId,
					SourceInstanceId: instance.SourceInstanceId,
					SourceClusterId:  entities.GenerateUUID(),
					TargetInstanceId: instance.TargetInstanceId,
					TargetClusterId:  entities.GenerateUUID(),
					InboundName:      entities.GenerateUUID(),
					OutboundName:     entities.GenerateUUID(),
				},
			}
			for _, link := range toAdd {
				_ = provider.AddConnectionInstanceLink(link)
			}
			links, err := provider.ListConnectionInstanceLinks(instance.ConnectionId)
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
				InboundName:        "",
				OutboundName:       "",
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
					InboundName:      entities.GenerateUUID(),
					OutboundName:     entities.GenerateUUID(),
				},
				{
					OrganizationId:   instance.OrganizationId,
					ConnectionId:     instance.ConnectionId,
					SourceInstanceId: instance.SourceInstanceId,
					SourceClusterId:  entities.GenerateUUID(),
					TargetInstanceId: instance.TargetInstanceId,
					TargetClusterId:  entities.GenerateUUID(),
					InboundName:      entities.GenerateUUID(),
					OutboundName:     entities.GenerateUUID(),
				},
			}
			for _, link := range toAdd {
				_ = provider.AddConnectionInstanceLink(link)
			}
			err := provider.RemoveConnectionInstanceLinks(instance.ConnectionId)
			gomega.Expect(err).To(gomega.Succeed())
			links, err := provider.ListConnectionInstanceLinks(instance.ConnectionId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(links).To(gomega.BeEmpty())
		})
	*/
}
