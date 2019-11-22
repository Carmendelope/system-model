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
			toAddA := entities.AddLogRequest{
				OrganizationId:         entities.GenerateUUID(),
				AppDescriptorId:        entities.GenerateUUID(),
				AppInstanceId:          entities.GenerateUUID(),
				ServiceGroupId:         entities.GenerateUUID(),
				ServiceGroupInstanceId: entities.GenerateUUID(),
				ServiceId:              entities.GenerateUUID(),
				ServiceInstanceId:      entities.GenerateUUID(),
				Created:                entities.GenerateInt64(),
			}
			err := provider.Add(&toAddA)
			gomega.Expect(err).To(gomega.BeNil())
			exists, err := provider.ExistsServiceInstanceLog(
				toAddA.OrganizationId,
				toAddA.AppInstanceId,
				toAddA.ServiceGroupInstanceId,
				toAddA.ServiceInstanceId,
			)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(exists).To(gomega.BeTrue())

			toUpdateA := entities.UpdateLogRequest{
				OrganizationId:    toAddA.OrganizationId,
				AppInstanceId:     toAddA.AppInstanceId,
				ServiceInstanceId: toAddA.ServiceInstanceId,
				Terminated:        toAddA.Created + 100,
			}
			err = provider.Update(&toUpdateA)
			gomega.Expect(err).To(gomega.BeNil())
			exists, err = provider.ExistsServiceInstanceLog(
				toAddA.OrganizationId,
				toAddA.AppInstanceId,
				toAddA.ServiceGroupInstanceId,
				toAddA.ServiceInstanceId,
			)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(exists).To(gomega.BeTrue())

			Query0 := entities.SearchLogsRequest{
				OrganizationId: toAddA.OrganizationId,
				From:           toAddA.Created + 25,
				To:             toAddA.Created + 75,
			}
			err, logResponse := provider.Search(&Query0)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(logResponse.OrganizationId).To(gomega.Equal(toAddA.OrganizationId))

			Query1 := entities.SearchLogsRequest{
				OrganizationId: toAddA.OrganizationId,
				From:           toAddA.Created - 100,
				To:             toAddA.Created + 200,
			}
			err, logResponse = provider.Search(&Query1)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(logResponse.OrganizationId).To(gomega.Equal(toAddA.OrganizationId))

			Query2 := entities.SearchLogsRequest{
				OrganizationId: toAddA.OrganizationId,
				From:           toAddA.Created + 50,
				To:             toAddA.Created + 200,
			}
			err, logResponse = provider.Search(&Query2)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(logResponse.OrganizationId).To(gomega.Equal(toAddA.OrganizationId))

			Query3 := entities.SearchLogsRequest{
				OrganizationId: toAddA.OrganizationId,
				From:           toAddA.Created - 100,
				To:             toAddA.Created + 50,
			}
			err, logResponse = provider.Search(&Query3)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(logResponse.OrganizationId).To(gomega.Equal(toAddA.OrganizationId))

			Query4 := entities.SearchLogsRequest{
				OrganizationId: toAddA.OrganizationId,
				From:           toAddA.Created - 100,
				To:             toAddA.Created - 50,
			}
			_, logResponse = provider.Search(&Query4)
			gomega.Expect(logResponse).To(gomega.BeNil())

			Query5 := entities.SearchLogsRequest{
				OrganizationId: toAddA.OrganizationId,
				From:           toAddA.Created + 200,
				To:             toAddA.Created + 300,
			}
			_, logResponse = provider.Search(&Query5)
			gomega.Expect(logResponse).To(gomega.BeNil())

			_ = provider.Clear()

			toAddB := entities.AddLogRequest{
				OrganizationId:         entities.GenerateUUID(),
				AppDescriptorId:        entities.GenerateUUID(),
				AppInstanceId:          entities.GenerateUUID(),
				ServiceGroupId:         entities.GenerateUUID(),
				ServiceGroupInstanceId: entities.GenerateUUID(),
				ServiceId:              entities.GenerateUUID(),
				ServiceInstanceId:      entities.GenerateUUID(),
				Created:                entities.GenerateInt64(),
			}
			err = provider.Add(&toAddB)
			exists, err = provider.ExistsServiceInstanceLog(
				toAddA.OrganizationId,
				toAddA.AppInstanceId,
				toAddA.ServiceGroupInstanceId,
				toAddA.ServiceInstanceId,
			)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(exists).To(gomega.BeFalse())

			Query6 := entities.SearchLogsRequest{
				OrganizationId: toAddB.OrganizationId,
				From:           toAddB.Created - 100,
				To:             toAddB.Created + 100,
			}
			err, logResponse = provider.Search(&Query6)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(logResponse.OrganizationId).To(gomega.Equal(toAddB.OrganizationId))

			Query7 := entities.SearchLogsRequest{
				OrganizationId: toAddB.OrganizationId,
				From:           toAddB.Created + 50,
				To:             toAddB.Created + 100,
			}
			err, logResponse = provider.Search(&Query7)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(logResponse.OrganizationId).To(gomega.Equal(toAddB.OrganizationId))

			Query8 := entities.SearchLogsRequest{
				OrganizationId: toAddB.OrganizationId,
				From:           toAddB.Created - 100,
				To:             toAddB.Created - 50,
			}
			err, logResponse = provider.Search(&Query8)
			gomega.Expect(logResponse).To(gomega.BeNil())
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