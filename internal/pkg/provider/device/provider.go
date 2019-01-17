package device

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities/device"
)

type Provider interface {

	// AddDeviceGroup adds a new device group
	AddDeviceGroup (deviceGroup device.DeviceGroup) derrors.Error
	// ExistsDeviceGroup checks if a group exists on the system.
	ExistsDeviceGroup(organizationID string, deviceGroupID string) (bool, derrors.Error)
	// GetDeviceGroup returns a device Group.
	GetDeviceGroup(organizationID string, deviceGroupID string) (* device.DeviceGroup, derrors.Error)
	// ListDeviceGroups returns a list of device groups in a organization.
	ListDeviceGroups(organizationID string) ([]device.DeviceGroup, derrors.Error)
	// Remove a device group
	RemoveDeviceGroup(organizationID string, deviceGroup string) derrors.Error

	// AddDevice adds a new device group
	AddDevice (device device.Device) derrors.Error
	// ExistsDevice checks if a device exists on the system.
	ExistsDevice(organizationID string, deviceGroupID string, deviceID string) (bool, derrors.Error)
	// GetDevice returns a device .
	GetDevice(organizationID string, deviceGroupID string, deviceID string) (* device.Device, derrors.Error)
	// ListDevice returns a list of device in a group.
	ListDevice(organizationID string, deviceGroupID string) ([]device.Device, derrors.Error)
	// Remove a device
	RemoveDevice(organizationID string, deviceGroupID string, deviceID string) derrors.Error

}