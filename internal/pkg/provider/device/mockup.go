package device

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities/device"
	"sync"
)

type MockupDeviceProvider struct {
	sync.Mutex
	// deviceGroup indexed by organization_id -> device_group_id
	deviceGroups map[string]map[string]device.DeviceGroup
	// devices indexed by (organization_id#device_group_id) -> device_id
	devices 	 map[string]map[string]device.Device

}

func NewMockupDeviceProvider () * MockupDeviceProvider {
	return &MockupDeviceProvider{
		deviceGroups:	make (map[string]map[string]device.DeviceGroup, 0),
		devices: 		make (map[string]map[string]device.Device, 0),
	}
}

func (m * MockupDeviceProvider) unsafeExistsGroup (organizationID string, deviceGroupID string) bool {
	groups, exists := m.deviceGroups[organizationID]

	if ! exists {
		return false
	}

	_, exists = groups[deviceGroupID]

	return exists
}

func (m * MockupDeviceProvider) unsafeExistsOrganization (organizationID string) bool {
	_, exists := m.deviceGroups[organizationID]
	return exists
}

func (m * MockupDeviceProvider) AddDeviceGroup (deviceGroup device.DeviceGroup) derrors.Error {
	m.Lock()
	defer m.Unlock()

	if !m.unsafeExistsOrganization(deviceGroup.OrganizationId){
		device := make(map[string]device.DeviceGroup)
		device[deviceGroup.DeviceGroupId] = deviceGroup
		m.deviceGroups[deviceGroup.OrganizationId] = device
	} else if !m.unsafeExistsGroup(deviceGroup.OrganizationId, deviceGroup.DeviceGroupId) {
			device := m.deviceGroups[deviceGroup.OrganizationId]
			device[deviceGroup.DeviceGroupId] = deviceGroup
			m.deviceGroups[deviceGroup.OrganizationId] = device
	} else{
		return derrors.NewAlreadyExistsError("Add device group").WithParams(deviceGroup.OrganizationId, deviceGroup.DeviceGroupId)
	}
	return nil
}

func (m * MockupDeviceProvider) ExistsDeviceGroup(organizationID string, deviceGroupID string) (bool, derrors.Error){
	m.Lock()
	defer m.Unlock()

	return m.unsafeExistsGroup(organizationID, deviceGroupID), nil
}


func (m * MockupDeviceProvider) GetDeviceGroup(organizationID string, deviceGroupID string) (* device.DeviceGroup, derrors.Error) {

	m.Lock()
	defer m.Unlock()

	groups, exists := m.deviceGroups[organizationID]
	if  exists {
		group , exists := groups[deviceGroupID]
		if exists{
			return &group, nil
		}
	}

	return  nil, derrors.NewNotFoundError("device group").WithParams(organizationID, deviceGroupID)
}

func (m * MockupDeviceProvider) ListDeviceGroups(organizationID string) ([]device.DeviceGroup, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	groups, exists := m.deviceGroups[organizationID]
	list := make([]device.DeviceGroup, 0)

	if  ! exists {
		return list, nil
	}

	for _, group := range groups{
		list = append(list, group)
	}
	return list, nil

}
func (m * MockupDeviceProvider) RemoveDeviceGroup(organizationID string, deviceGroupID string) derrors.Error {

	m.Lock()
	defer m.Unlock()

	groups, exists := m.deviceGroups[organizationID]
	if  exists {
		group , exists := groups[deviceGroupID]
		if exists{
			delete(groups, group.DeviceGroupId)
			return nil
		}
	}
	return derrors.NewNotFoundError("device group").WithParams(organizationID, deviceGroupID)

}

// ----------------------------------------------------------------------------------------------------

func CreateDeviceIndex(organizationID string, deviceGroupID string) string{
	return fmt.Sprintf("%s#%s", organizationID, deviceGroupID)
}

func (m * MockupDeviceProvider) unsafeExistDevicesInGroup (organizationID string, deviceGroupID string) bool {
	key := CreateDeviceIndex(organizationID, deviceGroupID)

	_, exists := m.devices[key]

	return exists
}

func (m * MockupDeviceProvider) unsafeExistsDevice (organizationID string, deviceGroupID string, deviceID string) bool {
	key := CreateDeviceIndex(organizationID, deviceGroupID)

	list, exists := m.devices[key]

	if ! exists {
		return false
	}

	_, exists = list[deviceID]

	return exists
}

func (m * MockupDeviceProvider) AddDevice (dev device.Device) derrors.Error {
	m.Lock()
	defer m.Unlock()

	key := CreateDeviceIndex(dev.OrganizationId, dev.DeviceGroupId)

	if !m.unsafeExistDevicesInGroup(dev.OrganizationId, dev.DeviceGroupId){
		device := make(map[string]device.Device)
		device[dev.DeviceId] = dev
		m.devices[key] = device
	} else if !m.unsafeExistsDevice(dev.OrganizationId, dev.DeviceGroupId, dev.DeviceId) {
		device := m.devices[key]
		device[dev.DeviceId] = dev
	} else{
		return derrors.NewAlreadyExistsError("Add device ").WithParams(dev.OrganizationId, dev.DeviceGroupId, dev.DeviceId)
	}
	return nil
}

func (m * MockupDeviceProvider) ExistsDevice(organizationID string, deviceGroupID string, deviceID string) (bool, derrors.Error){

	m.Lock()
	defer m.Unlock()

	return m.unsafeExistsDevice(organizationID, deviceGroupID, deviceID), nil
}

func (m * MockupDeviceProvider) GetDevice(organizationID string, deviceGroupID string, deviceID string) (* device.Device, derrors.Error){

	m.Lock()
	defer m.Unlock()

	key := CreateDeviceIndex(organizationID, deviceGroupID)

	devices, exists := m.devices[key]
	if exists {
		dev, exists := devices[deviceID]
		if exists{
			return &dev, nil
		}
	}

	return nil, derrors.NewNotFoundError("device").WithParams(organizationID, deviceGroupID, deviceID)

}
func (m * MockupDeviceProvider) ListDevice(organizationID string, deviceGroupID string) ([]device.Device, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	key := CreateDeviceIndex(organizationID, deviceGroupID)

	devices, exists := m.devices[key]

	if exists{
		devList := make([]device.Device, 0)
		for _, dev := range devices{
			devList = append(devList, dev)
		}
		return devList, nil
	}

	return nil, derrors.NewNotFoundError("devices list").WithParams(organizationID, deviceGroupID)

}
func (m * MockupDeviceProvider) RemoveDevice(organizationID string, deviceGroupID string, deviceID string) derrors.Error{
	m.Lock()
	defer m.Unlock()

	key := CreateDeviceIndex(organizationID, deviceGroupID)

	devices, exists := m.devices[key]
	if  exists {
		dev , exists := devices[deviceID]
		if exists{
			delete(devices, dev.DeviceId)
			return nil
		}
	}
	return derrors.NewNotFoundError("device").WithParams(organizationID, deviceGroupID, deviceID)
}