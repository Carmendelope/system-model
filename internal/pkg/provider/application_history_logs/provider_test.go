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
			gomega.Expect(err).To(gomega.Succeed())
			exists, err := provider.ExistsServiceInstanceLog(
				toAdd.OrganizationId,
				toAdd.AppInstanceId,
				toAdd.ServiceGroupInstanceId,
				toAdd.ServiceInstanceId,
			)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())
		})
	})

	//ginkgo.Context("UpdateServiceInstanceLog", func() {
	//	ginkgo.It("should be able to update a ServiceInstanceLog", func() {
	//		var toUpdate = entities.UpdateLogRequest{
	//			OrganizationId:    entities.GenerateUUID(),
	//			AppInstanceId:     entities.GenerateUUID(),
	//			ServiceInstanceId: entities.GenerateUUID(),
	//			Terminated:        entities.GenerateInt64(),
	//		}
	//		err := provider.Update(&toUpdate)
	//		gomega.Expect(err).To(gomega.Succeed())
	//		exists, err := provider.ExistsServiceInstanceLog(
	//			toUpdate.OrganizationId,
	//			toUpdate.AppInstanceId,
	//			toUpdate.Terminated,
	//			toUpdate.ServiceInstanceId,
	//		)
	//		gomega.Expect(err).To(gomega.Succeed())
	//		gomega.Expect(exists).To(gomega.BeFalse())
	//	})
	//})
}