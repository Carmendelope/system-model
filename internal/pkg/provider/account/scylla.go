/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package account

import (
	"github.com/nalej/derrors"
	"github.com/nalej/scylladb-utils/pkg/scylladb"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
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

func NewScyllaAccountProvider(address string, port int, keyspace string) *ScyllaAccountProvider {
	provider := ScyllaAccountProvider{
		ScyllaDB: scylladb.ScyllaDB{
			Address:  address,
			Port:     port,
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
func (sp *ScyllaAccountProvider) Add(account entities.Account) derrors.Error {
	sp.Lock()
	defer sp.Unlock()
	return sp.UnsafeAdd(AccountTable, AccountTablePK, account.AccountId, allAccountColumns, account)
}

// Update the information of an account.
func (sp *ScyllaAccountProvider) Update(account entities.Account) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeUpdate(AccountTable, AccountTablePK, account.AccountId, allAccountColumnsNoPK, account)
}

// Exists checks if an account exists on the system.
func (sp *ScyllaAccountProvider) Exists(accountID string) (bool, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeGenericExist(AccountTable, AccountTablePK, accountID)
}

func (sp *ScyllaAccountProvider) ExistsByName(accountName string) (bool, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeGenericExist(AccountTable, "name", accountName)
}

// Get an account.
func (sp *ScyllaAccountProvider) Get(accountID string) (*entities.Account, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	var account interface{} = &entities.Account{}

	err := sp.UnsafeGet(AccountTable, AccountTablePK, accountID, allAccountColumns, &account)
	if err != nil {
		return nil, err
	}
	return account.(*entities.Account), nil
}

func (sp *ScyllaAccountProvider) List() ([]entities.Account, derrors.Error) {
	sp.Lock()
	defer sp.Unlock()

	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select(AccountTable).Columns(allAccountColumns...).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names)

	accounts := make([]entities.Account, 0)
	cqlErr := gocqlx.Select(&accounts, q.Query)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list accounts")
	}

	return accounts, nil

}

// Remove an account
func (sp *ScyllaAccountProvider) Remove(accountID string) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeRemove(AccountTable, AccountTablePK, accountID)
}

// Clear all accounts
func (sp *ScyllaAccountProvider) Clear() derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeClear([]string{AccountTable})
}
