/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package asset

import "github.com/onsi/ginkgo"

var _ = ginkgo.Describe("Mockup Asset provider", func() {
	provider := NewMockupAssetProvider()
	RunTest(provider)
})
