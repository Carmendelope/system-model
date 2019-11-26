/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package application_network

import (
	"context"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-application-network-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/application"
	"github.com/nalej/system-model/internal/pkg/provider/application_network"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"sort"
)

func addOrganization(organizationProvider organization.Provider) entities.Organization {
	org := entities.Organization{
		ID:      entities.GenerateUUID(),
		Name:    "AppNet test org",
		Created: 0,
	}
	err := organizationProvider.Add(org)
	gomega.Expect(err).To(gomega.Succeed())
	return org
}

func addSourceInstance(organizationId string, outboundRequired bool, applicationProvider application.Provider) entities.AppInstance {
	sourceInstance := entities.AppInstance{
		OrganizationId:       organizationId,
		AppInstanceId:        entities.GenerateUUID(),
		InboundNetInterfaces: nil,
		OutboundNetInterfaces: []entities.OutboundNetworkInterface{{
			Name:     "source-outbound",
			Required: outboundRequired,
		}},
	}
	err := applicationProvider.AddInstance(sourceInstance)
	gomega.Expect(err).To(gomega.Succeed())
	return sourceInstance
}

func addTargetInstance(organizationId string, applicationProvider application.Provider) entities.AppInstance {
	sourceInstance := entities.AppInstance{
		OrganizationId:        organizationId,
		AppInstanceId:         entities.GenerateUUID(),
		InboundNetInterfaces:  []entities.InboundNetworkInterface{{Name: "target-inbound"}},
		OutboundNetInterfaces: nil,
	}
	err := applicationProvider.AddInstance(sourceInstance)
	gomega.Expect(err).To(gomega.Succeed())
	return sourceInstance
}

func addInstance(organizationId string, applicationProvider application.Provider) entities.AppInstance {
	sourceInstance := entities.AppInstance{
		OrganizationId: organizationId,
		AppInstanceId:  entities.GenerateUUID(),
		Groups: []entities.ServiceGroupInstance{
			{
				ServiceInstances: []entities.ServiceInstance{
					{
						ServiceId: entities.GenerateUUID()},
				}},
		},
	}
	err := applicationProvider.AddInstance(sourceInstance)
	gomega.Expect(err).To(gomega.Succeed())
	return sourceInstance
}

var _ = ginkgo.Describe("Application Network service", func() {
	var (
		server   *grpc.Server
		listener *bufconn.Listener
		client   grpc_application_network_go.ApplicationNetworkClient

		organizationProvider organization.Provider
		applicationProvider  application.Provider
		appNetProvider       application_network.Provider
	)

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		// Register the service
		organizationProvider = organization.NewMockupOrganizationProvider()
		applicationProvider = application.NewMockupApplicationProvider()
		appNetProvider = application_network.NewMockupApplicationNetworkProvider()
		manager := NewManager(organizationProvider, applicationProvider, appNetProvider)
		handler := NewHandler(manager)
		grpc_application_network_go.RegisterApplicationNetworkServer(server, handler)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())
		client = grpc_application_network_go.NewApplicationNetworkClient(conn)

		test.LaunchServer(server, listener)

	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		_ = listener.Close()
	})

	ginkgo.BeforeEach(func() {
		_ = organizationProvider.Clear()
		_ = applicationProvider.Clear()
		_ = appNetProvider.Clear()
	})

	ginkgo.Context("when adding a connection", func() {

		ginkgo.It("should support adding a new connection", func() {
			organization := addOrganization(organizationProvider)
			sourceInstance := addSourceInstance(organization.ID, false, applicationProvider)
			targetInstance := addTargetInstance(organization.ID, applicationProvider)

			addConnectionRequest := &grpc_application_network_go.AddConnectionRequest{
				OrganizationId:   organization.ID,
				SourceInstanceId: sourceInstance.AppInstanceId,
				TargetInstanceId: targetInstance.AppInstanceId,
				InboundName:      targetInstance.InboundNetInterfaces[0].Name,
				OutboundName:     sourceInstance.OutboundNetInterfaces[0].Name,
				IpRange:          entities.GenerateUUID(),
				ZtNetworkId:      entities.GenerateUUID(),
			}
			connectionInstance, err := client.AddConnection(context.Background(), addConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(connectionInstance).ToNot(gomega.BeNil())
			gomega.Expect(connectionInstance.OrganizationId).To(gomega.Equal(addConnectionRequest.OrganizationId))
			gomega.Expect(connectionInstance.SourceInstanceId).To(gomega.Equal(addConnectionRequest.SourceInstanceId))
			gomega.Expect(connectionInstance.TargetInstanceId).To(gomega.Equal(addConnectionRequest.TargetInstanceId))
			gomega.Expect(connectionInstance.InboundName).To(gomega.Equal(addConnectionRequest.InboundName))
			gomega.Expect(connectionInstance.OutboundName).To(gomega.Equal(addConnectionRequest.OutboundName))
		})

		ginkgo.It("should fail if an equivalent connection already exists", func() {
			organization := addOrganization(organizationProvider)
			sourceInstance := addSourceInstance(organization.ID, false, applicationProvider)
			targetInstance := addTargetInstance(organization.ID, applicationProvider)

			addConnectionRequest := &grpc_application_network_go.AddConnectionRequest{
				OrganizationId:   organization.ID,
				SourceInstanceId: sourceInstance.AppInstanceId,
				TargetInstanceId: targetInstance.AppInstanceId,
				InboundName:      targetInstance.InboundNetInterfaces[0].Name,
				OutboundName:     sourceInstance.OutboundNetInterfaces[0].Name,
				IpRange:          entities.GenerateUUID(),
				ZtNetworkId:      entities.GenerateUUID(),
			}
			_, err := client.AddConnection(context.Background(), addConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			connectionInstance, err := client.AddConnection(context.Background(), addConnectionRequest)
			gomega.Expect(err).ToNot(gomega.Succeed())
			gomega.Expect(connectionInstance).To(gomega.BeNil())
		})
	})

	ginkgo.Context("when getting a connection", func() {

		ginkgo.It("should retrieve a previously added connection using the composite PK", func() {
			organization := addOrganization(organizationProvider)
			sourceInstance := addSourceInstance(organization.ID, false, applicationProvider)
			targetInstance := addTargetInstance(organization.ID, applicationProvider)

			addConnectionRequest := &grpc_application_network_go.AddConnectionRequest{
				OrganizationId:   organization.ID,
				SourceInstanceId: sourceInstance.AppInstanceId,
				TargetInstanceId: targetInstance.AppInstanceId,
				InboundName:      targetInstance.InboundNetInterfaces[0].Name,
				OutboundName:     sourceInstance.OutboundNetInterfaces[0].Name,
				IpRange:          entities.GenerateUUID(),
				ZtNetworkId:      entities.GenerateUUID(),
			}
			connectionAdded, err := client.AddConnection(context.Background(), addConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(connectionAdded).ToNot(gomega.BeNil())

			connectionInstance, err := client.GetConnection(context.Background(), &grpc_application_network_go.ConnectionInstanceId{
				OrganizationId:   addConnectionRequest.OrganizationId,
				SourceInstanceId: addConnectionRequest.SourceInstanceId,
				TargetInstanceId: addConnectionRequest.TargetInstanceId,
				InboundName:      addConnectionRequest.InboundName,
				OutboundName:     addConnectionRequest.OutboundName,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(connectionInstance).ToNot(gomega.BeNil())
			gomega.Expect(connectionInstance).To(gomega.Equal(connectionAdded))
		})

		ginkgo.It("should return a not found error when trying to get a non existent connection", func() {
			connectionInstance, err := client.GetConnection(context.Background(), &grpc_application_network_go.ConnectionInstanceId{
				OrganizationId:   "",
				SourceInstanceId: "",
				TargetInstanceId: "",
				InboundName:      "",
				OutboundName:     "",
			})
			gomega.Expect(err).ToNot(gomega.Succeed())
			gomega.Expect(connectionInstance).To(gomega.BeNil())
		})

		ginkgo.It("should retrieve a previously added connection using the ZT network id", func() {
			organization := addOrganization(organizationProvider)
			sourceInstance := addSourceInstance(organization.ID, false, applicationProvider)
			targetInstance := addTargetInstance(organization.ID, applicationProvider)

			addConnectionRequest := &grpc_application_network_go.AddConnectionRequest{
				OrganizationId:   organization.ID,
				SourceInstanceId: sourceInstance.AppInstanceId,
				TargetInstanceId: targetInstance.AppInstanceId,
				InboundName:      targetInstance.InboundNetInterfaces[0].Name,
				OutboundName:     sourceInstance.OutboundNetInterfaces[0].Name,
				IpRange:          entities.GenerateUUID(),
				ZtNetworkId:      entities.GenerateUUID(),
			}
			connectionAdded, err := client.AddConnection(context.Background(), addConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(connectionAdded).ToNot(gomega.BeNil())

			connectionInstance, err := client.GetConnectionByZtNetworkId(context.Background(), &grpc_application_network_go.ZTNetworkId{
				OrganizationId: addConnectionRequest.OrganizationId,
				ZtNetworkId:    addConnectionRequest.ZtNetworkId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(connectionInstance).ToNot(gomega.BeNil())
			gomega.Expect(connectionInstance).To(gomega.Equal(connectionAdded))
		})

		ginkgo.It("should return a not found error when trying to get a non existent connection", func() {
			connectionInstance, err := client.GetConnectionByZtNetworkId(context.Background(), &grpc_application_network_go.ZTNetworkId{
				OrganizationId: "",
				ZtNetworkId:    "",
			})
			gomega.Expect(err).ToNot(gomega.Succeed())
			gomega.Expect(connectionInstance).To(gomega.BeNil())
		})
	})

	ginkgo.Context("when updating a connection", func() {
		ginkgo.It("should be able to update the status of the connection", func() {
			organization := addOrganization(organizationProvider)
			sourceInstance := addSourceInstance(organization.ID, false, applicationProvider)
			targetInstance := addTargetInstance(organization.ID, applicationProvider)

			addConnectionRequest := &grpc_application_network_go.AddConnectionRequest{
				OrganizationId:   organization.ID,
				SourceInstanceId: sourceInstance.AppInstanceId,
				TargetInstanceId: targetInstance.AppInstanceId,
				InboundName:      targetInstance.InboundNetInterfaces[0].Name,
				OutboundName:     sourceInstance.OutboundNetInterfaces[0].Name,
				IpRange:          entities.GenerateUUID(),
				ZtNetworkId:      entities.GenerateUUID(),
			}
			_, err := client.AddConnection(context.Background(), addConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			updateConnectionRequest := &grpc_application_network_go.UpdateConnectionRequest{
				OrganizationId:    addConnectionRequest.OrganizationId,
				SourceInstanceId:  addConnectionRequest.SourceInstanceId,
				TargetInstanceId:  addConnectionRequest.TargetInstanceId,
				InboundName:       addConnectionRequest.InboundName,
				OutboundName:      addConnectionRequest.OutboundName,
				UpdateStatus:      true,
				Status:            grpc_application_network_go.ConnectionStatus_ESTABLISHED,
				UpdateIpRange:     false,
				IpRange:           "",
				UpdateZtNetworkId: false,
				ZtNetworkId:       "",
			}
			success, err := client.UpdateConnection(context.Background(), updateConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ToNot(gomega.BeNil())
			connections, err := client.ListConnections(context.Background(), &grpc_organization_go.OrganizationId{OrganizationId: organization.ID})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(connections).ToNot(gomega.BeNil())
			gomega.Expect(connections.Connections).To(gomega.HaveLen(1))
			gomega.Expect(connections.Connections[0].Status).To(gomega.Equal(updateConnectionRequest.Status))
		})

		ginkgo.It("should be able to update the range IP of the connection", func() {
			organization := addOrganization(organizationProvider)
			sourceInstance := addSourceInstance(organization.ID, false, applicationProvider)
			targetInstance := addTargetInstance(organization.ID, applicationProvider)

			addConnectionRequest := &grpc_application_network_go.AddConnectionRequest{
				OrganizationId:   organization.ID,
				SourceInstanceId: sourceInstance.AppInstanceId,
				TargetInstanceId: targetInstance.AppInstanceId,
				InboundName:      targetInstance.InboundNetInterfaces[0].Name,
				OutboundName:     sourceInstance.OutboundNetInterfaces[0].Name,
				IpRange:          entities.GenerateUUID(),
				ZtNetworkId:      entities.GenerateUUID(),
			}
			_, err := client.AddConnection(context.Background(), addConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			newRange := "172.16.0.1-172.16.0.255"
			updateConnectionRequest := &grpc_application_network_go.UpdateConnectionRequest{
				OrganizationId:    addConnectionRequest.OrganizationId,
				SourceInstanceId:  addConnectionRequest.SourceInstanceId,
				TargetInstanceId:  addConnectionRequest.TargetInstanceId,
				InboundName:       addConnectionRequest.InboundName,
				OutboundName:      addConnectionRequest.OutboundName,
				UpdateStatus:      false,
				Status:            0,
				UpdateIpRange:     true,
				IpRange:           newRange,
				UpdateZtNetworkId: false,
				ZtNetworkId:       "",
			}
			success, err := client.UpdateConnection(context.Background(), updateConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ToNot(gomega.BeNil())
			connections, err := client.ListConnections(context.Background(), &grpc_organization_go.OrganizationId{OrganizationId: organization.ID})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(connections).ToNot(gomega.BeNil())
			gomega.Expect(connections.Connections).To(gomega.HaveLen(1))
			gomega.Expect(connections.Connections[0].IpRange).To(gomega.Equal(updateConnectionRequest.IpRange))
		})

		ginkgo.It("should be able to update the ZtNetwork ID of the connection", func() {
			organization := addOrganization(organizationProvider)
			sourceInstance := addSourceInstance(organization.ID, false, applicationProvider)
			targetInstance := addTargetInstance(organization.ID, applicationProvider)

			addConnectionRequest := &grpc_application_network_go.AddConnectionRequest{
				OrganizationId:   organization.ID,
				SourceInstanceId: sourceInstance.AppInstanceId,
				TargetInstanceId: targetInstance.AppInstanceId,
				InboundName:      targetInstance.InboundNetInterfaces[0].Name,
				OutboundName:     sourceInstance.OutboundNetInterfaces[0].Name,
				IpRange:          entities.GenerateUUID(),
				ZtNetworkId:      entities.GenerateUUID(),
			}
			_, err := client.AddConnection(context.Background(), addConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			updateConnectionRequest := &grpc_application_network_go.UpdateConnectionRequest{
				OrganizationId:    addConnectionRequest.OrganizationId,
				SourceInstanceId:  addConnectionRequest.SourceInstanceId,
				TargetInstanceId:  addConnectionRequest.TargetInstanceId,
				InboundName:       addConnectionRequest.InboundName,
				OutboundName:      addConnectionRequest.OutboundName,
				UpdateStatus:      false,
				Status:            0,
				UpdateIpRange:     false,
				IpRange:           "",
				UpdateZtNetworkId: true,
				ZtNetworkId:       entities.GenerateUUID(),
			}
			success, err := client.UpdateConnection(context.Background(), updateConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ToNot(gomega.BeNil())
			connections, err := client.ListConnections(context.Background(), &grpc_organization_go.OrganizationId{OrganizationId: organization.ID})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(connections).ToNot(gomega.BeNil())
			gomega.Expect(connections.Connections).To(gomega.HaveLen(1))
			gomega.Expect(connections.Connections[0].ZtNetworkId).To(gomega.Equal(updateConnectionRequest.ZtNetworkId))
		})

		ginkgo.It("should be able to update the range IP and the status of the connection", func() {
			organization := addOrganization(organizationProvider)
			sourceInstance := addSourceInstance(organization.ID, false, applicationProvider)
			targetInstance := addTargetInstance(organization.ID, applicationProvider)

			addConnectionRequest := &grpc_application_network_go.AddConnectionRequest{
				OrganizationId:   organization.ID,
				SourceInstanceId: sourceInstance.AppInstanceId,
				TargetInstanceId: targetInstance.AppInstanceId,
				InboundName:      targetInstance.InboundNetInterfaces[0].Name,
				OutboundName:     sourceInstance.OutboundNetInterfaces[0].Name,
				IpRange:          entities.GenerateUUID(),
				ZtNetworkId:      entities.GenerateUUID(),
			}
			_, err := client.AddConnection(context.Background(), addConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			newRange := "172.16.0.1-172.16.0.255"
			updateConnectionRequest := &grpc_application_network_go.UpdateConnectionRequest{
				OrganizationId:    addConnectionRequest.OrganizationId,
				SourceInstanceId:  addConnectionRequest.SourceInstanceId,
				TargetInstanceId:  addConnectionRequest.TargetInstanceId,
				InboundName:       addConnectionRequest.InboundName,
				OutboundName:      addConnectionRequest.OutboundName,
				UpdateStatus:      true,
				Status:            grpc_application_network_go.ConnectionStatus_ESTABLISHED,
				UpdateIpRange:     true,
				IpRange:           newRange,
				UpdateZtNetworkId: true,
				ZtNetworkId:       entities.GenerateUUID(),
			}
			success, err := client.UpdateConnection(context.Background(), updateConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ToNot(gomega.BeNil())
			connections, err := client.ListConnections(context.Background(), &grpc_organization_go.OrganizationId{OrganizationId: organization.ID})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(connections).ToNot(gomega.BeNil())
			gomega.Expect(connections.Connections).To(gomega.HaveLen(1))
			gomega.Expect(connections.Connections[0].Status).To(gomega.Equal(updateConnectionRequest.Status))
			gomega.Expect(connections.Connections[0].IpRange).To(gomega.Equal(updateConnectionRequest.IpRange))
			gomega.Expect(connections.Connections[0].ZtNetworkId).To(gomega.Equal(updateConnectionRequest.ZtNetworkId))
		})
	})

	ginkgo.Context("when removing a connection", func() {

		ginkgo.It("should support removing an existing connection with no required outbound", func() {
			organization := addOrganization(organizationProvider)
			sourceInstance := addSourceInstance(organization.ID, false, applicationProvider)
			targetInstance := addTargetInstance(organization.ID, applicationProvider)

			addConnectionRequest := &grpc_application_network_go.AddConnectionRequest{
				OrganizationId:   organization.ID,
				SourceInstanceId: sourceInstance.AppInstanceId,
				TargetInstanceId: targetInstance.AppInstanceId,
				InboundName:      targetInstance.InboundNetInterfaces[0].Name,
				OutboundName:     sourceInstance.OutboundNetInterfaces[0].Name,
			}
			connectionInstance, _ := client.AddConnection(context.Background(), addConnectionRequest)
			removeConnectionRequest := &grpc_application_network_go.RemoveConnectionRequest{
				OrganizationId:   connectionInstance.OrganizationId,
				SourceInstanceId: connectionInstance.SourceInstanceId,
				TargetInstanceId: connectionInstance.TargetInstanceId,
				InboundName:      connectionInstance.InboundName,
				OutboundName:     connectionInstance.OutboundName,
				UserConfirmation: false,
			}
			success, err := client.RemoveConnection(context.Background(), removeConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ToNot(gomega.BeNil())
		})

		ginkgo.It("should support removing an existing connection with required outbound and with positive user confirmation", func() {
			organization := addOrganization(organizationProvider)
			sourceInstance := addSourceInstance(organization.ID, true, applicationProvider)
			targetInstance := addTargetInstance(organization.ID, applicationProvider)

			addConnectionRequest := &grpc_application_network_go.AddConnectionRequest{
				OrganizationId:   organization.ID,
				SourceInstanceId: sourceInstance.AppInstanceId,
				TargetInstanceId: targetInstance.AppInstanceId,
				InboundName:      targetInstance.InboundNetInterfaces[0].Name,
				OutboundName:     sourceInstance.OutboundNetInterfaces[0].Name,
			}
			connectionInstance, _ := client.AddConnection(context.Background(), addConnectionRequest)
			removeConnectionRequest := &grpc_application_network_go.RemoveConnectionRequest{
				OrganizationId:   connectionInstance.OrganizationId,
				SourceInstanceId: connectionInstance.SourceInstanceId,
				TargetInstanceId: connectionInstance.TargetInstanceId,
				InboundName:      connectionInstance.InboundName,
				OutboundName:     connectionInstance.OutboundName,
				UserConfirmation: true,
			}
			success, err := client.RemoveConnection(context.Background(), removeConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ToNot(gomega.BeNil())
		})

		ginkgo.It("should fail when removing a not existing connection", func() {
			organization := addOrganization(organizationProvider)
			sourceInstance := addSourceInstance(organization.ID, false, applicationProvider)
			targetInstance := addTargetInstance(organization.ID, applicationProvider)

			removeConnectionRequest := &grpc_application_network_go.RemoveConnectionRequest{
				OrganizationId:   organization.ID,
				SourceInstanceId: sourceInstance.AppInstanceId,
				TargetInstanceId: targetInstance.AppInstanceId,
				InboundName:      targetInstance.InboundNetInterfaces[0].Name,
				OutboundName:     sourceInstance.OutboundNetInterfaces[0].Name,
				UserConfirmation: false,
			}
			success, err := client.RemoveConnection(context.Background(), removeConnectionRequest)
			gomega.Expect(err).ToNot(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())
		})

		ginkgo.It("should fail when removing a connection with required outbound but negative user confirmation", func() {
			organization := addOrganization(organizationProvider)
			sourceInstance := addSourceInstance(organization.ID, true, applicationProvider)
			targetInstance := addTargetInstance(organization.ID, applicationProvider)

			addConnectionRequest := &grpc_application_network_go.AddConnectionRequest{
				OrganizationId:   organization.ID,
				SourceInstanceId: sourceInstance.AppInstanceId,
				TargetInstanceId: targetInstance.AppInstanceId,
				InboundName:      targetInstance.InboundNetInterfaces[0].Name,
				OutboundName:     sourceInstance.OutboundNetInterfaces[0].Name,
			}
			connectionInstance, err := client.AddConnection(context.Background(), addConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			removeConnectionRequest := &grpc_application_network_go.RemoveConnectionRequest{
				OrganizationId:   connectionInstance.OrganizationId,
				SourceInstanceId: connectionInstance.SourceInstanceId,
				TargetInstanceId: connectionInstance.TargetInstanceId,
				InboundName:      connectionInstance.InboundName,
				OutboundName:     connectionInstance.OutboundName,
				UserConfirmation: false,
			}
			success, err := client.RemoveConnection(context.Background(), removeConnectionRequest)
			gomega.Expect(err).ToNot(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())
		})
	})

	ginkgo.Context("when listing connections", func() {

		ginkgo.It("should list all the connections in an organization", func() {
			organization := addOrganization(organizationProvider)
			instanceA := addSourceInstance(organization.ID, false, applicationProvider)
			instanceB := addSourceInstance(organization.ID, false, applicationProvider)
			instanceC := addSourceInstance(organization.ID, false, applicationProvider)

			instance1 := addTargetInstance(organization.ID, applicationProvider)
			instance2 := addTargetInstance(organization.ID, applicationProvider)
			instance3 := addTargetInstance(organization.ID, applicationProvider)

			addConnectionRequests := []*grpc_application_network_go.AddConnectionRequest{
				{
					OrganizationId:   organization.ID,
					SourceInstanceId: instanceA.AppInstanceId,
					TargetInstanceId: instance1.AppInstanceId,
					InboundName:      instance1.InboundNetInterfaces[0].Name,
					OutboundName:     instanceA.OutboundNetInterfaces[0].Name,
				},
				{
					OrganizationId:   organization.ID,
					SourceInstanceId: instanceA.AppInstanceId,
					TargetInstanceId: instance2.AppInstanceId,
					InboundName:      instance2.InboundNetInterfaces[0].Name,
					OutboundName:     instanceA.OutboundNetInterfaces[0].Name,
				},
				{
					OrganizationId:   organization.ID,
					SourceInstanceId: instanceA.AppInstanceId,
					TargetInstanceId: instance3.AppInstanceId,
					InboundName:      instance3.InboundNetInterfaces[0].Name,
					OutboundName:     instanceA.OutboundNetInterfaces[0].Name,
				},
				{
					OrganizationId:   organization.ID,
					SourceInstanceId: instanceB.AppInstanceId,
					TargetInstanceId: instance1.AppInstanceId,
					InboundName:      instance1.InboundNetInterfaces[0].Name,
					OutboundName:     instanceB.OutboundNetInterfaces[0].Name,
				},
				{
					OrganizationId:   organization.ID,
					SourceInstanceId: instanceB.AppInstanceId,
					TargetInstanceId: instance2.AppInstanceId,
					InboundName:      instance2.InboundNetInterfaces[0].Name,
					OutboundName:     instanceB.OutboundNetInterfaces[0].Name,
				},
				{
					OrganizationId:   organization.ID,
					SourceInstanceId: instanceC.AppInstanceId,
					TargetInstanceId: instance1.AppInstanceId,
					InboundName:      instance1.InboundNetInterfaces[0].Name,
					OutboundName:     instanceC.OutboundNetInterfaces[0].Name,
				},
			}
			for _, addConnectionRequest := range addConnectionRequests {
				_, err := client.AddConnection(context.Background(), addConnectionRequest)
				gomega.Expect(err).To(gomega.Succeed())
			}

			connections, err := client.ListConnections(context.Background(), &grpc_organization_go.OrganizationId{OrganizationId: organization.ID})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(connections.Connections).To(gomega.HaveLen(len(addConnectionRequests)))
			sort.Sort(connectionRequests(addConnectionRequests))
			sort.Sort(connectionInstances(connections.Connections))
			for i, connection := range connections.Connections {
				gomega.Expect(connection.OrganizationId).To(gomega.Equal(addConnectionRequests[i].OrganizationId))
				gomega.Expect(connection.SourceInstanceId).To(gomega.Equal(addConnectionRequests[i].SourceInstanceId))
				gomega.Expect(connection.TargetInstanceId).To(gomega.Equal(addConnectionRequests[i].TargetInstanceId))
				gomega.Expect(connection.InboundName).To(gomega.Equal(addConnectionRequests[i].InboundName))
				gomega.Expect(connection.OutboundName).To(gomega.Equal(addConnectionRequests[i].OutboundName))
			}
		})

		ginkgo.It("should retrieve an empty list if there are no connections on the organization", func() {
			organization := addOrganization(organizationProvider)
			connections, err := client.ListConnections(context.Background(), &grpc_organization_go.OrganizationId{OrganizationId: organization.ID})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(connections.Connections).To(gomega.BeEmpty())
		})
	})

	ginkgo.Context("when listing inbound connections", func() {

		ginkgo.It("should list all the inbound connections of an instance", func() {
			organization := addOrganization(organizationProvider)
			instanceA := addSourceInstance(organization.ID, false, applicationProvider)
			instanceB := addSourceInstance(organization.ID, false, applicationProvider)
			instanceC := addSourceInstance(organization.ID, false, applicationProvider)

			instance1 := addTargetInstance(organization.ID, applicationProvider)
			instance2 := addTargetInstance(organization.ID, applicationProvider)

			addConnectionRequests := []*grpc_application_network_go.AddConnectionRequest{
				{
					OrganizationId:   organization.ID,
					SourceInstanceId: instanceA.AppInstanceId,
					TargetInstanceId: instance1.AppInstanceId,
					InboundName:      instance1.InboundNetInterfaces[0].Name,
					OutboundName:     instanceA.OutboundNetInterfaces[0].Name,
				},
				{
					OrganizationId:   organization.ID,
					SourceInstanceId: instanceB.AppInstanceId,
					TargetInstanceId: instance1.AppInstanceId,
					InboundName:      instance1.InboundNetInterfaces[0].Name,
					OutboundName:     instanceB.OutboundNetInterfaces[0].Name,
				},
				{
					OrganizationId:   organization.ID,
					SourceInstanceId: instanceC.AppInstanceId,
					TargetInstanceId: instance1.AppInstanceId,
					InboundName:      instance1.InboundNetInterfaces[0].Name,
					OutboundName:     instanceC.OutboundNetInterfaces[0].Name,
				},
				{
					OrganizationId:   organization.ID,
					SourceInstanceId: instanceC.AppInstanceId,
					TargetInstanceId: instance2.AppInstanceId,
					InboundName:      instance2.InboundNetInterfaces[0].Name,
					OutboundName:     instanceC.OutboundNetInterfaces[0].Name,
				},
			}
			for _, addConnectionRequest := range addConnectionRequests {
				_, err := client.AddConnection(context.Background(), addConnectionRequest)
				gomega.Expect(err).To(gomega.Succeed())
			}

			connections, err := client.ListInboundConnections(context.Background(),
				&grpc_application_go.AppInstanceId{
					OrganizationId: organization.ID,
					AppInstanceId:  instance1.AppInstanceId,
				})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(connections).NotTo(gomega.BeNil())
			gomega.Expect(len(connections.Connections)).Should(gomega.Equal(3))
		})

		ginkgo.It("should not be able to list inbound connections if the organization does not exist", func() {
			_, err := client.ListInboundConnections(context.Background(),
				&grpc_application_go.AppInstanceId{
					OrganizationId: entities.GenerateUUID(),
					AppInstanceId:  entities.GenerateUUID(),
				})
			gomega.Expect(err).NotTo(gomega.Succeed())

		})

		ginkgo.It("should not be able to list inbound connections if the instance does not exist", func() {
			organization := addOrganization(organizationProvider)

			_, err := client.ListInboundConnections(context.Background(),
				&grpc_application_go.AppInstanceId{
					OrganizationId: organization.ID,
					AppInstanceId:  entities.GenerateUUID(),
				})
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
	})

	ginkgo.Context("when listing outbound connections", func() {

		ginkgo.It("should list all the outbound connections of an instance", func() {
			organization := addOrganization(organizationProvider)
			instanceA := addTargetInstance(organization.ID, applicationProvider)
			instanceB := addTargetInstance(organization.ID, applicationProvider)
			instanceC := addTargetInstance(organization.ID, applicationProvider)

			instance1 := addSourceInstance(organization.ID, false, applicationProvider)
			instance2 := addSourceInstance(organization.ID, false, applicationProvider)

			addConnectionRequests := []*grpc_application_network_go.AddConnectionRequest{
				{
					OrganizationId:   organization.ID,
					TargetInstanceId: instanceA.AppInstanceId,
					SourceInstanceId: instance1.AppInstanceId,
					OutboundName:     instance1.OutboundNetInterfaces[0].Name,
					InboundName:      instanceA.InboundNetInterfaces[0].Name,
				},
				{
					OrganizationId:   organization.ID,
					TargetInstanceId: instanceB.AppInstanceId,
					SourceInstanceId: instance1.AppInstanceId,
					OutboundName:     instance1.OutboundNetInterfaces[0].Name,
					InboundName:      instanceB.InboundNetInterfaces[0].Name,
				},
				{
					OrganizationId:   organization.ID,
					TargetInstanceId: instanceC.AppInstanceId,
					SourceInstanceId: instance1.AppInstanceId,
					OutboundName:     instance1.OutboundNetInterfaces[0].Name,
					InboundName:      instanceC.InboundNetInterfaces[0].Name,
				},
				{
					OrganizationId:   organization.ID,
					TargetInstanceId: instanceC.AppInstanceId,
					SourceInstanceId: instance2.AppInstanceId,
					OutboundName:     instance2.OutboundNetInterfaces[0].Name,
					InboundName:      instanceC.InboundNetInterfaces[0].Name,
				},
			}
			for _, addConnectionRequest := range addConnectionRequests {
				_, err := client.AddConnection(context.Background(), addConnectionRequest)
				gomega.Expect(err).To(gomega.Succeed())
			}

			connections, err := client.ListOutboundConnections(context.Background(),
				&grpc_application_go.AppInstanceId{
					OrganizationId: organization.ID,
					AppInstanceId:  instance1.AppInstanceId,
				})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(connections).NotTo(gomega.BeNil())
			gomega.Expect(len(connections.Connections)).Should(gomega.Equal(3))
		})

		ginkgo.It("should not be able to list outbouns connections if the organization does not exist", func() {
			_, err := client.ListOutboundConnections(context.Background(),
				&grpc_application_go.AppInstanceId{
					OrganizationId: entities.GenerateUUID(),
					AppInstanceId:  entities.GenerateUUID(),
				})
			gomega.Expect(err).NotTo(gomega.Succeed())

		})

		ginkgo.It("should not be able to list inbound connections if the instance does not exist", func() {
			organization := addOrganization(organizationProvider)

			_, err := client.ListOutboundConnections(context.Background(),
				&grpc_application_go.AppInstanceId{
					OrganizationId: organization.ID,
					AppInstanceId:  entities.GenerateUUID(),
				})
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
	})

	ginkgo.Context("when adding a ztConnection", func() {
		ginkgo.It("should be able to add when organization and instance exists", func() {
			organization := addOrganization(organizationProvider)
			instance := addInstance(organization.ID, applicationProvider)
			toAdd := &grpc_application_network_go.ZTNetworkConnection{
				OrganizationId: organization.ID,
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  instance.AppInstanceId,
				ServiceId:      instance.Groups[0].ServiceInstances[0].ServiceId,
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           grpc_application_network_go.ConnectionSide_SIDE_OUTBOUND,
			}
			added, err := client.AddZTNetworkConnection(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).NotTo(gomega.BeNil())
		})
		ginkgo.It("should not be able to add when organization does not exists", func() {
			toAdd := &grpc_application_network_go.ZTNetworkConnection{
				OrganizationId: entities.GenerateUUID(),
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  entities.GenerateUUID(),
				ServiceId:      entities.GenerateUUID(),
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           grpc_application_network_go.ConnectionSide_SIDE_OUTBOUND,
			}
			_, err := client.AddZTNetworkConnection(context.Background(), toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		ginkgo.It("should not be able to add a ztConnection twice", func() {
			organization := addOrganization(organizationProvider)
			instance := addInstance(organization.ID, applicationProvider)
			toAdd := &grpc_application_network_go.ZTNetworkConnection{
				OrganizationId: organization.ID,
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  instance.AppInstanceId,
				ServiceId:      instance.Groups[0].ServiceInstances[0].ServiceId,
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           grpc_application_network_go.ConnectionSide_SIDE_OUTBOUND,
			}
			_, err := client.AddZTNetworkConnection(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			_, err = client.AddZTNetworkConnection(context.Background(), toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})

	ginkgo.Context("when listing ztConnections", func() {
		ginkgo.It("Should be able to list all the ZT-Connecitions in a zt-Network", func() {
			num := 5
			organization := addOrganization(organizationProvider)
			toAdd := &grpc_application_network_go.ZTNetworkConnection{
				OrganizationId: organization.ID,
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  entities.GenerateUUID(),
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           grpc_application_network_go.ConnectionSide_SIDE_OUTBOUND,
			}
			for i := 0; i < num; i++ {
				instance := addInstance(organization.ID, applicationProvider)
				toAdd.AppInstanceId = instance.AppInstanceId
				toAdd.ServiceId = instance.Groups[0].ServiceInstances[0].ServiceId

				_, err := client.AddZTNetworkConnection(context.Background(), toAdd)
				gomega.Expect(err).To(gomega.Succeed())
			}
			list, err := client.ListZTNetworkConnection(context.Background(), &grpc_application_network_go.ZTNetworkId{
				OrganizationId: toAdd.OrganizationId,
				ZtNetworkId:    toAdd.ZtNetworkId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(list).NotTo(gomega.BeNil())
			gomega.Expect(len(list.Connections)).Should(gomega.Equal(num))

		})
		ginkgo.It("Should be able to return an empty list if there is no connections", func() {
			organization := addOrganization(organizationProvider)
			list, err := client.ListZTNetworkConnection(context.Background(), &grpc_application_network_go.ZTNetworkId{
				OrganizationId: organization.ID,
				ZtNetworkId:    entities.GenerateUUID(),
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(list).NotTo(gomega.BeNil())
			gomega.Expect(len(list.Connections)).Should(gomega.Equal(0))

		})
		ginkgo.It("Should not be able to return a list if there the organization does not exist", func() {
			_, err := client.ListZTNetworkConnection(context.Background(), &grpc_application_network_go.ZTNetworkId{
				OrganizationId: entities.GenerateUUID(),
				ZtNetworkId:    entities.GenerateUUID(),
			})
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})
	ginkgo.Context("when updating ztConnections", func() {
		ginkgo.It("should be able to update a ztconnection", func() {
			organization := addOrganization(organizationProvider)
			instance := addInstance(organization.ID, applicationProvider)
			toAdd := &grpc_application_network_go.ZTNetworkConnection{
				OrganizationId: organization.ID,
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  instance.AppInstanceId,
				ServiceId:      instance.Groups[0].ServiceInstances[0].ServiceId,
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           grpc_application_network_go.ConnectionSide_SIDE_OUTBOUND,
			}
			added, err := client.AddZTNetworkConnection(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).NotTo(gomega.BeNil())

			toUpdate := grpc_application_network_go.UpdateZTNetworkConnectionRequest{
				OrganizationId: toAdd.OrganizationId,
				ZtNetworkId:    toAdd.ZtNetworkId,
				AppInstanceId:  toAdd.AppInstanceId,
				ServiceId:      toAdd.ServiceId,
				ClusterId:      toAdd.ClusterId,
				UpdateZtIp:     true,
				ZtIp:           "yyy.yyy.yyy.yyy",
			}
			success, err := client.UpdateZTNetworkConnection(context.Background(), &toUpdate)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())

			list, err := client.ListZTNetworkConnection(context.Background(), &grpc_application_network_go.ZTNetworkId{
				OrganizationId: toAdd.OrganizationId,
				ZtNetworkId:    toAdd.ZtNetworkId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(list).NotTo(gomega.BeNil())
			gomega.Expect(len(list.Connections)).Should(gomega.Equal(1))
			gomega.Expect(list.Connections[0].ZtIp).Should(gomega.Equal(toUpdate.ZtIp))
			gomega.Expect(list.Connections[0].Side).Should(gomega.Equal(toAdd.Side))

		})
		ginkgo.It("should not be able to update a non existing ztconnection", func() {
			organization := addOrganization(organizationProvider)
			instance := addInstance(organization.ID, applicationProvider)
			toUpdate := grpc_application_network_go.UpdateZTNetworkConnectionRequest{
				OrganizationId: organization.ID,
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  instance.AppInstanceId,
				ServiceId:      instance.Groups[0].ServiceInstances[0].ServiceId,
				UpdateZtIp:     true,
				ZtIp:           "yyy.yyy.yyy.yyy",
			}
			success, err := client.UpdateZTNetworkConnection(context.Background(), &toUpdate)
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())

		})

	})
	ginkgo.Context("when removing ztConnections", func() {
		ginkgo.It("should be able to remove a ztconnection", func() {
			organization := addOrganization(organizationProvider)
			instance := addInstance(organization.ID, applicationProvider)
			toAdd := &grpc_application_network_go.ZTNetworkConnection{
				OrganizationId: organization.ID,
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  instance.AppInstanceId,
				ServiceId:      instance.Groups[0].ServiceInstances[0].ServiceId,
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           grpc_application_network_go.ConnectionSide_SIDE_OUTBOUND,
			}
			added, err := client.AddZTNetworkConnection(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).NotTo(gomega.BeNil())

			instance2 := addInstance(organization.ID, applicationProvider)
			toAdd.AppInstanceId = instance2.AppInstanceId
			toAdd.ServiceId = instance2.Groups[0].ServiceInstances[0].ServiceId

			added, err = client.AddZTNetworkConnection(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).NotTo(gomega.BeNil())

			success, err := client.RemoveZTNetworkConnection(context.Background(), &grpc_application_network_go.ZTNetworkConnectionId{
				OrganizationId: toAdd.OrganizationId,
				ZtNetworkId:    toAdd.ZtNetworkId,
				AppInstanceId:  toAdd.AppInstanceId,
				ServiceId:      toAdd.ServiceId,
				ClusterId:      toAdd.ClusterId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())
			// check there is no ztconnections
			list, err := client.ListZTNetworkConnection(context.Background(), &grpc_application_network_go.ZTNetworkId{
				OrganizationId: toAdd.OrganizationId,
				ZtNetworkId:    toAdd.ZtNetworkId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(list).NotTo(gomega.BeNil())
			gomega.Expect(len(list.Connections)).Should(gomega.Equal(1))

		})
		ginkgo.It("should be able to remove all the ztconnection of a network", func() {
			organization := addOrganization(organizationProvider)
			instance := addInstance(organization.ID, applicationProvider)
			toAdd := &grpc_application_network_go.ZTNetworkConnection{
				OrganizationId: organization.ID,
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  instance.AppInstanceId,
				ServiceId:      instance.Groups[0].ServiceInstances[0].ServiceId,
				ZtMember:       entities.GenerateUUID(),
				ZtIp:           "xxx.xxx.xxx.xxx",
				ClusterId:      entities.GenerateUUID(),
				Side:           grpc_application_network_go.ConnectionSide_SIDE_OUTBOUND,
			}
			added, err := client.AddZTNetworkConnection(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).NotTo(gomega.BeNil())

			instance2 := addInstance(organization.ID, applicationProvider)
			toAdd.AppInstanceId = instance2.AppInstanceId
			toAdd.ServiceId = instance2.Groups[0].ServiceInstances[0].ServiceId

			added, err = client.AddZTNetworkConnection(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).NotTo(gomega.BeNil())

			success, err := client.RemoveZTNetworkConnectionByNetworkId(context.Background(), &grpc_application_network_go.ZTNetworkId{
				OrganizationId: toAdd.OrganizationId,
				ZtNetworkId:    toAdd.ZtNetworkId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())
			// check there is no ztconnections
			list, err := client.ListZTNetworkConnection(context.Background(), &grpc_application_network_go.ZTNetworkId{
				OrganizationId: toAdd.OrganizationId,
				ZtNetworkId:    toAdd.ZtNetworkId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(list).NotTo(gomega.BeNil())
			gomega.Expect(len(list.Connections)).Should(gomega.Equal(0))

		})
		ginkgo.It("should not be able to remove a non existing ztconnection", func() {
			organization := addOrganization(organizationProvider)
			success, err := client.RemoveZTNetworkConnection(context.Background(), &grpc_application_network_go.ZTNetworkConnectionId{
				OrganizationId: organization.ID,
				ZtNetworkId:    entities.GenerateUUID(),
				AppInstanceId:  entities.GenerateUUID(),
				ServiceId:      entities.GenerateUUID(),
				ClusterId:      entities.GenerateUUID(),
			})
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())
		})
	})

})

type connectionInstances []*grpc_application_network_go.ConnectionInstance

func (c connectionInstances) Len() int      { return len(c) }
func (c connectionInstances) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c connectionInstances) Less(i, j int) bool {
	leftHash := c[i].SourceInstanceId + c[i].TargetInstanceId + c[i].InboundName + c[i].OutboundName
	rightHash := c[j].SourceInstanceId + c[j].TargetInstanceId + c[j].InboundName + c[j].OutboundName
	return leftHash < rightHash
}

type connectionRequests []*grpc_application_network_go.AddConnectionRequest

func (c connectionRequests) Len() int      { return len(c) }
func (c connectionRequests) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c connectionRequests) Less(i, j int) bool {
	leftHash := c[i].SourceInstanceId + c[i].TargetInstanceId + c[i].InboundName + c[i].OutboundName
	rightHash := c[j].SourceInstanceId + c[j].TargetInstanceId + c[j].InboundName + c[j].OutboundName
	return leftHash < rightHash
}
