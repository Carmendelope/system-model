/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package project

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("Mockup Project provider", func(){
	provider := NewMockupProjectProvider()
	RunTest(provider)
})

