package device

import (
	"context"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/entities/device"
	"github.com/rs/zerolog/log"
)

// Handler structure for the application requests.
type Handler struct {
	Manager Manager

}// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}

// AddDeviceGroup adds a new device group to the system.
func (h *Handler) AddDeviceGroup(ctx context.Context, addRequest *grpc_device_go.AddDeviceGroupRequest) (*grpc_device_go.DeviceGroup, error){

	err := device.ValidAddDeviceGroupRequest(addRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Interface("addRequest", addRequest).Msg("Adding device group")
	added, err := h.Manager.AddDeviceGroup(addRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return added.ToGRPC(), nil
}
// ListDeviceGroups obtains a list of device groups in an organization.
func (h *Handler) ListDeviceGroups(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_device_go.DeviceGroupList, error){

	err := entities.ValidOrganizationID(organizationID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	groups, err := h.Manager.ListDeviceGroups(organizationID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	toReturn := make([]*grpc_device_go.DeviceGroup, 0)
	for _, c := range groups {
		toReturn = append(toReturn, c.ToGRPC())
	}
	result := &grpc_device_go.DeviceGroupList{
		Groups: toReturn,
	}
	return result, nil
}
// GetDeviceGroup retrieves a given device group in an organization.
func (h *Handler) GetDeviceGroup(ctx context.Context, DeviceGroupID *grpc_device_go.DeviceGroupId) (*grpc_device_go.DeviceGroup, error){
	err := device.ValidDeviceGroupId(DeviceGroupID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	group, err := h.Manager.GetDeviceGroup(DeviceGroupID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return group.ToGRPC(), nil
}
// RemoveDeviceGroup removes a device group
func (h *Handler) RemoveDeviceGroup(ctx context.Context, removeRequest *grpc_device_go.RemoveDeviceGroupRequest) (*grpc_common_go.Success, error){

	err := device.ValidRemoveDeviceGroupRequest(removeRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.RemoveDeviceGroup(removeRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// ------------------------------------------------------------------------------------------------------------------

// AddDevice adds a new group to the system
func (h *Handler) AddDevice(ctx context.Context, addRequest *grpc_device_go.AddDeviceRequest) (*grpc_device_go.Device, error){
	err := device.ValidAddDeviceRequest(addRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Interface("addRequest", addRequest).Msg("Adding device")
	added, err := h.Manager.AddDevice(addRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return added.ToGRPC(), nil
}
// ListDevice obtains a list of devices in a device_group
func (h *Handler) ListDevices(ctx context.Context, deviceGroupRequest *grpc_device_go.DeviceGroupId) (*grpc_device_go.DeviceList, error){

	err := device.ValidDeviceGroupId(deviceGroupRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	devices, err := h.Manager.ListDevices(deviceGroupRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	toReturn := make([]*grpc_device_go.Device, 0)
	for _, c := range devices {
		toReturn = append(toReturn, c.ToGRPC())
	}
	result := &grpc_device_go.DeviceList{
		Devices: toReturn,
	}
	return result, nil
}
// GetDevice retrieves a given device in an organization.
func (h *Handler) GetDevice(ctx context.Context, deviceRequest *grpc_device_go.DeviceId) (*grpc_device_go.Device, error){
	err := device.ValidDeviceID(deviceRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	device, err := h.Manager.GetDevice(deviceRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return device.ToGRPC(), nil
}
// RemoveDevice removes a given device
func (h *Handler) RemoveDevice(ctx context.Context, removeRequest *grpc_device_go.RemoveDeviceRequest) (*grpc_common_go.Success, error) {

	err := device.ValidRemoveDeviceRequest(removeRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.RemoveDevice(removeRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}
// UpdateDEvice updates the device info (labels)
func (h *Handler) UpdateDevice(ctx context.Context, deviceRequest *grpc_device_go.UpdateDeviceRequest) (*grpc_device_go.Device, error){
	err := device.ValidUpdateDeviceRequest(deviceRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	device, err := h.Manager.UpdateDevice(deviceRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return device.ToGRPC(), nil

}
