/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package node

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestNodePackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Node package suite")
}
