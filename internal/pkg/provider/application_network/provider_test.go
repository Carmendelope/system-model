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
	ginkgo.Context("ConnectionInstance", func() {
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
				Status:             entities.ConnectionStatusWaiting,
				IpRange:            entities.GenerateUUID(),
				ZtNetworkId:        entities.GenerateUUID(),
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
				Status:             entities.ConnectionStatusWaiting,
				IpRange:            entities.GenerateUUID(),
				ZtNetworkId:        entities.GenerateUUID(),
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

		ginkgo.It("should be able to retrieve a previously inserted ConnectionInstance using the zt network id", func() {
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
				Status:             entities.ConnectionStatusWaiting,
				IpRange:            entities.GenerateUUID(),
				ZtNetworkId:        entities.GenerateUUID(),
			}
			err := provider.AddConnectionInstance(toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			connectionInstance, err := provider.GetConnectionByZtNetworkId(
				toAdd.ZtNetworkId,
			)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(connectionInstance).NotTo(gomega.BeEmpty())
			gomega.Expect(connectionInstance[0]).To(gomega.Equal(toAdd))
		})

		ginkgo.It("should be able to update a ConnectionInstance", func() {
			connectionInstance := entities.ConnectionInstance{
				OrganizationId:     entities.GenerateUUID(),
				ConnectionId:       entities.GenerateUUID(),
				SourceInstanceId:   entities.GenerateUUID(),
				SourceInstanceName: entities.GenerateUUID(),
				TargetInstanceId:   entities.GenerateUUID(),
				TargetInstanceName: entities.GenerateUUID(),
				InboundName:        entities.GenerateUUID(),
				OutboundName:       entities.GenerateUUID(),
				OutboundRequired:   false,
				Status:             entities.ConnectionStatusWaiting,
				IpRange:            entities.GenerateUUID(),
				ZtNetworkId:        entities.GenerateUUID(),
			}
			err := provider.AddConnectionInstance(connectionInstance)
			gomega.Expect(err).To(gomega.Succeed())
			connectionInstance.Status = entities.ConnectionStatusEstablished
			connectionInstance.IpRange = "172.16.0.1-172.16.0.254"
			connectionInstance.ZtNetworkId = entities.GenerateUUID()
			err = provider.UpdateConnectionInstance(connectionInstance)
			gomega.Expect(err).To(gomega.Succeed())
			updatedInstance, err := provider.GetConnectionInstance(
				connectionInstance.OrganizationId,
				connectionInstance.SourceInstanceId,
				connectionInstance.TargetInstanceId,
				connectionInstance.InboundName,
				connectionInstance.OutboundName,
			)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(*updatedInstance).To(gomega.Equal(connectionInstance))
			//gomega.Expect(updatedInstance.Status).To(gomega.Equal(connectionInstance.Status))
			//gomega.Expect(updatedInstance.IpRange).To(gomega.Equal(connectionInstance.IpRange))
			//gomega.Expect(updatedInstance.ZtNetworkId).To(gomega.Equal(connectionInstance.ZtNetworkId))
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
					Status:             entities.ConnectionStatusWaiting,
					IpRange:            entities.GenerateUUID(),
					ZtNetworkId:        entities.GenerateUUID(),
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
					Status:             entities.ConnectionStatusWaiting,
					IpRange:            entities.GenerateUUID(),
					ZtNetworkId:        entities.GenerateUUID(),
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
					Status:             entities.ConnectionStatusWaiting,
					IpRange:            entities.GenerateUUID(),
					ZtNetworkId:        entities.GenerateUUID(),
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
					Status:             entities.ConnectionStatusWaiting,
					IpRange:            entities.GenerateUUID(),
					ZtNetworkId:        entities.GenerateUUID(),
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
					Status:             entities.ConnectionStatusWaiting,
					IpRange:            entities.GenerateUUID(),
					ZtNetworkId:        entities.GenerateUUID(),
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
					Status:             entities.ConnectionStatusWaiting,
					IpRange:            entities.GenerateUUID(),
					ZtNetworkId:        entities.GenerateUUID(),
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
	})

	// Connection Instance Link
	// ------------------------
	ginkgo.Context("Connection Instance Link", func() {
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
				Status:             entities.ConnectionStatusWaiting,
				IpRange:            "",
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
				Status:           entities.ConnectionStatusWaiting,
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
				Status:             entities.ConnectionStatusWaiting,
				IpRange:            "",
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
				Status:           entities.ConnectionStatusWaiting,
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
				Status:             entities.ConnectionStatusWaiting,
				IpRange:            "",
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
					Status:           entities.ConnectionStatusWaiting,
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
					Status:           entities.ConnectionStatusWaiting,
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
				Status:             entities.ConnectionStatusWaiting,
				IpRange:            "",
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
					Status:           entities.ConnectionStatusWaiting,
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
					Status:           entities.ConnectionStatusWaiting,
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
	})

	ginkgo.Context("ZTNetworkConnection", func() {
		ginkgo.It("Should be able to add a ztnetworkConnection", func() {
			toAdd := &entities.ZTNetworkConnection{
				OrganizationId: entities.GenerateUUID(),
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  entities.GenerateUUID(),
				ServiceId:      entities.GenerateUUID(),
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           entities.ConnectionSideOutbound,
			}
			err := provider.AddZTConnection(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("Should be not able to add a ztnetworkConnection twice", func() {
			toAdd := &entities.ZTNetworkConnection{
				OrganizationId: entities.GenerateUUID(),
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  entities.GenerateUUID(),
				ServiceId:      entities.GenerateUUID(),
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           entities.ConnectionSideOutbound,
			}
			err := provider.AddZTConnection(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			err = provider.AddZTConnection(*toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		ginkgo.It("should be able to determinate if a ztnetwork connection exists", func() {
			toAdd := &entities.ZTNetworkConnection{
				OrganizationId: entities.GenerateUUID(),
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  entities.GenerateUUID(),
				ServiceId:      entities.GenerateUUID(),
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           entities.ConnectionSideOutbound,
			}
			err := provider.AddZTConnection(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			exits, err := provider.ExistsZTConnection(toAdd.OrganizationId, toAdd.ZtNetworkId, toAdd.AppInstanceId, toAdd.ServiceId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exits).To(gomega.BeTrue())

		})
		ginkgo.It("should be able to determinate if a ztnetwork connection does not exist", func() {
			exits, err := provider.ExistsZTConnection(entities.GenerateUUID(), entities.GenerateUUID(), entities.GenerateUUID(), entities.GenerateUUID())
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exits).NotTo(gomega.BeTrue())

		})
		ginkgo.It("should be able to get  a ztnetwork connection", func() {
			toAdd := &entities.ZTNetworkConnection{
				OrganizationId: entities.GenerateUUID(),
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  entities.GenerateUUID(),
				ServiceId:      entities.GenerateUUID(),
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           entities.ConnectionSideOutbound,
			}
			err := provider.AddZTConnection(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			retrieve, err := provider.GetZTConnection(toAdd.OrganizationId, toAdd.ZtNetworkId, toAdd.AppInstanceId, toAdd.ServiceId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieve).NotTo(gomega.BeNil())
			gomega.Expect(retrieve).Should(gomega.Equal(toAdd))
		})
		ginkgo.It("should not be able to get a ztnetwork connection when it does not exist", func() {
			_, err := provider.GetZTConnection(entities.GenerateUUID(), entities.GenerateUUID(), entities.GenerateUUID(), entities.GenerateUUID())
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("should be able to list the ztnetwork connections of a networkId", func() {
			toAdd := &entities.ZTNetworkConnection{
				OrganizationId: entities.GenerateUUID(),
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  entities.GenerateUUID(),
				ServiceId:      entities.GenerateUUID(),
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           entities.ConnectionSideOutbound,
			}
			err := provider.AddZTConnection(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			toAdd.AppInstanceId = entities.GenerateUUID()
			err = provider.AddZTConnection(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			list, err := provider.ListZTConnections(toAdd.OrganizationId, toAdd.ZtNetworkId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(list).NotTo(gomega.BeNil())
			gomega.Expect(len(list)).Should(gomega.Equal(2))
		})
		ginkgo.It("should be able to list an empty list of ztnetwork connections", func() {
			list, err := provider.ListZTConnections(entities.GenerateUUID(), entities.GenerateUUID())
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(list).NotTo(gomega.BeNil())
			gomega.Expect(len(list)).Should(gomega.Equal(0))
		})
		ginkgo.It("Should be able to update a ztNetwork connection", func() {
			toAdd := &entities.ZTNetworkConnection{
				OrganizationId: entities.GenerateUUID(),
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  entities.GenerateUUID(),
				ServiceId:      entities.GenerateUUID(),
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           entities.ConnectionSideOutbound,
			}
			err := provider.AddZTConnection(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			toAdd.ZtIp = "yyy.yyy.yyy.yyy"
			err = provider.UpdateZTConnection(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			retrieve, err := provider.GetZTConnection(toAdd.OrganizationId, toAdd.ZtNetworkId, toAdd.AppInstanceId, toAdd.ServiceId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieve.ZtIp).Should(gomega.Equal(toAdd.ZtIp))

		})
		ginkgo.It("Should not be able to update a non existing ztNetwork connection", func() {
			toAdd := &entities.ZTNetworkConnection{
				OrganizationId: entities.GenerateUUID(),
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  entities.GenerateUUID(),
				ServiceId:      entities.GenerateUUID(),
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           entities.ConnectionSideOutbound,
			}

			err := provider.UpdateZTConnection(*toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})

		ginkgo.It("should be able to remove a ztnetwork connections ", func() {
			toAdd := &entities.ZTNetworkConnection{
				OrganizationId: entities.GenerateUUID(),
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  entities.GenerateUUID(),
				ServiceId:      entities.GenerateUUID(),
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           entities.ConnectionSideOutbound,
			}
			err := provider.AddZTConnection(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.RemoveZTConnection(toAdd.OrganizationId, toAdd.ZtNetworkId)
			gomega.Expect(err).To(gomega.Succeed())

		})
		ginkgo.It("should not be able to remove a ztnetwork connections if it does not exist", func() {
			err := provider.RemoveZTConnection(entities.GenerateUUID(), entities.GenerateUUID())
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("should be able to remove a ztnetwork connections ", func() {
			toAdd := &entities.ZTNetworkConnection{
				OrganizationId: entities.GenerateUUID(),
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  entities.GenerateUUID(),
				ServiceId:      entities.GenerateUUID(),
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           entities.ConnectionSideOutbound,
			}
			err := provider.AddZTConnection(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			toAdd.AppInstanceId = entities.GenerateUUID()
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.RemoveZTConnection(toAdd.OrganizationId, toAdd.ZtNetworkId)
			gomega.Expect(err).To(gomega.Succeed())

		})
	})

}
