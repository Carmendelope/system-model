/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package user

import (
	"context"
	"fmt"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/account"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	uProvider "github.com/nalej/system-model/internal/pkg/provider/user"
	"github.com/nalej/system-model/internal/pkg/server/testhelpers"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func createAddUserRequest(organizationID string, email string) * grpc_user_go.AddUserRequest{
	return &grpc_user_go.AddUserRequest{
		OrganizationId:       organizationID,
		Email:                email,
		Password:             "testPassword",
		Name:                 "test user",
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
	var accProvider account.Provider

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		organizationProvider = orgProvider.NewMockupOrganizationProvider()
		userProvider = uProvider.NewMockupUserProvider()
		accProvider = account.NewMockupAccountProvider()

		// Register the service
		manager := NewManager(organizationProvider, userProvider, accProvider)
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

	ginkgo.Context("User", func() {
		// Target organization.
		var targetOrganization * entities.Organization
		ginkgo.BeforeEach(func(){
			ginkgo.By("cleaning the mockups", func(){
				organizationProvider.(*orgProvider.MockupOrganizationProvider).Clear()
				userProvider.(*uProvider.MockupUserProvider).Clear()
				// Initial data
				targetOrganization = testhelpers.CreateOrganization(organizationProvider)
			})
		})

		ginkgo.It("should be able to add a new user", func(){
			toAdd := createAddUserRequest(targetOrganization.ID, "email@email.com")
			added, err := client.AddUser(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())
			gomega.Expect(added.Email).Should(gomega.Equal(toAdd.Email))
		})

		ginkgo.It("should be able to retrieve an existing user", func(){
			toAdd := createAddUserRequest(targetOrganization.ID, "email@email.com")
			added, err := client.AddUser(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			userID := &grpc_user_go.UserId{
				OrganizationId:       targetOrganization.ID,
				Email:                added.Email,
			}
			retrieved, err := client.GetUser(context.Background(), userID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			gomega.Expect(retrieved.Email).Should(gomega.Equal(added.Email))
			gomega.Expect(retrieved.Name).Should(gomega.Equal(added.Name))
		})

		ginkgo.It("should be able to retrieve the list of users", func(){
			numUsers := 10
			for i := 0; i < numUsers; i++ {
				email := fmt.Sprintf("email%d@email.com", i)
				toAdd := createAddUserRequest(targetOrganization.ID, email)
				_, err := client.AddUser(context.Background(), toAdd)
				gomega.Expect(err).To(gomega.Succeed())
			}
			organizationID := &grpc_organization_go.OrganizationId{
				OrganizationId:       targetOrganization.ID,
			}
			users, err := client.GetUsers(context.Background(), organizationID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(users).ToNot(gomega.BeNil())
			gomega.Expect(len(users.Users)).Should(gomega.Equal(numUsers))
		})

		ginkgo.It("should be able to remove an existing user", func(){
			toAdd := createAddUserRequest(targetOrganization.ID, "email@email.com")
			added, err := client.AddUser(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			removeRequest := &grpc_user_go.RemoveUserRequest{
				OrganizationId:       added.OrganizationId,
				Email:                added.Email,
			}

			success, err := client.RemoveUser(context.Background(), removeRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).ToNot(gomega.BeNil())
		})

		ginkgo.It("should be able to retrieve an existing user", func(){
			toAdd := createAddUserRequest(targetOrganization.ID, "email@email.com")
			added, err := client.AddUser(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			updateReq := &grpc_user_go.UpdateUserRequest{
				OrganizationId:       targetOrganization.ID,
				Email:                added.Email,
				Name: "newNameUpdate",
			}
			_,err = client.Update(context.Background(), updateReq)
			gomega.Expect(err).To(gomega.Succeed())

			userID := &grpc_user_go.UserId{
				OrganizationId:       targetOrganization.ID,
				Email:                added.Email,
			}
			retrieved, err := client.GetUser(context.Background(), userID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).ShouldNot(gomega.BeNil())
			gomega.Expect(retrieved.Email).Should(gomega.Equal(added.Email))
			gomega.Expect(retrieved.Name).Should(gomega.Equal("newNameUpdate"))
		})

		ginkgo.It("Should be able to update the contact info ", func(){
			// add User
			toAdd := testhelpers.CreateNewAddUser()
			added, err := client.AddUser(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(added).ShouldNot(gomega.BeNil())

			// update ContactInfo
			updateRequest := &grpc_user_go.UpdateContactInfoRequest{
				Email: toAdd.Email,
				FullName: "Full name updated",
				Phone: map[string]string{"mobile": "666.66.66.66"},
				Title: "title updated",
			}
			success, err := client.UpdateContactInfo(context.Background(), updateRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())

			// check the update works
			retrieved, err := client.GetUser(context.Background(), &grpc_user_go.UserId{
				Email: toAdd.Email,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).NotTo(gomega.BeNil())
			gomega.Expect(retrieved.ContactInfo.FullName).Should(gomega.Equal(updateRequest.FullName))
			gomega.Expect(retrieved.ContactInfo.Title).Should(gomega.Equal(updateRequest.Title))
			gomega.Expect(retrieved.ContactInfo.CompanyName).Should(gomega.Equal(toAdd.CompanyName))
			gomega.Expect(retrieved.Name).Should(gomega.Equal(toAdd.Name))

		})
		ginkgo.It("Should not be able to update the contact info if the user does not exists ", func(){
			// update ContactInfo
			updateRequest := &grpc_user_go.UpdateContactInfoRequest{
				Email: fmt.Sprintf("%s.nalej.com", entities.GenerateUUID()),
				FullName: "Full name updated",
				Phone: map[string]string{"mobile": "666.66.66.66"},
				Title: "title updated",
			}
			success, err := client.UpdateContactInfo(context.Background(), updateRequest)
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())

		})
	})

	ginkgo.Context("AccountUser", func(){
		var user *entities.User
		var acc *entities.Account

		ginkgo.BeforeEach(func(){
			userProvider.(*uProvider.MockupUserProvider).Clear()
			accProvider.(*account.MockupAccountProvider).Clear()

			user = testhelpers.AddUser(userProvider)
			acc = testhelpers.AddAccount(accProvider)

		})
		// Add
		ginkgo.It("should be able to add an accountUser", func() {

			accountUser, err := client.AddAccountUser(context.Background(), &grpc_user_go.AddAccountUserRequest{
				AccountId: acc.AccountId,
				Email: user.Email,
				RoleId: entities.GenerateUUID(),
				Internal: false,
				Status: grpc_user_go.UserStatus_ACTIVE,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(accountUser).NotTo(gomega.BeNil())
		})
		ginkgo.It("should not be able to add an accountUser without Role", func() {

			_, err := client.AddAccountUser(context.Background(), &grpc_user_go.AddAccountUserRequest{
				AccountId: acc.AccountId,
				Email: user.Email,
				RoleId: "",
				Internal: false,
				Status: grpc_user_go.UserStatus_ACTIVE,
			})
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("should not be able to add an accountUser if the user does not exit", func() {

			_, err := client.AddAccountUser(context.Background(), &grpc_user_go.AddAccountUserRequest{
				AccountId: acc.AccountId,
				Email: "invalid_user@nalej.com",
				RoleId: entities.GenerateUUID(),
				Internal: false,
				Status: grpc_user_go.UserStatus_ACTIVE,
			})
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("should not be able to add an accountUser without Role", func() {

			_, err := client.AddAccountUser(context.Background(), &grpc_user_go.AddAccountUserRequest{
				AccountId: entities.GenerateUUID(),
				Email: user.Email,
				RoleId: entities.GenerateUUID(),
				Internal: false,
				Status: grpc_user_go.UserStatus_ACTIVE,
			})
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		// Remove
		ginkgo.It("should be able to remove an accountUser", func(){
			accountUser, err := client.AddAccountUser(context.Background(), &grpc_user_go.AddAccountUserRequest{
				AccountId: acc.AccountId,
				Email: user.Email,
				RoleId: entities.GenerateUUID(),
				Internal: false,
				Status: grpc_user_go.UserStatus_ACTIVE,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(accountUser).NotTo(gomega.BeNil())

			// Remove it
			success, err := client.RemoveAccountUser(context.Background(), &grpc_user_go.AccountUserId{
				Email: user.Email,
				AccountId: acc.AccountId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())
		})
		ginkgo.It("should not be able to remove an accountUser if the user does not exists", func(){
			accountUser, err := client.AddAccountUser(context.Background(), &grpc_user_go.AddAccountUserRequest{
				AccountId: acc.AccountId,
				Email: user.Email,
				RoleId: entities.GenerateUUID(),
				Internal: false,
				Status: grpc_user_go.UserStatus_ACTIVE,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(accountUser).NotTo(gomega.BeNil())

			// Remove it
			_, err = client.RemoveAccountUser(context.Background(), &grpc_user_go.AccountUserId{
				Email: "invalid_email@nale.com",
				AccountId: acc.AccountId,
			})
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("should not be able to remove an accountUser if does not exists", func(){

			// Remove it
			_, err := client.RemoveAccountUser(context.Background(), &grpc_user_go.AccountUserId{
				Email: user.Email,
				AccountId: acc.AccountId,
			})
			gomega.Expect(err).NotTo(gomega.Succeed())
		})

		// Update
		ginkgo.It("should be able to update an accountUser", func(){
			accountUser, err := client.AddAccountUser(context.Background(), &grpc_user_go.AddAccountUserRequest{
				AccountId: acc.AccountId,
				Email: user.Email,
				RoleId: entities.GenerateUUID(),
				Internal: false,
				Status: grpc_user_go.UserStatus_ACTIVE,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(accountUser).NotTo(gomega.BeNil())

			// update it
			updated, err := client.UpdateAccountUser(context.Background(), &grpc_user_go.AccountUserUpdateRequest{
				Email: user.Email,
				AccountId: acc.AccountId,
				UpdateRoleId: false,
				UpdateStatus: true,
				Status: grpc_user_go.UserStatus_DEACTIVATED,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(updated).NotTo(gomega.BeNil())

		})
		ginkgo.It("should not be able to update an accountUser if the user does not exists", func(){
			// update it
			_, err := client.UpdateAccountUser(context.Background(), &grpc_user_go.AccountUserUpdateRequest{
				Email: "invalid_mail@nalej.com",
				AccountId: acc.AccountId,
				UpdateRoleId: false,
				UpdateStatus: true,
				Status: grpc_user_go.UserStatus_DEACTIVATED,
			})
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("should not be able to update a non existing accountUser", func(){

			// update it
			_, err := client.UpdateAccountUser(context.Background(), &grpc_user_go.AccountUserUpdateRequest{
				Email: user.Email,
				AccountId: acc.AccountId,
				UpdateRoleId: false,
				UpdateStatus: true,
				Status: grpc_user_go.UserStatus_DEACTIVATED,
			})
			gomega.Expect(err).NotTo(gomega.Succeed())

		})

		// List
	})

	ginkgo.Context("AccountUserInvites", func(){

	})

	ginkgo.Context("ProjectUser", func(){

	})

	ginkgo.Context("ProjectUserInvites", func(){

	})
})


