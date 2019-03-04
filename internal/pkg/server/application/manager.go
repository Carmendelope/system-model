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
	"github.com/nalej/system-model/internal/pkg/provider/device"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/rs/zerolog/log"
)

// Manager structure with the required providers for application operations.
type Manager struct {
	OrgProvider organization.Provider
	AppProvider application.Provider
	DevProvider device.Provider
}

// NewManager creates a Manager using a set of providers.
func NewManager(orgProvider organization.Provider, appProvider application.Provider, devProvider device.Provider) Manager {
	return Manager{orgProvider, appProvider, devProvider}
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

	// -----------------
	// check if the descriptor has device_names in the rules
	// we need to convert deviceGroupNames into deviceGroupIds
	names := make(map[string]bool, 0) // uses a map to avoid insert a device group twice
	for _, rules := range addRequest.Rules{
		if len(rules.DeviceGroupNames) > 0 {
			for _, name := range rules.DeviceGroupNames {
				names[name] = true
			}
		}
	}
	// map to array
	keys := make([]string, len(names))
	i:=0
	for key, _  := range names{
		keys[i] = key
		i += 1
	}

	deviceGroupIds := make (map[string]string, 0) // map of deviceGroupIds indexed by deviceGroupNames
	if len(keys) > 0 {
		deviceGroups, err := m.DevProvider.GetDeviceGroupsByName(addRequest.OrganizationId, keys)
		if err != nil {
			return nil, err
		}
		for _,  deviceGroup := range deviceGroups {
			deviceGroupIds[deviceGroup.Name] = deviceGroup.DeviceGroupId
		}
	}
	// ---------------------

	descriptor := entities.NewAppDescriptorFromGRPC(addRequest, deviceGroupIds)
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

// GetDescriptor retrieves a single application 0,descriptor.
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

// UpdateAppDescriptor allows the user to update the information of a registered descriptor.
func (m *Manager) UpdateAppDescriptor(request *grpc_application_go.UpdateAppDescriptorRequest) (*entities.AppDescriptor, derrors.Error) {
	exists, err := m.OrgProvider.Exists(request.OrganizationId)
	if err != nil {
		return nil, err
	}
	if ! exists{
		return nil, derrors.NewNotFoundError("organizationID").WithParams(request.OrganizationId)
	}
	old, err := m.AppProvider.GetDescriptor(request.AppDescriptorId)
	if err != nil{
		return nil, err
	}
	old.ApplyUpdate(*request)
	err = m.AppProvider.UpdateDescriptor(*old)
	if err != nil{
		return nil, err
	}
	return old, nil
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
// TODO: wait for the conductor to be implemented
func (m * Manager) UpdateService(updateRequest * grpc_application_go.UpdateServiceStatusRequest) error {

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

    aux := toUpdate

    // find the service instance
    for indexGroup, g := range toUpdate.Groups {
    	// find the group
        if g.ServiceGroupInstanceId == updateRequest.ServiceGroupInstanceId {
        	// find the service
			changed := false
        	for indexService, serviceInstance := range g.ServiceInstances {
        		if serviceInstance.ServiceInstanceId == updateRequest.ServiceInstanceId {
        			// found and updated
        			// build the endpoint instances
        			endpoints := make([]entities.EndpointInstance,len(updateRequest.Endpoints))
        			for i, ep := range updateRequest.Endpoints {
        				endpoints[i] = entities.EndpointInstanceFromGRPC(ep)
					}
        			aux.Groups[indexGroup].ServiceInstances[indexService].Status = entities.ServiceStatusFromGRPC[updateRequest.Status]
					aux.Groups[indexGroup].ServiceInstances[indexService].Endpoints = endpoints
					aux.Groups[indexGroup].ServiceInstances[indexService].DeployedOnClusterId = updateRequest.DeployedOnClusterId
					changed = true
				}
			}
        	if !changed {
				return derrors.NewInternalError("update service failed. Not all the entries were found.")
			}
		}
    }


	err = m.AppProvider.UpdateInstance(*aux)
	if err != nil {
		return derrors.NewInternalError("impossible to update instance").CausedBy(err)
	}

	return nil

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

func (m * Manager) AddServiceGroupInstance(request *grpc_application_go.AddServiceGroupInstanceRequest) (*entities.ServiceGroupInstance, derrors.Error){

	// check if the app instance exists (for this organization)
	exists, err := m.OrgProvider.InstanceExists(request.OrganizationId, request.AppInstanceId)
	if err != nil {
		return nil, err
	}
	if ! exists {
		return nil, derrors.NewNotFoundError("appInstanceId").WithParams(request.OrganizationId, request.AppInstanceId)
	}

	// Check if the app descriptor exists (for this organization)
	exists, err = m.OrgProvider.DescriptorExists(request.OrganizationId, request.AppDescriptorId)
	if err != nil {
		return nil, err
	}
	if ! exists {
		return nil, derrors.NewNotFoundError("appDescriptorId").WithParams(request.OrganizationId, request.AppDescriptorId)
	}

	// get the app_descriptor
	appDesc, err := m.AppProvider.GetDescriptor(request.AppDescriptorId)
	if err != nil {
		return nil, err
	}

	// get the service_group in the app descriptor
	var serviceGroup *entities.ServiceGroup
	for _, sg := range appDesc.Groups {
		if sg.ServiceGroupId == request.ServiceGroupId {
			serviceGroup = &sg
			break
		}
	}
	if serviceGroup == nil {
		return nil, derrors.NewNotFoundError("ServiceGroupId").WithParams(request.ServiceGroupId)
	}

	// serviceGroupInstance
	sgInst := serviceGroup.ToEmptyServiceGroupInstance(request.AppInstanceId)

	// get the app instance
	retrieved, err := m.AppProvider.GetInstance(request.AppInstanceId)
	if err != nil {
		return nil, err
	}

	// add the new service group into the instance groups
	retrieved.Groups = append (retrieved.Groups, *sgInst)

	// update
	err = m.AppProvider.UpdateInstance(*retrieved)
	if err != nil {
		return nil, err
	}

	return sgInst, nil
}

func (m * Manager) AddServiceInstance(request *grpc_application_go.AddServiceInstanceRequest) (*entities.ServiceInstance, derrors.Error) {
	// Check if the app descriptor exists (for this organization)
	exists, err := m.OrgProvider.DescriptorExists(request.OrganizationId, request.AppDescriptorId)
	if err != nil {
		return nil, err
	}
	if ! exists {
		return nil, derrors.NewNotFoundError("appDescriptorId").WithParams(request.OrganizationId, request.AppDescriptorId)
	}

	// get the app_descriptor
	appDesc, err := m.AppProvider.GetDescriptor(request.AppDescriptorId)
	if err != nil {
		return nil, err
	}

	// get the service_group in the app descriptor
	var serviceGroup *entities.ServiceGroup
	for _, sg := range appDesc.Groups {
		if sg.ServiceGroupId == request.ServiceGroupId {
			serviceGroup = &sg
			break
		}
	}
	if serviceGroup == nil {
		return nil, derrors.NewNotFoundError("ServiceGroupId").WithParams(request.ServiceGroupId)
	}

	// get the service in the service_group
	var service *entities.Service
	for _, serv := range serviceGroup.Services {
		if serv.ServiceId == request.ServiceId {
			service = &serv
			break
		}
	}
	if service == nil {
		return nil, derrors.NewNotFoundError("serviceID").WithParams(request.ServiceId)
	}

	// Instance creation
	serviceInstance := service.ToServiceInstance(request.AppInstanceId, request.ServiceGroupInstanceId)

	// get the instance
	retrieved, err := m.AppProvider.GetInstance(request.AppInstanceId)
	if err != nil {
		return nil, err
	}

	// look for the service_group_instance and add the new service into service group
	found := false // boolean to control if the service group has been found
	for i:= 0; i < len(retrieved.Groups); i++ {
		if retrieved.Groups[i].ServiceGroupId == request.ServiceGroupId &&
			retrieved.Groups[i].ServiceGroupInstanceId == request.ServiceGroupInstanceId{
			retrieved.Groups[i].ServiceInstances = append(retrieved.Groups[i].ServiceInstances, *serviceInstance)
			found = true
			break
		}
	}
	if ! found {
		return nil, derrors.NewNotFoundError("ServiceGroupInstanceId").WithParams(request.ServiceGroupInstanceId)
	}

	// update the instance
	err = m.AppProvider.UpdateInstance(*retrieved)
	if err != nil {
		return nil, err
	}

	return serviceInstance, nil
}
