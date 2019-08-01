/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package project

import (
	"context"
	"fmt"
	"github.com/nalej/grpc-account-go"
	"github.com/nalej/grpc-project-go"
	"github.com/nalej/grpc-utils/pkg/test"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/account"
	"github.com/nalej/system-model/internal/pkg/provider/project"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"time"
)

func AddAccount(provider account.Provider) *entities.Account{
	account := &entities.Account{
		AccountId: entities.GenerateUUID(),
		Name: fmt.Sprintf("Account-%s", entities.GenerateUUID()),
		Created: time.Now().Unix(),
		State: entities.AccountState_Active,
	}
	err := provider.Add(*account)
	gomega.Expect(err).To(gomega.Succeed())

	return account
}

func CreateAddProjectRequest(accountId string) *grpc_project_go.AddProjectRequest {
	return &grpc_project_go.AddProjectRequest{
		AccountId:accountId,
		Name: fmt.Sprintf("Project-%s", entities.GenerateUUID()),
	}
}

var _ = ginkgo.Describe("Project service", func() {
	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client grpc_project_go.ProjectsClient

	// Providers
	var accountProvider account.Provider
	var projectProvider project.Provider

	var targetAccount *entities.Account

	ginkgo.BeforeSuite(func() {
		listener = test.GetDefaultListener()
		server = grpc.NewServer()
		test.LaunchServer(server, listener)

		// Register the service
		accountProvider = account.NewMockupAccountProvider()
		projectProvider = project.NewMockupProjectProvider()
		manager := NewManager(accountProvider, projectProvider)
		handler := NewHandler(manager)
		grpc_project_go.RegisterProjectsServer(server, handler)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).Should(gomega.Succeed())
		client = grpc_project_go.NewProjectsClient(conn)
	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
	})

	ginkgo.BeforeEach(func() {
		targetAccount = AddAccount(accountProvider)
	})
 	//	AddProject(context.Context, *AddProjectRequest) (*Project, error)
	ginkgo.Context("Adding project", func() {
		ginkgo.It("should be able to add a new project", func() {

			toAdd := CreateAddProjectRequest(targetAccount.AccountId)

			project , err := client.AddProject(context.Background(),toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(project).NotTo(gomega.BeNil())
			gomega.Expect(project.ProjectId).NotTo(gomega.BeEmpty())
		})
		ginkgo.It("should not be able to add a project in a wrong account", func() {
			toAdd := CreateAddProjectRequest(entities.GenerateUUID())

			_ , err := client.AddProject(context.Background(),toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
	})
	// 	GetProject(context.Context, *ProjectId) (*Project, error)
	ginkgo.Context("getting project", func() {
		ginkgo.It("should be able to get a project", func() {

			toAdd := CreateAddProjectRequest(targetAccount.AccountId)

			project , err := client.AddProject(context.Background(),toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(project).NotTo(gomega.BeNil())
			gomega.Expect(project.ProjectId).NotTo(gomega.BeEmpty())

			retrieved, err := client.GetProject(context.Background(), &grpc_project_go.ProjectId{
				AccountId: project.OwnerAccountId,
				ProjectId: project.ProjectId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).NotTo(gomega.BeNil())
			gomega.Expect(retrieved).Should(gomega.Equal(project))

		})
		ginkgo.It("should not be able to get a project of a wrong account", func() {

			_ , err := client.GetProject(context.Background(),&grpc_project_go.ProjectId{
				AccountId: entities.GenerateUUID(),
				ProjectId: entities.GenerateUUID(),
			})
			gomega.Expect(err).NotTo(gomega.Succeed())


		})
		ginkgo.It("should not be able to get a non existing project", func() {

			_ , err := client.GetProject(context.Background(),&grpc_project_go.ProjectId{
				AccountId: targetAccount.AccountId,
				ProjectId: entities.GenerateUUID(),
			})
			gomega.Expect(err).NotTo(gomega.Succeed())


		})
	})
	// RemoveProject(context.Context, *ProjectId) (*grpc_common_go.Success, error)
	ginkgo.Context("Removing project", func() {
		ginkgo.It("should be able to remove a project", func(){
			toAdd := CreateAddProjectRequest(targetAccount.AccountId)

			project , err := client.AddProject(context.Background(),toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(project).NotTo(gomega.BeNil())

			success, err := client.RemoveProject(context.Background(), &grpc_project_go.ProjectId{
				AccountId: targetAccount.AccountId,
				ProjectId: project.ProjectId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())

		})
		ginkgo.It("should not be able to remove a project if it does not exist", func(){
			success, err := client.RemoveProject(context.Background(), &grpc_project_go.ProjectId{
				AccountId: targetAccount.AccountId,
				ProjectId: entities.GenerateUUID(),
			})
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())

		})
	})
	// 	ListAccountProjects(context.Context, *grpc_account_go.AccountId) (*ProjectList, error)
	ginkgo.Context("listing projects", func() {
		ginkgo.It("should be able to list the projects of an account", func(){
			numProjects := 10
			for i:= 0; i<numProjects; i++ {
				toAdd := CreateAddProjectRequest(targetAccount.AccountId)

				project , err := client.AddProject(context.Background(),toAdd)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(project).NotTo(gomega.BeNil())
			}

			projects, err := client.ListAccountProjects(context.Background(), &grpc_account_go.AccountId{
				AccountId: targetAccount.AccountId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(projects).NotTo(gomega.BeNil())
			gomega.Expect(len(projects.Projects)).Should(gomega.Equal(numProjects))

		})
		ginkgo.It("should be able to return an empty list of projects of an existing account", func(){
			projects, err := client.ListAccountProjects(context.Background(), &grpc_account_go.AccountId{
				AccountId: targetAccount.AccountId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(projects).NotTo(gomega.BeNil())
			gomega.Expect(len(projects.Projects)).Should(gomega.Equal(0))

		})
		ginkgo.It("should not be able to return a list of projects of a non existing account", func(){
			_, err := client.ListAccountProjects(context.Background(), &grpc_account_go.AccountId{
				AccountId: entities.GenerateUUID(),
			})
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
	})
    // 	UpdateProject(context.Context, *UpdateProjectRequest) (*grpc_common_go.Success, error)
	ginkgo.Context("Updating project", func() {
		ginkgo.It("should be able to update a project", func() {

			toAdd := CreateAddProjectRequest(targetAccount.AccountId)

			project , err := client.AddProject(context.Background(),toAdd)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(project).NotTo(gomega.BeNil())

			success, err := client.UpdateProject(context.Background(), &grpc_project_go.UpdateProjectRequest{
				AccountId:targetAccount.AccountId,
				ProjectId: project.ProjectId,
				UpdateName: true,
				Name: "name updated",
				UpdateState: true,
				State: grpc_project_go.ProjectState_DEACTIVATED,
				UpdateStateInfo: true,
				StateInfo: "deactivated for test",
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(success).NotTo(gomega.BeNil())

			// check if the update works
			retrieved, err := client.GetProject(context.Background(), &grpc_project_go.ProjectId{
				AccountId:targetAccount.AccountId,
				ProjectId: project.ProjectId,
			})
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).NotTo(gomega.BeNil())
			gomega.Expect(retrieved.Name).Should(gomega.Equal("name updated"))
			gomega.Expect(retrieved.State).Should(gomega.Equal(grpc_project_go.ProjectState_DEACTIVATED))
			gomega.Expect(retrieved.StateInfo).Should(gomega.Equal("deactivated for test"))

		})
		ginkgo.It("should not be able to update a non existing project", func() {
			success, err := client.UpdateProject(context.Background(), &grpc_project_go.UpdateProjectRequest{
				AccountId:targetAccount.AccountId,
				ProjectId: entities.GenerateUUID(),
				UpdateName: true,
				Name: "name updated",
				UpdateState: true,
				State: grpc_project_go.ProjectState_DEACTIVATED,
				UpdateStateInfo: true,
				StateInfo: "deactivated for test",
			})
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(success).To(gomega.BeNil())

		})
	})

})