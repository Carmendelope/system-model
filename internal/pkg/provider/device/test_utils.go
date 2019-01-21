package device

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/entities/device"
	"math/rand"
)


type DeviceTestHelper struct {

}

func NewDeviceTestHepler() *DeviceTestHelper {
	return &DeviceTestHelper{}
}

func (d * DeviceTestHelper) CreateDeviceGroup() * device.DeviceGroup{

	labels := make(map[string]string, 0)
	labels["lab1"] = "etiq_1"
	labels["lab2"] = "etiq_2"

	return &device.DeviceGroup{
		OrganizationId: entities.GenerateUUID(),
		DeviceGroupId: entities.GenerateUUID(),
		Name: fmt.Sprintf("Test-%d Device Group", rand.Intn(10)),
		Created: rand.Int63(),
		Labels: labels,
	}
}
func (d * DeviceTestHelper) CreateOrganizationDeviceGroup(organizationID string) * device.DeviceGroup{

	labels := make(map[string]string, 0)
	labels["lab1"] = "etiq_1"
	labels["lab2"] = "etiq_2"

	return &device.DeviceGroup{
		OrganizationId: organizationID,
		DeviceGroupId: entities.GenerateUUID(),
		Name: fmt.Sprintf("Test-%d Device Group", rand.Intn(10)),
		Created: rand.Int63(),
		Labels: labels,
	}
}

func (d * DeviceTestHelper) CreateDevice() *device.Device {

	labels := make(map[string]string, 0)
	labels["lab1"] = "etiq_1"
	labels["lab2"] = "etiq_2"

	return &device.Device{
		OrganizationId: entities.GenerateUUID(),
		DeviceGroupId: entities.GenerateUUID(),
		DeviceId: entities.GenerateUUID(),
		RegisterSince: rand.Int63(),
		Labels: labels,
	}
}

func (d * DeviceTestHelper) CreateGroupDevices(organizationID string, deviceGroupID string) *device.Device {

	labels := make(map[string]string, 0)
	labels["lab1"] = "etiq_1"
	labels["lab2"] = "etiq_2"

	return &device.Device{
		OrganizationId: organizationID,
		DeviceGroupId: deviceGroupID,
		DeviceId: entities.GenerateUUID(),
		RegisterSince: rand.Int63(),
		Labels: labels,
	}
}