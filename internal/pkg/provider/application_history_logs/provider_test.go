package application_history_logs

import (
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest (provider Provider) {
	ginkgo.AfterEach(func() {
		_ = provider.Clear()
	})

	ginkgo.Context("AddServiceInstanceLog", func() {
		ginkgo.It("should be able to add a ServiceInstanceLog", func() {
			toAdd := entities.AddLogRequest{
				OrganizationId:         entities.GenerateUUID(),
				AppDescriptorId:        entities.GenerateUUID(),
				AppInstanceId:          entities.GenerateUUID(),
				ServiceGroupId:         entities.GenerateUUID(),
				ServiceGroupInstanceId: entities.GenerateUUID(),
				ServiceId:              entities.GenerateUUID(),
				ServiceInstanceId:      entities.GenerateUUID(),
				Created:                entities.GenerateInt64(),
			}
			err := provider.Add(&toAdd)
			gomega.Expect(err).To(gomega.BeNil())
			exists, err := provider.ExistsServiceInstanceLog(
				toAdd.OrganizationId,
				toAdd.AppInstanceId,
				toAdd.ServiceGroupInstanceId,
				toAdd.ServiceInstanceId,
			)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(exists).To(gomega.BeTrue())
		})
	})

	ginkgo.Context("UpdateServiceInstanceLog", func() {
		ginkgo.It("should be able to update a ServiceInstanceLog", func() {
			toAdd := entities.AddLogRequest{
				OrganizationId:         entities.GenerateUUID(),
				AppDescriptorId:        entities.GenerateUUID(),
				AppInstanceId:          entities.GenerateUUID(),
				ServiceGroupId:         entities.GenerateUUID(),
				ServiceGroupInstanceId: entities.GenerateUUID(),
				ServiceId:              entities.GenerateUUID(),
				ServiceInstanceId:      entities.GenerateUUID(),
				Created:                entities.GenerateInt64(),
			}
			err := provider.Add(&toAdd)
			gomega.Expect(err).To(gomega.BeNil())
			exists, err := provider.ExistsServiceInstanceLog(
				toAdd.OrganizationId,
				toAdd.AppInstanceId,
				toAdd.ServiceGroupInstanceId,
				toAdd.ServiceInstanceId,
			)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(exists).To(gomega.BeTrue())


			toUpdate := entities.UpdateLogRequest{
				OrganizationId:    toAdd.OrganizationId,
				AppInstanceId:     toAdd.AppInstanceId,
				ServiceInstanceId: toAdd.ServiceInstanceId,
				Terminated:        toAdd.Created + 100,
			}
			err = provider.Update(&toUpdate)
			gomega.Expect(err).To(gomega.BeNil())
			exists, err = provider.ExistsServiceInstanceLog(
				toAdd.OrganizationId,
				toAdd.AppInstanceId,
				toAdd.ServiceGroupInstanceId,
				toAdd.ServiceInstanceId,
			)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(exists).To(gomega.BeTrue())
		})
	})

	ginkgo.Context("SearchServiceInstanceLog", func() {
		ginkgo.It("should be able to search for some ServiceInstanceLogs", func() {
			toAdd := entities.AddLogRequest{
				OrganizationId:         entities.GenerateUUID(),
				AppDescriptorId:        entities.GenerateUUID(),
				AppInstanceId:          entities.GenerateUUID(),
				ServiceGroupId:         entities.GenerateUUID(),
				ServiceGroupInstanceId: entities.GenerateUUID(),
				ServiceId:              entities.GenerateUUID(),
				ServiceInstanceId:      entities.GenerateUUID(),
				Created:                entities.GenerateInt64(),
			}
			err := provider.Add(&toAdd)
			gomega.Expect(err).To(gomega.BeNil())
			exists, err := provider.ExistsServiceInstanceLog(
				toAdd.OrganizationId,
				toAdd.AppInstanceId,
				toAdd.ServiceGroupInstanceId,
				toAdd.ServiceInstanceId,
			)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(exists).To(gomega.BeTrue())

			toSearch := entities.SearchLogsRequest{
				OrganizationId: toAdd.OrganizationId,
				From:           toAdd.Created - 100,
				To:             toAdd.Created + 100,
			}
			err, logResponse := provider.Search(&toSearch)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(logResponse.OrganizationId).To(gomega.Equal(toAdd.OrganizationId))
		})
	})

	ginkgo.Context("RemoveServiceInstanceLog", func() {
		ginkgo.It("should be able to remove a ServiceInstanceLog", func() {
			toAdd := entities.AddLogRequest{
				OrganizationId:         entities.GenerateUUID(),
				AppDescriptorId:        entities.GenerateUUID(),
				AppInstanceId:          entities.GenerateUUID(),
				ServiceGroupId:         entities.GenerateUUID(),
				ServiceGroupInstanceId: entities.GenerateUUID(),
				ServiceId:              entities.GenerateUUID(),
				ServiceInstanceId:      entities.GenerateUUID(),
				Created:                entities.GenerateInt64(),
			}
			err := provider.Add(&toAdd)
			gomega.Expect(err).To(gomega.BeNil())
			exists, err := provider.ExistsServiceInstanceLog(
				toAdd.OrganizationId,
				toAdd.AppInstanceId,
				toAdd.ServiceGroupInstanceId,
				toAdd.ServiceInstanceId,
			)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(exists).To(gomega.BeTrue())

			toRemove := entities.RemoveLogRequest{
				OrganizationId: toAdd.OrganizationId,
				AppInstanceId:  toAdd.AppInstanceId,
			}
			err = provider.Remove(&toRemove)
			gomega.Expect(err).To(gomega.BeNil())
			exists, err = provider.ExistsServiceInstanceLog(
				toAdd.OrganizationId,
				toAdd.AppInstanceId,
				toAdd.ServiceGroupInstanceId,
				toAdd.ServiceInstanceId,
			)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(exists).To(gomega.BeFalse())
		})
	})
}