package entities

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
	//
	AppDescriptorId string `json:"app_descriptor_id,omitempty" cql:"app_descriptor_id"`
	//
	AppInstanceId string `json:"app_instance_id,omitempty" cql:"app_instance_id"`
	//
	ServiceGroupId string `json:"service_group_id,omitempty" cql:"service_group_id"`
	//
	ServiceGroupInstanceId string `json:"service_group_instance_id,omitempty" cql:"service_group_instance_id"`
	//
	ServiceId string `json:"service_id,omitempty" cql:"service_id"`
	//
	ServiceInstanceId string `json:"service_instance_id,omitempty" cql:"service_instance_id"`
	//
	Created int64 `json:"created,omitempty" cql:"created"`
	//
	Terminated int64 `json:"terminated,omitempty" cql:"terminated"`
}

type AddLogRequest struct {

}

type UpdateLogRequest struct {

}

type SearchLogsRequest struct {

}

type RemoveLogRequest struct {

}
