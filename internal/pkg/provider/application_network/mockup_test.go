/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package application_network

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("Mockup Application Network provider", func() {
	provider := NewMockupApplicationNetworkProvider()
	RunTest(provider)
})
