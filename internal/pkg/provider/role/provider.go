/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package role

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for application
type Provider interface {
	// Add a new role to the system.
	Add(role entities.Role) derrors.Error
	// Update an existing role in the system
	Update(role entities.Role) derrors.Error
	// Exists checks if a role exists on the system.
	Exists(roleID string) (bool, derrors.Error)
	// Get a role.
	Get(roleID string) (* entities.Role, derrors.Error)
	// Remove a role
	Remove(roleID string) derrors.Error
}
