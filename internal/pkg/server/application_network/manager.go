/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package application_network

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-network-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/application"
	"github.com/nalej/system-model/internal/pkg/provider/application_network"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
)

type Manager struct {
	OrganizationProvider organization.Provider
	ApplicationProvider  application.Provider
	AppNetProvider       application_network.Provider
}

func NewManager(organizationProvider organization.Provider, applicationProvider application.Provider, appNetProvider application_network.Provider) Manager {
	return Manager{
		OrganizationProvider: organizationProvider,
		ApplicationProvider:  applicationProvider,
		AppNetProvider:       appNetProvider,
	}
}

// AddConnectionInstance Adds a new connection instance to an organization
func (manager *Manager) AddConnectionInstance(addConnectionRequest *grpc_application_network_go.AddConnectionRequest) (*entities.ConnectionInstance, derrors.Error) {
	err := manager.validOrganization(addConnectionRequest.OrganizationId)
	if err != nil {
		return nil, err
	}

	sourceInstance, err := manager.ApplicationProvider.GetInstance(addConnectionRequest.SourceInstanceId)
	if err != nil {
		return nil, derrors.NewNotFoundError("sourceInstanceID", err).WithParams(addConnectionRequest.SourceInstanceId)
	}
	var outboundRequired bool
	found := false
	for _, iface := range sourceInstance.OutboundNetInterfaces {
		if iface.Name == addConnectionRequest.OutboundName {
			found = true
			outboundRequired = iface.Required
		}
	}
	if !found {
		return nil, derrors.NewNotFoundError("outboundName").WithParams(addConnectionRequest.OutboundName)
	}

	targetInstance, err := manager.ApplicationProvider.GetInstance(addConnectionRequest.TargetInstanceId)
	if err != nil {
		return nil, derrors.NewNotFoundError("targetInstanceID", err).WithParams(addConnectionRequest.TargetInstanceId)
	}
	found = false
	for _, iface := range targetInstance.InboundNetInterfaces {
		if iface.Name == addConnectionRequest.InboundName {
			found = true
		}
	}
	if !found {
		return nil, derrors.NewNotFoundError("inboundName").WithParams(addConnectionRequest.OutboundName)
	}

	instance := entities.NewConnectionInstanceFromGRPC(
		*addConnectionRequest,
		sourceInstance.Name,
		targetInstance.Name,
		outboundRequired,
	)

	if err = manager.AppNetProvider.AddConnectionInstance(*instance); err != nil {
		return nil, err
	}
	return instance, nil
}

// RemoveConnectionInstance Removes the given connection instance
func (manager *Manager) RemoveConnectionInstance(removeConnectionRequest *grpc_application_network_go.RemoveConnectionRequest) derrors.Error {
	err := manager.validOrganization(removeConnectionRequest.OrganizationId)
	if err != nil {
		return err
	}

	sourceInstance, err := manager.ApplicationProvider.GetInstance(removeConnectionRequest.SourceInstanceId)
	if err != nil {
		return derrors.NewNotFoundError("sourceInstanceID", err).WithParams(removeConnectionRequest.SourceInstanceId)
	}
	var outboundRequired bool
	found := false
	for _, iface := range sourceInstance.OutboundNetInterfaces {
		if iface.Name == removeConnectionRequest.OutboundName {
			found = true
			outboundRequired = iface.Required
		}
	}
	if !found {
		return derrors.NewNotFoundError("outboundName").WithParams(removeConnectionRequest.OutboundName)
	}
	if outboundRequired && !removeConnectionRequest.UserConfirmation {
		return derrors.NewGenericError("outbound connection is required but user did not grant confirmation")
	}

	targetInstance, err := manager.ApplicationProvider.GetInstance(removeConnectionRequest.TargetInstanceId)
	if err != nil {
		return derrors.NewNotFoundError("targetInstanceID", err).WithParams(removeConnectionRequest.TargetInstanceId)
	}
	found = false
	for _, iface := range targetInstance.InboundNetInterfaces {
		if iface.Name == removeConnectionRequest.InboundName {
			found = true
		}
	}
	if !found {
		return derrors.NewNotFoundError("inboundName").WithParams(removeConnectionRequest.OutboundName)
	}

	links, err := manager.AppNetProvider.ListConnectionInstanceLinks(
		removeConnectionRequest.OrganizationId,
		removeConnectionRequest.SourceInstanceId,
		removeConnectionRequest.TargetInstanceId,
		removeConnectionRequest.InboundName,
		removeConnectionRequest.OutboundName,
	)
	if err != nil {
		return err
	}
	if len(links) > 0 {
		return derrors.NewFailedPreconditionError("the connectionInstance still has links associated")
	}

	err = manager.AppNetProvider.RemoveConnectionInstance(
		removeConnectionRequest.OrganizationId,
		removeConnectionRequest.SourceInstanceId,
		removeConnectionRequest.TargetInstanceId,
		removeConnectionRequest.InboundName,
		removeConnectionRequest.OutboundName,
	)
	if err != nil {
		return err
	}
	return nil
}

// ListConnectionInstances Retrieves a list of all the connection instances linked to the given organization
func (manager *Manager) ListConnectionInstances(organizationId *grpc_organization_go.OrganizationId) ([]entities.ConnectionInstance, derrors.Error) {
	if err := manager.validOrganization(organizationId.OrganizationId); err != nil {
		return nil, err
	}

	listConnectionInstances, err := manager.AppNetProvider.ListConnectionInstances(organizationId.OrganizationId)
	if err != nil {
		return nil, err
	}
	return listConnectionInstances, nil
}

// validOrganization check if the organization ID corresponds to an existing organization
func (manager *Manager) validOrganization(orgID string) derrors.Error {
	exists, err := manager.OrganizationProvider.Exists(orgID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organizationID").WithParams(orgID)
	}
	return nil
}
