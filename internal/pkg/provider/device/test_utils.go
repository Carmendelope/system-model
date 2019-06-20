package device

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/entities/devices"
	"math/rand"
)


type DeviceTestHelper struct {

}

func NewDeviceTestHepler() *DeviceTestHelper {
	return &DeviceTestHelper{}
}

func (d * DeviceTestHelper) CreateDeviceGroup() *devices.DeviceGroup {

	labels := make(map[string]string, 0)
	labels["lab1"] = "etiq_1"
	labels["lab2"] = "etiq_2"

	return &devices.DeviceGroup{
		OrganizationId: entities.GenerateUUID(),
		DeviceGroupId: entities.GenerateUUID(),
		Name: fmt.Sprintf("Test-%d Device Group", rand.Intn(10)),
		Created: rand.Int63(),
		Labels: labels,
	}
}
func (d * DeviceTestHelper) CreateOrganizationDeviceGroup(organizationID string) *devices.DeviceGroup {

	labels := make(map[string]string, 0)
	labels["lab1"] = "etiq_1"
	labels["lab2"] = "etiq_2"

	return &devices.DeviceGroup{
		OrganizationId: organizationID,
		DeviceGroupId: entities.GenerateUUID(),
		Name: fmt.Sprintf("Test-%d Device Group", rand.Intn(10)),
		Created: rand.Int63(),
		Labels: labels,
	}
}

func (d *DeviceTestHelper) CreateOs () *entities.OperatingSystemInfo {
	return &entities.OperatingSystemInfo{
		Name: "Linux ubuntu",
		Version: "3.0.0",
		Architecture: "arch",
		Class: entities.DARWIN,
	}
}
func (d *DeviceTestHelper) CreateHardware () *entities.HardwareInfo {
	return &entities.HardwareInfo{
		// []string{"authService1", "authService2"}
		Cpus: []* entities.CPUInfo {
			{
				Manufacturer: 	"man_1",
				Model: 			"model1",
				Architecture:   "arch_1",
				NumCores:       2,
			},
			{
				Manufacturer: 	"man_2",
				Model: 			"model2",
				Architecture:   "arch_2",
				NumCores:       2,
			},
		},
		InstalledRam: int64(2000),
		NetInterfaces: []*entities.NetworkingHardwareInfo {
			{
				Type: "type",
				LinkCapacity: int64(8000),
			},
		},
	}
}

func (d *DeviceTestHelper) CreateStorage() []*entities.StorageHardwareInfo  {
	return []*entities.StorageHardwareInfo{
		{
			Type: "shi_type",
			TotalCapacity: int64(25000),
		},
	}
}

func (d * DeviceTestHelper) CreateDevice() *devices.Device {

	labels := map[string]string {
		"lab1": "value_1",
		"lab2": "value_2",
	}
	return &devices.Device{
		OrganizationId: entities.GenerateUUID(),
		DeviceGroupId: 	entities.GenerateUUID(),
		DeviceId: 		entities.GenerateUUID(),
		RegisterSince:  rand.Int63(),
		Labels: 		labels,
		Os: 			d.CreateOs(),
		Hardware: 		d.CreateHardware(),
		Storage:		d.CreateStorage(),
	}
}

func (d * DeviceTestHelper) CreateGroupDevices(organizationID string, deviceGroupID string) *devices.Device {

	labels := make(map[string]string, 0)
	labels["lab1"] = "etiq_1"
	labels["lab2"] = "etiq_2"

	return &devices.Device{
		OrganizationId: organizationID,
		DeviceGroupId: deviceGroupID,
		DeviceId: entities.GenerateUUID(),
		RegisterSince: rand.Int63(),
		Labels: labels,
	}
}