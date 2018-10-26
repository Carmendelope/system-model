/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package role

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"testing"
)

func TestRolePackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Role package suite")
}
