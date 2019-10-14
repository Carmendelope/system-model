package application

import (
	"github.com/google/uuid"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func RunTest(provider Provider) {

	ginkgo.AfterEach(func() {
		provider.Clear()
	})

	ginkgo.Context("Descriptor", func() {
		// AddDescriptor
		ginkgo.It("Should be able to add a descriptor", func() {

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
		ginkgo.It("Should not be able to get the descriptor", func() {
			app, err := provider.GetDescriptor(uuid.New().String())
			gomega.Expect(err).NotTo(gomega.Succeed())
			gomega.Expect(app).To(gomega.BeNil())
		})

		// DescriptorExists
		ginkgo.It("Should be able to find the descriptor", func() {

			descriptor := CreateTestApplicationDescriptor(uuid.New().String())

			// add the application
			err := provider.AddDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())

			// find it
			exists, err := provider.DescriptorExists(descriptor.AppDescriptorId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())
		})
		ginkgo.It("Should not be able to find the descriptor", func() {
			exists, err := provider.DescriptorExists(uuid.New().String())
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).NotTo(gomega.BeTrue())
		})

		ginkgo.It("should be able to update a descriptor", func() {
			descriptor := CreateTestApplicationDescriptor(uuid.New().String())
			// add the application
			err := provider.AddDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())
			// update
			descriptor.Name = "newName"
			descriptor.InboundNetInterfaces = []entities.InboundNetworkInterface{{Name: "inbound1mod"}, {Name: "inbound2mod"}}
			err = provider.UpdateDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())
			// check the update
			descriptorAux, err := provider.GetDescriptor(descriptor.AppDescriptorId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(descriptor).NotTo(gomega.BeNil())
			gomega.Expect(descriptor.Name).Should(gomega.Equal(descriptorAux.Name))
			gomega.Expect(descriptor.InboundNetInterfaces).Should(gomega.Equal(descriptorAux.InboundNetInterfaces))
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
		ginkgo.It("Should be able to find the appInstance", func() {

			app := CreateTestApplication(uuid.New().String(), uuid.New().String())

			// add the application
			err := provider.AddInstance(*app)
			gomega.Expect(err).To(gomega.Succeed())

			// find it
			exists, err := provider.InstanceExists(app.AppInstanceId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(exists).To(gomega.BeTrue())
		})
		ginkgo.It("Should not be able to find the appInstance", func() {
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
		ginkgo.It("Should not be able to get the appInstance", func() {
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

	ginkgo.Context("App EntryPoints", func() {
		ginkgo.It("should be able to add an appEndPoint", func() {
			entrypoint := CreateAppEndPoint()
			err := provider.AddAppEndpoint(*entrypoint)
			gomega.Expect(err).To(gomega.Succeed())

		})
		ginkgo.It("should be able to add an appEndPoint twice", func() {
			entrypoint := CreateAppEndPoint()
			err := provider.AddAppEndpoint(*entrypoint)
			gomega.Expect(err).To(gomega.Succeed())

			entrypoint.Protocol = entities.HTTPS
			err = provider.AddAppEndpoint(*entrypoint)
			gomega.Expect(err).To(gomega.Succeed())

		})
		ginkgo.It("should be able to get EndPoints by name", func() {
			entrypoint := CreateAppEndPoint()
			err := provider.AddAppEndpoint(*entrypoint)
			gomega.Expect(err).To(gomega.Succeed())

			retrieved, err := provider.GetAppEndpointByFQDN(entrypoint.GlobalFqdn)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).NotTo(gomega.BeEmpty())
			gomega.Expect(retrieved[0].OrganizationId).Should(gomega.Equal(entrypoint.OrganizationId))

		})
		ginkgo.It("should be able to get EndPoint list by name", func() {
			endpoint := CreateAppEndPoint()
			err := provider.AddAppEndpoint(*endpoint)
			gomega.Expect(err).To(gomega.Succeed())

			endpoint.OrganizationId = uuid.New().String()
			err = provider.AddAppEndpoint(*endpoint)
			gomega.Expect(err).To(gomega.Succeed())

			retrieved, err := provider.GetAppEndpointByFQDN(endpoint.GlobalFqdn)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(retrieved).NotTo(gomega.BeEmpty())
			gomega.Expect(len(retrieved)).Should(gomega.Equal(2))

		})
		ginkgo.It("should be able to delete an appEndpoint", func() {
			endpoint := CreateAppEndPoint()
			err := provider.AddAppEndpoint(*endpoint)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.DeleteAppEndpoints(endpoint.OrganizationId, endpoint.AppInstanceId)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("should be able to delete all the EndPoints in a application", func() {
			endpoint := CreateAppEndPoint()
			err := provider.AddAppEndpoint(*endpoint)
			gomega.Expect(err).To(gomega.Succeed())

			endpoint.ServiceInstanceId = uuid.New().String()
			err = provider.AddAppEndpoint(*endpoint)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.DeleteAppEndpoints(endpoint.OrganizationId, endpoint.AppInstanceId)
			gomega.Expect(err).To(gomega.Succeed())
		})

	})

	ginkgo.Context("Instance Parameters", func() {
		ginkgo.It("Should be able to add instance parameters", func() {

			parameters := []entities.InstanceParameter{
				{"param1", "value1"},
				{"param2", "value2"},
			}
			err := provider.AddInstanceParameters(uuid.New().String(), parameters)
			gomega.Expect(err).To(gomega.Succeed())

		})
		ginkgo.It("Should not be able to add instance parameters twice", func() {
			instanceID := uuid.New().String()
			parameters := []entities.InstanceParameter{
				{"param1", "value1"},
				{"param2", "value2"},
			}
			err := provider.AddInstanceParameters(instanceID, parameters)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.AddInstanceParameters(instanceID, parameters)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("Should be able to retrieve the params of an instance", func() {
			instanceID := uuid.New().String()
			parameters := []entities.InstanceParameter{
				{"param1", "value1"},
				{"param2", "value2"},
			}
			err := provider.AddInstanceParameters(instanceID, parameters)
			gomega.Expect(err).To(gomega.Succeed())

			params, err := provider.GetInstanceParameters(instanceID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(params).NotTo(gomega.BeNil())
			gomega.Expect(len(params)).Should(gomega.Equal(2))
		})
		ginkgo.It("Should be able to retrieve an empty list if the instance has no params", func() {
			instanceID := uuid.New().String()

			params, err := provider.GetInstanceParameters(instanceID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(params).NotTo(gomega.BeNil())
			gomega.Expect(len(params)).Should(gomega.Equal(0))
		})
		ginkgo.It("should be able to remove the params of an instance", func() {
			instanceID := uuid.New().String()
			parameters := []entities.InstanceParameter{
				{"param1", "value1"},
				{"param2", "value2"},
			}
			err := provider.AddInstanceParameters(instanceID, parameters)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.DeleteInstanceParameters(instanceID)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("should not fail when deleting the parameters of an instance (which do not exist)", func() {
			instanceID := uuid.New().String()

			err := provider.DeleteInstanceParameters(instanceID)
			gomega.Expect(err).To(gomega.Succeed())
		})
	})

	ginkgo.Context("Descriptor Parameters", func() {
		ginkgo.It("should be able to retrieves descriptor parameters", func() {
			appDescriptorID := uuid.New().String()
			descriptor := CreateApplicationDescriptorWithParameters(appDescriptorID)

			err := provider.AddDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())

			params, err := provider.GetDescriptorParameters(descriptor.AppDescriptorId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(params).NotTo(gomega.BeEmpty())
		})
		ginkgo.It("should be able to retrieves an empty list when the descriptor has no parameters", func() {
			appDescriptorID := uuid.New().String()
			descriptor := CreateTestApplicationDescriptor(appDescriptorID)
			descriptor.Parameters = nil

			err := provider.AddDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())

			params, err := provider.GetDescriptorParameters(descriptor.AppDescriptorId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(params).NotTo(gomega.BeNil())
			gomega.Expect(params).To(gomega.BeEmpty())
		})
	})

	ginkgo.Context("Parametrized Descriptor", func() {
		ginkgo.It("Should be able to add a parametrized descriptor", func() {

			descriptor := CreateParametrizedDescriptor(uuid.New().String())
			err := provider.AddParametrizedDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())
		})
		ginkgo.It("Should not be able to add a parametrized descriptor twice", func() {

			descriptor := CreateParametrizedDescriptor(uuid.New().String())
			err := provider.AddParametrizedDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.AddParametrizedDescriptor(*descriptor)
			gomega.Expect(err).NotTo(gomega.Succeed())
		})
		ginkgo.It("Should be able to get a parametrized descriptor", func() {

			descriptor := CreateParametrizedDescriptor(uuid.New().String())
			err := provider.AddParametrizedDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())

			parametrized, err := provider.GetParametrizedDescriptor(descriptor.AppInstanceId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(parametrized).NotTo(gomega.BeNil())

		})
		ginkgo.It("Should not be able to get a non-existent parametrized descriptor", func() {

			_, err := provider.GetParametrizedDescriptor(uuid.New().String())
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
		ginkgo.It("Should be able to determinate if a parametrized descriptor exists", func() {

			descriptor := CreateParametrizedDescriptor(uuid.New().String())
			err := provider.AddParametrizedDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())

			exists, err := provider.ParametrizedDescriptorExists(descriptor.AppInstanceId)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(*exists).To(gomega.BeTrue())

		})
		ginkgo.It("Should be able to determinate a parametrized descriptor does not exist", func() {

			exists, err := provider.ParametrizedDescriptorExists(uuid.New().String())
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(*exists).NotTo(gomega.BeTrue())

		})
		ginkgo.It("Should be able to delete a parametrized descriptor", func() {

			descriptor := CreateParametrizedDescriptor(uuid.New().String())
			err := provider.AddParametrizedDescriptor(*descriptor)
			gomega.Expect(err).To(gomega.Succeed())

			err = provider.DeleteParametrizedDescriptor(descriptor.AppInstanceId)
			gomega.Expect(err).To(gomega.Succeed())

		})
		ginkgo.It("Should not be able to delete a non-existent parametrized descriptor", func() {

			err := provider.DeleteParametrizedDescriptor(uuid.New().String())
			gomega.Expect(err).NotTo(gomega.Succeed())

		})
	})
}
