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
 *
 */

package user

import (
	"context"
	"fmt"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/nalej/system-model/internal/pkg/entities"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	uProvider "github.com/nalej/system-model/internal/pkg/provider/user"
	"github.com/nalej/system-model/internal/pkg/server/testhelpers"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func createAddUserRequest(organizationID string, email string) *grpc_user_go.AddUserRequest {
	return &grpc_user_go.AddUserRequest{
		OrganizationId: organizationID,
		Email:          email,
		Password:       "testPassword",
		Name:           "test user",
	}
}

var _ = ginkgo.Describe("User service", func() {

	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client grpc_user_go.UsersClient

	// Providers
	var organizationProvider orgProvider.Provider
	var userProvider uProvider.Provider

	// Target organization.
	var targetOrganization *entities.Organization

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		organizationProvider = orgProvider.NewMockupOrganizationProvider()
		userProvider = uProvider.NewMockupUserProvider()

		// Register the service
		manager := NewManager(organizationProvider, userProvider)
		handler := NewHandler(manager)
		grpc_user_go.RegisterUsersServer(server, handler)

		test.LaunchServer(server, listener)


		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_user_go.NewUsersClient(conn)
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
	})

	ginkgo.BeforeEach(func() {
		ginkgo.By("cleaning the mockups", func() {
			organizationProvider.(*orgProvider.MockupOrganizationProvider).Clear()
			userProvider.(*uProvider.MockupUserProvider).Clear()
			// Initial data
			targetOrganization = testhelpers.CreateOrganization(organizationProvider)
		})
	})

	ginkgo.It("should be able to add a new user", func() {
		toAdd := createAddUserRequest(targetOrganization.ID, "email@email.com")
		added, err := client.AddUser(context.Background(), toAdd)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(added).ShouldNot(gomega.BeNil())
		gomega.Expect(added.Email).Should(gomega.Equal(toAdd.Email))
	})

	ginkgo.It("should be able to retrieve an existing user", func() {
		toAdd := createAddUserRequest(targetOrganization.ID, "email@email.com")
		added, err := client.AddUser(context.Background(), toAdd)
		gomega.Expect(err).To(gomega.Succeed())

		userID := &grpc_user_go.UserId{
			OrganizationId: targetOrganization.ID,
			Email:          added.Email,
		}
		retrieved, err := client.GetUser(context.Background(), userID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
		gomega.Expect(retrieved.Email).Should(gomega.Equal(added.Email))
		gomega.Expect(retrieved.Name).Should(gomega.Equal(added.Name))
	})

	ginkgo.It("should be able to retrieve the list of users", func() {
		numUsers := 10
		for i := 0; i < numUsers; i++ {
			email := fmt.Sprintf("email%d@email.com", i)
			toAdd := createAddUserRequest(targetOrganization.ID, email)
			_, err := client.AddUser(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
		}
		organizationID := &grpc_organization_go.OrganizationId{
			OrganizationId: targetOrganization.ID,
		}
		users, err := client.GetUsers(context.Background(), organizationID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(users).ToNot(gomega.BeNil())
		gomega.Expect(len(users.Users)).Should(gomega.Equal(numUsers))
	})

	ginkgo.It("should be able to remove an existing user", func() {
		toAdd := createAddUserRequest(targetOrganization.ID, "email@email.com")
		added, err := client.AddUser(context.Background(), toAdd)
		gomega.Expect(err).To(gomega.Succeed())

		removeRequest := &grpc_user_go.RemoveUserRequest{
			OrganizationId: added.OrganizationId,
			Email:          added.Email,
		}

		success, err := client.RemoveUser(context.Background(), removeRequest)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(success).ToNot(gomega.BeNil())
	})

	ginkgo.It("should be able to retrieve an existing user", func() {
		toAdd := createAddUserRequest(targetOrganization.ID, "email@email.com")
		added, err := client.AddUser(context.Background(), toAdd)
		gomega.Expect(err).To(gomega.Succeed())

		updateReq := &grpc_user_go.UpdateUserRequest{
			OrganizationId: targetOrganization.ID,
			Email:          added.Email,
			Name:           "newNameUpdate",
		}
		_, err = client.Update(context.Background(), updateReq)
		gomega.Expect(err).To(gomega.Succeed())

		userID := &grpc_user_go.UserId{
			OrganizationId: targetOrganization.ID,
			Email:          added.Email,
		}
		retrieved, err := client.GetUser(context.Background(), userID)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
		gomega.Expect(retrieved.Email).Should(gomega.Equal(added.Email))
		gomega.Expect(retrieved.Name).Should(gomega.Equal("newNameUpdate"))
	})

})
