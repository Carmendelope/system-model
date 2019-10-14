/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package eic

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("Mockup EIC provider", func() {
	provider := NewMockupEICProvider()
	RunTest(provider)
})
