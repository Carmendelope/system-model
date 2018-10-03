/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package application

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/application"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
)

type Manager struct {
	OrgProvider organization.Provider
	AppProvider application.Provider
}

func NewManager(orgProvider organization.Provider, appProvider application.Provider) Manager {
	return Manager{orgProvider, appProvider}
}

func (m * Manager) AddAppDescriptor(addRequest * grpc_application_go.AddAppDescriptorRequest) (* entities.AppDescriptor, derrors.Error) {
	exists := m.OrgProvider.Exists(addRequest.OrganizationId)
	if !exists{
		return nil, derrors.NewNotFoundError("organizationID").WithParams(addRequest.OrganizationId)
	}
	descriptor := entities.NewAppDescriptorFromGRPC(addRequest)
	err := m.AppProvider.AddDescriptor(*descriptor)
	if err != nil {
		return nil, err
	}
	err = m.OrgProvider.AddDescriptor(descriptor.OrganizationId, descriptor.AppDescriptorId)
	if err != nil {
	    return nil, err
	}

	return descriptor, nil
}

func (m * Manager) ListDescriptors(orgID * grpc_organization_go.OrganizationId) ([] entities.AppDescriptor, derrors.Error) {
	if !m.OrgProvider.Exists(orgID.OrganizationId){
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

func (m * Manager) GetDescriptor(appDescID * grpc_application_go.AppDescriptorId) (* entities.AppDescriptor, derrors.Error){
	if ! m.OrgProvider.Exists(appDescID.OrganizationId){
		return nil, derrors.NewNotFoundError("organizationID").WithParams(appDescID.OrganizationId)
	}

	if !m.OrgProvider.DescriptorExists(appDescID.OrganizationId, appDescID.AppDescriptorId){
		return nil, derrors.NewNotFoundError("appDescriptorID").WithParams(appDescID.OrganizationId, appDescID.AppDescriptorId)
	}
	return m.AppProvider.GetDescriptor(appDescID.AppDescriptorId)
}

func (m * Manager) AddAppInstance(addRequest * grpc_application_go.AddAppInstanceRequest) (* entities.AppInstance, derrors.Error) {

	if !m.OrgProvider.Exists(addRequest.OrganizationId){
		return nil, derrors.NewNotFoundError("organizationID").WithParams(addRequest.OrganizationId)
	}
	if !m.AppProvider.DescriptorExists(addRequest.AppDescriptorId){
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

func (m * Manager) ListInstances(orgID * grpc_organization_go.OrganizationId) ([] entities.AppInstance, derrors.Error) {
	if !m.OrgProvider.Exists(orgID.OrganizationId){
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

func (m * Manager) GetInstance(appInstID * grpc_application_go.AppInstanceId) (* entities.AppInstance, derrors.Error){
	if ! m.OrgProvider.Exists(appInstID.OrganizationId){
		return nil, derrors.NewNotFoundError("organizationID").WithParams(appInstID.OrganizationId)
	}

	if !m.OrgProvider.InstanceExists(appInstID.OrganizationId, appInstID.AppInstanceId){
		return nil, derrors.NewNotFoundError("appInstanceID").WithParams(appInstID.OrganizationId, appInstID.AppInstanceId)
	}
	return m.AppProvider.GetInstance(appInstID.AppInstanceId)
}
