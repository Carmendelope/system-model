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
	"github.com/nalej/grpc-account-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/account"
)

type Manager struct {
	AccountProvider account.Provider
}

func NewManager(accProvider account.Provider) Manager {
	return Manager{
		AccountProvider: accProvider,
	}
}

func (m *Manager) AddAccount(request *grpc_account_go.AddAccountRequest) (*entities.Account, derrors.Error) {

	// Check if there is an account with the same name (no repeated names are allowed)
	exists, err := m.AccountProvider.ExistsByName(request.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, derrors.NewInvalidArgumentError("Account Name already exists").WithParams(request.Name)
	}

	// add the account
	toAdd := entities.NewAccountFromGRPC(request)
	err = m.AccountProvider.Add(*toAdd)
	if err != nil {
		return nil, err
	}

	return toAdd, nil
}

func (m *Manager) GetAccount(request *grpc_account_go.AccountId) (*entities.Account, derrors.Error) {

	account, err := m.AccountProvider.Get(request.AccountId)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (m *Manager) ListAccounts() ([]entities.Account, derrors.Error) {

	return m.AccountProvider.List()

}

func (m *Manager) UpdateAccount(request *grpc_account_go.UpdateAccountRequest) derrors.Error {

	oldAccount, err := m.AccountProvider.Get(request.AccountId)
	if err != nil {
		return err
	}
	// if the name is being udpated, we need to confirm there is no other account with this name
	if request.UpdateName {
		exists, err := m.AccountProvider.ExistsByName(request.Name)
		if err != nil {
			return err
		}
		if exists {
			return derrors.NewInvalidArgumentError("Account Name already exists").WithParams(request.Name)
		}
	}
	oldAccount.ApplyUpdate(request)

	return m.AccountProvider.Update(*oldAccount)

}

func (m *Manager) UpdateAccountBillingInfo(request *grpc_account_go.UpdateAccountBillingInfoRequest) derrors.Error {

	oldAccount, err := m.AccountProvider.Get(request.AccountId)
	if err != nil {
		return err
	}
	oldAccount.ApplyUpdateBillingInfo(request)

	return m.AccountProvider.Update(*oldAccount)
}
