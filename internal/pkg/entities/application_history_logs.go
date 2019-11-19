package entities

import (
	"github.com/nalej/derrors"
	grpc_application_history_logs_go "github.com/nalej/grpc-application-history-logs-go"
)

type LogResponse struct {
	// OrganizationId with the organization identifier
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	//
	AvailableFrom int64 `json:"available_from,omitempty" cql:"available_from"`
	//
	AvailableTo int64 `json:"available_to,omitempty" cql:"available_to"`
	//
	Events []ServiceInstanceLog `json:"events,omitempty" cql:"events"`
}

type ServiceInstanceLog struct {
	// OrganizationId with the organization identifier
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// ApplicationDescriptorId
	AppDescriptorId string `json:"app_descriptor_id,omitempty" cql:"app_descriptor_id"`
	// ApplicationInstanceId
	AppInstanceId string `json:"app_instance_id,omitempty" cql:"app_instance_id"`
	// ServiceGroupId
	ServiceGroupId string `json:"service_group_id,omitempty" cql:"service_group_id"`
	// ServiceGroupInstanceId
	ServiceGroupInstanceId string `json:"service_group_instance_id,omitempty" cql:"service_group_instance_id"`
	// ServiceId
	ServiceId string `json:"service_id,omitempty" cql:"service_id"`
	// ServiceInstanceId
	ServiceInstanceId string `json:"service_instance_id,omitempty" cql:"service_instance_id"`
	// Timestamp when the information of when this service instance was created
	Created int64 `json:"created,omitempty" cql:"created"`
	// Timestamp when the information of when this service instance was terminated
	Terminated int64 `json:"terminated,omitempty" cql:"terminated"`
}

type AddLogRequest struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// ApplicationInstanceId
	AppInstanceId string `json:"app_instance_id,omitempty" cql:"app_instance_id"`
	// ApplicationDescriptorId
	AppDescriptorId string `json:"app_descriptor_id,omitempty" cql:"app_descriptor_id"`
	// ServiceGroupId
	ServiceGroupId string `json:"service_group_id,omitempty" cql:"service_group_id"`
	// ServiceGroupInstanceId
	ServiceGroupInstanceId string `json:"service_group_instance_id,omitempty" cql:"service_group_instance_id"`
	// ServiceId
	ServiceId string `json:"service_id,omitempty" cql:"service_id"`
	// ServiceInstanceId
	ServiceInstanceId string `json:"service_instance_id,omitempty" cql:"service_instance_id"`
	// Timestamp when the information of when this service instance was created
	Created int64 `json:"created,omitempty" cql:"created"`
}

type UpdateLogRequest struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// ApplicationInstanceId
	AppInstanceId string `json:"app_instance_id,omitempty" cql:"app_instance_id"`
	// ServiceInstanceId
	ServiceInstanceId string `json:"service_instance_id,omitempty" cql:"service_instance_id"`
	// Timestamp when the information of when this service instance was terminated
	Terminated int64 `json:"terminated,omitempty" cql:"terminated"`
}

type SearchLogsRequest struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// AvailableFrom contains the timestamp from which a service instance was available
	AvailableFrom int64 `json:"available_from,omitempty" cql:"available_from"`
	// AvailableTo contains the timestamp to which a service instance was available
	AvailableTo int64 `json:"available_to,omitempty" cql:"available_to"`
}

type RemoveLogRequest struct {
	// OrganizationId with the organization identifier.
	OrganizationId string `json:"organization_id,omitempty" cql:"organization_id"`
	// ApplicationInstanceId
	AppInstanceId string `json:"app_instance_id,omitempty" cql:"app_instance_id"`
}

func ValidAddLogRequest (addLogRequest *grpc_application_history_logs_go.AddLogRequest) derrors.Error {
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

func ValidUpdateLogRequest (updateLogRequest *grpc_application_history_logs_go.UpdateLogRequest) derrors.Error {
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

func ValidSearchLogRequest (searchLogRequest *grpc_application_history_logs_go.SearchLogRequest) derrors.Error {
	if searchLogRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	return nil
}

func ValidRemoveLogRequest (removeLogRequest *grpc_application_history_logs_go.RemoveLogsRequest) derrors.Error {
	if removeLogRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if removeLogRequest.AppInstanceId == "" {
		return derrors.NewInvalidArgumentError(emptyAppInstanceId)
	}
	return nil
}