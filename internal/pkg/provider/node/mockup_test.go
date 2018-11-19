package node

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("Mockup node provider", func(){


	sp := NewMockupNodeProvider()
	RunTest(sp)

})
