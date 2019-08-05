/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */
package account

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestAccountPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Account package suite")
}