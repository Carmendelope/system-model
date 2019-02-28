package device

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"time"
)

//  Device model the information available regarding a Device of an organization
type Device struct {
	OrganizationId	string
	DeviceGroupId 	string
	DeviceId 		string
	RegisterSince	int64
	Labels			map[string]string
}

type DeviceGroup struct {
	OrganizationId 	string
	DeviceGroupId	string
	Name 			string
	Created  		int64
	Labels			map[string]string
}

func NewDeviceGroup (organizationID string, deviceGroupID string, name string, labels map[string]string) * DeviceGroup {
	return &DeviceGroup{
		OrganizationId:organizationID,
		DeviceGroupId: deviceGroupID,
		Name : name,
		Created: time.Now().Unix(),
		Labels: labels,
	}
}


// ----------- Device Group ----------- //
func NewDeviceGroupFromGRPC (addRequest * grpc_device_go.AddDeviceGroupRequest) * DeviceGroup {
	if addRequest == nil {
		return nil
	}

	return &DeviceGroup{
		OrganizationId: addRequest.OrganizationId,
		DeviceGroupId: entities.GenerateUUID(),
		Name: addRequest.Name,
		Labels: addRequest.Labels,
		Created: time.Now().Unix(),
	}

}

func (d * DeviceGroup) ToGRPC() * grpc_device_go.DeviceGroup {
	return &grpc_device_go.DeviceGroup{
		OrganizationId: d.OrganizationId,
		DeviceGroupId:	d.DeviceGroupId,
		Name:			d.Name,
		Labels:			d.Labels,
		Created:		d.Created,
	}
}

func ValidAddDeviceGroupRequest (addRequest * grpc_device_go.AddDeviceGroupRequest) derrors.Error {

	if addRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}

	return nil
}

func ValidDeviceGroupId (deviceGroup * grpc_device_go.DeviceGroupId) derrors.Error {

	if deviceGroup.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if deviceGroup.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError("device_group_id cannot be empty")
	}

	return nil
}

func ValidRemoveDeviceGroupRequest (removeRequest * grpc_device_go.RemoveDeviceGroupRequest) derrors.Error {
	if removeRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if removeRequest.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError("device_group_id cannot be empty")
	}

	return nil
}

// ----------- Device ----------- //
func NewDeviceFromGRPC (addRequest * grpc_device_go.AddDeviceRequest) * Device{
	return &Device{
		OrganizationId: addRequest.OrganizationId,
		DeviceGroupId: addRequest.DeviceGroupId,
		DeviceId:     addRequest.DeviceId,
		Labels:			addRequest.Labels,
		RegisterSince: time.Now().Unix(),
	}
}

func (d * Device) ToGRPC() *grpc_device_go.Device {
	return &grpc_device_go.Device{
		OrganizationId: d.OrganizationId,
		DeviceGroupId: d.DeviceGroupId,
		DeviceId: d.DeviceId,
		RegisterSince: d.RegisterSince,
		Labels:d.Labels,
	}
}

func (d *Device) ApplyUpdate(updateRequest grpc_device_go.UpdateDeviceRequest) {

	if updateRequest.AddLabels {
		if d.Labels == nil {
			d.Labels = make(map[string]string, 0)
		}
		for k, v := range updateRequest.Labels {
			d.Labels[k] = v
		}
	}
	if updateRequest.RemoveLabels {
		for k, _ := range updateRequest.Labels {
			delete(d.Labels, k)
		}
	}
}

func ValidDeviceID (device * grpc_device_go.DeviceId) derrors.Error {

	if device.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if device.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError("device_group_id cannot be empty")
	}
	if device.DeviceId == "" {
		return derrors.NewInvalidArgumentError("device_id cannot be empty")
	}

	return nil
}

func ValidAddDeviceRequest (request * grpc_device_go.AddDeviceRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError("device_group_id cannot be empty")
	}
	if request.DeviceId == "" {
		return derrors.NewInvalidArgumentError("device_id cannot be empty")
	}
	return nil
}

func ValidRemoveDeviceRequest (request * grpc_device_go.RemoveDeviceRequest) derrors.Error {

	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError("device_group_id cannot be empty")
	}
	if request.DeviceId == "" {
		return derrors.NewInvalidArgumentError("device_id cannot be empty")
	}
	return nil
 }

func ValidUpdateDeviceRequest(request * grpc_device_go.UpdateDeviceRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError("device_group_id cannot be empty")
	}
	if request.DeviceId == "" {
		return derrors.NewInvalidArgumentError("device_id cannot be empty")
	}
	if request.AddLabels && request.RemoveLabels {
		return derrors.NewInvalidArgumentError("add_labels and remove_labels can not be true at the same time")
	}
	if ! request.AddLabels && !request.RemoveLabels {
		return derrors.NewInvalidArgumentError("add_labels and remove_labels cannot be false at the same time")
	}
	if request == nil || len(request.Labels) == 0{
		return derrors.NewInvalidArgumentError("labels cannot be empty")
	}

	return nil
}

