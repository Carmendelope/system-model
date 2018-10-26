/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package testhelpers

import (
	"fmt"
	orgProvider "github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"math/rand"
)

func CreateOrganization(orgProvider orgProvider.Provider) * entities.Organization {
	toAdd := entities.NewOrganization(fmt.Sprintf("org-%d-%d", ginkgo.GinkgoRandomSeed(), rand.Int()))
	err := orgProvider.Add(*toAdd)
	gomega.Expect(err).To(gomega.Succeed())
	return toAdd
}

