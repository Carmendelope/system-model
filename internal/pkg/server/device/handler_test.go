package device

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/nalej/grpc-device-go"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/nalej/system-model/internal/pkg/entities"
	devEntities "github.com/nalej/system-model/internal/pkg/entities/device"
	"github.com/nalej/system-model/internal/pkg/provider/device"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/nalej/system-model/internal/pkg/server/testhelpers"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"math/rand"
)

func GenerateAddDeviceGroup(organizationID string) * grpc_device_go.AddDeviceGroupRequest  {

	labels := make(map[string]string, 0)
	for i := 0; i < rand.Intn(5); i++ {
		labels[fmt.Sprintf("eti_%d", i)] = fmt.Sprintf("Label_%d", i)
	}


	return &grpc_device_go.AddDeviceGroupRequest{
		RequestId: "request_id",
		OrganizationId: organizationID,
		Name: fmt.Sprintf("organization_test-%s", uuid.New().String()),
		Labels:labels,
	}
}

func CreateAssetInfo () *grpc_inventory_go.AssetInfo {
	return &grpc_inventory_go.AssetInfo{
		Os: & grpc_inventory_go.OperatingSystemInfo{
			Name: 		"Linux ubuntu",
			Version: 	"3.0.0",
			Class: 		grpc_inventory_go.OperatingSystemClass_WINDOWS,
			Architecture: "arch",
		},
		Hardware: &grpc_inventory_go.HardwareInfo{
			Cpus: []* grpc_inventory_go.CPUInfo {
				{
					Manufacturer: 	"man_1",
					Model: 			"model1",
					Architecture:   "arch_1",
					NumCores:       2,
				},
			},
			InstalledRam: int64(2000),
			NetInterfaces: []*grpc_inventory_go.NetworkingHardwareInfo {
				{
					Type: "type",
					LinkCapacity: int64(8000),
				},
			},
		},
		Storage: []*grpc_inventory_go.StorageHardwareInfo{
			{
				Type:          "shi_type",
				TotalCapacity: int64(25000),
			},
		},
	}
}

func GenerateAddDevice(organizationID string, deviceGroupID string) * grpc_device_go.AddDeviceRequest  {

	labels := make(map[string]string, 0)

	for i := 0; i < rand.Intn(5); i++ {
		labels[fmt.Sprintf("eti_%d", i)] = fmt.Sprintf("Label_%d", i)
	}


	return &grpc_device_go.AddDeviceRequest{
		OrganizationId: organizationID,
		DeviceGroupId: 	deviceGroupID,
		DeviceId: 		entities.GenerateUUID(),
		Labels: 		labels,
		AssetInfo:      CreateAssetInfo(),
	}
}

var _ = ginkgo.Describe("Applications", func(){

	// gRPC server
	var server * grpc.Server
	// grpc test listener
	var listener * bufconn.Listener
	// client
	var client grpc_device_go.DevicesClient

	// Target organization.
	var targetOrganization * entities.Organization
	var targetDeviceGroup * devEntities.DeviceGroup

	// Organization Provider
	var organizationProvider organization.Provider
	var deviceProvider device.Provider

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()


		// Create providers
		organizationProvider = organization.NewMockupOrganizationProvider()
		deviceProvider = device.NewMockupDeviceProvider()

		manager := NewManager(deviceProvider, organizationProvider)
		handler := NewHandler(manager)
		grpc_device_go.RegisterDevicesServer(server, handler)

		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_device_go.NewDevicesClient(conn)
	})

	ginkgo.AfterSuite(func(){
		server.Stop()
		listener.Close()
	})

	ginkgo.BeforeEach(func(){
		ginkgo.By("cleaning the mockups", func(){
			targetOrganization = testhelpers.CreateOrganization(organizationProvider)
		})
	})
	ginkgo.AfterEach(func() {
		testhelpers.DeleteGroups(deviceProvider,targetOrganization.ID)
	})

	ginkgo.Context("Device Group", func(){
		ginkgo.Context("adding device group", func(){
			ginkgo.It("should add an device group", func(){
				toAdd := GenerateAddDeviceGroup(targetOrganization.ID)
				group, err := client.AddDeviceGroup(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(group).ShouldNot(gomega.BeNil())
				gomega.Expect(group.DeviceGroupId).ShouldNot(gomega.BeNil())
				gomega.Expect(group.Name).Should(gomega.Equal(toAdd.Name))
			})
			ginkgo.It("should fail on an empty request", func(){
				toAdd := &grpc_device_go.AddDeviceGroupRequest{}
				group, err := client.AddDeviceGroup(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(group).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a non existing organization", func(){
				toAdd := GenerateAddDeviceGroup(targetOrganization.ID)
				toAdd.OrganizationId = "does not exists"
				group, err := client.AddDeviceGroup(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(group).Should(gomega.BeNil())
			})
		})
		ginkgo.Context("get device group", func(){
			ginkgo.It("should get an existing device group", func(){
				toAdd := GenerateAddDeviceGroup(targetOrganization.ID)
				group, err := client.AddDeviceGroup(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(group).ShouldNot(gomega.BeNil())
				retrieved, err := client.GetDeviceGroup(context.Background(), &grpc_device_go.DeviceGroupId{
					OrganizationId: group.OrganizationId,
					DeviceGroupId: group.DeviceGroupId,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(retrieved.Name).Should(gomega.Equal(group.Name))
			})
			ginkgo.It("should fail on a non existing device group", func(){
				retrieved, err := client.GetDeviceGroup(context.Background(), &grpc_device_go.DeviceGroupId{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a non existing organization", func(){
				retrieved, err := client.GetDeviceGroup(context.Background(), &grpc_device_go.DeviceGroupId{
					OrganizationId: "does not exists",
					DeviceGroupId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
		})
		ginkgo.Context("listing device groups", func(){
			ginkgo.It("should device groups on an existing organization", func(){
				numGroups := 3
				for i := 0; i < numGroups; i ++ {
					toAdd := GenerateAddDeviceGroup(targetOrganization.ID)
					group, err := client.AddDeviceGroup(context.Background(), toAdd)
					gomega.Expect(err).Should(gomega.Succeed())
					gomega.Expect(group).ShouldNot(gomega.BeNil())
				}
				retrieved, err := client.ListDeviceGroups(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: targetOrganization.ID,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				//gomega.Expect(len(retrieved.Groups)).Should(gomega.Equal(numGroups))
			})
			ginkgo.It("should fail on a non existing organization", func(){
				retrieved, err := client.ListDeviceGroups(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
			ginkgo.It("should work on an organization without groups", func(){

				retrieved, err := client.ListDeviceGroups(context.Background(), &grpc_organization_go.OrganizationId{
					OrganizationId: targetOrganization.ID,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(len(retrieved.Groups)).Should(gomega.Equal(0))
			})
		})
		ginkgo.Context("removing device group", func() {
			ginkgo.It("Should remove a group", func(){
				toAdd := GenerateAddDeviceGroup(targetOrganization.ID)
				group, err := client.AddDeviceGroup(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())

				// remove group
				removed, err := client.RemoveDeviceGroup(context.Background(), &grpc_device_go.RemoveDeviceGroupRequest{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: group.DeviceGroupId,
				})
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(removed).ShouldNot(gomega.BeNil())

			})
			ginkgo.It("Should not be able to remove a group on a non existing organization", func(){

				// remove group
				removed, err := client.RemoveDeviceGroup(context.Background(), &grpc_device_go.RemoveDeviceGroupRequest{
					OrganizationId: "does not exists",
					DeviceGroupId: "device_id",
				})
				gomega.Expect(err).NotTo(gomega.Succeed())
				gomega.Expect(removed).Should(gomega.BeNil())

			})
			ginkgo.It("Should not be able to remove a non existing group", func(){

				// remove group
				removed, err := client.RemoveDeviceGroup(context.Background(), &grpc_device_go.RemoveDeviceGroupRequest{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: "does not exists",
				})
				gomega.Expect(err).NotTo(gomega.Succeed())
				gomega.Expect(removed).Should(gomega.BeNil())

			})
		})
	})
	ginkgo.Context("Devices", func() {
		ginkgo.BeforeEach(func() {
			targetDeviceGroup = testhelpers.CreateDeviceGroup(deviceProvider, targetOrganization.ID, "dgName")
		})
		ginkgo.Context("adding device", func() {
			ginkgo.It("should add an device group", func() {
				toAdd := GenerateAddDevice(targetOrganization.ID, targetDeviceGroup.DeviceGroupId)
				group, err := client.AddDevice(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(group).ShouldNot(gomega.BeNil())
				gomega.Expect(group.DeviceId).ShouldNot(gomega.BeNil())
			})
			ginkgo.It("should fail on an empty request", func() {
				toAdd := &grpc_device_go.AddDeviceRequest{}
				group, err := client.AddDevice(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(group).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a non existing organization", func() {
				toAdd := GenerateAddDevice(targetOrganization.ID, "does not exists")
				toAdd.OrganizationId = "does not exists"
				group, err := client.AddDevice(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(group).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a non existing group", func() {
				toAdd := GenerateAddDevice(targetOrganization.ID, "does not exists")
				group, err := client.AddDevice(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(group).Should(gomega.BeNil())
			})
		})
		ginkgo.Context("get device", func(){
			ginkgo.It("should get an existing device", func(){
				toAdd := GenerateAddDevice(targetOrganization.ID, targetDeviceGroup.DeviceGroupId)
				device, err := client.AddDevice(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(device).ShouldNot(gomega.BeNil())
				retrieved, err := client.GetDevice(context.Background(), &grpc_device_go.DeviceId{
					OrganizationId: device.OrganizationId,
					DeviceGroupId: device.DeviceGroupId,
					DeviceId: device.DeviceId,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			})
			ginkgo.It("should fail on a non existing device", func(){
				retrieved, err := client.GetDevice(context.Background(), &grpc_device_go.DeviceId{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: targetDeviceGroup.DeviceGroupId,
					DeviceId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a non existing organization", func(){
				retrieved, err := client.GetDevice(context.Background(), &grpc_device_go.DeviceId{
					OrganizationId: "does not exists",
					DeviceGroupId: "does not exists",
					DeviceId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a non existing group", func(){
				retrieved, err := client.GetDevice(context.Background(), &grpc_device_go.DeviceId{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: "does not exists",
					DeviceId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
		})
		ginkgo.Context("listing devices", func(){
			ginkgo.It("should devices on an existing group", func(){
				numGroups := rand.Intn(4) +1
				for i := 0; i < numGroups; i ++ {
					toAdd := GenerateAddDevice(targetOrganization.ID, targetDeviceGroup.DeviceGroupId)
					group, err := client.AddDevice(context.Background(), toAdd)
					gomega.Expect(err).Should(gomega.Succeed())
					gomega.Expect(group).ShouldNot(gomega.BeNil())
				}
				retrieved, err := client.ListDevices(context.Background(), &grpc_device_go.DeviceGroupId{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: targetDeviceGroup.DeviceGroupId,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(len(retrieved.Devices)).Should(gomega.Equal(numGroups))
			})
			ginkgo.It("should fail on a non existing organization", func(){
				retrieved, err := client.ListDevices(context.Background(), &grpc_device_go.DeviceGroupId{
					OrganizationId: "does not exists",
					DeviceGroupId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
			ginkgo.It("should fail on a non existing group", func(){
				retrieved, err := client.ListDevices(context.Background(), &grpc_device_go.DeviceGroupId{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: "does not exists",
				})
				gomega.Expect(err).Should(gomega.HaveOccurred())
				gomega.Expect(retrieved).Should(gomega.BeNil())
			})
			ginkgo.It("should work on an organization without devices", func(){

				retrieved, err := client.ListDevices(context.Background(), &grpc_device_go.DeviceGroupId{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: targetDeviceGroup.DeviceGroupId,
				})
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(len(retrieved.Devices)).Should(gomega.Equal(0))
			})
		})
		ginkgo.Context("removing device", func() {
			ginkgo.It("Should remove a device", func(){
				toAdd := GenerateAddDevice(targetOrganization.ID, targetDeviceGroup.DeviceGroupId)
				group, err := client.AddDevice(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())

				// remove group
				removed, err := client.RemoveDevice(context.Background(), &grpc_device_go.RemoveDeviceRequest{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: group.DeviceGroupId,
					DeviceId: toAdd.DeviceId,
				})
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(removed).ShouldNot(gomega.BeNil())

			})
			ginkgo.It("Should not be able to remove a device on a non existing organization", func(){

				// remove group
				removed, err := client.RemoveDevice(context.Background(), &grpc_device_go.RemoveDeviceRequest{
					OrganizationId: "does not exists",
					DeviceGroupId: "does not exists",
					DeviceId: "does not exists",
				})
				gomega.Expect(err).NotTo(gomega.Succeed())
				gomega.Expect(removed).Should(gomega.BeNil())

			})
			ginkgo.It("Should not be able to remove a device on a non existing group", func(){

				// remove group
				removed, err := client.RemoveDevice(context.Background(), &grpc_device_go.RemoveDeviceRequest{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: "does not exists",
				})
				gomega.Expect(err).NotTo(gomega.Succeed())
				gomega.Expect(removed).Should(gomega.BeNil())

			})
			ginkgo.It("Should not be able to remove a non existing device", func(){

				// remove group
				removed, err := client.RemoveDevice(context.Background(), &grpc_device_go.RemoveDeviceRequest{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: targetDeviceGroup.DeviceGroupId,
				})
				gomega.Expect(err).NotTo(gomega.Succeed())
				gomega.Expect(removed).Should(gomega.BeNil())

			})
		})
		ginkgo.Context("update device group", func() {
			ginkgo.It("Should update a device", func(){
				toAdd := GenerateAddDevice(targetOrganization.ID, targetDeviceGroup.DeviceGroupId)
				toAdd.Labels = nil
				group, err := client.AddDevice(context.Background(), toAdd)
				gomega.Expect(err).Should(gomega.Succeed())

				// update device (add label)
				updated, err := client.UpdateDevice(context.Background(), &grpc_device_go.UpdateDeviceRequest{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: group.DeviceGroupId,
					DeviceId: toAdd.DeviceId,
					AddLabels: true,
					RemoveLabels: false,
					Labels: map[string]string{"label1":"value1"},

				})
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(updated).ShouldNot(gomega.BeNil())

				// get the device to check the updated works
				retrieved, err := client.GetDevice(context.Background(), &grpc_device_go.DeviceId{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: group.DeviceGroupId,
					DeviceId: toAdd.DeviceId,
				})
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(len(retrieved.Labels)).Should(gomega.Equal(1))

				// removeLabel
				updated, err = client.UpdateDevice(context.Background(), &grpc_device_go.UpdateDeviceRequest{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: group.DeviceGroupId,
					DeviceId: toAdd.DeviceId,
					AddLabels: false,
					RemoveLabels: true,
					Labels: map[string]string{"label1":"value1"},

				})
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(updated).ShouldNot(gomega.BeNil())

				// get the device to check the updated works
				retrieved, err = client.GetDevice(context.Background(), &grpc_device_go.DeviceId{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: group.DeviceGroupId,
					DeviceId: toAdd.DeviceId,
				})
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
				gomega.Expect(retrieved.Labels).Should(gomega.BeNil())


			})
			ginkgo.It("Should not be able to update a device on a non existing organization", func(){

				_, err := client.UpdateDevice(context.Background(), &grpc_device_go.UpdateDeviceRequest{
					OrganizationId: targetOrganization.ID,
					DeviceGroupId: uuid.New().String(),
					DeviceId: uuid.New().String(),
					AddLabels: false,
					RemoveLabels: true,
					Labels: map[string]string{"label1":"value1"},

				})
				gomega.Expect(err).NotTo(gomega.Succeed())

			})
		})
	})

})
