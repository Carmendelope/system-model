/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package user

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for application
type Provider interface {
	// Add a new user to the system.
	Add(user entities.User) derrors.Error
	// Update an existing user in the system
	Update(user entities.User) derrors.Error
	// Exists checks if a user exists on the system.
	Exists(email string) bool
	// Get a user.
	Get(email string) (* entities.User, derrors.Error)
	// Remove a user.
	Remove(email string) derrors.Error
}
