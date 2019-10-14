package device

import (
	"github.com/google/uuid"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {

	ginkgo.AfterEach(func() {
		provider.Clear()
	})

	ginkgo.Context("device group tests", func() {
		ginkgo.It("Should be able to add device group", func() {

			toAdd := NewDeviceTestHepler().CreateDeviceGroup()

			err := provider.AddDeviceGroup(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

		})
		ginkgo.It("Should not be able to add the same device group twice", func() {
			toAdd := NewDeviceTestHepler().CreateDeviceGroup()

			err := provider.AddDeviceGroup(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.AddDeviceGroup(*toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})

		ginkgo.It("Should be able to delete a device group", func() {
			toAdd := NewDeviceTestHepler().CreateDeviceGroup()
			err := provider.AddDeviceGroup(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.RemoveDeviceGroup(toAdd.OrganizationId, toAdd.DeviceGroupId)
			gomega.Expect(err).To(gomega.Succeed())

		})
		ginkgo.It("Should not be able to delete a device group", func() {
			toAdd := NewDeviceTestHepler().CreateDeviceGroup()

			err := provider.RemoveDeviceGroup(toAdd.OrganizationId, toAdd.DeviceGroupId)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})

		ginkgo.It("Should be able to find a device group", func() {
			toAdd := NewDeviceTestHepler().CreateDeviceGroup()
			err := provider.AddDeviceGroup(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			exists, err := provider.ExistsDeviceGroup(toAdd.OrganizationId, toAdd.DeviceGroupId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())

		})
		ginkgo.It("Should not be able to find a device group", func() {
			toAdd := NewDeviceTestHepler().CreateDeviceGroup()

			exists, err := provider.ExistsDeviceGroup(toAdd.OrganizationId, toAdd.DeviceGroupId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())

		})
		ginkgo.It("Should be able to find a device group by name", func() {
			toAdd := NewDeviceTestHepler().CreateDeviceGroup()
			err := provider.AddDeviceGroup(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			exists, err := provider.ExistsDeviceGroupByName(toAdd.OrganizationId, toAdd.Name)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())

		})
		ginkgo.It("Should not be able to find a device group by name", func() {
			toAdd := NewDeviceTestHepler().CreateDeviceGroup()

			exists, err := provider.ExistsDeviceGroup(toAdd.OrganizationId, toAdd.Name)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())

		})
		ginkgo.It("Should be able to list a device groups", func() {
			helper := NewDeviceTestHepler()

			toAdd := helper.CreateDeviceGroup()
			err := provider.AddDeviceGroup(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			toAdd = helper.CreateOrganizationDeviceGroup(toAdd.OrganizationId)
			err = provider.AddDeviceGroup(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			toAdd = helper.CreateOrganizationDeviceGroup(toAdd.OrganizationId)
			err = provider.AddDeviceGroup(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			list, err := provider.ListDeviceGroups(toAdd.OrganizationId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list)).To(gomega.Equal(3))

		})
		ginkgo.It("Should not be able to list a device groups", func() {
			helper := NewDeviceTestHepler()

			toAdd := helper.CreateDeviceGroup()

			list, err := provider.ListDeviceGroups(toAdd.OrganizationId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(list).To(gomega.BeEmpty())
		})

		ginkgo.It("should be able to get a device group", func() {
			helper := NewDeviceTestHepler()

			toAdd := helper.CreateDeviceGroup()
			err := provider.AddDeviceGroup(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			received, err := provider.GetDeviceGroup(toAdd.OrganizationId, toAdd.DeviceGroupId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(received.Name).To(gomega.Equal(toAdd.Name))

		})

		ginkgo.It("should not be able to get a device group", func() {
			helper := NewDeviceTestHepler()

			toAdd := helper.CreateDeviceGroup()

			_, err := provider.GetDeviceGroup(toAdd.OrganizationId, toAdd.DeviceGroupId)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		ginkgo.It("should be able to get devices groups by name", func() {
			names := make([]string, 0)
			helper := NewDeviceTestHepler()
			organizationID := uuid.New().String()

			for i := 0; i < 5; i++ {
				toAdd := helper.CreateOrganizationDeviceGroup(organizationID)
				err := provider.AddDeviceGroup(*toAdd)
				gomega.Expect(err).To(gomega.Succeed())
				names = append(names, toAdd.Name)
			}

			deviceGroups, err := provider.GetDeviceGroupsByName(organizationID, names)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(deviceGroups).NotTo(gomega.BeEmpty())
			gomega.Expect(len(deviceGroups)).Should(gomega.Equal(5))

		})
		ginkgo.It("should be able to get an empty devices groups by name", func() {
			names := make([]string, 0)
			organizationID := uuid.New().String()

			deviceGroups, err := provider.GetDeviceGroupsByName(organizationID, names)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(deviceGroups).To(gomega.BeEmpty())

		})

	})
	ginkgo.Context("Device tests", func() {
		ginkgo.It("Should be able to add a device", func() {
			toAdd := NewDeviceTestHepler().CreateDevice()

			err := provider.AddDevice(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

		})
		ginkgo.It("Should not be able to add a device twice", func() {
			toAdd := NewDeviceTestHepler().CreateDevice()

			err := provider.AddDevice(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.AddDevice(*toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("Should be able to find a device", func() {
			toAdd := NewDeviceTestHepler().CreateDevice()

			err := provider.AddDevice(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			exists, err := provider.ExistsDevice(toAdd.OrganizationId, toAdd.DeviceGroupId, toAdd.DeviceId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())
		})
		ginkgo.It("Should not be able to find a device", func() {
			toAdd := NewDeviceTestHepler().CreateDevice()

			exists, err := provider.ExistsDevice(toAdd.OrganizationId, toAdd.DeviceGroupId, toAdd.DeviceId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())

		})
		ginkgo.It("Should be able to get device info", func() {
			toAdd := NewDeviceTestHepler().CreateDevice()

			err := provider.AddDevice(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			dev, err := provider.GetDevice(toAdd.OrganizationId, toAdd.DeviceGroupId, toAdd.DeviceId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(dev.RegisterSince).To(gomega.Equal(toAdd.RegisterSince))
		})
		ginkgo.It("Should not be able to get device info", func() {
			toAdd := NewDeviceTestHepler().CreateDevice()

			_, err := provider.GetDevice(toAdd.OrganizationId, toAdd.DeviceGroupId, toAdd.DeviceId)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("Should be able to get the devices of a group", func() {
			helper := NewDeviceTestHepler()

			toAdd := helper.CreateDevice()
			err := provider.AddDevice(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			toAdd = helper.CreateGroupDevices(toAdd.OrganizationId, toAdd.DeviceGroupId)
			err = provider.AddDevice(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			toAdd = helper.CreateGroupDevices(toAdd.OrganizationId, toAdd.DeviceGroupId)
			err = provider.AddDevice(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			list, err := provider.ListDevices(toAdd.OrganizationId, toAdd.DeviceGroupId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(list).To(gomega.HaveLen(3))

		})
		ginkgo.It("Should be able to get empty list of devices of a group ", func() {
			helper := NewDeviceTestHepler()

			toAdd := helper.CreateDevice()
			list, err := provider.ListDevices(toAdd.OrganizationId, toAdd.DeviceGroupId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list)).Should(gomega.Equal(0))

		})
		ginkgo.It("Should be able to remove a device", func() {

			toAdd := NewDeviceTestHepler().CreateDevice()
			err := provider.AddDevice(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.RemoveDevice(toAdd.OrganizationId, toAdd.DeviceGroupId, toAdd.DeviceId)
			gomega.Expect(err).To(gomega.Succeed())

		})
		ginkgo.It("Should not be able to remove a device", func() {

			toAdd := NewDeviceTestHepler().CreateDevice()
			err := provider.RemoveDevice(toAdd.OrganizationId, toAdd.DeviceGroupId, toAdd.DeviceId)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})

		ginkgo.It("Should be able to update a device removing all labels", func() {
			toAdd := NewDeviceTestHepler().CreateDevice()

			err := provider.AddDevice(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			// remove labels
			toAdd.Labels = nil
			err = provider.UpdateDevice(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			// check the update
			retrieve, err := provider.GetDevice(toAdd.OrganizationId, toAdd.DeviceGroupId, toAdd.DeviceId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieve).NotTo(gomega.BeNil())
			gomega.Expect(retrieve.Labels).To(gomega.BeNil())

		})

		ginkgo.It("Should be able to update a device adding labels", func() {
			toAdd := NewDeviceTestHepler().CreateDevice()

			err := provider.AddDevice(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			// remove labels
			toAdd.Labels["label3"] = "value3"
			err = provider.UpdateDevice(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			// check the update
			retrieve, err := provider.GetDevice(toAdd.OrganizationId, toAdd.DeviceGroupId, toAdd.DeviceId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieve).NotTo(gomega.BeNil())
			gomega.Expect(len(retrieve.Labels)).Should(gomega.Equal(3))

		})

		ginkgo.It("Should be able to update a device adding location", func() {
			toAdd := NewDeviceTestHepler().CreateDevice()

			err := provider.AddDevice(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			// remove labels
			toAdd.Location = &entities.InventoryLocation{
				Geolocation: "Spain",
				Geohash:     "SP",
			}
			err = provider.UpdateDevice(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			// check the update
			retrieve, err := provider.GetDevice(toAdd.OrganizationId, toAdd.DeviceGroupId, toAdd.DeviceId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieve).NotTo(gomega.BeNil())
			gomega.Expect(retrieve.Location).ShouldNot(gomega.BeNil())
			gomega.Expect(retrieve.Location.Geolocation).Should(gomega.Equal(toAdd.Location.Geolocation))
			gomega.Expect(retrieve.Location.Geohash).Should(gomega.Equal(toAdd.Location.Geohash))
		})

		ginkgo.It("Should not be able to update a non existing device", func() {
			toAdd := NewDeviceTestHepler().CreateDevice()

			err := provider.UpdateDevice(*toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})

	})
}
