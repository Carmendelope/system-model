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

package role

import (
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {

	var roleOK = "RoleID-1"
	var roleKO = "RoleID-2"

	ginkgo.BeforeEach(func() {
		provider.Clear()
	})

	// AddUser
	ginkgo.It("Should be able to add role", func() {

		role := &entities.Role{OrganizationId: "org",
			RoleId:  roleOK,
			Name:    "name",
			Created: 1}

		err := provider.Add(*role)
		gomega.Expect(err).To(gomega.Succeed())

	})

	//	Update
	ginkgo.It("Should be able to update role", func() {

		// insert a role
		role := &entities.Role{OrganizationId: "org",
			RoleId:  roleOK,
			Name:    "name",
			Created: 1}

		err := provider.Add(*role)
		gomega.Expect(err).To(gomega.Succeed())

		role.OrganizationId = "organization_MOD"
		err = provider.Update(*role)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to update role", func() {

		role := &entities.Role{OrganizationId: "org",
			RoleId:  roleOK,
			Name:    "name modified",
			Created: 1}

		err := provider.Update(*role)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	//	Exists
	ginkgo.It("Should be able to find role", func() {

		// insert a role
		role := &entities.Role{OrganizationId: "org",
			RoleId:  roleOK,
			Name:    "name modified",
			Created: 1}

		err := provider.Add(*role)
		gomega.Expect(err).To(gomega.Succeed())

		// ask if exists
		exists, err := provider.Exists(roleOK)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())

	})
	ginkgo.It("Should not be able to find role", func() {

		exists, err := provider.Exists(roleKO)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).NotTo(gomega.BeTrue())

	})
	//	Get
	ginkgo.It("Should be able to return role", func() {

		// insert a role
		role := &entities.Role{OrganizationId: "org",
			RoleId:  roleOK,
			Name:    "Name",
			Created: 1}

		err := provider.Add(*role)
		gomega.Expect(err).To(gomega.Succeed())

		returnedRole, err := provider.Get(roleOK)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(returnedRole).NotTo(gomega.BeNil())
	})
	ginkgo.It("Should not be able to return role", func() {

		_, err := provider.Get(roleKO)
		gomega.Expect(err).NotTo(gomega.Succeed())
	})

	//	Remove
	ginkgo.It("Should be able to remove role", func() {

		// insert a role
		role := &entities.Role{OrganizationId: "org",
			RoleId:  roleOK,
			Name:    "Name",
			Created: 1}

		err := provider.Add(*role)
		gomega.Expect(err).To(gomega.Succeed())

		// remove it
		err = provider.Remove(role.RoleId)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to remove role", func() {

		err := provider.Remove(roleKO)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

}
