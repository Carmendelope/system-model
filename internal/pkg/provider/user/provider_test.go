package user

import (
	"fmt"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	//"github.com/onsi/gomega"
)


func RunTest(provider Provider) {

	ginkgo.BeforeEach(func(){
		provider.Clear()
	})

	ginkgo.Context("testing User", func() {
		// AddUser
		ginkgo.It("Should be able to add user", func(){
			mail := fmt.Sprintf("%s@nalej.com", entities.GenerateUUID())
			user := CreateUser(mail)
			err := provider.Add(*user)
			gomega.Expect(err).To(gomega.Succeed())
		})

		// Update
		ginkgo.It("Should be able to update user", func(){
			mail := fmt.Sprintf("%s@nalej.com", entities.GenerateUUID())
			user := CreateUser(mail)
			err := provider.Add(*user)
			gomega.Expect(err).To(gomega.Succeed())

			user.Name = "new name"
			user.ContactInfo.CompanyName = "new company name"
			err = provider.Update(*user)
			gomega.Expect(err).To(gomega.Succeed())

			// check the update works
			retrieved, err := provider.Get(user.Email)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).Should(gomega.Equal(user))
		})
		ginkgo.It("Should not be able to update user", func(){
			mail := fmt.Sprintf("%s@nalej.com", entities.GenerateUUID())
			user := CreateUser(mail)

			err := provider.Update(*user)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})

		// Exists
		ginkgo.It("Should be able to find the user", func(){

			mail := fmt.Sprintf("%s@nalej.com", entities.GenerateUUID())
			user := CreateUser(mail)
			err := provider.Add(*user)
			gomega.Expect(err).To(gomega.Succeed())

			exists, err := provider.Exists(user.Email)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())
		})
		ginkgo.It("Should not be able to find the user", func(){

			exists, err := provider.Exists("invalid_email@nalej.com")
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())
		})

		// Get
		ginkgo.It("Should be able to return the user", func(){

			mail := fmt.Sprintf("%s@nalej.com", entities.GenerateUUID())
			user := CreateUser(mail)
			err := provider.Add(*user)
			gomega.Expect(err).To(gomega.Succeed())

			user, err = provider.Get(user.Email)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(user).NotTo(gomega.BeNil())
		})
		ginkgo.It("Should not be able to return the user", func(){

			exists, err := provider.Exists("invalid_email@nalej.com")
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())
		})

		// Remove
		ginkgo.It("Should be able to remove the user", func(){

			user := CreateUser(entities.GenerateUUID())

			err := provider.Add(*user)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.Remove(user.Email)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("Should not be able to remove the user", func(){

			err := provider.Remove("invalid_email@nalej.com")
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})

	ginkgo.Context("AccountUser", func() {

		var user *entities.User

		ginkgo.BeforeEach(func(){
			user, _ = CreateAddUser(provider)
			gomega.Expect(user).NotTo(gomega.BeNil())
		})
		// Add
		ginkgo.It("Should be able to add an accountUser", func() {
			accUser := CreateAccountUser(user.Email)
			err := provider.AddAccountUser(*accUser)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("Should not be able to add an accountUser twice", func() {
			accUser := CreateAccountUser(user.Email)

			err := provider.AddAccountUser(*accUser)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.AddAccountUser(*accUser)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("should not be able to add an accountUser if the user does not exist", func(){
			accUser := CreateAccountUser("invalid_email@nalej.com")

			err := provider.AddAccountUser(*accUser)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		// Update
		ginkgo.It("Should be able to update an accountUser", func() {
			accUser := CreateAccountUser(user.Email)
			err := provider.AddAccountUser(*accUser)
			gomega.Expect(err).To(gomega.Succeed())

			accUser.Status = 0
			accUser.RoleId = entities.GenerateUUID()
			err = provider.UpdateAccountUser(*accUser)
			gomega.Expect(err).To(gomega.Succeed())

		})
		ginkgo.It("Should not be able to update an accountUser if it does not exist", func() {
			accUser := CreateAccountUser(user.Email)
			err := provider.UpdateAccountUser(*accUser)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		//Remove
		ginkgo.It("Should be able to remove an accountUser", func() {
			accUser := CreateAccountUser(user.Email)
			err := provider.AddAccountUser(*accUser)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.RemoveAccountUser(accUser.AccountId, accUser.Email)
			gomega.Expect(err).To(gomega.Succeed())

		})
		ginkgo.It("Should not be able to remove an accountUser if it does not exist", func() {
			accUser := CreateAccountUser(user.Email)
			err := provider.RemoveAccountUser(accUser.AccountId, accUser.Email)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		// List
		ginkgo.It("Should be able to list the accountUser of a user", func() {
			num := 10
			for i:= 0; i < num; i++ {
				accUser := CreateAccountUser(user.Email)
				err := provider.AddAccountUser(*accUser)
				gomega.Expect(err).To(gomega.Succeed())
			}

			list, err := provider.ListAccountUser(user.Email)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list)).Should(gomega.Equal(10))

		})
		ginkgo.It("Should be able to return an empty list of accountUser", func() {

			list, err := provider.ListAccountUser(user.Email)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list)).Should(gomega.Equal(0))

		})
		ginkgo.It("Should not be able to return a list of accountUser if the user does not exist", func() {

			_, err := provider.ListAccountUser("invalid_user@nalej.com")
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
	})

	ginkgo.Context("AccountUserInvite", func() {
		var user *entities.User

		ginkgo.BeforeEach(func(){
			user, _ = CreateAddUser(provider)
			gomega.Expect(user).NotTo(gomega.BeNil())
		})

		// AddAccountUserInvite
		ginkgo.It("Should be able to add an accountUserInvite", func() {
			invite := CreateAccountUserInvite(user.Email)
			err := provider.AddAccountUserInvite(*invite)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("Should not be able to add an accountUserInvite twice", func() {
			invite := CreateAccountUserInvite(user.Email)

			err := provider.AddAccountUserInvite(*invite)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.AddAccountUserInvite(*invite)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("should not be able to add an accountUserInvite if the user does not exist", func(){
			invite := CreateAccountUserInvite("invalid_email@nalej.com")

			err := provider.AddAccountUserInvite(*invite)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		// GetAccountUserInvite
		ginkgo.It("Should be able to get an accountUserInvite", func() {
			invite := CreateAccountUserInvite(user.Email)
			err := provider.AddAccountUserInvite(*invite)
			gomega.Expect(err).To(gomega.Succeed())

			retrieve, err := provider.GetAccountUserInvite(invite.AccountId, invite.Email)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieve).Should(gomega.Equal(invite))

		})
		ginkgo.It("Should not be able to get an accountUserInvite if it does not exist", func() {
			_ , err  := provider.GetAccountUserInvite(entities.GenerateUUID(), user.Email)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		// RemoveAccountUserInvite
		ginkgo.It("Should be able to remove an accountUserInvite", func() {
			invite := CreateAccountUserInvite(user.Email)
			err := provider.AddAccountUserInvite(*invite)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.RemoveAccountUserInvite(invite.AccountId, invite.Email)
			gomega.Expect(err).To(gomega.Succeed())

		})
		ginkgo.It("Should not be able to remove an accountUserInvite if it does not exist", func() {
			err  := provider.RemoveAccountUserInvite(entities.GenerateUUID(), user.Email)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		// List
		ginkgo.It("Should be able to list the accountUserInvite of a user", func() {
			num := 10
			for i:= 0; i < num; i++ {
				invite := CreateAccountUserInvite(user.Email)
				err := provider.AddAccountUserInvite(*invite)
				gomega.Expect(err).To(gomega.Succeed())
			}

			list, err := provider.ListAccountUserInvites(user.Email)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list)).Should(gomega.Equal(10))

		})
		ginkgo.It("Should be able to return an empty list of accountUserInvite", func() {

			list, err := provider.ListAccountUserInvites(user.Email)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list)).Should(gomega.Equal(0))

		})
		ginkgo.It("Should not be able to return a list of accountUserInvite if the user does not exist", func() {

			_, err := provider.ListAccountUserInvites("invalid_user@nalej.com")
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
	})

	ginkgo.Context("ProjectUser", func(){
		var user *entities.User

		ginkgo.BeforeEach(func(){
			user, _ = CreateAddUser(provider)
			gomega.Expect(user).NotTo(gomega.BeNil())
		})
		// Add
		ginkgo.It("Should be able to add a projectUser", func() {
			projUser := CreateProjectUser(entities.GenerateUUID(), entities.GenerateUUID(), user.Email)
			err := provider.AddProjectUser(*projUser)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("Should not be able to add a projectUser twice", func() {
			projUser := CreateProjectUser(entities.GenerateUUID(), entities.GenerateUUID(), user.Email)

			err := provider.AddProjectUser(*projUser)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.AddProjectUser(*projUser)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("should not be able to add a projectUser if the user does not exist", func(){
			projUser := CreateProjectUser(entities.GenerateUUID(), entities.GenerateUUID(), "invalid_email@nalej.com")

			err := provider.AddProjectUser(*projUser)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		// Update
		ginkgo.It("Should be able to update a projectUser", func() {
			projUser := CreateProjectUser(entities.GenerateUUID(), entities.GenerateUUID(), user.Email)
			err := provider.AddProjectUser(*projUser)
			gomega.Expect(err).To(gomega.Succeed())

			projUser.Status = 0
			projUser.RoleId = entities.GenerateUUID()
			err = provider.UpdateProjectUser(*projUser)
			gomega.Expect(err).To(gomega.Succeed())

		})
		ginkgo.It("Should not be able to update a projectUser if it does not exist", func() {
			projUser := CreateProjectUser(entities.GenerateUUID(), entities.GenerateUUID(), user.Email)
			err := provider.UpdateProjectUser(*projUser)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		//Remove
		ginkgo.It("Should be able to remove a projectUser", func() {
			projUser := CreateProjectUser(entities.GenerateUUID(),entities.GenerateUUID(), user.Email)
			err := provider.AddProjectUser(*projUser)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.RemoveProjectUser(projUser.AccountId, projUser.ProjectId, projUser.Email)
			gomega.Expect(err).To(gomega.Succeed())

		})
		ginkgo.It("Should not be able to remove a projectUserUser if it does not exist", func() {
			projUser := CreateProjectUser(entities.GenerateUUID(), entities.GenerateUUID(), user.Email)
			err := provider.RemoveProjectUser(projUser.AccountId, projUser.ProjectId, projUser.Email)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		// List
		ginkgo.It("Should be able to list the projectUser of a project", func() {
			num := 10
			accountID := entities.GenerateUUID()
			projectID := entities.GenerateUUID()
			for i:= 0; i < num; i++ {
				email := generateMail()
				err := AddUser(provider, CreateUser(email))

				projUser := CreateProjectUser(accountID, projectID, email)
				err = provider.AddProjectUser(*projUser)
				gomega.Expect(err).To(gomega.Succeed())
			}

			list, err := provider.ListProjectUser(accountID, projectID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list)).Should(gomega.Equal(10))

		})
		ginkgo.It("Should be able to return an empty list of projectUser", func() {

			list, err := provider.ListProjectUser(entities.GenerateUUID(), entities.GenerateUUID())
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list)).Should(gomega.Equal(0))

		})
	})

	ginkgo.Context("ProjectUserInvite", func() {

		// CreateProjectUserInvite

		var user *entities.User
		var accountId string
		var projectId string

		ginkgo.BeforeEach(func(){
			user, _ = CreateAddUser(provider)
			gomega.Expect(user).NotTo(gomega.BeNil())

			accountId = entities.GenerateUUID()
			projectId = entities.GenerateUUID()

		})

		// AddProjectUserInvite
		ginkgo.It("Should be able to add an projectUserInvite", func() {
			invite := CreateProjectUserInvite(accountId, projectId, user.Email)
			err := provider.AddProjectUserInvite(*invite)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("Should not be able to add a projectUserInvite twice", func() {
			invite := CreateProjectUserInvite(accountId, projectId, user.Email)

			err := provider.AddProjectUserInvite(*invite)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.AddProjectUserInvite(*invite)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("should not be able to add a projectUserInvite if the user does not exist", func(){
			invite := CreateProjectUserInvite(accountId, projectId, "invalid_user@nalej.com")

			err := provider.AddProjectUserInvite(*invite)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		// GetProjectUserInvite
		ginkgo.It("Should be able to get a projectUserInvite", func() {
			invite := CreateProjectUserInvite(accountId, projectId, user.Email)
			err := provider.AddProjectUserInvite(*invite)
			gomega.Expect(err).To(gomega.Succeed())

			retrieve, err := provider.GetProjectUserInvite(invite.AccountId, invite.ProjectId, invite.Email)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieve).Should(gomega.Equal(invite))

		})
		ginkgo.It("Should not be able to get a projectUserInvite if it does not exist", func() {
			_ , err  := provider.GetProjectUserInvite(entities.GenerateUUID(),entities.GenerateUUID(), user.Email)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		// RemoveProjectUserInvite
		ginkgo.It("Should be able to remove a projectUserInvite", func() {
			invite := CreateProjectUserInvite(accountId, projectId, user.Email)
			err := provider.AddProjectUserInvite(*invite)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.RemoveProjectUserInvite(invite.AccountId, invite.ProjectId, invite.Email)
			gomega.Expect(err).To(gomega.Succeed())

		})
		ginkgo.It("Should not be able to remove a projectUserInvite if it does not exist", func() {
			err  := provider.RemoveProjectUserInvite(entities.GenerateUUID(), entities.GenerateUUID(), user.Email)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		// List
		ginkgo.It("Should be able to list the projectUserInvites of a user", func() {
			num := 10
			for i:= 0; i < num; i++ {

				invite := CreateProjectUserInvite(entities.GenerateUUID(), entities.GenerateUUID(), user.Email)
				err := provider.AddProjectUserInvite(*invite)
				gomega.Expect(err).To(gomega.Succeed())
			}
			list, err := provider.ListProjectUserInvites(user.Email)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list)).Should(gomega.Equal(10))

		})
		ginkgo.It("Should be able to return an empty list of projectUserInvites", func() {

			list, err := provider.ListProjectUserInvites(user.Email)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list)).Should(gomega.Equal(0))

		})
		ginkgo.It("Should not be able to return a list of projectUserInvites if the user does not exist", func() {

			_, err := provider.ListProjectUserInvites("invalid_user@nalej.com")
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
	})
}

