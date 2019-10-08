/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package application_network

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-go"
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

func (manager *Manager) UpdateConnectionInstance(updateConnectionRequest *grpc_application_network_go.UpdateConnectionRequest) derrors.Error {
	err := manager.validOrganization(updateConnectionRequest.OrganizationId)
	if err != nil {
		return err
	}

	sourceInstance, err := manager.ApplicationProvider.GetInstance(updateConnectionRequest.SourceInstanceId)
	if err != nil {
		return derrors.NewNotFoundError("sourceInstanceID", err).WithParams(updateConnectionRequest.SourceInstanceId)
	}
	found := false
	for _, iface := range sourceInstance.OutboundNetInterfaces {
		if iface.Name == updateConnectionRequest.OutboundName {
			found = true
		}
	}
	if !found {
		return derrors.NewNotFoundError("outboundName").WithParams(updateConnectionRequest.OutboundName)
	}

	targetInstance, err := manager.ApplicationProvider.GetInstance(updateConnectionRequest.TargetInstanceId)
	if err != nil {
		return derrors.NewNotFoundError("targetInstanceID", err).WithParams(updateConnectionRequest.TargetInstanceId)
	}
	found = false
	for _, iface := range targetInstance.InboundNetInterfaces {
		if iface.Name == updateConnectionRequest.InboundName {
			found = true
		}
	}
	if !found {
		return derrors.NewNotFoundError("inboundName").WithParams(updateConnectionRequest.OutboundName)
	}

	connectionInstance, err := manager.AppNetProvider.GetConnectionInstance(
		updateConnectionRequest.OrganizationId,
		updateConnectionRequest.SourceInstanceId,
		updateConnectionRequest.TargetInstanceId,
		updateConnectionRequest.InboundName,
		updateConnectionRequest.OutboundName)
	if err != nil {
		return derrors.NewNotFoundError("connectionInstance").WithParams(updateConnectionRequest)
	}

	connectionInstance.ApplyUpdate(updateConnectionRequest)

	if err = manager.AppNetProvider.UpdateConnectionInstance(*connectionInstance); err != nil {
		return err
	}
	return nil
}
/*
func (manager *Manager) GetConnectionInstance(connectionId grpc_application_network_go.ConnectionInstanceId) (entities.ConnectionInstance, derrors.Error) {
	err := manager.validOrganization(connectionId.OrganizationId)
	if err != nil {
		return err
	}

	return  manager.AppNetProvider.GetConnectionInstance(connectionId.OrganizationId, connectionId.SourceInstanceId, connectionId.TargetInstanceId,
		connectionId.InboundName, connectionId.OutboundName)

}
*/
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

// validOrganization check if the appInstanceId ID corresponds to an existing instance
func (manager *Manager) validInstance(appInstanceId string) derrors.Error {
	exists, err := manager.ApplicationProvider.InstanceExists(appInstanceId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("appInstanceId").WithParams(appInstanceId)
	}
	return nil
}

// ListInboundConnections retrieves a list with all the connections where the appInstanceId is the target
func (manager *Manager) ListInboundConnections(appInstanceID *grpc_application_go.AppInstanceId) ([]entities.ConnectionInstance, derrors.Error) {

	err := manager.validOrganization(appInstanceID.OrganizationId)
	if err != nil {
		return nil, err
	}

	_, err = manager.ApplicationProvider.GetInstance(appInstanceID.AppInstanceId)
	if err != nil {
		return nil, derrors.NewNotFoundError("appInstance", err).WithParams(appInstanceID.AppInstanceId)
	}

	inboundList, err := manager.AppNetProvider.ListInboundConnections(appInstanceID.OrganizationId, appInstanceID.AppInstanceId)
	if err != nil {
		return nil, err
	}
	return inboundList, nil
}

// ListOutboundConnections retrieves a list with all the connections where the appInstanceId is the source
func (manager *Manager) ListOutboundConnections(appInstanceID *grpc_application_go.AppInstanceId) ([]entities.ConnectionInstance, derrors.Error) {

	err := manager.validOrganization(appInstanceID.OrganizationId)
	if err != nil {
		return nil, err
	}

	_, err = manager.ApplicationProvider.GetInstance(appInstanceID.AppInstanceId)
	if err != nil {
		return nil, derrors.NewNotFoundError("appInstance", err).WithParams(appInstanceID.AppInstanceId)
	}

	outboundList, err := manager.AppNetProvider.ListOutboundConnections(appInstanceID.OrganizationId, appInstanceID.AppInstanceId)
	if err != nil {
		return nil, err
	}
	return outboundList, nil
}


// check if the instance has a service wich its identifier is serviceId
func (manager *Manager) checkServiceId (inst *entities.AppInstance, serviceId string) derrors.Error {
	found := false
	for  i:=0; i<len(inst.Groups) && !found ; i++ { // , group := range inst.Groups {
		for j:=0; j<len(inst.Groups[i].ServiceInstances) && ! found; j++ { //, service := range group.ServiceInstances {
			if inst.Groups[i].ServiceInstances[j].ServiceId == serviceId { // .ServiceId == serviceId{
				found = true
			}
		}
	}
	if !found {
		return derrors.NewNotFoundError("no service found in the instance").WithParams(serviceId, inst.AppInstanceId)
	}
	return nil
}

func (manager *Manager) AddZTNetworkConnection(addRequest *grpc_application_network_go.ZTNetworkConnection) (*entities.ZTNetworkConnection, derrors.Error){

	// check if the organization exists
	err := manager.validOrganization(addRequest.OrganizationId)
	if err != nil {
		return nil, err
	}

	inst, err := manager.ApplicationProvider.GetInstance(addRequest.AppInstanceId)
	if err != nil {
		return nil, err
	}

	err = manager.checkServiceId(inst, addRequest.ServiceId)
	if err != nil {
		return nil, err
	}

	toAdd := entities.NewZTNetworkConnectionFromGRPC(addRequest)
	err = manager.AppNetProvider.AddZTConnection(*toAdd)
	if err != nil {
		return nil, err
	}
	return toAdd, nil
}
// ListZTNetworkConnection lists the connections in one zt network (one inbound and one outbound)
func (manager *Manager) ListZTNetworkConnection(ztNetworkId *grpc_application_network_go.ZTNetworkConnectionId) ([]entities.ZTNetworkConnection, derrors.Error){
	// check if the organization exists
	err := manager.validOrganization(ztNetworkId.OrganizationId)
	if err != nil {
		return nil, err
	}
	list, err := manager.AppNetProvider.ListZTConnections(ztNetworkId.OrganizationId, ztNetworkId.ZtNetworkId)
	if err != nil {
		return nil, err
	}
	return list, nil
}
// UpdateZTNetworkConnection updates an existing zt connection
func (manager *Manager) UpdateZTNetworkConnection(updateRequest *grpc_application_network_go.UpdateZTNetworkConnectionRequest) derrors.Error{
	// check if the organization exists
	err := manager.validOrganization(updateRequest.OrganizationId)
	if err != nil {
		return err
	}

	// check if the instance exists
	err = manager.validInstance(updateRequest.AppInstanceId)
	if err != nil {
		return err
	}

	conn, err := manager.AppNetProvider.GetZTConnection(updateRequest.OrganizationId, updateRequest.ZtNetworkId, updateRequest.AppInstanceId, updateRequest.ServiceId)
	if err != nil {
		return err
	}

	conn.ApplyUpdate(updateRequest)

	return manager.AppNetProvider.UpdateZTConnection(*conn)
}

// Remove ZTNetwork removes the ztNetworkConnection (the inbound and the outbound)
func (manager *Manager) RemoveZTNetworkConnection(ztNetworkId *grpc_application_network_go.ZTNetworkConnectionId) derrors.Error{
	// check if the organization exists
	err := manager.validOrganization(ztNetworkId.OrganizationId)
	if err != nil {
		return err
	}

	return manager.AppNetProvider.RemoveZTConnection(ztNetworkId.OrganizationId, ztNetworkId.ZtNetworkId)

}