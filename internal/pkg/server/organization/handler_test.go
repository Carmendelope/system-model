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

package organization

import (
	"context"
	"github.com/google/uuid"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/nalej/system-model/internal/pkg/provider/organization_setting"
	"github.com/nalej/system-model/internal/pkg/server/testhelpers"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var _ = ginkgo.Describe("Organization service", func() {
	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client grpc_organization_go.OrganizationsClient

	var orgProvider organization.Provider
	var settingProvider organization_setting.Provider

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		// Register the service
		orgProvider = organization.NewMockupOrganizationProvider()
		settingProvider = organization_setting.NewMockupOrganizationSettingProvider()
		manager := NewManager(orgProvider, settingProvider)
		handler := NewHandler(manager)
		grpc_organization_go.RegisterOrganizationsServer(server, handler)

		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_organization_go.NewOrganizationsClient(conn)
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
	})

	ginkgo.BeforeEach(func() {
		orgProvider.Clear()
	})

	ginkgo.Context("adding organization", func() {
		ginkgo.It("should support adding a new organization", func() {
			toAdd := testhelpers.CreateAddOrganizationRequest()
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())
			gomega.Expect(org.Name).To(gomega.Equal(toAdd.Name))
			gomega.Expect(org.OrganizationId).ShouldNot(gomega.BeNil())

			retrieved, err := client.GetOrganization(context.Background(), &grpc_organization_go.OrganizationId{
				OrganizationId: org.OrganizationId})
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(*retrieved).Should(gomega.Equal(*org))

		})

		ginkgo.It("should fail if the organization name is not specified", func() {
			toAdd := &grpc_organization_go.AddOrganizationRequest{}
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(org).Should(gomega.BeNil())
		})

		ginkgo.It("should fail if the organization name already exists", func() {
			toAdd := testhelpers.CreateAddOrganizationRequest()
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			//sameNameOrg := createOrganization("org_test")
			_, err = client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})

	})

	ginkgo.Context("retrieve organization", func() {
		ginkgo.It("should work on existing organization", func() {
			toAdd := testhelpers.CreateAddOrganizationRequest()
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			toGet := grpc_organization_go.OrganizationId{
				OrganizationId: org.OrganizationId,
			}
			retrieved, err := client.GetOrganization(context.Background(), &toGet)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			gomega.Expect(retrieved).Should(gomega.Equal(org))
		})

		ginkgo.It("should recover a list of organizations", func() {
			toAdd := testhelpers.CreateAddOrganizationRequest()
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			toAdd = testhelpers.CreateAddOrganizationRequest()
			org, err = client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			retrieved, err := client.ListOrganizations(context.Background(), new(grpc_common_go.Empty))
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			gomega.Expect(retrieved.Organizations).ShouldNot(gomega.BeNil())
			gomega.Expect(retrieved.Organizations).Should(gomega.HaveLen(2))
		})

		ginkgo.It("should recover a list empty", func() {
			retrieved, err := client.ListOrganizations(context.Background(), new(grpc_common_go.Empty))
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			gomega.Expect(retrieved.Organizations).Should(gomega.HaveLen(0))
		})

		ginkgo.It("should fail on none existing organization", func() {
			toGet := grpc_organization_go.OrganizationId{
				OrganizationId: "notFound",
			}
			retrieved, err := client.GetOrganization(context.Background(), &toGet)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(retrieved).Should(gomega.BeNil())
		})

		ginkgo.It("should fail on empty request", func() {
			toGet := grpc_organization_go.OrganizationId{}
			retrieved, err := client.GetOrganization(context.Background(), &toGet)
			gomega.Expect(err).Should(gomega.HaveOccurred())
			gomega.Expect(retrieved).Should(gomega.BeNil())
		})
	})

	ginkgo.Context("update organization", func() {
		ginkgo.It("should support updating an existing organization", func() {
			toAdd := testhelpers.CreateAddOrganizationRequest()
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			toUpdate := testhelpers.CreateUpdateOrganizationRequest(org.OrganizationId, false, "")
			success, err := client.UpdateOrganization(context.Background(), toUpdate)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())

			retrieved, err := client.GetOrganization(context.Background(), &grpc_organization_go.OrganizationId{OrganizationId: org.OrganizationId})
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(retrieved).NotTo(gomega.Equal(org))

		})

		ginkgo.It("should fail on non-existing organization", func() {
			toUpdate := testhelpers.CreateUpdateOrganizationRequest(uuid.New().String(), false, "")
			success, err := client.UpdateOrganization(context.Background(), toUpdate)
			gomega.Expect(err).ShouldNot(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())
		})

		ginkgo.It("should fail when removing the name of an organization", func() {
			toAdd := testhelpers.CreateAddOrganizationRequest()
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			toUpdate := testhelpers.CreateUpdateOrganizationRequest(org.OrganizationId, true, "")
			success, err := client.UpdateOrganization(context.Background(), toUpdate)
			gomega.Expect(err).ShouldNot(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())
		})
		ginkgo.It("Should fail when updating the name of an organization if there is another with that name", func() {
			toAdd1 := testhelpers.CreateAddOrganizationRequest()
			org1, err := client.AddOrganization(context.Background(), toAdd1)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org1).ShouldNot(gomega.BeNil())

			toAdd2 := testhelpers.CreateAddOrganizationRequest()
			org2, err := client.AddOrganization(context.Background(), toAdd2)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org2).ShouldNot(gomega.BeNil())

			toUpdate := testhelpers.CreateUpdateOrganizationRequest(org1.OrganizationId, true, org2.Name)
			success, err := client.UpdateOrganization(context.Background(), toUpdate)
			gomega.Expect(err).ShouldNot(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())

		})
	})

	ginkgo.PContext("remove organization", func() {
		ginkgo.PIt("should support removing an existing organization", func() {

		})
		ginkgo.PIt("should fail on non-existing organization", func() {

		})
	})

	ginkgo.Context("adding setting", func() {
		ginkgo.It("Should be able to add a setting", func() {
			// add Organization
			toAdd := testhelpers.CreateAddOrganizationRequest()
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			settingToAdd := testhelpers.CreateAddSettingRequest(org.OrganizationId)
			setting, err := client.AddSetting(context.Background(), settingToAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(setting).NotTo(gomega.BeNil())

		})
		ginkgo.It("Should not be able to add a setting if the organization does not exist", func() {
			settingToAdd := testhelpers.CreateAddSettingRequest("organization-test")
			_, err := client.AddSetting(context.Background(), settingToAdd)
			gomega.Expect(err).ShouldNot(gomega.Succeed())

		})
	})
	ginkgo.Context("removing setting", func() {

		ginkgo.It("Should be able to remove a setting", func() {
			// add Organization
			toAdd := testhelpers.CreateAddOrganizationRequest()
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			settingToAdd := testhelpers.CreateAddSettingRequest(org.OrganizationId)
			setting, err := client.AddSetting(context.Background(), settingToAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(setting).NotTo(gomega.BeNil())

			success, err := client.RemoveSetting(context.Background(), &grpc_organization_go.SettingKey{
				OrganizationId: setting.OrganizationId,
				Key:            setting.Key,
			})
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(success).ShouldNot(gomega.BeNil())

		})
		ginkgo.It("Should not be able to remove a non existing setting", func() {
			// add Organization
			toAdd := testhelpers.CreateAddOrganizationRequest()
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			success, err := client.RemoveSetting(context.Background(), &grpc_organization_go.SettingKey{
				OrganizationId: org.OrganizationId,
				Key:            "key",
			})
			gomega.Expect(err).ShouldNot(gomega.Succeed())
			gomega.Expect(success).Should(gomega.BeNil())

		})
	})
	ginkgo.Context("updating setting", func() {
		ginkgo.It("Should be able to update a setting", func() {
			// add Organization
			toAdd := testhelpers.CreateAddOrganizationRequest()
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			settingToAdd := testhelpers.CreateAddSettingRequest(org.OrganizationId)
			setting, err := client.AddSetting(context.Background(), settingToAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(setting).NotTo(gomega.BeNil())

			toUpdate := testhelpers.CreateUpdateSettingRequest(setting)
			success, err := client.UpdateSetting(context.Background(), toUpdate)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())

		})
		ginkgo.It("Should not be able to update a non existing setting", func() {
			// add Organization
			toAdd := testhelpers.CreateAddOrganizationRequest()
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			settingToAdd := testhelpers.CreateAddSettingRequest(org.OrganizationId)
			setting, err := client.AddSetting(context.Background(), settingToAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(setting).NotTo(gomega.BeNil())

			toUpdate := testhelpers.CreateUpdateSettingRequest(setting)
			toUpdate.Key = "wrong key"
			success, err := client.UpdateSetting(context.Background(), toUpdate)
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())

		})
	})
	ginkgo.Context("listing setting", func() {
		ginkgo.It("Should be able to add a setting", func() {
			// add Organization
			toAdd := testhelpers.CreateAddOrganizationRequest()
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			numSettings := 5
			for i := 0; i < numSettings; i++ {
				settingToAdd := testhelpers.CreateAddSettingRequest(org.OrganizationId)
				setting, err := client.AddSetting(context.Background(), settingToAdd)
				gomega.Expect(err).Should(gomega.Succeed())
				gomega.Expect(setting).NotTo(gomega.BeNil())
			}
			list, err := client.ListSettings(context.Background(), &grpc_organization_go.OrganizationId{
				OrganizationId: org.OrganizationId,
			})
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(list).NotTo(gomega.BeNil())
			gomega.Expect(len(list.Settings)).Should(gomega.Equal(numSettings))

		})
		ginkgo.It("Should be able to return an empty list of settings", func() {
			// add Organization
			toAdd := testhelpers.CreateAddOrganizationRequest()
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			list, err := client.ListSettings(context.Background(), &grpc_organization_go.OrganizationId{
				OrganizationId: org.OrganizationId,
			})
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(list).NotTo(gomega.BeNil())
			gomega.Expect(len(list.Settings)).Should(gomega.Equal(0))

		})
		ginkgo.It("Should not be able to return a list of settings if the organization does not exist", func() {

			_, err := client.ListSettings(context.Background(), &grpc_organization_go.OrganizationId{
				OrganizationId: uuid.New().String(),
			})
			gomega.Expect(err).ShouldNot(gomega.Succeed())
		})
	})
	ginkgo.Context("removing setting", func() {
		ginkgo.It("Should be able to remove a setting", func() {
			// add Organization
			toAdd := testhelpers.CreateAddOrganizationRequest()
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			settingToAdd := testhelpers.CreateAddSettingRequest(org.OrganizationId)
			setting, err := client.AddSetting(context.Background(), settingToAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(setting).NotTo(gomega.BeNil())

			success, err := client.RemoveSetting(context.Background(), &grpc_organization_go.SettingKey{
				OrganizationId: org.OrganizationId,
				Key:            setting.Key,
			})
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())

		})
		ginkgo.It("Should not be able to remove a non existing setting", func() {
			// add Organization
			toAdd := testhelpers.CreateAddOrganizationRequest()
			org, err := client.AddOrganization(context.Background(), toAdd)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(org).ShouldNot(gomega.BeNil())

			success, err := client.RemoveSetting(context.Background(), &grpc_organization_go.SettingKey{
				OrganizationId: org.OrganizationId,
				Key:            uuid.New().String(),
			})
			gomega.Expect(err).ShouldNot(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())

		})
	})

})
