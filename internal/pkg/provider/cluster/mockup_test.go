package cluster

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("Mockup Cluster provider", func() {

	sp := NewMockupClusterProvider()
	RunTest(sp)

})
