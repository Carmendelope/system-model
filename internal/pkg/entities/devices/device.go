package devices

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-device-go"
	grpc_device_manager_go "github.com/nalej/grpc-device-manager-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"time"
)

//  Device model the information available regarding a Device of an organization
type Device struct {
	OrganizationId	string                           `json:"organization_id,omitempty"`
	DeviceGroupId 	string                           `json:"device_group_id,omitempty"`
	DeviceId 		string                           `json:"device_id,omitempty"`
	RegisterSince	int64                            `json:"register_since,omitempty"`
	Labels			map[string]string                `json:"labels,omitempty"`
	Os 				*entities.OperatingSystemInfo    `json:"os,omitempty" cql:"os"`
	Hardware 		*entities.HardwareInfo           `json:"hardware,omitempty" cql:"hardware"`
	Storage 		[]*entities.StorageHardwareInfo  `json:"storage,omitempty" cql:"storage"`
	Location        *entities.InventoryLocation      `json:"location,omitempty"`
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
		DeviceGroupId:  entities.GenerateUUID(),
		Name:           addRequest.Name,
		Labels:         addRequest.Labels,
		Created:        time.Now().Unix(),
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

func ValidGetDeviceGroupsRequest (request *grpc_device_go.GetDeviceGroupsRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if len(request.DeviceGroupNames) == 0{
		return derrors.NewInvalidArgumentError("device_group_names cannot be empty")
	}

	return nil
}

// ----------- Device ----------- //
func NewDeviceFromGRPC (addRequest * grpc_device_go.AddDeviceRequest) * Device{

	var os *entities.OperatingSystemInfo
	var hardware 		*entities.HardwareInfo
	var storage 		[]*entities.StorageHardwareInfo
	storage = make ([]*entities.StorageHardwareInfo, 0)

	if addRequest.AssetInfo != nil {
		os = entities.NewOperatingSystemInfoFromGRPC(addRequest.AssetInfo.Os)
		hardware = entities.NewHardwareInfoFromGRPC(addRequest.AssetInfo.Hardware)
		for _, sto := range addRequest.AssetInfo.Storage {
			storage = append(storage, entities.NewStorageHardwareInfoFromGRPC(sto))
		}
	}

	return &Device{
		OrganizationId:	addRequest.OrganizationId,
		DeviceGroupId: 	addRequest.DeviceGroupId,
		DeviceId:     	addRequest.DeviceId,
		Labels:			addRequest.Labels,
		RegisterSince: 	time.Now().Unix(),
		Os: 		  	os,
		Hardware: 		hardware,
		Storage: 		storage,
	}
}

func (d * Device) ToGRPC() *grpc_device_go.Device {

	storage := make ([]*grpc_inventory_go.StorageHardwareInfo, 0)
	for _, sto := range d.Storage {
		storage = append(storage, sto.ToGRPC())
	}

	return &grpc_device_go.Device{
		OrganizationId: d.OrganizationId,
		DeviceGroupId: d.DeviceGroupId,
		DeviceId: d.DeviceId,
		RegisterSince: d.RegisterSince,
		Labels:d.Labels,
		AssetInfo: &grpc_inventory_go.AssetInfo{
			Os: d.Os.ToGRPC(),
			Hardware: d.Hardware.ToGRPC(),
			Storage:storage,
		},
	}
}

func (d * Device) ToGRPCDeviceManager() *grpc_device_manager_go.Device {

	storage := make ([]*grpc_inventory_go.StorageHardwareInfo, 0)
	for _, sto := range d.Storage {
		storage = append(storage, sto.ToGRPC())
	}

	return &grpc_device_manager_go.Device{
		OrganizationId: d.OrganizationId,
		DeviceGroupId: d.DeviceGroupId,
		DeviceId: d.DeviceId,
		RegisterSince: d.RegisterSince,
		Labels:d.Labels,
		AssetInfo: &grpc_inventory_go.AssetInfo{
			Os:       d.Os.ToGRPC(),
			Hardware: d.Hardware.ToGRPC(),
			Storage:  storage,
		},
		Location: d.Location.ToGRPC(),
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

func (d * Device) ApplyLocationUpdate (request *grpc_device_manager_go.UpdateDeviceLocationRequest) {
	if request.UpdateLocation {
		d.Location.Geolocation = request.Location.Geolocation
		d.Location.Geohash = request.Location.Geohash
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
	if request == nil || len(request.Labels) == 0 || request.Location == nil {
		return derrors.NewInvalidArgumentError("request cannot be empty")
	}

	return nil
}

func ValidUpdateDeviceLocationRequest(request * grpc_device_manager_go.UpdateDeviceLocationRequest) derrors.Error {
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id cannot be empty")
	}
	if request.DeviceGroupId == "" {
		return derrors.NewInvalidArgumentError("device_group_id cannot be empty")
	}
	if request.DeviceId == "" {
		return derrors.NewInvalidArgumentError("device_id cannot be empty")
	}
	if request.Location != nil && request.Location.Geolocation == "" {
		return derrors.NewInvalidArgumentError("location cannot be empty")
	}

	return nil
}
