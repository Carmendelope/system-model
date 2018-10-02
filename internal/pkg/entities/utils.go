/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/satori/go.uuid"
)

// GenerateUUID generates a new UUID.
func GenerateUUID() string {
	return uuid.NewV4().String()
}
