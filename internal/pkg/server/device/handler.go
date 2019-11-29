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

package device

import (
	"context"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/entities/devices"
	"github.com/rs/zerolog/log"
)

// Handler structure for the application requests.
type Handler struct {
	Manager Manager
} // NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

// AddDeviceGroup adds a new device group to the system.
func (h *Handler) AddDeviceGroup(ctx context.Context, addRequest *grpc_device_go.AddDeviceGroupRequest) (*grpc_device_go.DeviceGroup, error) {
	err := devices.ValidAddDeviceGroupRequest(addRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid add device group request")
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Interface("addRequest", addRequest).Msg("Adding device group")
	added, err := h.Manager.AddDeviceGroup(addRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add device group")
		return nil, conversions.ToGRPCError(err)
	}
	return added.ToGRPC(), nil
}

// ListDeviceGroups obtains a list of device groups in an organization.
func (h *Handler) ListDeviceGroups(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_device_go.DeviceGroupList, error) {
	err := entities.ValidOrganizationID(organizationID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid organization identifier")
		return nil, conversions.ToGRPCError(err)
	}
	groups, err := h.Manager.ListDeviceGroups(organizationID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot list device groups")
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
func (h *Handler) GetDeviceGroup(ctx context.Context, DeviceGroupID *grpc_device_go.DeviceGroupId) (*grpc_device_go.DeviceGroup, error) {
	err := devices.ValidDeviceGroupId(DeviceGroupID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid device group identifier")
		return nil, conversions.ToGRPCError(err)
	}
	group, err := h.Manager.GetDeviceGroup(DeviceGroupID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get device group")
		return nil, conversions.ToGRPCError(err)
	}
	return group.ToGRPC(), nil
}

// RemoveDeviceGroup removes a device group
func (h *Handler) RemoveDeviceGroup(ctx context.Context, removeRequest *grpc_device_go.RemoveDeviceGroupRequest) (*grpc_common_go.Success, error) {
	err := devices.ValidRemoveDeviceGroupRequest(removeRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid remove device request")
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.RemoveDeviceGroup(removeRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove device group")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// GetDeviceGroupsByNames obtains a list the device groups .
func (h *Handler) GetDeviceGroupsByNames(ctx context.Context, request *grpc_device_go.GetDeviceGroupsRequest) (*grpc_device_go.DeviceGroupList, error) {
	err := devices.ValidGetDeviceGroupsRequest(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid get device groups request")
		return nil, conversions.ToGRPCError(err)
	}
	groups, err := h.Manager.GetDeviceGroupsByNames(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get device groups by name")
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

// ------------------------------------------------------------------------------------------------------------------

// AddDevice adds a new group to the system
func (h *Handler) AddDevice(ctx context.Context, addRequest *grpc_device_go.AddDeviceRequest) (*grpc_device_go.Device, error) {
	err := devices.ValidAddDeviceRequest(addRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid add device request")
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Interface("addRequest", addRequest).Msg("Adding device")
	added, err := h.Manager.AddDevice(addRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add device")
		return nil, conversions.ToGRPCError(err)
	}
	return added.ToGRPC(), nil
}

// ListDevice obtains a list of devices in a device_group
func (h *Handler) ListDevices(ctx context.Context, deviceGroupRequest *grpc_device_go.DeviceGroupId) (*grpc_device_go.DeviceList, error) {
	err := devices.ValidDeviceGroupId(deviceGroupRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid device group identifier")
		return nil, conversions.ToGRPCError(err)
	}
	devices, err := h.Manager.ListDevices(deviceGroupRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot list devices")
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
func (h *Handler) GetDevice(ctx context.Context, deviceRequest *grpc_device_go.DeviceId) (*grpc_device_go.Device, error) {
	err := devices.ValidDeviceID(deviceRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid device identifier")
		return nil, conversions.ToGRPCError(err)
	}
	device, err := h.Manager.GetDevice(deviceRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get device")
		return nil, conversions.ToGRPCError(err)
	}
	return device.ToGRPC(), nil
}

// RemoveDevice removes a given device
func (h *Handler) RemoveDevice(ctx context.Context, removeRequest *grpc_device_go.RemoveDeviceRequest) (*grpc_common_go.Success, error) {
	err := devices.ValidRemoveDeviceRequest(removeRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid remove device request")
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.RemoveDevice(removeRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove device")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// UpdateDevice updates the device info (labels)
func (h *Handler) UpdateDevice(ctx context.Context, request *grpc_device_go.UpdateDeviceRequest) (*grpc_device_go.Device, error) {
	err := devices.ValidUpdateDeviceRequest(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid update device request")
		return nil, conversions.ToGRPCError(err)
	}
	updated, err := h.Manager.UpdateDevice(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot update device")
		return nil, conversions.ToGRPCError(err)
	}
	return updated.ToGRPC(), nil
}
