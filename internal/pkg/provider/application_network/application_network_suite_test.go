/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package application_network_test

import (
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestApplicationNetwork(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "ApplicationNetwork provider Suite")
}
