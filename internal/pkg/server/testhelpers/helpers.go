/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package testhelpers

import (
	"fmt"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/entities/devices"
	"github.com/nalej/system-model/internal/pkg/provider/account"
	devProvider "github.com/nalej/system-model/internal/pkg/provider/device"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/nalej/system-model/internal/pkg/provider/user"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"math/rand"
	"time"
)

func CreateOrganization(orgProvider orgProvider.Provider) * entities.Organization {
	toAdd := entities.NewOrganization(fmt.Sprintf("org-%d-%d", ginkgo.GinkgoRandomSeed(), rand.Int()))
	err := orgProvider.Add(*toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	return toAdd
}

func CreateDeviceGroup(devProvider devProvider.Provider, organizationID string, deviceGroupName string) *devices.DeviceGroup {
	labels := make(map[string]string, 0)
	toAdd := devices.NewDeviceGroup( organizationID, entities.GenerateUUID(), deviceGroupName,labels)
	err := devProvider.AddDeviceGroup(*toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	return toAdd
}

func DeleteGroups(devProvider devProvider.Provider, organizationID string){

	groups, err := devProvider.ListDeviceGroups(organizationID)
	gomega.Expect(err).To(gomega.Succeed())

	for _, group := range groups {
		list, err := devProvider.ListDevices(organizationID, group.DeviceGroupId)
		gomega.Expect(err).To(gomega.Succeed())

		for _, device := range list{
			err = devProvider.RemoveDevice(organizationID, group.DeviceGroupId, device.DeviceId)
			gomega.Expect(err).To(gomega.Succeed())
		}

		err = devProvider.RemoveDeviceGroup(organizationID, group.DeviceGroupId)
		gomega.Expect(err).To(gomega.Succeed())
	}
}

// --------------
// -- user helper
// --------------
func CreateNewAddUser()  *grpc_user_go.AddUserRequest{
	return &grpc_user_go.AddUserRequest{
		Email: 		fmt.Sprintf("%s.nalej.com", entities.GenerateUUID()),
		Password: 	"******",
		Name: 		"user test name",
		PhotoUrl: 	"",
		FullName: 	"user test full name",
		Address:    "address",
		Phone: 		map[string]string{"home": "00.000.00.00", "mobile": "606.11.22.33"},
		AltEmail:	fmt.Sprintf("%s.nalej.com", entities.GenerateUUID()),
		CompanyName:"Company name",
		Title: 		"title",
	}
}

func AddUser (provider user.Provider) *entities.User{
	toAdd := entities.User{
		Email: fmt.Sprintf("%s.nalej.com", entities.GenerateUUID()),
		Name: 		"user test name",
		PhotoUrl: 	"",
		MemberSince: time.Now().Unix(),
		ContactInfo: &entities.UserContactInfo{
			FullName: "user test full name",
		},
	}
	err := provider.Add(toAdd)
	gomega.Expect(err).To(gomega.Succeed())

	return &toAdd
}

// -----------------
// -- account helper
// -----------------

func AddAccount (provider account.Provider) * entities.Account {
	toAdd := entities.Account{
		AccountId: entities.GenerateUUID(),
		Name: "test account name",
		Created: time.Now().Unix(),
		State: entities.AccountState_Active,
		StateInfo: "Active for test",
	}
	err := provider.Add(toAdd)
	gomega.Expect(err).To(gomega.Succeed())

	return &toAdd
}

