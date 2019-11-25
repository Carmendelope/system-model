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

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-history-logs-go"
)

type LogResponse struct {
	// OrganizationId with the organization identifier
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// From defines the timestamp from which the request will be taken into account
	From int64 `json:"from,omitempty" cql:"from"`
	// To defines the timestamp to which the request will be taken into account
	To int64 `json:"to,omitempty" cql:"to"`
	// Events contains the entries of the service instance history result of the query
	Events []ServiceInstanceLog `json:"events,omitempty" cql:"events"`
}

type ServiceInstanceLog struct {
	// OrganizationId with the organization identifier
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// ApplicationDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty" cql:"app_descriptor_id"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty" cql:"app_instance_id"`
	// ServiceGroupId with the group identifier.
	ServiceGroupId string `json:"service_group_id,omitempty" cql:"service_group_id"`
	// ServiceGroupInstanceId  with the group instance identifier.
	ServiceGroupInstanceId string `json:"service_group_instance_id,omitempty" cql:"service_group_instance_id"`
	// ServiceId with the service identifier.
	ServiceId string `json:"service_id,omitempty" cql:"service_id"`
	// ServiceInstanceId with the service instance identifier.
	ServiceInstanceId string `json:"service_instance_id,omitempty" cql:"service_instance_id"`
	// Timestamp when the information of when this service instance was created
	Created int64 `json:"created,omitempty" cql:"created"`
	// Timestamp when the information of when this service instance was terminated
	Terminated int64 `json:"terminated,omitempty" cql:"terminated"`
}

type AddLogRequest struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty" cql:"app_instance_id"`
	// ApplicationDescriptorId with the application descriptor identifier.
	AppDescriptorId string `json:"app_descriptor_id,omitempty" cql:"app_descriptor_id"`
	// ServiceGroupId with the group identifier.
	ServiceGroupId string `json:"service_group_id,omitempty" cql:"service_group_id"`
	// ServiceGroupInstanceId  with the group instance identifier.
	ServiceGroupInstanceId string `json:"service_group_instance_id,omitempty" cql:"service_group_instance_id"`
	// ServiceId with the service identifier.
	ServiceId string `json:"service_id,omitempty" cql:"service_id"`
	// ServiceInstanceId with the service instance identifier.
	ServiceInstanceId string `json:"service_instance_id,omitempty" cql:"service_instance_id"`
	// Created with the timestamp of when the information of when this service instance was created
	Created int64 `json:"created,omitempty" cql:"created"`
}

type UpdateLogRequest struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty" cql:"app_instance_id"`
	// ServiceInstanceId with the service instance identifier.
	ServiceInstanceId string `json:"service_instance_id,omitempty" cql:"service_instance_id"`
	// Timestamp when the information of when this service instance was terminated
	Terminated int64 `json:"terminated,omitempty" cql:"terminated"`
}

type SearchLogsRequest struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// From contains the timestamp from which a service instance was available
	From int64 `json:"available_from,omitempty" cql:"available_from"`
	// To contains the timestamp to which a service instance was available
	To int64 `json:"available_to,omitempty" cql:"available_to"`
}

type RemoveLogRequest struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// AppInstanceId with the application instance identifier.
	AppInstanceId string `json:"app_instance_id,omitempty" cql:"app_instance_id"`
}

func ValidAddLogRequest(addLogRequest *grpc_application_history_logs_go.AddLogRequest) derrors.Error {
	if addLogRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if addLogRequest.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyAppInstanceId)
	}
	if addLogRequest.AppDescriptorId == "" {
		return derrors.NewInvalidArgumentError(emptyAppDescriptorId)
	}
	if addLogRequest.ServiceGroupId == "" {
		return derrors.NewInvalidArgumentError(emptyServiceGroupId)
	}
	if addLogRequest.ServiceGroupInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyServiceGroupInstanceId)
	}
	if addLogRequest.ServiceInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyServiceInstanceId)
	}
	return nil
}

func ValidUpdateLogRequest(updateLogRequest *grpc_application_history_logs_go.UpdateLogRequest) derrors.Error {
	if updateLogRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if updateLogRequest.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyAppInstanceId)
	}
	if updateLogRequest.ServiceInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyServiceInstanceId)
	}
	return nil
}

func ValidSearchLogRequest(searchLogRequest *grpc_application_history_logs_go.SearchLogRequest) derrors.Error {
	if searchLogRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	return nil
}

func ValidRemoveLogRequest(removeLogRequest *grpc_application_history_logs_go.RemoveLogsRequest) derrors.Error {
	if removeLogRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if removeLogRequest.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyAppInstanceId)
	}
	return nil
}

func ToAddLogRequest(addLogRequest grpc_application_history_logs_go.AddLogRequest) AddLogRequest {
	return AddLogRequest{
		OrganizationId:         addLogRequest.OrganizationId,
		AppInstanceId:          addLogRequest.AppInstanceId,
		AppDescriptorId:        addLogRequest.AppDescriptorId,
		ServiceGroupId:         addLogRequest.ServiceGroupId,
		ServiceGroupInstanceId: addLogRequest.ServiceGroupInstanceId,
		ServiceId:              addLogRequest.ServiceId,
		ServiceInstanceId:      addLogRequest.ServiceInstanceId,
		Created:                addLogRequest.Created,
	}
}

func ToUpdateLogRequest(updateLogRequest grpc_application_history_logs_go.UpdateLogRequest) UpdateLogRequest {
	return UpdateLogRequest{
		OrganizationId:    updateLogRequest.OrganizationId,
		AppInstanceId:     updateLogRequest.AppInstanceId,
		ServiceInstanceId: updateLogRequest.ServiceInstanceId,
		Terminated:        updateLogRequest.Terminated,
	}
}

func ToSearchLogsRequest(searchLogsRequest grpc_application_history_logs_go.SearchLogRequest) SearchLogsRequest {
	return SearchLogsRequest{
		OrganizationId: searchLogsRequest.OrganizationId,
		From:           searchLogsRequest.From,
		To:             searchLogsRequest.To,
	}
}

func ToRemoveLogRequest(removeLogRequest grpc_application_history_logs_go.RemoveLogsRequest) RemoveLogRequest {
	return RemoveLogRequest{
		OrganizationId: removeLogRequest.OrganizationId,
		AppInstanceId:  removeLogRequest.AppInstanceId,
	}
}

func ToGRPCLogRequest(logResponse LogResponse) grpc_application_history_logs_go.LogResponse {
	events := make([]*grpc_application_history_logs_go.ServiceInstanceLog, 0)
	for _, event := range logResponse.Events {
		events = append(events, ToGRPCServiceInstanceLog(event))
	}

	return grpc_application_history_logs_go.LogResponse{
		OrganizationId: logResponse.OrganizationId,
		From:           0,
		To:             0,
		Events:         events,
	}
}

func ToGRPCServiceInstanceLog(serviceInstanceLog ServiceInstanceLog) *grpc_application_history_logs_go.ServiceInstanceLog {
	return &grpc_application_history_logs_go.ServiceInstanceLog{
		OrganizationId:         serviceInstanceLog.OrganizationId,
		AppDescriptorId:        serviceInstanceLog.AppDescriptorId,
		AppInstanceId:          serviceInstanceLog.AppInstanceId,
		ServiceGroupId:         serviceInstanceLog.ServiceGroupId,
		ServiceGroupInstanceId: serviceInstanceLog.ServiceGroupInstanceId,
		ServiceId:              serviceInstanceLog.ServiceId,
		ServiceInstanceId:      serviceInstanceLog.ServiceInstanceId,
		Created:                serviceInstanceLog.Created,
		Terminated:             serviceInstanceLog.Terminated,
	}
}
