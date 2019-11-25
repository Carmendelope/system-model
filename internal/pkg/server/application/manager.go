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

package application

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/application"
	"github.com/nalej/system-model/internal/pkg/provider/device"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/rs/zerolog/log"
	"strings"
)

// Manager structure with the required providers for application operations.
type Manager struct {
	OrgProvider      organization.Provider
	AppProvider      application.Provider
	DevProvider      device.Provider
	PublicHostDomain string
}

// NewManager creates a Manager using a set of providers.
func NewManager(orgProvider organization.Provider, appProvider application.Provider, devProvider device.Provider, publicHostDomain string) Manager {
	return Manager{orgProvider, appProvider, devProvider, publicHostDomain}
}

func (m *Manager) extractGroupIds(organizationID string, rules []*grpc_application_go.SecurityRule) (map[string]string, derrors.Error) {
	// -----------------
	// check if the descriptor has device_names in the rules
	// we need to convert deviceGroupNames into deviceGroupIds
	names := make(map[string]bool, 0) // uses a map to avoid insert a device group twice
	for _, rules := range rules {
		if len(rules.DeviceGroupNames) > 0 {
			for _, name := range rules.DeviceGroupNames {
				names[name] = true
			}
		}
	}
	// map to array
	keys := make([]string, len(names))
	i := 0
	for key, _ := range names {
		keys[i] = key
		i += 1
	}

	deviceGroupIds := make(map[string]string, 0) // map of deviceGroupIds indexed by deviceGroupNames
	if len(keys) > 0 {
		deviceGroups, err := m.DevProvider.GetDeviceGroupsByName(organizationID, keys)
		if err != nil {
			return nil, err
		}

		for _, deviceGroup := range deviceGroups {
			deviceGroupIds[deviceGroup.Name] = deviceGroup.DeviceGroupId
		}

		// check the devices number returned (it should be the the same as deviceNames)
		if len(deviceGroupIds) != len(keys) {
			return nil, derrors.NewNotFoundError("device group names").WithParams(keys)
		}

	}
	// ---------------------
	return deviceGroupIds, nil

}

// AddAppDescriptor adds a new application descriptor to a given organization.
func (m *Manager) AddAppDescriptor(addRequest *grpc_application_go.AddAppDescriptorRequest) (*entities.AppDescriptor, derrors.Error) {
	exists, err := m.OrgProvider.Exists(addRequest.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(addRequest.OrganizationId)
	}

	descriptor, err := entities.NewAppDescriptorFromGRPC(addRequest)
	if err != nil {
		return nil, err
	}

	// Validate AppDescriptor
	err = entities.ValidateDescriptor(*descriptor)
	if err != nil {
		return nil, err
	}

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
func (m *Manager) ListDescriptors(orgID *grpc_organization_go.OrganizationId) ([]entities.AppDescriptor, derrors.Error) {
	exists, err := m.OrgProvider.Exists(orgID.OrganizationId)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(orgID.OrganizationId)
	}
	descriptors, err := m.OrgProvider.ListDescriptors(orgID.OrganizationId)
	if err != nil {
		return nil, err
	}
	result := make([]entities.AppDescriptor, 0)
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
func (m *Manager) GetDescriptor(appDescID *grpc_application_go.AppDescriptorId) (*entities.AppDescriptor, derrors.Error) {
	exists, err := m.OrgProvider.Exists(appDescID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(appDescID.OrganizationId)
	}
	exists, err = m.OrgProvider.DescriptorExists(appDescID.OrganizationId, appDescID.AppDescriptorId)
	if err != nil {
		return nil, err
	}
	if !exists {
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
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(request.OrganizationId)
	}
	old, err := m.AppProvider.GetDescriptor(request.AppDescriptorId)
	if err != nil {
		return nil, err
	}
	old.ApplyUpdate(*request)
	err = m.AppProvider.UpdateDescriptor(*old)
	if err != nil {
		return nil, err
	}
	return old, nil
}

// RemoveAppDescriptor removes an application descriptor.
func (m *Manager) RemoveAppDescriptor(appDescID *grpc_application_go.AppDescriptorId) derrors.Error {
	exists, err := m.OrgProvider.Exists(appDescID.OrganizationId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organizationID").WithParams(appDescID.OrganizationId)
	}
	exists, err = m.OrgProvider.DescriptorExists(appDescID.OrganizationId, appDescID.AppDescriptorId)
	if err != nil {
		return err
	}
	if !exists {
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

func (m *Manager) GetDescriptorAppParameters(request *grpc_application_go.AppDescriptorId) ([]entities.Parameter, derrors.Error) {
	exists, err := m.OrgProvider.Exists(request.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(request.OrganizationId)
	}
	exists, err = m.OrgProvider.DescriptorExists(request.OrganizationId, request.AppDescriptorId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("appDescriptorID").WithParams(request.OrganizationId, request.AppDescriptorId)
	}
	return m.AppProvider.GetDescriptorParameters(request.AppDescriptorId)
}

// AddAppInstance adds a new application instance to a given organization.
func (m *Manager) AddAppInstance(addRequest *grpc_application_go.AddAppInstanceRequest) (*entities.AppInstance, derrors.Error) {

	exists, err := m.OrgProvider.Exists(addRequest.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(addRequest.OrganizationId)
	}
	exists, err = m.AppProvider.DescriptorExists(addRequest.AppDescriptorId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("descriptorID").WithParams(addRequest.OrganizationId, addRequest.AppDescriptorId)
	}

	descriptor, err := m.AppProvider.GetDescriptor(addRequest.AppDescriptorId)
	if err != nil {
		return nil, err
	}

	instance := entities.NewAppInstanceFromAddInstanceRequestGRPC(addRequest, descriptor)
	err = m.AppProvider.AddInstance(*instance)
	if err != nil {
		return nil, err
	}
	err = m.OrgProvider.AddInstance(instance.OrganizationId, instance.AppInstanceId)
	if err != nil {
		return nil, err
	}

	// add parameters
	if addRequest.Parameters != nil {
		parameters := make([]entities.InstanceParameter, 0)
		for _, param := range addRequest.Parameters.Parameters {
			parameters = append(parameters, *entities.NewInstanceParamFromGRPC(param))
		}
		err = m.AppProvider.AddInstanceParameters(instance.AppInstanceId, parameters)
		if err != nil {
			log.Error().Str("instance_id", instance.AppInstanceId).Str("trace", err.DebugReport()).Msg("error saving instance parameters.")
			// if error storing instance parameters -> delete instance and return the error
			rollBackErr := m.AppProvider.DeleteInstance(instance.AppInstanceId)
			if rollBackErr != nil {
				log.Error().Str("instance_id", instance.AppInstanceId).Str("trace", rollBackErr.DebugReport()).Msg("Error removing instance")
			}
			return nil, err
		}
	}

	return instance, nil
}

// ListInstances retrieves the list of instances associated with an organization.
func (m *Manager) ListInstances(orgID *grpc_organization_go.OrganizationId) ([]entities.AppInstance, derrors.Error) {
	exists, err := m.OrgProvider.Exists(orgID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(orgID.OrganizationId)
	}
	instances, err := m.OrgProvider.ListInstances(orgID.OrganizationId)
	if err != nil {
		return nil, err
	}

	result := make([]entities.AppInstance, 0)
	for _, instID := range instances {
		toAdd, err := m.AppProvider.GetInstance(instID)
		if err != nil {
			// NP-1593.
			// It can happen, that while an  instance is being undeploying, a list of the instances is requested (and the join fails)
			log.Warn().Str("instance", instID).Msg("not found!!")
		} else {
			// Fill Global FQdn
			err = m.fillGlobalFqdn(toAdd)
			if err != nil {
				return nil, err
			}
			result = append(result, *toAdd)
		}
	}

	return result, nil
}

func (m *Manager) fillGlobalFqdn(instance *entities.AppInstance) derrors.Error {
	// Load ServiceGroup GlobalFqn
	for i := 0; i < len(instance.Groups); i++ {
		globalFQDN, err := m.AppProvider.GetAppEndpointList(instance.Groups[i].OrganizationId, instance.Groups[i].AppInstanceId, instance.Groups[i].ServiceGroupInstanceId)
		if err != nil {
			return err
		}

		// map to avoid repeated fqdn
		if globalFQDN != nil && len(globalFQDN) > 0 {
			fqdns := make(map[string]bool, 0)
			for _, fqdn := range globalFQDN {
				fqdns[fqdn.GlobalFqdn] = true
			}
			// map to list
			instance.Groups[i].GlobalFqdn = make([]string, 0)
			for key, _ := range fqdns {
				instance.Groups[i].GlobalFqdn = append(instance.Groups[i].GlobalFqdn, fmt.Sprintf("%s.ep.%s", key, m.PublicHostDomain))
			}
		}

	}
	return nil
}

// GetInstance retrieves a single instance.
func (m *Manager) GetInstance(appInstID *grpc_application_go.AppInstanceId) (*entities.AppInstance, derrors.Error) {
	exists, err := m.OrgProvider.Exists(appInstID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(appInstID.OrganizationId)
	}

	exists, err = m.OrgProvider.InstanceExists(appInstID.OrganizationId, appInstID.AppInstanceId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("appInstanceID").WithParams(appInstID.OrganizationId, appInstID.AppInstanceId)
	}
	instance, err := m.AppProvider.GetInstance(appInstID.AppInstanceId)
	if err != nil {
		return nil, err
	}

	err = m.fillGlobalFqdn(instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

// UpdateInstance updates the information of a given instance.
func (m *Manager) UpdateInstance(updateRequest *grpc_application_go.UpdateAppStatusRequest) derrors.Error {
	exists, err := m.OrgProvider.InstanceExists(updateRequest.OrganizationId, updateRequest.AppInstanceId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("appInstanceID").WithParams(updateRequest.OrganizationId, updateRequest.AppInstanceId)
	}

	toUpdate, err := m.AppProvider.GetInstance(updateRequest.AppInstanceId)
	if err != nil {
		return derrors.NewInternalError("impossible to get old instance", err)
	}

	toUpdate.Status = entities.AppStatusFromGRPC[updateRequest.Status]
	if updateRequest.Info != "" {
		toUpdate.Info = updateRequest.Info
	}

	err = m.AppProvider.UpdateInstance(*toUpdate)
	if err != nil {
		return derrors.NewInternalError("impossible to update instance").CausedBy(err)
	}

	return nil
}

// UpdateService updates an application service.
// TODO: wait for the conductor to be implemented
func (m *Manager) UpdateService(updateRequest *grpc_application_go.UpdateServiceStatusRequest) derrors.Error {

	exists, err := m.OrgProvider.InstanceExists(updateRequest.OrganizationId, updateRequest.AppInstanceId)

	if err != nil {
		return err
	}
	if !exists {
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
					endpoints := make([]entities.EndpointInstance, len(updateRequest.Endpoints))
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

func (m *Manager) UpdateAppInstance(appInstance *grpc_application_go.AppInstance) derrors.Error {
	localEntity := entities.NewAppInstanceFromGRPC(appInstance)

	err := m.AppProvider.UpdateInstance(*localEntity)
	if err != nil {
		return derrors.NewInternalError("impossible to update application instance").CausedBy(err)
	}
	return nil
}

// RemoveAppInstance removes an application instance
func (m *Manager) RemoveAppInstance(appInstID *grpc_application_go.AppInstanceId) derrors.Error {
	exists, err := m.OrgProvider.Exists(appInstID.OrganizationId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organizationID").WithParams(appInstID.OrganizationId)
	}
	exists, err = m.OrgProvider.InstanceExists(appInstID.OrganizationId, appInstID.AppInstanceId)
	if err != nil {
		return err
	}
	if !exists {
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
	} else { // delete parameters (if exist)
		instErr := m.AppProvider.DeleteInstanceParameters(appInstID.AppInstanceId)
		if instErr != nil {
			log.Error().Str("instanceID", appInstID.AppInstanceId).Str("trace", instErr.DebugReport()).Msg("Error removing parameters")
		}
	}
	return err
}

func (m *Manager) GetInstanceParameters(request *grpc_application_go.AppInstanceId) ([]entities.InstanceParameter, derrors.Error) {
	exists, err := m.OrgProvider.Exists(request.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(request.OrganizationId)
	}
	exists, err = m.OrgProvider.InstanceExists(request.OrganizationId, request.AppInstanceId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("AppInstanceId").WithParams(request.AppInstanceId)
	}
	parameters, err := m.AppProvider.GetInstanceParameters(request.AppInstanceId)
	if err != nil {
		return nil, err
	}
	return parameters, nil
}

func (m *Manager) AddServiceGroupInstances(request *grpc_application_go.AddServiceGroupInstancesRequest) ([]entities.ServiceGroupInstance, derrors.Error) {

	// check if the app instance exists (for this organization)
	exists, err := m.OrgProvider.InstanceExists(request.OrganizationId, request.AppInstanceId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("appInstanceId").WithParams(request.OrganizationId, request.AppInstanceId)
	}

	// Check if the app descriptor exists (for this organization)
	exists, err = m.OrgProvider.DescriptorExists(request.OrganizationId, request.AppDescriptorId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("appDescriptorId").WithParams(request.OrganizationId, request.AppDescriptorId)
	}

	// get the app_descriptor
	appDesc, err := m.AppProvider.GetParametrizedDescriptor(request.AppInstanceId)
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

	// Generate as many service group instances as required
	result := make([]entities.ServiceGroupInstance, request.NumInstances)
	for numReplica := int32(0); numReplica < request.NumInstances; numReplica++ {
		// create the group
		sgi := serviceGroup.ToServiceGroupInstance(request.AppInstanceId)
		// fill the metadata
		sgi.FillMetadata(int(request.NumInstances))
		result[numReplica] = *sgi
	}

	// get the app instance
	retrieved, err := m.AppProvider.GetInstance(request.AppInstanceId)
	if err != nil {
		return nil, err
	}

	// set the new values for these service group instances
	retrieved.Groups = append(retrieved.Groups, result...)

	// update
	err = m.AppProvider.UpdateInstance(*retrieved)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *Manager) RemoveServiceGroupInstances(removeRequest *grpc_application_go.RemoveServiceGroupInstancesRequest) derrors.Error {
	// Get the corresponding instance
	appInst, err := m.AppProvider.GetInstance(removeRequest.AppInstanceId)
	if err != nil {
		return err
	}

	appInst.Groups = nil

	// update
	err = m.AppProvider.UpdateInstance(*appInst)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) GetServiceGroupInstanceMetadata(request *grpc_application_go.GetServiceGroupInstanceMetadataRequest) (*entities.InstanceMetadata, derrors.Error) {
	// Get the corresponding instance
	appInst, err := m.AppProvider.GetInstance(request.AppInstanceId)
	if err != nil {
		return nil, err
	}

	// Find the service group instance
	for _, groupInst := range appInst.Groups {
		if groupInst.ServiceGroupInstanceId == request.ServiceGroupInstanceId {
			return groupInst.Metadata, nil
		}
	}

	// Not found
	return nil, derrors.NewNotFoundError(fmt.Sprintf("service group instance %s not found", request.ServiceGroupInstanceId))
}

func (m *Manager) UpdateServiceGroupInstanceMetadata(request *grpc_application_go.InstanceMetadata) derrors.Error {
	// Get the corresponding instance
	appInst, err := m.AppProvider.GetInstance(request.AppInstanceId)
	if err != nil {
		return err
	}

	// Find the service group instance and update it
	targetGroupIndex := 0
	found := false
	for i, groupInst := range appInst.Groups {
		if groupInst.ServiceGroupInstanceId == request.MonitoredInstanceId {
			targetGroupIndex = i
			found = true
			break
		}
	}

	if !found {
		return derrors.NewNotFoundError(fmt.Sprintf("service group instance %s not found", request.MonitoredInstanceId))
	}

	//update the corresponding application instance
	appInst.Groups[targetGroupIndex].Metadata = entities.NewMetadataFromGRPC(request)
	err = m.AppProvider.UpdateInstance(*appInst)
	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) AddServiceInstance(request *grpc_application_go.AddServiceInstanceRequest) (*entities.ServiceInstance, derrors.Error) {
	// Check if the app descriptor exists (for this organization)
	exists, err := m.OrgProvider.DescriptorExists(request.OrganizationId, request.AppDescriptorId)
	if err != nil {
		return nil, err
	}
	if !exists {
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
	for i := 0; i < len(retrieved.Groups); i++ {
		if retrieved.Groups[i].ServiceGroupId == request.ServiceGroupId &&
			retrieved.Groups[i].ServiceGroupInstanceId == request.ServiceGroupInstanceId {
			retrieved.Groups[i].ServiceInstances = append(retrieved.Groups[i].ServiceInstances, *serviceInstance)
			found = true
			break
		}
	}
	if !found {
		return nil, derrors.NewNotFoundError("ServiceGroupInstanceId").WithParams(request.ServiceGroupInstanceId)
	}

	// update the instance
	err = m.AppProvider.UpdateInstance(*retrieved)
	if err != nil {
		return nil, err
	}

	return serviceInstance, nil
}

// AddAppEndPoint adds a new App Endpoint to a given service instance
func (m *Manager) AddAppEndpoint(appEndpoint *grpc_application_go.AddAppEndpointRequest) derrors.Error {

	endpoint, err := entities.NewAppEndpointFromGRPC(appEndpoint)
	if err != nil {
		return err
	}

	err = m.AppProvider.AddAppEndpoint(*endpoint)
	if err != nil {
		return err
	}

	return nil
}

// GetAppEndPoint retrieves an appEndpoint
func (m *Manager) GetAppEndpoint(request *grpc_application_go.GetAppEndPointRequest) (*grpc_application_go.AppEndpointList, derrors.Error) {

	split := strings.Split(request.Fqdn, ".")
	globalFqdn := fmt.Sprintf("%s.%s.%s.%s", split[0], split[1], split[2], split[3])

	list, err := m.AppProvider.GetAppEndpointByFQDN(globalFqdn)
	if err != nil {
		return nil, err
	}

	endpointList := make([]*grpc_application_go.AppEndpoint, 0)

	// check that there are only end points of an organization
	if len(list) > 0 {
		organizationID := list[0].OrganizationId
		for _, endpoint := range list {
			if endpoint.OrganizationId != organizationID {
				return nil, derrors.NewInternalError("Unable to return app end points, several organizations have the same endpoint")
			}
			endpointList = append(endpointList, endpoint.ToGRPC())
		}
	}

	return &grpc_application_go.AppEndpointList{
		AppEndpoints: endpointList,
	}, nil
}

func (m *Manager) RemoveAppEndpoints(removeRequest *grpc_application_go.RemoveAppEndpointRequest) derrors.Error {
	return m.AppProvider.DeleteAppEndpoints(removeRequest.OrganizationId, removeRequest.AppInstanceId)
}

func (m *Manager) AddZtNetwork(request *grpc_application_go.AddAppZtNetworkRequest) derrors.Error {
	return m.AppProvider.AddAppZtNetwork(entities.AppZtNetwork{OrganizationId: request.OrganizationId,
		AppInstanceId: request.AppInstanceId, ZtNetworkId: request.NetworkId, VSAList: request.VsaList})
}

func (m *Manager) RemoveZtNetwork(request *grpc_application_go.RemoveAppZtNetworkRequest) derrors.Error {
	return m.AppProvider.RemoveAppZtNetwork(request.OrganizationId, request.AppInstanceId)
}

func (m *Manager) AddZtNetworkProxy(request *entities.ServiceProxy) derrors.Error {
	return m.AppProvider.AddZtNetworkProxy(*request)
}

func (m *Manager) RemoveZtNetworkProxy(organizationId string, appInstanceId string, fqdn string, clusterId string,
	serviceGroupInstanceId string, serviceInstanceId string) derrors.Error {
	return m.AppProvider.RemoveZtNetworkProxy(organizationId, appInstanceId, fqdn, clusterId, serviceGroupInstanceId, serviceInstanceId)
}

func (m *Manager) GetAppZtNetwork(request *grpc_application_go.GetAppZtNetworkRequest) (*entities.AppZtNetwork, derrors.Error) {
	return m.AppProvider.GetAppZtNetwork(request.OrganizationId, request.AppInstanceId)
}

func (m *Manager) AddAppZtNetworkMember(request *grpc_application_go.AddAuthorizedZtNetworkMemberRequest) (*entities.AppZtNetworkMembers, derrors.Error) {
	return m.AppProvider.AddAppZtNetworkMember(*entities.NewAppZtNetworkMemberFromGRPC(request))
}

func (m *Manager) RemoveAppZtNetworkMember(organizationId string, appInstanceId string, serviceGroupInstanceId string, serviceApplicationInstanceId string, ztNetworkId string) derrors.Error {
	return m.AppProvider.RemoveAppZtNetworkMember(organizationId, appInstanceId, serviceGroupInstanceId, serviceApplicationInstanceId, ztNetworkId)
}

func (m *Manager) GetAppZtNetworkMember(organizationId string, appInstanceId string, serviceGroupInstanceId string, serviceApplicationInstanceId string) (*entities.AppZtNetworkMembers, derrors.Error) {
	return m.AppProvider.GetAppZtNetworkMember(organizationId, appInstanceId, serviceGroupInstanceId, serviceApplicationInstanceId)
}

func (m *Manager) ListAuthorizedZTNetworkMembers(organizationId string, appInstanceId string, ztNetworkId string) (*grpc_application_go.ZtNetworkMembers, derrors.Error) {
	retrieved, err := m.AppProvider.ListAppZtNetworkMembers(organizationId, appInstanceId, ztNetworkId)

	if err != nil {
		return nil, err
	}

	list := make([]*grpc_application_go.ZtNetworkMember, 0)

	for _, ret := range retrieved {
		for _, member := range ret.Members {
			list = append(list, &grpc_application_go.ZtNetworkMember{
				OrganizationId:               ret.OrganizationId,
				NetworkId:                    ret.ZtNetworkId,
				MemberId:                     member.MemberId,
				AppInstanceId:                ret.AppInstanceId,
				ServiceGroupInstanceId:       ret.ServiceGroupInstanceId,
				ServiceApplicationInstanceId: ret.ServiceApplicationInstanceId,
				IsProxy:                      member.IsProxy,
				CreatedAt:                    member.CreatedAt,
			})
		}
	}

	return &grpc_application_go.ZtNetworkMembers{
		Members: list,
	}, nil
}

func (m *Manager) fillDeviceGroupIds(desc *entities.ParametrizedDescriptor) derrors.Error {
	// -----------------
	// check if the descriptor has device_names in the rules
	// we need to convert deviceGroupNames into deviceGroupIds
	names := make(map[string]bool, 0) // uses a map to avoid insert a device group twice
	for _, rules := range desc.Rules {
		if len(rules.DeviceGroupNames) > 0 {
			for _, name := range rules.DeviceGroupNames {
				names[name] = true
			}
		}
	}
	// map to array
	keys := make([]string, len(names))
	i := 0
	for key := range names {
		keys[i] = key
		i += 1
	}

	deviceGroupIds := make(map[string]string, 0) // map of deviceGroupIds indexed by deviceGroupNames
	if len(keys) > 0 {
		deviceGroups, err := m.DevProvider.GetDeviceGroupsByName(desc.OrganizationId, keys)
		if err != nil {
			return err
		}

		for _, deviceGroup := range deviceGroups {
			deviceGroupIds[deviceGroup.Name] = deviceGroup.DeviceGroupId
		}

		// check the devices number returned (it should be the the same as deviceNames)
		if len(deviceGroupIds) != len(keys) {
			return derrors.NewNotFoundError("device group names").WithParams(keys)
		}

	}

	// once we have all the ids of the devices groups, we add them to the descriptor
	for i := 0; i < len(desc.Rules); i++ {
		ids := make([]string, 0)

		for j := 0; j < len(desc.Rules[i].DeviceGroupNames); j++ {

			id, exists := deviceGroupIds[desc.Rules[i].DeviceGroupNames[j]]
			if !exists {
				log.Error().Str("deviceName", desc.Rules[i].DeviceGroupNames[j]).Msg("Device id not found")
				return derrors.NewNotFoundError("device group id").WithParams(desc.Rules[i].DeviceGroupNames[j])
			}
			ids = append(ids, id)
		}

		desc.Rules[i].DeviceGroupIds = ids
	}

	return nil

}

func (m *Manager) getDeviceGroupIds(organizationID string, deviceNames []string) (map[string]string, derrors.Error) {

	deviceGroupIds := make(map[string]string, 0) // map of deviceGroupIds indexed by deviceGroupNames
	if len(deviceNames) > 0 {
		deviceGroups, err := m.DevProvider.GetDeviceGroupsByName(organizationID, deviceNames)
		if err != nil {
			return nil, err
		}

		for _, deviceGroup := range deviceGroups {
			deviceGroupIds[deviceGroup.Name] = deviceGroup.DeviceGroupId
		}

		// check the devices number returned (it should be the the same as deviceNames)
		if len(deviceGroupIds) != len(deviceNames) {
			return nil, derrors.NewNotFoundError("device group names").WithParams(deviceNames)
		}

	}
	// ---------------------
	return deviceGroupIds, nil
}

// AddParametrizedDescriptor adds a parametrized descriptor to a given descriptor
func (m *Manager) AddParametrizedDescriptor(descriptor *grpc_application_go.ParametrizedDescriptor) (*entities.ParametrizedDescriptor, derrors.Error) {

	// check if the organization exists
	exists, err := m.OrgProvider.Exists(descriptor.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(descriptor.OrganizationId)
	}
	// check if the descriptor exists
	exists, err = m.AppProvider.DescriptorExists(descriptor.AppDescriptorId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("descriptorID").WithParams(descriptor.OrganizationId, descriptor.AppDescriptorId)
	}

	// check if the instance exists
	exists, err = m.AppProvider.InstanceExists(descriptor.AppInstanceId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("instanceID").WithParams(descriptor.OrganizationId, descriptor.AppInstanceId)
	}

	// Convert to ParametrizedDescriptor
	newDesc := entities.NewParametrizedDescriptorFromGRPC(descriptor)

	// fill deviceGroupIds
	err = m.fillDeviceGroupIds(newDesc)
	if err != nil {
		return nil, err
	}
	err = m.AppProvider.AddParametrizedDescriptor(*newDesc)
	if err != nil {
		return nil, err
	}

	return newDesc, nil
}

// GetParametrizedDescriptor retrieves the parametrized descriptor associated with an instance
func (m *Manager) GetParametrizedDescriptor(request *grpc_application_go.AppInstanceId) (*entities.ParametrizedDescriptor, derrors.Error) {
	// check if the organization exists
	exists, err := m.OrgProvider.Exists(request.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(request.OrganizationId)
	}

	// check if the instance exists
	exists, err = m.AppProvider.InstanceExists(request.AppInstanceId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("instanceID").WithParams(request.OrganizationId, request.AppInstanceId)
	}

	descriptor, err := m.AppProvider.GetParametrizedDescriptor(request.AppInstanceId)
	if err != nil {
		return nil, err
	}

	return descriptor, nil
}

// RemoveParametrizedDescriptor removes the parametrized descriptor associated with an instance
func (m *Manager) RemoveParametrizedDescriptor(request *grpc_application_go.AppInstanceId) derrors.Error {
	// check if the organization exists
	exists, err := m.OrgProvider.Exists(request.OrganizationId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organizationID").WithParams(request.OrganizationId)
	}

	err = m.AppProvider.DeleteParametrizedDescriptor(request.AppInstanceId)
	if err != nil {
		return err
	}
	return nil

}
