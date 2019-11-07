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
 *
 */

package device

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities/devices"
	"sync"
)

type MockupDeviceProvider struct {
	sync.Mutex
	// deviceGroup indexed by organization_id -> device_group_id
	deviceGroups map[string]map[string]devices.DeviceGroup
	// deviceGroupsByName indexed by organization_id#device_group_name
	deviceGroupsByName map[string]devices.DeviceGroup
	// devices indexed by (organization_id#device_group_id) -> device_id
	devices map[string]map[string]devices.Device
}

func NewMockupDeviceProvider() *MockupDeviceProvider {
	return &MockupDeviceProvider{
		deviceGroups:       make(map[string]map[string]devices.DeviceGroup, 0),
		deviceGroupsByName: make(map[string]devices.DeviceGroup, 0),
		devices:            make(map[string]map[string]devices.Device, 0),
	}
}

// ----------------------------------------------------------------------------------------------------

func (m *MockupDeviceProvider) unsafeExistsGroup(organizationID string, deviceGroupID string) bool {
	groups, exists := m.deviceGroups[organizationID]

	if !exists {
		return false
	}

	_, exists = groups[deviceGroupID]

	return exists
}

func (m *MockupDeviceProvider) unsafeExistsOrganization(organizationID string) bool {
	_, exists := m.deviceGroups[organizationID]
	return exists
}

func (m *MockupDeviceProvider) AddDeviceGroup(deviceGroup devices.DeviceGroup) derrors.Error {
	m.Lock()
	defer m.Unlock()

	if !m.unsafeExistsOrganization(deviceGroup.OrganizationId) {
		device := make(map[string]devices.DeviceGroup)
		device[deviceGroup.DeviceGroupId] = deviceGroup
		m.deviceGroups[deviceGroup.OrganizationId] = device
	} else if !m.unsafeExistsGroup(deviceGroup.OrganizationId, deviceGroup.DeviceGroupId) {
		device := m.deviceGroups[deviceGroup.OrganizationId]
		device[deviceGroup.DeviceGroupId] = deviceGroup
		m.deviceGroups[deviceGroup.OrganizationId] = device
	} else {
		return derrors.NewAlreadyExistsError("Add device group").WithParams(deviceGroup.OrganizationId, deviceGroup.DeviceGroupId)
	}

	//
	m.deviceGroupsByName[fmt.Sprintf("%s#%s", deviceGroup.OrganizationId, deviceGroup.Name)] = deviceGroup
	return nil
}

func (m *MockupDeviceProvider) ExistsDeviceGroup(organizationID string, deviceGroupID string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	return m.unsafeExistsGroup(organizationID, deviceGroupID), nil
}

func (m *MockupDeviceProvider) ExistsDeviceGroupByName(organizationID string, name string) (bool, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	_, exists := m.deviceGroupsByName[fmt.Sprintf("%s#%s", organizationID, name)]

	return exists, nil

}

func (m *MockupDeviceProvider) GetDeviceGroup(organizationID string, deviceGroupID string) (*devices.DeviceGroup, derrors.Error) {

	m.Lock()
	defer m.Unlock()

	groups, exists := m.deviceGroups[organizationID]
	if exists {
		group, exists := groups[deviceGroupID]
		if exists {
			return &group, nil
		}
	}

	return nil, derrors.NewNotFoundError("device group").WithParams(organizationID, deviceGroupID)
}

func (m *MockupDeviceProvider) GetDeviceGroupsByName(organizationID string, groupNames []string) ([]devices.DeviceGroup, derrors.Error) {

	m.Lock()
	defer m.Unlock()

	deviceGroups := make([]devices.DeviceGroup, 0)

	for _, name := range groupNames {
		group, exists := m.deviceGroupsByName[fmt.Sprintf("%s#%s", organizationID, name)]

		if exists {
			deviceGroups = append(deviceGroups, group)
		}

	}

	return deviceGroups, nil
}

func (m *MockupDeviceProvider) ListDeviceGroups(organizationID string) ([]devices.DeviceGroup, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	groups, exists := m.deviceGroups[organizationID]
	list := make([]devices.DeviceGroup, 0)

	if !exists {
		return list, nil
	}

	for _, group := range groups {
		list = append(list, group)
	}
	return list, nil

}

func (m *MockupDeviceProvider) RemoveDeviceGroup(organizationID string, deviceGroupID string) derrors.Error {

	m.Lock()
	defer m.Unlock()

	groups, exists := m.deviceGroups[organizationID]
	if exists {
		group, exists := groups[deviceGroupID]
		if exists {
			delete(m.deviceGroupsByName, fmt.Sprintf("%s#%s", organizationID, group.Name))

			if len(groups) == 1 {
				delete(m.deviceGroups, organizationID)
			} else {
				delete(groups, group.DeviceGroupId)
			}
			return nil
		}
	}

	return derrors.NewNotFoundError("device group").WithParams(organizationID, deviceGroupID)

}

// ----------------------------------------------------------------------------------------------------

func CreateDeviceIndex(organizationID string, deviceGroupID string) string {
	return fmt.Sprintf("%s#%s", organizationID, deviceGroupID)
}

func (m *MockupDeviceProvider) unsafeExistDevicesInGroup(organizationID string, deviceGroupID string) bool {
	key := CreateDeviceIndex(organizationID, deviceGroupID)

	_, exists := m.devices[key]

	return exists
}

func (m *MockupDeviceProvider) unsafeExistsDevice(organizationID string, deviceGroupID string, deviceID string) bool {
	key := CreateDeviceIndex(organizationID, deviceGroupID)

	list, exists := m.devices[key]

	if !exists {
		return false
	}

	_, exists = list[deviceID]

	return exists
}

func (m *MockupDeviceProvider) AddDevice(dev devices.Device) derrors.Error {
	m.Lock()
	defer m.Unlock()

	key := CreateDeviceIndex(dev.OrganizationId, dev.DeviceGroupId)

	if !m.unsafeExistDevicesInGroup(dev.OrganizationId, dev.DeviceGroupId) {
		device := make(map[string]devices.Device)
		device[dev.DeviceId] = dev
		m.devices[key] = device
	} else if !m.unsafeExistsDevice(dev.OrganizationId, dev.DeviceGroupId, dev.DeviceId) {
		device := m.devices[key]
		device[dev.DeviceId] = dev
	} else {
		return derrors.NewAlreadyExistsError("Add device ").WithParams(dev.OrganizationId, dev.DeviceGroupId, dev.DeviceId)
	}
	return nil
}

func (m *MockupDeviceProvider) ExistsDevice(organizationID string, deviceGroupID string, deviceID string) (bool, derrors.Error) {

	m.Lock()
	defer m.Unlock()

	return m.unsafeExistsDevice(organizationID, deviceGroupID, deviceID), nil
}

func (m *MockupDeviceProvider) GetDevice(organizationID string, deviceGroupID string, deviceID string) (*devices.Device, derrors.Error) {

	m.Lock()
	defer m.Unlock()

	key := CreateDeviceIndex(organizationID, deviceGroupID)

	devices, exists := m.devices[key]
	if exists {
		dev, exists := devices[deviceID]
		if exists {
			return &dev, nil
		}
	}

	return nil, derrors.NewNotFoundError("device").WithParams(organizationID, deviceGroupID, deviceID)

}

func (m *MockupDeviceProvider) ListDevices(organizationID string, deviceGroupID string) ([]devices.Device, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	key := CreateDeviceIndex(organizationID, deviceGroupID)
	devList := make([]devices.Device, 0)

	devices, exists := m.devices[key]

	if exists {
		for _, dev := range devices {
			devList = append(devList, dev)
		}

	}
	return devList, nil

	// return nil, derrors.NewNotFoundError("devices list").WithParams(organizationID, deviceGroupID)

}

func (m *MockupDeviceProvider) RemoveDevice(organizationID string, deviceGroupID string, deviceID string) derrors.Error {
	m.Lock()
	defer m.Unlock()

	key := CreateDeviceIndex(organizationID, deviceGroupID)

	devices, exists := m.devices[key]
	if exists {
		dev, exists := devices[deviceID]
		if exists {
			if len(devices) == 1 {
				delete(m.devices, key)
			} else {
				delete(devices, dev.DeviceId)
			}
			return nil
		}
	}
	return derrors.NewNotFoundError("device").WithParams(organizationID, deviceGroupID, deviceID)
}

func (m *MockupDeviceProvider) UpdateDevice(device devices.Device) derrors.Error {

	m.Lock()
	defer m.Unlock()

	if !m.unsafeExistsDevice(device.OrganizationId, device.DeviceGroupId, device.DeviceId) {
		return derrors.NewNotFoundError("device").WithParams(device.OrganizationId, device.DeviceGroupId, device.DeviceId)
	}
	key := CreateDeviceIndex(device.OrganizationId, device.DeviceGroupId)
	devices := m.devices[key]
	devices[device.DeviceId] = device

	return nil
}

// ----------------------------------------------------------------------------------------------------

func (m *MockupDeviceProvider) Clear() derrors.Error {
	m.devices = make(map[string]map[string]devices.Device, 0)
	m.deviceGroups = make(map[string]map[string]devices.DeviceGroup, 0)

	return nil
}
