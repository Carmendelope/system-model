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

type Handler struct {
	Manager Manager
}

func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}

func (h *Handler) AddAppDescriptor(ctx context.Context, addRequest *grpc_application_go.AddAppDescriptorRequest) (*grpc_application_go.AppDescriptor, error) {
	err := entities.ValidAddAppDescriptorRequest(addRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Interface("addRequest", addRequest).Msg("Adding application descriptor")
	added, err := h.Manager.AddAppDescriptor(addRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return added.ToGRPC(), nil
}

func (h *Handler) GetAppDescriptors(ctx context.Context, orgID *grpc_organization_go.OrganizationId) (*grpc_application_go.AppDescriptorList, error) {
	descriptors, err := h.Manager.ListDescriptors(orgID)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}

	toReturn := make([]*grpc_application_go.AppDescriptor, 0)
	for _, d := range descriptors {
		toReturn = append(toReturn, d.ToGRPC())
	}
	result := &grpc_application_go.AppDescriptorList{
		Descriptors:          toReturn,
	}
	return result, nil
}

func (h *Handler) GetAppDescriptor(ctx context.Context, appDescID *grpc_application_go.AppDescriptorId) (*grpc_application_go.AppDescriptor, error) {
	descriptor, err := h.Manager.GetDescriptor(appDescID)
	if err != nil {
	    return nil, conversions.ToGRPCError(err)
	}
	return descriptor.ToGRPC(), nil
}


func (h *Handler) AddAppInstance(ctx context.Context, addInstanceRequest *grpc_application_go.AddAppInstanceRequest) (*grpc_application_go.AppInstance, error) {
	panic("implement me")
}

func (h *Handler) GetAppInstances(ctx context.Context, orgID *grpc_organization_go.OrganizationId) (*grpc_application_go.AppInstanceList, error) {
	panic("implement me")
}

func (h *Handler) GetAppInstance(ctx context.Context, appInstID *grpc_application_go.AppInstanceId) (*grpc_application_go.AppInstance, error) {
	panic("implement me")
}

func (h *Handler) UpdateAppStatus(ctx context.Context, updateAppStatus *grpc_application_go.UpdateAppStatusRequest) (*grpc_common_go.Success, error) {
	panic("implement me")
}

func (h *Handler) UpdateServiceStatus(ctx context.Context, updateServiceStatus *grpc_application_go.UpdateServiceStatusRequest) (*grpc_common_go.Success, error) {
	panic("implement me")
}


