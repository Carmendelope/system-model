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
 *
 */

package organization_setting

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"math/rand"
)

func CreateOrganizationSetting(organizationID string) *entities.OrganizationSetting {
	key := rand.Int()
	return entities.NewOrganizationSetting(organizationID, fmt.Sprintf("key_%d", key),
		fmt.Sprintf("%d", key), "description")
}

func RunTest(provider Provider) {

	ginkgo.BeforeEach(func() {
		provider.Clear()
	})

	ginkgo.Context("adding a setting", func() {
		ginkgo.It("should be able to add a setting", func() {
			orgId := uuid.New().String()
			setting := CreateOrganizationSetting(orgId)
			err := provider.Add(*setting)
			gomega.Expect(err).Should(gomega.Succeed())
		})
		ginkgo.It("should not be able to add a setting twice", func() {
			orgId := uuid.New().String()
			setting := CreateOrganizationSetting(orgId)
			err := provider.Add(*setting)
			gomega.Expect(err).Should(gomega.Succeed())

			err = provider.Add(*setting)
			gomega.Expect(err).ShouldNot(gomega.Succeed())
		})
	})
	ginkgo.Context("getting a setting", func() {
		ginkgo.It("should be able to get a setting", func() {

			orgId := uuid.New().String()
			setting := CreateOrganizationSetting(orgId)
			err := provider.Add(*setting)
			gomega.Expect(err).Should(gomega.Succeed())

			retrieved, err := provider.Get(setting.OrganizationId, setting.Key)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(retrieved).NotTo(gomega.BeNil())
			gomega.Expect(retrieved).Should(gomega.Equal(setting))

		})
		ginkgo.It("should not be able to get a non existing setting", func() {
			orgId := uuid.New().String()
			setting := CreateOrganizationSetting(orgId)

			_, err := provider.Get(setting.OrganizationId, setting.Key)
			gomega.Expect(err).ShouldNot(gomega.Succeed())
		})
	})
	ginkgo.Context("listing settings", func() {
		ginkgo.It("should be able to return the settings of an organization", func() {

			orgId := uuid.New().String()
			numSettings := 2
			for i := 0; i < numSettings; i++ {
				setting := CreateOrganizationSetting(orgId)
				err := provider.Add(*setting)
				gomega.Expect(err).Should(gomega.Succeed())
			}

			list, err := provider.List(orgId)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(list).NotTo(gomega.BeNil())
			gomega.Expect(len(list)).Should(gomega.Equal(numSettings))

			// add other organization with settings and check when listing the settings, only its settings are being returned

			orgId2 := uuid.New().String()
			numSettings2 := 5
			for i := 0; i < numSettings2; i++ {
				setting := CreateOrganizationSetting(orgId2)
				err := provider.Add(*setting)
				gomega.Expect(err).Should(gomega.Succeed())
			}

			list2, err := provider.List(orgId2)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(list2).NotTo(gomega.BeNil())
			gomega.Expect(len(list2)).Should(gomega.Equal(numSettings2))
		})
		ginkgo.It("should be able to return an empty list of settings of an organization", func() {
			orgId := uuid.New().String()
			list, err := provider.List(orgId)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(list).NotTo(gomega.BeNil())
			gomega.Expect(len(list)).Should(gomega.Equal(0))
		})

	})
	ginkgo.Context("updating a setting", func() {
		ginkgo.It("should be able to update a setting", func() {
			orgId := uuid.New().String()
			setting := CreateOrganizationSetting(orgId)
			err := provider.Add(*setting)
			gomega.Expect(err).Should(gomega.Succeed())

			setting.Value = "New value"
			setting.Description = "New Description"

			err = provider.Update(*setting)
			gomega.Expect(err).Should(gomega.Succeed())

			retrieved, err := provider.Get(orgId, setting.Key)
			gomega.Expect(err).Should(gomega.Succeed())
			gomega.Expect(retrieved).Should(gomega.Equal(setting))

		})
		ginkgo.It("should not be able to update a non existing setting", func() {
			orgId := uuid.New().String()
			setting := CreateOrganizationSetting(orgId)
			err := provider.Update(*setting)
			gomega.Expect(err).ShouldNot(gomega.Succeed())
		})
	})
	ginkgo.Context("removing a setting", func() {
		ginkgo.It("should be able to remove a setting", func() {
			orgId := uuid.New().String()
			setting := CreateOrganizationSetting(orgId)
			err := provider.Add(*setting)
			gomega.Expect(err).Should(gomega.Succeed())

			err = provider.Remove(orgId, setting.Key)
			gomega.Expect(err).Should(gomega.Succeed())
		})
		ginkgo.It("should not be able to remove a non existing setting", func() {
			orgId := uuid.New().String()
			setting := CreateOrganizationSetting(orgId)

			err := provider.Remove(orgId, setting.Key)
			gomega.Expect(err).ShouldNot(gomega.Succeed())
		})
	})
}
