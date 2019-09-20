/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package application

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("Mockup Application provider", func() {

	sp := NewMockupApplicationProvider()
	RunTest(sp)

})
