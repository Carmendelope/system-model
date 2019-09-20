/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package application_network

import (
	"context"
	"github.com/nalej/grpc-application-network-go"
	grpc_organization_go "github.com/nalej/grpc-organization-go"
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
		test.LaunchServer(server, listener)

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
			}
			_, err := client.AddConnection(context.Background(), addConnectionRequest)
			gomega.Expect(err).To(gomega.Succeed())
			connectionInstance, err := client.AddConnection(context.Background(), addConnectionRequest)
			gomega.Expect(err).ToNot(gomega.Succeed())
			gomega.Expect(connectionInstance).To(gomega.BeNil())
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
