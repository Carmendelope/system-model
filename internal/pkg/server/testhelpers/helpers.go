/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package testhelpers

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/entities/device"
	devProvider "github.com/nalej/system-model/internal/pkg/provider/device"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"math/rand"
)

func CreateOrganization(orgProvider orgProvider.Provider) * entities.Organization {
	toAdd := entities.NewOrganization(fmt.Sprintf("org-%d-%d", ginkgo.GinkgoRandomSeed(), rand.Int()))
	err := orgProvider.Add(*toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	return toAdd
}

func CreateDeviceGroup(devProvider devProvider.Provider, organizationID string) *device.DeviceGroup  {
	labels := make(map[string]string, 0)
	toAdd := device.NewDeviceGroup( organizationID, entities.GenerateUUID(), "test device group",labels)
	err := devProvider.AddDeviceGroup(*toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	return toAdd
}

func DeleteGroups(devProvider devProvider.Provider, organizationID string){

	groups, err := devProvider.ListDeviceGroups(organizationID)
	gomega.Expect(err).To(gomega.Succeed())

	for _, group := range groups {
		list, err := devProvider.ListDevice(organizationID, group.DeviceGroupId)
		gomega.Expect(err).To(gomega.Succeed())

		for _, device := range list{
			err = devProvider.RemoveDevice(organizationID, group.DeviceGroupId, device.DeviceId)
			gomega.Expect(err).To(gomega.Succeed())
		}

		err = devProvider.RemoveDeviceGroup(organizationID, group.DeviceGroupId)
		gomega.Expect(err).To(gomega.Succeed())
	}


}