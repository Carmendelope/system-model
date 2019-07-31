/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package account

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for account
type Provider interface {
	// Add a new account to the system.
	Add(account entities.Account) derrors.Error
	// Update the information of an account.
	Update(account entities.Account) derrors.Error
	// Exists checks if an account exists on the system.
	Exists(accountID string) (bool, derrors.Error)
	// Get an account.
	Get(accountID string) (*entities.Account, derrors.Error)
	// Remove an account
	Remove(accountID string) derrors.Error
	// Clear all accounts
	Clear() derrors.Error
}
