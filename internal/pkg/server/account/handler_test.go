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

package account

import (
	"context"
	"fmt"
	"github.com/nalej/grpc-account-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/account"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func CreateAddAccountRequest() *grpc_account_go.AddAccountRequest {
	return &grpc_account_go.AddAccountRequest{

		Name:           fmt.Sprintf("Account-%s", entities.GenerateUUID()),
		FullName:       "Account Full Name",
		CompanyName:    "Company",
		Address:        "Address 10",
		AdditionalInfo: "Additional info of account",
	}
}

var _ = ginkgo.Describe("Account service", func() {
	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client grpc_account_go.AccountsClient

	// Providers
	var accountProvider account.Provider

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()

		// Register the service
		accountProvider = account.NewMockupAccountProvider()
		manager := NewManager(accountProvider)
		handler := NewHandler(manager)
		grpc_account_go.RegisterAccountsServer(server, handler)

		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_account_go.NewAccountsClient(conn)
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
	})

	//AddAccount(ctx context.Context, in *AddAccountRequest, opts ...grpc.CallOption) (*Account, error)
	ginkgo.Context("Adding account", func() {
		ginkgo.It("should be able to add a new account", func() {
			toAdd := CreateAddAccountRequest()
			account, err := client.AddAccount(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(account).NotTo(gomega.BeNil())
			gomega.Expect(account.AccountId).NotTo(gomega.BeEmpty())
		})
		ginkgo.It("should not be able to add an account without name", func() {
			toAdd := CreateAddAccountRequest()
			toAdd.Name = ""
			_, err := client.AddAccount(context.Background(), toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		ginkgo.It("should not be able to add two accounts with the same name", func() {
			toAdd := CreateAddAccountRequest()
			account, err := client.AddAccount(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(account).NotTo(gomega.BeNil())

			_, err = client.AddAccount(context.Background(), toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
	})
	// GetAccount(ctx context.Context, in *AccountId, opts ...grpc.CallOption) (*Account, error)
	ginkgo.Context("Getting account", func() {
		ginkgo.It("should be able to get an account", func() {
			toAdd := CreateAddAccountRequest()
			account, err := client.AddAccount(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(account).NotTo(gomega.BeNil())

			retrieved, err := client.GetAccount(context.Background(), &grpc_account_go.AccountId{
				AccountId: account.AccountId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).NotTo(gomega.BeNil())
			gomega.Expect(retrieved).Should(gomega.Equal(account))

		})
		ginkgo.It("should not be able to get a non existing account", func() {
			_, err := client.GetAccount(context.Background(), &grpc_account_go.AccountId{
				AccountId: entities.GenerateUUID(),
			})
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
	})
	// 	ListAccounts(ctx context.Context, in *grpc_common_go.Empty, opts ...grpc.CallOption) (*AccountList, error)
	ginkgo.Context("Getting account", func() {
		ginkgo.It("should be able to get an account", func() {
			numAccounts := 10
			for i := 0; i < numAccounts; i++ {
				toAdd := CreateAddAccountRequest()
				account, err := client.AddAccount(context.Background(), toAdd)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(account).NotTo(gomega.BeNil())
			}

			list, err := client.ListAccounts(context.Background(), &grpc_common_go.Empty{})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(list).NotTo(gomega.BeNil())
			gomega.Expect(len(list.Accounts)).ShouldNot(gomega.BeZero())

		})
	})
	//UpdateAccount(ctx context.Context, in *UpdateAccountRequest, opts ...grpc.CallOption) (*grpc_common_go.Success, error)
	ginkgo.Context("Updating account", func() {
		ginkgo.It("Should be able to update an account", func() {

			toAdd := CreateAddAccountRequest()
			account, err := client.AddAccount(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(account).NotTo(gomega.BeNil())

			success, err := client.UpdateAccount(context.Background(), &grpc_account_go.UpdateAccountRequest{
				AccountId:  account.AccountId,
				UpdateName: true,
				Name:       "updated name",
			})
			gomega.Expect(success).ShouldNot(gomega.BeNil())
			gomega.Expect(err).To(gomega.Succeed())

			// check the account to check the udpate works
			retrieved, err := client.GetAccount(context.Background(), &grpc_account_go.AccountId{
				AccountId: account.AccountId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved.Name).Should(gomega.Equal("updated name"))

		})
		ginkgo.It("Should not be able to update the name of an account if the name already exists", func() {

			toAdd1 := CreateAddAccountRequest()
			account1, err := client.AddAccount(context.Background(), toAdd1)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(account1).NotTo(gomega.BeNil())

			toAdd2 := CreateAddAccountRequest()
			account2, err := client.AddAccount(context.Background(), toAdd2)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(account2).NotTo(gomega.BeNil())

			success, err := client.UpdateAccount(context.Background(), &grpc_account_go.UpdateAccountRequest{
				AccountId:  account1.AccountId,
				UpdateName: true,
				Name:       account2.Name,
			})
			gomega.Expect(success).Should(gomega.BeNil())
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
	})
	//UpdateAccountBillingInfo(ctx context.Context, in *UpdateAccountBillingInfoRequest, opts ...grpc.CallOption) (*grpc_common_go.Success, error)
	ginkgo.Context("updating billing info of an account", func() {
		ginkgo.It("Should be able to update the billing information of account", func() {

			toAdd := CreateAddAccountRequest()
			account, err := client.AddAccount(context.Background(), toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(account).NotTo(gomega.BeNil())

			success, err := client.UpdateAccountBillingInfo(context.Background(), &grpc_account_go.UpdateAccountBillingInfoRequest{
				AccountId:      account.AccountId,
				UpdateFullName: true,
				FullName:       "full name updated",
			})
			gomega.Expect(success).ShouldNot(gomega.BeNil())
			gomega.Expect(err).To(gomega.Succeed())

			// check the account to check the udpate works
			retrieved, err := client.GetAccount(context.Background(), &grpc_account_go.AccountId{
				AccountId: account.AccountId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved.BillingInfo.FullName).Should(gomega.Equal("full name updated"))
		})
		ginkgo.It("Should not be able to update the billing information of a non existing account", func() {

			success, err := client.UpdateAccountBillingInfo(context.Background(), &grpc_account_go.UpdateAccountBillingInfoRequest{
				AccountId:      entities.GenerateUUID(),
				UpdateFullName: true,
				FullName:       "full name updated",
			})
			gomega.Expect(success).Should(gomega.BeNil())
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
	})
})
