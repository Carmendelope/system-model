package user

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("Mockup user provider", func(){

	sp := NewMockupUserProvider()
	RunTest(sp)

})