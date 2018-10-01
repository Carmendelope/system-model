/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/satori/go.uuid"
)

const OrganizationPrefix = "o-"
const ClusterPrefix = "c-"
const NodePrefix = "n-"
const DevicePrefix = "d-"
const AppDescPrefix = "app-"
const AppInstPrefix = "inst-"
const SecurityRulePrefix = "sr-"
const ServicePrefix = "s-"
const ServiceGroupPrefix = "sg-"
const ConfigFilePrefix = "cf-"

// GenerateUUID generates a new UUID.
//   params:
//     prefix The UUID prefix.
//   returns:
//     A new UUID.
func GenerateUUID(prefix string) string {
	return prefix + uuid.NewV4().String()
}
