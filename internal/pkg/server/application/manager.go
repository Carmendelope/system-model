/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package application

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/application"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/rs/zerolog/log"
)

// Manager structure with the required providers for application operations.
type Manager struct {
	OrgProvider organization.Provider
	AppProvider application.Provider
}

// NewManager creates a Manager using a set of providers.
func NewManager(orgProvider organization.Provider, appProvider application.Provider) Manager {
	return Manager{orgProvider, appProvider}
}

// AddAppDescriptor adds a new application descriptor to a given organization.
func (m * Manager) AddAppDescriptor(addRequest * grpc_application_go.AddAppDescriptorRequest) (* entities.AppDescriptor, derrors.Error) {
	exists, err := m.OrgProvider.Exists(addRequest.OrganizationId)
	if err != nil{
		return nil, err
	}
	if !exists{
		return nil, derrors.NewNotFoundError("organizationID").WithParams(addRequest.OrganizationId)
	}
	descriptor := entities.NewAppDescriptorFromGRPC(addRequest)
	err = m.AppProvider.AddDescriptor(*descriptor)
	if err != nil {
		return nil, err
	}
	err = m.OrgProvider.AddDescriptor(descriptor.OrganizationId, descriptor.AppDescriptorId)
	if err != nil {
	    return nil, err
	}

	return descriptor, nil
}

// ListDescriptors obtains a list of descriptors associated with an organization.
func (m * Manager) ListDescriptors(orgID * grpc_organization_go.OrganizationId) ([] entities.AppDescriptor, derrors.Error) {
	exists, err := m.OrgProvider.Exists(orgID.OrganizationId)
	if err != nil{
		return nil, err
	}

	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(orgID.OrganizationId)
	}
	descriptors, err := m.OrgProvider.ListDescriptors(orgID.OrganizationId)
	if err != nil {
		return nil, err
	}
	result := make([] entities.AppDescriptor, 0)
	for _, dID := range descriptors {
		toAdd, err := m.AppProvider.GetDescriptor(dID)
		if err != nil {
		    return nil, err
		}
		result = append(result, *toAdd)
	}
	return result, nil
}

// GetDescriptor retrieves a single application descriptor.
func (m * Manager) GetDescriptor(appDescID * grpc_application_go.AppDescriptorId) (* entities.AppDescriptor, derrors.Error){
	exists, err := m.OrgProvider.Exists(appDescID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if ! exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(appDescID.OrganizationId)
	}
	exists, err = m.OrgProvider.DescriptorExists(appDescID.OrganizationId, appDescID.AppDescriptorId)
	if err != nil {
		return nil, err
	}
	if !exists{
		return nil, derrors.NewNotFoundError("appDescriptorID").WithParams(appDescID.OrganizationId, appDescID.AppDescriptorId)
	}
	return m.AppProvider.GetDescriptor(appDescID.AppDescriptorId)
}

// RemoveAppDescriptor removes an application descriptor.
func (m * Manager) RemoveAppDescriptor(appDescID *grpc_application_go.AppDescriptorId) derrors.Error {
	exists, err := m.OrgProvider.Exists(appDescID.OrganizationId)
	if err != nil {
		return err
	}
	if ! exists {
		return derrors.NewNotFoundError("organizationID").WithParams(appDescID.OrganizationId)
	}
	exists, err = m.OrgProvider.DescriptorExists(appDescID.OrganizationId, appDescID.AppDescriptorId)
	if err != nil {
		return err
	}
	if ! exists {
		return derrors.NewNotFoundError("appDescriptorId").WithParams(appDescID.AppDescriptorId)
	}
	err = m.OrgProvider.DeleteDescriptor(appDescID.OrganizationId, appDescID.AppDescriptorId)
	if err != nil {
		return err
	}
	err = m.AppProvider.DeleteDescriptor(appDescID.AppDescriptorId)
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("Error removing user. Rollback!")
		rollbackError := m.OrgProvider.AddDescriptor(appDescID.OrganizationId, appDescID.AppDescriptorId)
		if rollbackError != nil {
			log.Error().Str("trace", conversions.ToDerror(rollbackError).DebugReport()).Msg("error in Rollback")
		}
	}
	return err
}

// AddAppInstance adds a new application instance to a given organization.
func (m * Manager) AddAppInstance(addRequest * grpc_application_go.AddAppInstanceRequest) (* entities.AppInstance, derrors.Error) {

	exists, err := m.OrgProvider.Exists(addRequest.OrganizationId)
	if err != nil {
		return nil, err
	}
	if ! exists{
		return nil, derrors.NewNotFoundError("organizationID").WithParams(addRequest.OrganizationId)
	}
	exists, err = m.AppProvider.DescriptorExists(addRequest.AppDescriptorId)
	if err != nil {
		return nil, err
	}
	if ! exists {
		return nil, derrors.NewNotFoundError("descriptorID").WithParams(addRequest.OrganizationId, addRequest.AppDescriptorId)
	}

	descriptor, err := m.AppProvider.GetDescriptor(addRequest.AppDescriptorId)
	if err != nil {
	    return nil, err
	}

	instance := entities.NewAppInstanceFromGRPC(addRequest, descriptor)
	err = m.AppProvider.AddInstance(*instance)
	if err != nil {
		return nil, err
	}
	err = m.OrgProvider.AddInstance(instance.OrganizationId, instance.AppInstanceId)
	if err != nil {
		return nil, err
	}
	return instance, nil
}

// ListInstances retrieves the list of instances associated with an organization.
func (m * Manager) ListInstances(orgID * grpc_organization_go.OrganizationId) ([] entities.AppInstance, derrors.Error) {
	exists, err := m.OrgProvider.Exists(orgID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists{
		return nil, derrors.NewNotFoundError("organizationID").WithParams(orgID.OrganizationId)
	}
	instances, err := m.OrgProvider.ListInstances(orgID.OrganizationId)
	if err != nil {
		return nil, err
	}
	result := make([] entities.AppInstance, 0)
	for _, instID := range instances {
		toAdd, err := m.AppProvider.GetInstance(instID)
		if err != nil {
			return nil, err
		}
		result = append(result, *toAdd)
	}
	return result, nil
}

// GetInstance retrieves a single instance.
func (m * Manager) GetInstance(appInstID * grpc_application_go.AppInstanceId) (* entities.AppInstance, derrors.Error){
	exists, err := m.OrgProvider.Exists(appInstID.OrganizationId)
	if err != nil{
		return nil, err
	}
	if ! exists{
		return nil, derrors.NewNotFoundError("organizationID").WithParams(appInstID.OrganizationId)
	}

	exists, err = m.OrgProvider.InstanceExists(appInstID.OrganizationId, appInstID.AppInstanceId)
	if err != nil {
		return nil, err
	}
	if !exists{
		return nil, derrors.NewNotFoundError("appInstanceID").WithParams(appInstID.OrganizationId, appInstID.AppInstanceId)
	}
	return m.AppProvider.GetInstance(appInstID.AppInstanceId)
}

// UpdateInstance updates the information of a given instance.
func (m * Manager) UpdateInstance(updateRequest * grpc_application_go.UpdateAppStatusRequest) error {
	exists, err := m.OrgProvider.InstanceExists(updateRequest.OrganizationId, updateRequest.AppInstanceId)
	if err != nil {
		return err
	}
	if !exists{
		return derrors.NewNotFoundError("appInstanceID").WithParams(updateRequest.OrganizationId, updateRequest.AppInstanceId)
	}

	toUpdate, err := m.AppProvider.GetInstance(updateRequest.AppInstanceId)
	if err != nil {
		return derrors.NewInternalError("impossible to get old instance", err)
	}

	toUpdate.Status = entities.AppStatusFromGRPC[updateRequest.Status]

	err = m.AppProvider.UpdateInstance(*toUpdate)
	if err != nil {
		return derrors.NewInternalError("impossible to update instance").CausedBy(err)
	}

	return nil
}

// UpdateService updates an application service.
// TODO: wait until JuanMa has the conductor implemented
func (m * Manager) UpdateService(updateRequest * grpc_application_go.UpdateServiceStatusRequest) error {
	/*
	exists, err := m.OrgProvider.InstanceExists(updateRequest.OrganizationId, updateRequest.AppInstanceId)

	if err != nil {
		return err
	}
    if !exists{
        return derrors.NewNotFoundError("appInstanceID").WithParams(updateRequest.OrganizationId, updateRequest.AppInstanceId)
    }
    toUpdate, err := m.AppProvider.GetInstance(updateRequest.AppInstanceId)
    if err != nil {
        return derrors.NewInternalError("impossible to get parent instance", err)
    }

    // find the service instance
    for index, s := range toUpdate.Groups {
        if s.ServiceId == updateRequest.ServiceId {
            toUpdate.Services[index].Status = entities.ServiceStatusFromGRPC[updateRequest.Status]
            toUpdate.Services[index].Endpoints = updateRequest.Endpoints
			toUpdate.Services[index].DeployedOnClusterId = updateRequest.DeployedOnClusterId
            err = m.AppProvider.UpdateInstance(*toUpdate)
            if err != nil {
                return derrors.NewInternalError("impossible to update instance").CausedBy(err)
            }
            return nil
        }
    }

    return derrors.NewInternalError("service not found")
    */
    return derrors.NewUnimplementedError("not implemented yet!")
}

// RemoveAppInstance removes an application instance
func (m * Manager) RemoveAppInstance(appInstID *grpc_application_go.AppInstanceId) derrors.Error {
	exists, err := m.OrgProvider.Exists(appInstID.OrganizationId)
	if err != nil{
		return err
	}
	if ! exists{
		return derrors.NewNotFoundError("organizationID").WithParams(appInstID.OrganizationId)
	}
	exists, err = m.OrgProvider.InstanceExists(appInstID.OrganizationId, appInstID.AppInstanceId)
	if err != nil{
		return err
	}
	if ! exists{
		return derrors.NewNotFoundError("AppInstanceId").WithParams(appInstID.AppInstanceId)
	}
	err = m.OrgProvider.DeleteInstance(appInstID.OrganizationId, appInstID.AppInstanceId)
	if err != nil {
		return err
	}
	err = m.AppProvider.DeleteInstance(appInstID.AppInstanceId)
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("Error removing app Instance. Rollback!")
		rollbackError := m.OrgProvider.AddInstance(appInstID.OrganizationId, appInstID.AppInstanceId)
		if rollbackError != nil {
			log.Error().Str("trace", conversions.ToDerror(rollbackError).DebugReport()).
				Str("appInstID.OrganizationId", appInstID.OrganizationId).
				Str("appInstID.AppInstanceId", appInstID.AppInstanceId).Msg("error in Rollback")
		}
	}
	return err
}
