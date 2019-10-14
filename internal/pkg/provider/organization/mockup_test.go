package organization

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("Organization Application provider", func() {

	sp := NewMockupOrganizationProvider()
	RunTest(sp)

})
