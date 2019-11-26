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
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities/devices"
)

type Provider interface {

	// AddDeviceGroup adds a new device group
	AddDeviceGroup(deviceGroup devices.DeviceGroup) derrors.Error
	// ExistsDeviceGroup checks if a group exists on the system.
	ExistsDeviceGroup(organizationID string, deviceGroupID string) (bool, derrors.Error)
	// ExistsDeviceGroupByName checks if a group exists on the system.
	ExistsDeviceGroupByName(organizationID string, name string) (bool, derrors.Error)
	// GetDeviceGroup returns a device Group.
	GetDeviceGroup(organizationID string, deviceGroupID string) (*devices.DeviceGroup, derrors.Error)
	// ListDeviceGroups returns a list of device groups in a organization.
	ListDeviceGroups(organizationID string) ([]devices.DeviceGroup, derrors.Error)
	// GetDeviceGroupsByName returns a list o devices which names are in groupName list
	GetDeviceGroupsByName(organizationID string, groupNames []string) ([]devices.DeviceGroup, derrors.Error)
	// Remove a device group
	RemoveDeviceGroup(organizationID string, deviceGroup string) derrors.Error

	// AddDevice adds a new device group
	AddDevice(device devices.Device) derrors.Error
	// ExistsDevice checks if a device exists on the system.
	ExistsDevice(organizationID string, deviceGroupID string, deviceID string) (bool, derrors.Error)
	// GetDevice returns a device .
	GetDevice(organizationID string, deviceGroupID string, deviceID string) (*devices.Device, derrors.Error)
	// ListDevice returns a list of device in a group.
	ListDevices(organizationID string, deviceGroupID string) ([]devices.Device, derrors.Error)
	// Remove a device
	RemoveDevice(organizationID string, deviceGroupID string, deviceID string) derrors.Error
	//UpdateDevice updates the device information
	UpdateDevice(device devices.Device) derrors.Error

	Clear() derrors.Error
}
