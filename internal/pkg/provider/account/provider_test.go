/*
 * Copyright 2020 Nalej
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

package account

import (
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {
	ginkgo.AfterEach(func() {
		provider.Clear()
	})
	ginkgo.Context("adding account", func() {
		ginkgo.It("should be able to add an account", func() {
			toAdd := CreateAccount()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("should not be able to add an account twice", func() {
			toAdd := CreateAccount()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.Add(*toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})
	ginkgo.Context("getting account", func() {
		ginkgo.It("should be able to get an account", func() {
			toAdd := CreateAccount()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			retrieve, err := provider.Get(toAdd.AccountId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieve).NotTo(gomega.BeNil())
			gomega.Expect(retrieve).Should(gomega.Equal(toAdd))
		})
		ginkgo.It("should not be able to get a non existing account", func() {
			_, err := provider.Get(entities.GenerateUUID())
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})
	ginkgo.Context("removing account", func() {
		ginkgo.It("should be able to remove an account", func() {
			toAdd := CreateAccount()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.Remove(toAdd.AccountId)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("should not be able to remove a non existing account", func() {
			err := provider.Remove(entities.GenerateUUID())
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})
	ginkgo.Context("updating account", func() {
		ginkgo.It("should be able to update an account", func() {
			toAdd := CreateAccount()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			// update Account
			toAdd.Name = "updated name"
			toAdd.BillingInfo.FullName = "full name updated"
			toAdd.State = entities.AccountState_Deactivated
			toAdd.StateInfo = "deactivated for test"

			err = provider.Update(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			// check the update works
			retrieve, err := provider.Get(toAdd.AccountId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieve).NotTo(gomega.BeNil())
			gomega.Expect(retrieve).Should(gomega.Equal(toAdd))

		})
		ginkgo.It("should not be able to update a non existing account", func() {
			toAdd := CreateAccount()

			err := provider.Update(*toAdd)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})
	ginkgo.Context("checking if exists account", func() {
		ginkgo.It("should be able to check an account exists", func() {
			toAdd := CreateAccount()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			exists, err := provider.Exists(toAdd.AccountId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())
		})
		ginkgo.It("should be able to check an account does not exist", func() {
			exists, err := provider.Exists(entities.GenerateUUID())
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())
		})
	})
	ginkgo.Context("checking if exists account by name", func() {
		ginkgo.It("should be able to check if a name of an account exists", func() {
			toAdd := CreateAccount()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			exists, err := provider.ExistsByName(toAdd.Name)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())
		})
		ginkgo.It("should be able to check that a name of an account does not exist", func() {
			exists, err := provider.ExistsByName(entities.GenerateUUID())
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())
		})
		ginkgo.It("should be able to check that a name of an account does not exist after delete it", func() {
			toAdd := CreateAccount()
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			// remove the account
			err = provider.Remove(toAdd.AccountId)
			gomega.Expect(err).To(gomega.Succeed())

			exists, err := provider.ExistsByName(toAdd.Name)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())
		})
		ginkgo.It("should be able to check that a name of an account does not exist after update it", func() {
			toAdd := CreateAccount()
			name := toAdd.Name
			err := provider.Add(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			// update the account
			toAdd.Name = "name updated"
			err = provider.Update(*toAdd)
			gomega.Expect(err).To(gomega.Succeed())

			exists, err := provider.ExistsByName(name)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())
		})
	})
	ginkgo.Context("listing accounts", func() {
		ginkgo.It("should be able to list accounts where there are", func() {
			numAccounts := 10
			for i := 0; i < numAccounts; i++ {
				toAdd := CreateAccount()
				err := provider.Add(*toAdd)
				gomega.Expect(err).To(gomega.Succeed())
			}
			list, err := provider.List()
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list)).Should(gomega.Equal(numAccounts))
		})
		ginkgo.It("should be able to return an empty list of accounts", func() {
			list, err := provider.List()
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(len(list)).Should(gomega.Equal(0))
		})
	})
}
