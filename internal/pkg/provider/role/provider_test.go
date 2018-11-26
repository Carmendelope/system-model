package role

import (
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {

	var roleOK = "RoleID-1"
	var roleKO = "RoleID-2"

	ginkgo.BeforeEach(func(){
		provider.Clear()
	})

	// AddUser
	ginkgo.It("Should be able to add role", func(){

		role := &entities.Role{ OrganizationId:"org",
			RoleId: roleOK,
			Name:"name",
			Created:1 }

		err := provider.Add(*role)
		gomega.Expect(err).To(gomega.Succeed())

	})

	//	Update
	ginkgo.It("Should be able to update role", func(){

		// insert a role
		role := &entities.Role{ OrganizationId:"org",
			RoleId: roleOK,
			Name:"name",
			Created:1 }

		err := provider.Add(*role)
		gomega.Expect(err).To(gomega.Succeed())

		role.OrganizationId = "organization_MOD"
		err = provider.Update(*role)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to update role", func(){

		role := &entities.Role{ OrganizationId:"org",
			RoleId: roleOK,
			Name:"name modified",
			Created:1 }

		err := provider.Update(*role)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	//	Exists
	ginkgo.It("Should be able to find role", func(){

		// insert a role
		role := &entities.Role{ OrganizationId:"org",
			RoleId: roleOK,
			Name:"name modified",
			Created:1 }

		err := provider.Add(*role)
		gomega.Expect(err).To(gomega.Succeed())

		// ask if exists
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

		// insert a role
		role := &entities.Role{ OrganizationId:"org",
			RoleId: roleOK,
			Name:"Name",
			Created:1 }

		err := provider.Add(*role)
		gomega.Expect(err).To(gomega.Succeed())

		returnedRole, err := provider.Get(roleOK)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(returnedRole).NotTo(gomega.BeNil())
	})
	ginkgo.It("Should not be able to return role", func(){

		_, err := provider.Get(roleKO)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})

	//	Remove
	ginkgo.It("Should be able to remove role", func(){

		// insert a role
		role := &entities.Role{ OrganizationId:"org",
			RoleId: roleOK,
			Name:"Name",
			Created:1 }

		err := provider.Add(*role)
		gomega.Expect(err).To(gomega.Succeed())

		// remove it
		err = provider.Remove(role.RoleId)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to remove role", func(){

		 err := provider.Remove(roleKO)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

}
