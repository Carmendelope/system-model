/*
 * Copyright 2019 Nalej
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

package role

import (
	"context"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-role-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/nalej/system-model/internal/pkg/entities"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	rProvider "github.com/nalej/system-model/internal/pkg/provider/role"
	"github.com/nalej/system-model/internal/pkg/server/testhelpers"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func createAddRoleRequest(organizationID string) *grpc_role_go.AddRoleRequest {
	return &grpc_role_go.AddRoleRequest{
		OrganizationId: organizationID,
		Name:           "name",
		Description:    "description",
		Internal:       false,
	}
}

var _ = ginkgo.Describe("Role service", func() {

	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client grpc_role_go.RolesClient

	// Providers
	var organizationProvider orgProvider.Provider
	var roleProvider rProvider.Provider

	// Target organization.
	var targetOrganization *entities.Organization

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()
		test.LaunchServer(server, listener)

		organizationProvider = orgProvider.NewMockupOrganizationProvider()
		roleProvider = rProvider.NewMockupRoleProvider()

		// Register the service
		manager := NewManager(organizationProvider, roleProvider)
		handler := NewHandler(manager)
		grpc_role_go.RegisterRolesServer(server, handler)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_role_go.NewRolesClient(conn)
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
	})

	ginkgo.BeforeEach(func() {
		ginkgo.By("cleaning the mockups", func() {
			organizationProvider.(*orgProvider.MockupOrganizationProvider).Clear()
			roleProvider.(*rProvider.MockupRoleProvider).Clear()
			// Initial data
			targetOrganization = testhelpers.CreateOrganization(organizationProvider)
		})
	})

	ginkgo.It("should be able to add a new role", func() {
		toAdd := createAddRoleRequest(targetOrganization.ID)
		added, err := client.AddRole(context.Background(), toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(added).ShouldNot(gomega.BeNil())
		gomega.Expect(added.RoleId).ShouldNot(gomega.BeEmpty())
		gomega.Expect(added.OrganizationId).Should(gomega.Equal(toAdd.OrganizationId))
	})

	ginkgo.It("should be able to retrieve an existing role", func() {
		toAdd := createAddRoleRequest(targetOrganization.ID)
		added, err := client.AddRole(context.Background(), toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		roleID := &grpc_role_go.RoleId{
			OrganizationId: added.OrganizationId,
			RoleId:         added.RoleId,
		}
		retrieved, err := client.GetRole(context.Background(), roleID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(retrieved).ToNot(gomega.BeNil())
		gomega.Expect(retrieved.RoleId).Should(gomega.Equal(added.RoleId))
	})

	ginkgo.It("should be able to list the existing roles", func() {
		numRoles := 10
		for i := 0; i < numRoles; i++ {
			toAdd := createAddRoleRequest(targetOrganization.ID)
			_, err := client.AddRole(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
		}
		organizationID := &grpc_organization_go.OrganizationId{
			OrganizationId: targetOrganization.ID,
		}
		roles, err := client.ListRoles(context.Background(), organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(roles).ToNot(gomega.BeNil())
		gomega.Expect(len(roles.Roles)).Should(gomega.Equal(numRoles))
	})

	ginkgo.It("should be able to remove a role", func() {
		toAdd := createAddRoleRequest(targetOrganization.ID)
		added, err := client.AddRole(context.Background(), toAdd)
		gomega.Expect(err).To(gomega.Succeed())

		removeRequest := &grpc_role_go.RemoveRoleRequest{
			OrganizationId: added.OrganizationId,
			RoleId:         added.RoleId,
		}

		success, err := client.RemoveRole(context.Background(), removeRequest)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(success).ToNot(gomega.BeNil())
	})

})
