/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package application

import (
	"context"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
)

// Handler structure for the application requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

// AddAppDescriptor adds a new application descriptor to a given organization.
func (h *Handler) AddAppDescriptor(ctx context.Context, addRequest *grpc_application_go.AddAppDescriptorRequest) (*grpc_application_go.AppDescriptor, error) {
	err := entities.ValidAddAppDescriptorRequest(addRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid add application descriptor request")
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Interface("addRequest", addRequest).Msg("Adding application descriptor")
	added, err := h.Manager.AddAppDescriptor(addRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add application descriptor")
		return nil, conversions.ToGRPCError(err)
	}
	return added.ToGRPC(), nil
}

// ListAppDescriptors retrieves a list of application descriptors.
func (h *Handler) ListAppDescriptors(ctx context.Context, orgID *grpc_organization_go.OrganizationId) (*grpc_application_go.AppDescriptorList, error) {
	descriptors, err := h.Manager.ListDescriptors(orgID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot list descriptors")
		return nil, conversions.ToGRPCError(err)
	}
	toReturn := make([]*grpc_application_go.AppDescriptor, 0)
	for _, d := range descriptors {
		toReturn = append(toReturn, d.ToGRPC())
	}
	result := &grpc_application_go.AppDescriptorList{
		Descriptors: toReturn,
	}
	return result, nil
}

// GetAppDescriptor retrieves a given application descriptor.
func (h *Handler) GetAppDescriptor(ctx context.Context, appDescID *grpc_application_go.AppDescriptorId) (*grpc_application_go.AppDescriptor, error) {
	err := entities.ValidAppDescriptorId(appDescID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid application descriptor identifier")
		return nil, conversions.ToGRPCError(err)
	}
	descriptor, err := h.Manager.GetDescriptor(appDescID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get application descriptor")
		return nil, conversions.ToGRPCError(err)
	}
	return descriptor.ToGRPC(), nil
}

// UpdateAppDescriptor allows the user to update the information of a registered descriptor.
func (h *Handler) UpdateAppDescriptor(ctx context.Context, request *grpc_application_go.UpdateAppDescriptorRequest) (*grpc_application_go.AppDescriptor, error) {
	err := entities.ValidUpdateAppDescriptorRequest(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid update app descriptor request")
		return nil, conversions.ToGRPCError(err)
	}
	updated, err := h.Manager.UpdateAppDescriptor(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot update application descriptor")
		return nil, conversions.ToGRPCError(err)
	}
	return updated.ToGRPC(), nil
}

// RemoveAppDescriptor removes an application descriptor.
func (h *Handler) RemoveAppDescriptor(ctx context.Context, appDescID *grpc_application_go.AppDescriptorId) (*grpc_common_go.Success, error) {
	err := entities.ValidAppDescriptorId(appDescID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid application descriptor identifier")
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.RemoveAppDescriptor(appDescID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove application descriptor")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// GetDescriptorAppParameters retrieves a list of application parameters of a descriptor
func (h *Handler) GetDescriptorAppParameters(ctx context.Context, request *grpc_application_go.AppDescriptorId) (*grpc_application_go.AppParameterList, error) {
	err := entities.ValidAppDescriptorId(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid application descriptor identifier")
		return nil, conversions.ToGRPCError(err)
	}
	parameters, err := h.Manager.GetDescriptorAppParameters(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get application descriptor")
		return nil, conversions.ToGRPCError(err)
	}
	toReturn := make([]*grpc_application_go.AppParameter, 0)
	for _, param := range parameters {
		toReturn = append(toReturn, param.ToGRPC())
	}
	return &grpc_application_go.AppParameterList{
		Parameters: toReturn,
	}, nil
}

// GetInstanceParameters retrieves a list of application parameters of an instance
func (h *Handler) GetInstanceParameters(ctx context.Context, request *grpc_application_go.AppInstanceId) (*grpc_application_go.InstanceParameterList, error) {
	err := entities.ValidAppInstanceId(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid application instance identifier")
		return nil, conversions.ToGRPCError(err)
	}
	instanceParams, err := h.Manager.GetInstanceParameters(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get instance parameters")
		return nil, conversions.ToGRPCError(err)
	}
	toReturn := make([]*grpc_application_go.InstanceParameter, 0)
	for _, param := range instanceParams {
		toReturn = append(toReturn, param.ToGRPC())
	}
	return &grpc_application_go.InstanceParameterList{
		Parameters: toReturn,
	}, nil
}

// AddAppInstance adds a new application instance to a given organization.
func (h *Handler) AddAppInstance(ctx context.Context, addInstanceRequest *grpc_application_go.AddAppInstanceRequest) (*grpc_application_go.AppInstance, error) {
	err := entities.ValidAddAppInstanceRequest(addInstanceRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid add application instance request")
		return nil, conversions.ToGRPCError(err)
	}
	added, err := h.Manager.AddAppInstance(addInstanceRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add application instance")
		return nil, conversions.ToGRPCError(err)
	}
	return added.ToGRPC(), nil
}

// UpdateAppInstance adds a new application instance to a given organization.
func (h *Handler) UpdateAppInstance(ctx context.Context, appInstance *grpc_application_go.AppInstance) (*grpc_common_go.Success, error) {
	// TODO validate the application instance
	err := h.Manager.UpdateAppInstance(appInstance)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot update application instance")
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}

// ListAppInstances retrieves a list of application instances.
func (h *Handler) ListAppInstances(ctx context.Context, orgID *grpc_organization_go.OrganizationId) (*grpc_application_go.AppInstanceList, error) {
	instances, err := h.Manager.ListInstances(orgID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot list instances")
		return nil, conversions.ToGRPCError(err)
	}
	toReturn := make([]*grpc_application_go.AppInstance, 0)
	for _, inst := range instances {
		toReturn = append(toReturn, inst.ToGRPC())
	}
	result := &grpc_application_go.AppInstanceList{
		Instances: toReturn,
	}
	return result, nil
}

// GetAppInstance retrieves a given application instance.
func (h *Handler) GetAppInstance(ctx context.Context, appInstID *grpc_application_go.AppInstanceId) (*grpc_application_go.AppInstance, error) {
	err := entities.ValidAppInstanceId(appInstID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid application instance identifier")
		return nil, conversions.ToGRPCError(err)
	}
	instance, err := h.Manager.GetInstance(appInstID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get application instance")
		return nil, conversions.ToGRPCError(err)
	}
	return instance.ToGRPC(), nil
}

// UpdateAppStatus updates the status of an application instance.
func (h *Handler) UpdateAppStatus(ctx context.Context, updateAppStatus *grpc_application_go.UpdateAppStatusRequest) (*grpc_common_go.Success, error) {
	err := entities.ValidUpdateAppStatusRequest(updateAppStatus)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid update application status request")
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.UpdateInstance(updateAppStatus)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot update instance status")
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}

// UpdateServiceStatus updates the status of an application instance service.
func (h *Handler) UpdateServiceStatus(ctx context.Context, updateServiceStatus *grpc_application_go.UpdateServiceStatusRequest) (*grpc_common_go.Success, error) {
	err := entities.ValidUpdateServiceStatusRequest(updateServiceStatus)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid update service status request")
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.UpdateService(updateServiceStatus)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot update service")
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}

// RemoveAppInstance removes an application instance
func (h *Handler) RemoveAppInstance(ctx context.Context, appInstID *grpc_application_go.AppInstanceId) (*grpc_common_go.Success, error) {
	err := entities.ValidAppInstanceId(appInstID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid application instance identifier")
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.RemoveAppInstance(appInstID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove application instance")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// AddServiceGroupInstance to an already existing application instance
func (h *Handler) AddServiceGroupInstances(ctx context.Context, addRequest *grpc_application_go.AddServiceGroupInstancesRequest) (*grpc_application_go.ServiceGroupInstancesList, error) {
	err := entities.ValidAddServiceGroupInstanceRequest(addRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid add service group instances request")
		return nil, conversions.ToGRPCError(err)
	}
	instances, err := h.Manager.AddServiceGroupInstances(addRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add service group instances")
		return nil, conversions.ToGRPCError(err)
	}
	result := make([]*grpc_application_go.ServiceGroupInstance, len(instances))
	for i, inst := range instances {
		result[i] = inst.ToGRPC()
	}
	toReturn := grpc_application_go.ServiceGroupInstancesList{
		ServiceGroupInstances: result,
	}
	return &toReturn, nil
}

func (h *Handler) RemoveServiceGroupInstances(ctx context.Context, removeRequest *grpc_application_go.RemoveServiceGroupInstancesRequest) (*grpc_common_go.Success, error) {
	err := entities.ValidateRemoveServiceGroupInstancesRequest(removeRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid remove service group instances request")
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.RemoveServiceGroupInstances(removeRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove service group instances")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// AddServiceInstance to an already existing service group instance
func (h *Handler) AddServiceInstance(ctx context.Context, addRequest *grpc_application_go.AddServiceInstanceRequest) (*grpc_application_go.ServiceInstance, error) {
	err := entities.ValidAddServiceInstanceRequest(addRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid add service instance request")
		return nil, conversions.ToGRPCError(err)
	}
	serviceInstance, err := h.Manager.AddServiceInstance(addRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add service instance")
		return nil, conversions.ToGRPCError(err)
	}
	return serviceInstance.ToGRPC(), nil
}

// GetServiceGroupInstanceMetadata returns the metadata entry of an existing ServiceGroupInstance
func (h *Handler) GetServiceGroupInstanceMetadata(ctx context.Context, getRequest *grpc_application_go.GetServiceGroupInstanceMetadataRequest) (*grpc_application_go.InstanceMetadata, error) {
	err := entities.ValidGetServiceGroupInstanceMetadataRequest(getRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid get service group instance metadata request")
		return nil, conversions.ToGRPCError(err)
	}
	metadata, err := h.Manager.GetServiceGroupInstanceMetadata(getRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get service group instance metadata")
		return nil, conversions.ToGRPCError(err)
	}
	return metadata.ToGRPC(), nil
}

// UpdateServiceGroupInstanceMetadata updates the value of an existing metadata instance
func (h *Handler) UpdateServiceGroupInstanceMetadata(ctx context.Context, updateMetadataRequest *grpc_application_go.InstanceMetadata) (*grpc_common_go.Success, error) {
	err := entities.ValidUpdateInstanceMetadata(updateMetadataRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid update metadata request")
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.UpdateServiceGroupInstanceMetadata(updateMetadataRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot update service group instance metadata")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// AddAppEndPoint adds a new App Endpoint to a given service instance
func (h *Handler) AddAppEndpoint(ctx context.Context, request *grpc_application_go.AddAppEndpointRequest) (*grpc_common_go.Success, error) {
	err := entities.ValidAddAppEndpointRequest(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid add application endpoint request")
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.AddAppEndpoint(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add app endpoint")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// GetAppEndPoint retrieves an appEndpoint
func (h *Handler) GetAppEndpoints(ctx context.Context, request *grpc_application_go.GetAppEndPointRequest) (*grpc_application_go.AppEndpointList, error) {
	err := entities.ValidGetAppEndPointRequest(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid get application endpoint request")
		return nil, conversions.ToGRPCError(err)
	}
	endpoint, err := h.Manager.GetAppEndpoint(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get app endpoint")
		return nil, conversions.ToGRPCError(err)
	}
	return endpoint, nil
}

func (h *Handler) RemoveAppEndpoints(ctx context.Context, request *grpc_application_go.RemoveAppEndpointRequest) (*grpc_common_go.Success, error) {
	err := entities.ValidRemoveEndpointRequest(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid remove app endpoint request")
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.RemoveAppEndpoints(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove app endpoints")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// AddAppZtNetwork add a new zerotier network for an existing app instance
func (h *Handler) AddAppZtNetwork(ctx context.Context, request *grpc_application_go.AddAppZtNetworkRequest) (*grpc_common_go.Success, error) {
	err := entities.ValidAddAppZtNetworkRequest(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid add application ZT network request")
		return nil, err
	}
	err = h.Manager.AddZtNetwork(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add ZT network")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// RemoveAppZtNetwork remove the network instance for an application instance
func (h *Handler) RemoveAppZtNetwork(ctx context.Context, request *grpc_application_go.RemoveAppZtNetworkRequest) (*grpc_common_go.Success, error) {
	err := entities.ValidRemoveAppZtNetworkRequest(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid remove application ZT network request")
		return nil, err
	}
	err = h.Manager.RemoveZtNetwork(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove ZT network")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// GetAppZtnetwork get an existing network instance associated with an application.
func (h *Handler) GetAppZtNetwork(ctx context.Context, request *grpc_application_go.GetAppZtNetworkRequest) (*grpc_application_go.AppZtNetwork, error) {
	err := entities.ValidGetAppZtNetworkRequest(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid get app ZT network request")
		return nil, err
	}
	retrieved, err := h.Manager.GetAppZtNetwork(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get application ZT network")
		return nil, conversions.ToGRPCError(err)
	}
	return retrieved.ToGRPC(), nil
}

// AddParametrizedDescriptor adds a parametrized descriptor to a given descriptor
func (h *Handler) AddParametrizedDescriptor(ctx context.Context, request *grpc_application_go.ParametrizedDescriptor) (*grpc_application_go.ParametrizedDescriptor, error) {
	// TODO validate entity
	param, err := h.Manager.AddParametrizedDescriptor(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add parametrized descriptor")
		return nil, conversions.ToGRPCError(err)
	}
	return param.ToGRPC(), nil
}

// GetParametrizedDescriptor retrieves the parametrized descriptor associated with an instance
func (h *Handler) GetParametrizedDescriptor(ctx context.Context, instanceID *grpc_application_go.AppInstanceId) (*grpc_application_go.ParametrizedDescriptor, error) {
	err := entities.ValidAppInstanceId(instanceID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid app instance identifier")
		return nil, conversions.ToGRPCError(err)
	}
	descriptor, err := h.Manager.GetParametrizedDescriptor(instanceID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get parametrized descriptor")
		return nil, conversions.ToGRPCError(err)
	}
	return descriptor.ToGRPC(), nil
}

// RemoveParametrizedDescriptor removes the parametrized descriptor associated with an instance
func (h *Handler) RemoveParametrizedDescriptor(ctx context.Context, instanceID *grpc_application_go.AppInstanceId) (*grpc_common_go.Success, error) {
	err := entities.ValidAppInstanceId(instanceID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid application instance identifier")
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.RemoveParametrizedDescriptor(instanceID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove parametrized descriptor")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// Add Zt member authorization data
func (h *Handler) AddAuthorizedZtNetworkMember(ctx context.Context, req *grpc_application_go.AddAuthorizedZtNetworkMemberRequest) (*grpc_application_go.ZtNetworkMember, error) {
	err := entities.ValidAddAuthorizedNetworkMemberRequest(req)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid add authorized ZT network member request")
		return nil, err
	}
	storedMember, err := h.Manager.AddAppZtNetworkMember(req)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add application ZT network member")
		return nil, conversions.ToGRPCError(err)
	}
	return storedMember.ToGRPC(), nil
}

// Delete Zt member authorization data
func (h *Handler) RemoveAuthorizedZtNetworkMember(ctx context.Context, req *grpc_application_go.RemoveAuthorizedZtNetworkMemberRequest) (*grpc_common_go.Success, error) {
	err := entities.ValidRemoveAuthorizedZtNetworkMemberRequest(req)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid remove authorized ZT network member request")
		return nil, err
	}
	err = h.Manager.RemoveAppZtNetworkMember(req.OrganizationId, req.AppInstanceId, req.ServiceGroupInstanceId, req.ServiceApplicationInstanceId, req.ZtNetworkId)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove application ZT network member")
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}

// Get the Zt authorized members
func (h *Handler) GetAuthorizedZtNetworkMember(ctx context.Context, req *grpc_application_go.GetAuthorizedZtNetworkMemberRequest) (*grpc_application_go.ZtNetworkMembers, error) {
	err := entities.ValidGetAuthorizedZtNetworkMemberRequest(req)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid get authorized ZT network member request")
		return nil, err
	}
	retrieved, err := h.Manager.GetAppZtNetworkMember(req.OrganizationId, req.AppInstanceId, req.ServiceGroupInstanceId, req.ServiceApplicationInstanceId)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get application ZT network member")
		return nil, err
	}
	return retrieved.ToArrayGRPC(), err
}

// ListAuthorizedZTNetworkMembers retrieves a list of ztMembers in a zero tier network
func (h *Handler)  ListAuthorizedZTNetworkMembers(ctx context.Context, req *grpc_application_go.ListAuthorizedZtNetworkMemberRequest) (*grpc_application_go.ZtNetworkMembers, error){

	vErr := entities.ValidListAuthorizedZtNetworkMemberRequest(req)
	if vErr != nil {
		log.Error().Str("trace", vErr.DebugReport()).Msg("invalid list authorized ZT network member request")
		return nil, conversions.ToGRPCError(vErr)
	}

	retrieved, err := h.Manager.ListAuthorizedZTNetworkMembers(req.OrganizationId, req.AppInstanceId, req.ZtNetworkId)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get application ZT network member")
		return nil, conversions.ToGRPCError(err)
	}

	return retrieved, err

}
// AddZtNetworkProxy adds a new network proxy for an existing private network
func (h *Handler) AddZtNetworkProxy(ctx context.Context, req *grpc_application_go.ServiceProxy) (*grpc_common_go.Success, error) {
	err := entities.ValidAddZtNetworkProxy(req)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid service proxy")
		return nil, err
	}
	err = h.Manager.AddZtNetworkProxy(&entities.ServiceProxy{
		OrganizationId: req.OrganizationId, AppInstanceId: req.AppInstanceId,
		ServiceGroupInstanceId: req.ServiceGroupInstanceId, ClusterId: req.ClusterId,
		FQDN: req.Fqdn, IP: req.Ip, ServiceInstanceId: req.ServiceInstanceId, ServiceGroupId: req.ServiceGroupId,
		ServiceId: req.ServiceId,
	})
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add ZT network proxy")
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}

// RemoveZtnetworkProxy removes a proxy from the list of available entries of a private network
func (h *Handler) RemoveZtNetworkProxy(ctx context.Context, req *grpc_application_go.RemoveAppZtNetworkProxy) (*grpc_common_go.Success, error) {
	err := entities.ValidRemoveZtNetworkProxy(req)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid remove application ZT network proxy")
		return nil, err
	}
	err = h.Manager.RemoveZtNetworkProxy(req.OrganizationId, req.AppInstanceId, req.Fqdn, req.ClusterId, req.ServiceGroupInstanceId, req.ServiceInstanceId)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove ZT network proxy")
		return nil, err
	}
	return nil, nil
}
