package application

import (
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {

	ginkgo.BeforeEach(func() {
		provider.Clear()
	})

	// AddDescriptor
	ginkgo.It("Should be able to add a descriptor", func() {

		descriptor := CreateTestApplicationDescriptor("XX001")

		err := provider.AddDescriptor(*descriptor)
		gomega.Expect(err).To(gomega.Succeed())
	})

	// GetDescriptors
	ginkgo.It("Should be able to get the Descriptor", func() {

		descriptor := CreateTestApplicationDescriptor("xx0001")

		// add the application
		err := provider.AddDescriptor(*descriptor)
		gomega.Expect(err).To(gomega.Succeed())

		// get it
		descriptor, err = provider.GetDescriptor(descriptor.AppDescriptorId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(descriptor).NotTo(gomega.BeNil())
	})
	ginkgo.It ("Should not be able to get the descriptor", func (){
		app, err := provider.GetDescriptor("xx0001")
		gomega.Expect(err).NotTo(gomega.Succeed())
		gomega.Expect(app).To(gomega.BeNil())
	})

	// DescriptorExists
	ginkgo.It("Should be able to find the descriptor", func(){

		descriptor := CreateTestApplicationDescriptor("xx0001")

		// add the application
		err := provider.AddDescriptor(*descriptor)
		gomega.Expect(err).To(gomega.Succeed())

		// find it
		exists, err := provider.DescriptorExists(descriptor.AppDescriptorId)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).To(gomega.BeTrue())
	})
	ginkgo.It("Should not be able to find the descriptor", func(){
		exists, err := provider.DescriptorExists("xx001")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(exists).NotTo(gomega.BeTrue())
	})

	// DeleteDescriptor
	ginkgo.It("Should be able to remove the descriptor", func() {

		descriptor := CreateTestApplicationDescriptor("xx0001")

		// add the application
		err := provider.AddDescriptor(*descriptor)
		gomega.Expect(err).To(gomega.Succeed())

		// delete it
		err = provider.DeleteDescriptor(descriptor.AppDescriptorId)
		gomega.Expect(err).To(gomega.Succeed())
	})
	ginkgo.It("Should not be able to remove the descriptor", func() {
		err := provider.DeleteDescriptor("xx0001")
		gomega.Expect(err).NotTo(gomega.Succeed())
	})

	// ---------------------------------------------------------------------------------------------------------------------

	// Add Application Instance
	ginkgo.It("Should be able to add an application", func() {

		app := CreateTestApplication("0001")

		err := provider.AddInstance(*app)
		gomega.Expect(err).To(gomega.Succeed())

	})

	// Update Application Instance
	ginkgo.It("Should be able to udpate an application", func() {

		app := CreateTestApplication("0001")

		// add the application
		err := provider.AddInstance(*app)
		gomega.Expect(err).To(gomega.Succeed())

		// modify some fields
		groups := make ([]entities.ServiceGroupInstance, 0)
		groups = append(groups, CreateTestServiceGroupInstance("XXXXX"))
		app.Groups = groups

		// and update it
		err = provider.UpdateInstance(*app)
		gomega.Expect(err).To(gomega.Succeed())

	})
	ginkgo.It("Should not be able to udpate an application", func() {

		app := CreateTestApplication("0001")

		// and update it
		err := provider.UpdateInstance(*app)
		gomega.Expect(err).NotTo(gomega.Succeed())

	})

	// ExistsInstance
	ginkgo.It("Should be able to find the appInstance", func(){

		app := CreateTestApplication("0001")

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

		app := CreateTestApplication("0001")

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

		app := CreateTestApplication("0001")

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

}
