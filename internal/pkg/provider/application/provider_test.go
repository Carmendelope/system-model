package application

import (
	"github.com/google/uuid"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {

	ginkgo.BeforeEach(func() {
		provider.Clear()
	})

	ginkgo.Context("Descriptor", func() {
		// AddDescriptor
		ginkgo.FIt("Should be able to add a descriptor", func() {

			descriptor := CreateTestApplicationDescriptor(uuid.New().String())

			err := provider.AddDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())
		})

		// GetDescriptors
		ginkgo.It("Should be able to get the Descriptor", func() {

			descriptorId := uuid.New().String()

			descriptor := CreateTestApplicationDescriptor(descriptorId)

			// add the application
			err := provider.AddDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())

			// get it
			descriptor, err = provider.GetDescriptor(descriptor.AppDescriptorId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(descriptor).NotTo(gomega.BeNil())
		})
		ginkgo.It ("Should not be able to get the descriptor", func (){
			app, err := provider.GetDescriptor(uuid.New().String())
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(app).To(gomega.BeNil())
		})

		// DescriptorExists
		ginkgo.It("Should be able to find the descriptor", func(){

			descriptor := CreateTestApplicationDescriptor(uuid.New().String())

			// add the application
			err := provider.AddDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())

			// find it
			exists, err := provider.DescriptorExists(descriptor.AppDescriptorId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())
		})
		ginkgo.It("Should not be able to find the descriptor", func(){
			exists, err := provider.DescriptorExists(uuid.New().String())
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())
		})

		ginkgo.It("should be able to update a descriptor", func(){
			descriptor := CreateTestApplicationDescriptor(uuid.New().String())
			// add the application
			err := provider.AddDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())
			// update
			descriptor.Name = "newName"
			err = provider.UpdateDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())
			// check the update
			descriptor, err = provider.GetDescriptor(descriptor.AppDescriptorId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(descriptor).NotTo(gomega.BeNil())
			gomega.Expect(descriptor.Name).Should(gomega.Equal(descriptor.Name))
		})

		// DeleteDescriptor
		ginkgo.It("Should be able to remove the descriptor", func() {

			descriptor := CreateTestApplicationDescriptor(uuid.New().String())

			// add the application
			err := provider.AddDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())

			// delete it
			err = provider.DeleteDescriptor(descriptor.AppDescriptorId)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("Should not be able to remove the descriptor", func() {
			err := provider.DeleteDescriptor(uuid.New().String())
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})

	// ---------------------------------------------------------------------------------------------------------------------

	ginkgo.Context("Instance", func() {
		// Add Application Instance
		ginkgo.It("Should be able to add an application", func() {

			app := CreateTestApplication(uuid.New().String(), uuid.New().String())

			err := provider.AddInstance(*app)
			gomega.Expect(err).To(gomega.Succeed())

		})

		// Update Application Instance
		ginkgo.It("Should be able to update an application", func() {
			app := CreateTestApplication(uuid.New().String(), uuid.New().String())

			err := provider.AddInstance(*app)
			gomega.Expect(err).To(gomega.Succeed())

			app.Status = entities.Deploying
			err = provider.UpdateInstance(*app)
			gomega.Expect(err).To(gomega.Succeed())

			recovered, err := provider.GetInstance(app.AppInstanceId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(recovered).NotTo(gomega.BeNil())
			gomega.Expect(recovered.Status).Should(gomega.Equal(entities.Deploying))

		})
		ginkgo.It("Should not be able to update an application", func() {
			app := CreateTestApplication(uuid.New().String(), uuid.New().String())
			err := provider.UpdateInstance(*app)
			gomega.Expect(err).NotTo(gomega.Succeed())

		})

		// ExistsInstance
		ginkgo.It("Should be able to find the appInstance", func(){

			app := CreateTestApplication(uuid.New().String(), uuid.New().String())

			// add the application
			err := provider.AddInstance(*app)
			gomega.Expect(err).To(gomega.Succeed())

			// find it
			exists, err := provider.InstanceExists(app.AppInstanceId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())
		})
		ginkgo.It("Should not be able to find the appInstance", func(){
			exists, err := provider.InstanceExists("application instance")
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())
		})

		// 	GetInstance
		ginkgo.It("Should be able to get the appInstance", func() {

			app := CreateTestApplication(uuid.New().String(), uuid.New().String())

			// add the application
			err := provider.AddInstance(*app)
			gomega.Expect(err).To(gomega.Succeed())

			// get it
			app, err = provider.GetInstance(app.AppInstanceId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(app).NotTo(gomega.BeNil())
		})
		ginkgo.It ("Should not be able to get the appInstance", func (){
			app, err := provider.GetInstance("application instance")
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(app).To(gomega.BeNil())
		})

		// DeleteInstance
		ginkgo.It("Should be able to remove the appInstance", func() {

			app := CreateTestApplication(uuid.New().String(), uuid.New().String())

			// add the application
			err := provider.AddInstance(*app)
			gomega.Expect(err).To(gomega.Succeed())

			// delete it
			err = provider.DeleteInstance(app.AppInstanceId)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("Should not be able to remove the appInstance", func() {
			err := provider.DeleteInstance("application instance")
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
	})


}
