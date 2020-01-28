/*
 * Copyright 2020 Nalej
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

package testhelpers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/entities/devices"
	devProvider "github.com/nalej/system-model/internal/pkg/provider/device"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"math/rand"
)

func CreateOrganization() *entities.Organization {
	return entities.NewOrganization(
		fmt.Sprintf(fmt.Sprintf("org-%d-%d", ginkgo.GinkgoRandomSeed(), rand.Int())),
		"test@email.com",
		"Nalej Test Address",
		"City Test",
		"State Test",
		"U.S.A",
		"XXX",
		"Photo")
}

func CreateAddOrganizationRequest() *grpc_organization_go.AddOrganizationRequest {
	return &grpc_organization_go.AddOrganizationRequest{
		Name:        fmt.Sprintf("Nalej-%s", uuid.New().String()),
		FullAddress: "Address",
		City:        "City",
		State:       "State",
		Country:     "Country",
		ZipCode:     "Zip Code",
		PhotoBase64: "Photo",
	}
}

func CreateUpdateOrganizationRequest(id string, updateName bool, newName string) *grpc_organization_go.UpdateOrganizationRequest {
	return &grpc_organization_go.UpdateOrganizationRequest{
		OrganizationId:    id,
		UpdateName:        updateName,
		Name:              newName,
		UpdateFullAddress: true,
		FullAddress:       "Address modified",
	}
}

func AddOrganization(orgProvider orgProvider.Provider) *entities.Organization {
	toAdd := CreateOrganization()
	err := orgProvider.Add(*toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	return toAdd
}

func CreateDeviceGroup(devProvider devProvider.Provider, organizationID string, deviceGroupName string) *devices.DeviceGroup {
	labels := make(map[string]string, 0)
	toAdd := devices.NewDeviceGroup(organizationID, entities.GenerateUUID(), deviceGroupName, labels)
	err := devProvider.AddDeviceGroup(*toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	return toAdd
}

func DeleteGroups(devProvider devProvider.Provider, organizationID string) {

	groups, err := devProvider.ListDeviceGroups(organizationID)
	gomega.Expect(err).To(gomega.Succeed())

	for _, group := range groups {
		list, err := devProvider.ListDevices(organizationID, group.DeviceGroupId)
		gomega.Expect(err).To(gomega.Succeed())

		for _, device := range list {
			err = devProvider.RemoveDevice(organizationID, group.DeviceGroupId, device.DeviceId)
			gomega.Expect(err).To(gomega.Succeed())
		}

		err = devProvider.RemoveDeviceGroup(organizationID, group.DeviceGroupId)
		gomega.Expect(err).To(gomega.Succeed())
	}

}

func CreateAddSettingRequest(organizationID string) *grpc_organization_go.AddSettingRequest {
	randValue := rand.Int()
	return &grpc_organization_go.AddSettingRequest{
		OrganizationId: organizationID,
		Key:            fmt.Sprintf("Key-%d", randValue),
		Value:          fmt.Sprintf("Value-%d", randValue),
		Description:    "Description",
	}
}

func CreateUpdateSettingRequest(setting *grpc_organization_go.OrganizationSetting) *grpc_organization_go.UpdateSettingRequest {
	return &grpc_organization_go.UpdateSettingRequest{
		OrganizationId:    setting.OrganizationId,
		Key:               setting.Key,
		UpdateValue:       true,
		Value:             "new value",
		UpdateDescription: true,
		Description:       "new description",
	}
}
