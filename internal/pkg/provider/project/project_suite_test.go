/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package project

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestProjectProviderPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Project Providers package suite")
}
