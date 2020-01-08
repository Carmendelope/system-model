/*
 * Copyright 2020 Nalej
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

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-account-go"
	"time"
)

type AccountState int32

const (
	AccountState_Active AccountState = iota + 1
	AccountState_Deactivated
)

var AccountStateToGRPC = map[AccountState]grpc_account_go.AccountState{
	AccountState_Active:      grpc_account_go.AccountState_ACTIVE,
	AccountState_Deactivated: grpc_account_go.AccountState_DEACTIVATED,
}
var AccountStateFromGRPC = map[grpc_account_go.AccountState]AccountState{
	grpc_account_go.AccountState_ACTIVE:      AccountState_Active,
	grpc_account_go.AccountState_DEACTIVATED: AccountState_Deactivated,
}

// AccountBillingInfo with the billing information of an account
type AccountBillingInfo struct {
	// AccountId with the account identifier
	AccountId string `json:"account_id,omitempty" cql:"account_id"`
	// FullName with the name of a person that receives the invoice
	FullName string `json:"full_name,omitempty" cql:"full_name"`
	// CompanyName with the name of the company
	CompanyName string `json:"company_name,omitempty" cql:"company_name"`
	// Address with the address of the company
	Address string `json:"address,omitempty" cql:"address"`
	// AdditionalInfo with more information of the company
	AdditionalInfo string `json:"additional_info,omitempty" cql:"additional_info"`
}

func NewAccountBillingInfoFromGRPC(accountId string, fullName string, companyName string, address string, additionalInfo string) *AccountBillingInfo {

	return &AccountBillingInfo{
		AccountId:      accountId,
		FullName:       fullName,
		CompanyName:    companyName,
		Address:        address,
		AdditionalInfo: additionalInfo,
	}
}

func (abi *AccountBillingInfo) ToGRPC() *grpc_account_go.AccountBillingInfo {
	return &grpc_account_go.AccountBillingInfo{
		AccountId:      abi.AccountId,
		FullName:       abi.FullName,
		CompanyName:    abi.CompanyName,
		Address:        abi.Address,
		AdditionalInfo: abi.AdditionalInfo,
	}
}

// Account model with the information related to a given account
type Account struct {
	// AccountId with the account identifier. This value is created by the system
	AccountId string `json:"account_id,omitempty"`
	// Name of the account. The name of the account must be unique
	Name string `json:"name,omitempty"`
	// Created timestamp
	Created int64 `json:"created,omitempty"`
	// BillingInfo with the billing information of the account
	BillingInfo *AccountBillingInfo `json:"billing_info,omitempty"`
	// State with the state of the account
	State AccountState `json:"state,omitempty"`
	// StateInfo in case the account is in a non active state,
	// it contains the information about the reason for this state
	StateInfo string `json:"state_info,omitempty"`
}

func NewAccountFromGRPC(account *grpc_account_go.AddAccountRequest) *Account {
	if account == nil {
		return nil
	}

	id := GenerateUUID()
	return &Account{
		AccountId:   id,
		Name:        account.Name,
		Created:     time.Now().Unix(),
		BillingInfo: NewAccountBillingInfoFromGRPC(id, account.FullName, account.CompanyName, account.Address, account.AdditionalInfo),
		State:       AccountState_Active,
		StateInfo:   "",
	}
}

func (a *Account) ToGRPC() *grpc_account_go.Account {
	if a == nil {
		return nil
	}
	return &grpc_account_go.Account{
		AccountId:   a.AccountId,
		Name:        a.Name,
		Created:     a.Created,
		BillingInfo: a.BillingInfo.ToGRPC(),
		State:       AccountStateToGRPC[a.State],
		StateInfo:   a.StateInfo,
	}
}

// -------------------
// apply update
// -------------------
func (a *Account) ApplyUpdate(update *grpc_account_go.UpdateAccountRequest) {

	if update.UpdateName {
		a.Name = update.Name
	}
	if update.UpdateState {
		a.State = AccountStateFromGRPC[update.State]
	}
	if update.UpdateStateInfo {
		a.StateInfo = update.StateInfo
	}

}

func (a *Account) ApplyUpdateBillingInfo(update *grpc_account_go.UpdateAccountBillingInfoRequest) {
	if a.BillingInfo == nil {
		a.BillingInfo = &AccountBillingInfo{}
	}

	if update.UpdateFullName {
		a.BillingInfo.FullName = update.FullName
	}
	if update.UpdateCompanyName {
		a.BillingInfo.CompanyName = update.CompanyName
	}
	if update.UpdateAddress {
		a.BillingInfo.Address = update.Address
	}
	if update.UpdateAdditionalInfo {
		a.BillingInfo.AdditionalInfo = update.AdditionalInfo
	}
}

// -------------------
// validation methods
// -------------------
func ValidateAddAccountRequest(request *grpc_account_go.AddAccountRequest) derrors.Error {

	if request.Name == "" {
		return derrors.NewInvalidArgumentError(emptyName)
	}
	return nil
}

func ValidateAccountId(request *grpc_account_go.AccountId) derrors.Error {
	if request.AccountId == "" {
		return derrors.NewInvalidArgumentError(emptyAccountId)
	}
	return nil
}

func ValidateUpdateAccountRequest(request *grpc_account_go.UpdateAccountRequest) derrors.Error {
	if request.AccountId == "" {
		return derrors.NewInvalidArgumentError(emptyAccountId)
	}
	return nil
}

func ValidateUpdateAccountBillingInfoRequest(request *grpc_account_go.UpdateAccountBillingInfoRequest) derrors.Error {
	if request.AccountId == "" {
		return derrors.NewInvalidArgumentError(emptyAccountId)
	}
	return nil
}
