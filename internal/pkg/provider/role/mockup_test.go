package role

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("Mockup role provider", func() {

	sp := NewMockupRoleProvider()
	RunTest(sp)

})
