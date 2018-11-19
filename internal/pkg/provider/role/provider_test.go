package role

import (
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {

	var roleOK = "RoleID-1"
	var roleKO = "RoleID-2"

	// AddUser
	ginkgo.It("Should be able to add role", func(){

		user := &entities.Role{ OrganizationId:"org",
			RoleId: roleOK,
			Name:"name",
			Created:1 }

		err := provider.Add(*user)
		gomega.Expect(err).To(gomega.Succeed())

	})

	//	Update
	ginkgo.It("Should be able to update role", func(){

		user := &entities.Role{ OrganizationId:"org",
			RoleId: roleOK,
			Name:"name modified",
			Created:1 }

		err := provider.Update(*user)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to update role", func(){

		user := &entities.Role{ OrganizationId:"org",
			RoleId: roleKO,
			Name:"name modified",
			Created:1 }

		err := provider.Update(*user)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	//	Exists
	ginkgo.It("Should be able to find role", func(){

		exists, err := provider.Exists(roleOK)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())

	})
	ginkgo.It("Should not be able to find role", func(){

		exists, err := provider.Exists(roleKO)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).NotTo(gomega.BeTrue())

	})
	//	Get
	ginkgo.It("Should be able to return role", func(){

		role, err := provider.Get(roleOK)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(role).NotTo(gomega.BeNil())
	})
	ginkgo.It("Should not be able to return role", func(){

		_, err := provider.Get(roleKO)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})

	//	Remove
	ginkgo.It("Should be able to remove role", func(){

		err := provider.Remove(roleOK)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to remove role", func(){

		 err := provider.Remove(roleKO)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

}
