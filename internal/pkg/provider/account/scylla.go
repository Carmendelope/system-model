/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package account

import (
	"github.com/nalej/derrors"
	"github.com/nalej/scylladb-utils/pkg/scylladb"
	"github.com/nalej/system-model/internal/pkg/entities"
	"sync"
)

const AccountTable = "Account"
const AccountTablePK = "account_id"
var allAccountColumns = []string{"account_id", "name", "created", "billing_info", "state", "state_info"}
var allAccountColumnsNoPK = []string{"name", "created", "billing_info", "state", "state_info"}

type ScyllaAccountProvider struct {
	scylladb.ScyllaDB
	sync.Mutex
}

func NewScyllaAccountProvider(address string, port int, keyspace string) * ScyllaAccountProvider{
	provider := ScyllaAccountProvider{
		ScyllaDB : scylladb.ScyllaDB{
			Address: address,
			Port : port,
			Keyspace: keyspace,
		},
	}
	provider.Connect()
	return &provider
}

// disconnect from the database
func (sp *ScyllaAccountProvider) Disconnect() {
	sp.Lock()
	defer sp.Unlock()
	sp.ScyllaDB.Disconnect()
}

// Add a new account to the system.
func (sp *ScyllaAccountProvider) Add(account entities.Account) derrors.Error{
	sp.Lock()
	defer sp.Unlock()
	return sp.UnsafeAdd(AccountTable, AccountTablePK, account.AccountId, allAccountColumns, account)
}

// Update the information of an account.
func (sp *ScyllaAccountProvider) Update(account entities.Account) derrors.Error{
	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeUpdate(AccountTable, AccountTablePK, account.AccountId, allAccountColumnsNoPK, account)
}

// Exists checks if an account exists on the system.
func (sp *ScyllaAccountProvider) Exists(accountID string) (bool, derrors.Error){
	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeGenericExist(AccountTable, AccountTablePK, accountID)
}

// Get an account.
func (sp *ScyllaAccountProvider) Get(accountID string) (*entities.Account, derrors.Error){
	sp.Lock()
	defer sp.Unlock()

	var account interface{} = &entities.Account{}

	err := sp.UnsafeGet(AccountTable, AccountTablePK, accountID, allAccountColumns, &account)
	if err != nil {
		return nil, err
	}
	return account.(*entities.Account), nil
}

// Remove an account
func (sp *ScyllaAccountProvider) Remove(accountID string) derrors.Error{
	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeRemove(AccountTable, AccountTablePK, accountID)
}

// Clear all accounts
func (sp *ScyllaAccountProvider) Clear() derrors.Error{
	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeClear([]string{AccountTable})
}