package user

import (
	//"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	//"github.com/onsi/gomega"
)

func RunTest(provider Provider) {

	var email = "email@company.eu"
	var email2 = "email_2@company.eu"

	ginkgo.BeforeEach(func() {
		provider.Clear()
	})

	// AddUser
	ginkgo.It("Should be able to add user", func() {

		user := &entities.User{OrganizationId: "org",
			Email:       email,
			Name:        "name",
			MemberSince: 1}

		err := provider.Add(*user)
		gomega.Expect(err).To(gomega.Succeed())

	})

	// Update
	ginkgo.It("Should be able to update user", func() {

		user := &entities.User{OrganizationId: "organization",
			Email:       email,
			Name:        "Name",
			MemberSince: 1}

		err := provider.Add(*user)
		gomega.Expect(err).To(gomega.Succeed())

		user.OrganizationId = "organization_mod"

		err = provider.Update(*user)
		gomega.Expect(err).To(gomega.Succeed())
	})
	ginkgo.It("Should not be able to update user", func() {

		user := &entities.User{OrganizationId: "org",
			Email:       email2,
			Name:        "name",
			MemberSince: 2}

		err := provider.Update(*user)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})

	// Exists
	ginkgo.It("Should be able to find the user", func() {

		user := &entities.User{OrganizationId: "organization",
			Email:       email,
			Name:        "Name",
			MemberSince: 1}

		err := provider.Add(*user)
		gomega.Expect(err).To(gomega.Succeed())

		exists, err := provider.Exists(email)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())
	})
	ginkgo.It("Should not be able to find the user", func() {

		exists, err := provider.Exists(email2)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).NotTo(gomega.BeTrue())
	})

	// Get
	ginkgo.It("Should be able to return the user", func() {

		user := &entities.User{OrganizationId: "organization",
			Email:       email,
			Name:        "Name",
			MemberSince: 1,
			PhotoUrl:    "../../photo"}

		err := provider.Add(*user)
		gomega.Expect(err).To(gomega.Succeed())

		user, err = provider.Get(email)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(user).NotTo(gomega.BeNil())
	})
	ginkgo.It("Should not be able to return the user", func() {

		exists, err := provider.Exists(email2)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).NotTo(gomega.BeTrue())
	})

	// Remove
	ginkgo.It("Should be able to remove the user", func() {

		user := &entities.User{OrganizationId: "organization",
			Email:       email,
			Name:        "Name",
			MemberSince: 1}

		err := provider.Add(*user)
		gomega.Expect(err).To(gomega.Succeed())

		err = provider.Remove(email)
		gomega.Expect(err).To(gomega.Succeed())
	})
	ginkgo.It("Should not be able to remove the user", func() {

		err := provider.Remove(email2)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})

}
