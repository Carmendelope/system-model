/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/grpc-account-go"
	"time"
)

type AccountState int32

const (
	AccountState_Active AccountState = iota + 1
	AccountState_Deactivated
)

var AccountStateGRPC = map[AccountState]grpc_account_go.AccountState{
	AccountState_Active: grpc_account_go.AccountState_ACTIVE,
	AccountState_Deactivated : grpc_account_go.AccountState_DEACTIVATED,
}
var AccountStateFromGRPC = map[grpc_account_go.AccountState]AccountState {
	grpc_account_go.AccountState_ACTIVE: AccountState_Active,
	grpc_account_go.AccountState_DEACTIVATED: AccountState_Deactivated,
}

// AccountBillingInfo with the billing information of an account
type AccountBillingInfo struct {
	// AccountId with the account identifier
	AccountId string	`json:"account_id,omitempty" cql:"account_id"`
	// FullName with the name of a person that receives the invoice
	FullName string 	`json:"full_name,omitempty" cql:"full_name"`
	// CompanyName with the name of the company
	CompanyName string 	`json:"company_name,omitempty" cql:"company_name"`
	// Address with the address of the company
	Address string 		`json:"address,omitempty" cql:"address"`
	// AdditionalInfo with more information of the company
	AdditionalInfo string `json:"additional_info,omitempty" cql:"additional_info"`
}

func NewAccountBillingInfoFromGRPC(accountId string, info *grpc_account_go.AccountBillingInfo) *AccountBillingInfo {
	if info == nil {
		return  nil
	}

	return &AccountBillingInfo{
		AccountId:	accountId,
		FullName: 	info.FullName,
		CompanyName:info.CompanyName,
		Address: 	info.Address,
		AdditionalInfo: info.AdditionalInfo,
	}
}

func (abi *AccountBillingInfo) ToGRPC () *grpc_account_go.AccountBillingInfo {
	return &grpc_account_go.AccountBillingInfo{
		AccountId:	abi.AccountId,
		FullName: 	abi.FullName,
		CompanyName:abi.CompanyName,
		Address: 	abi.Address,
		AdditionalInfo: abi.AdditionalInfo,
	}
}

// Account model with the information related to a given account
type Account struct {
	// AccountId with the account identifier. This value is created by the system
	AccountId string 	`json:"account_id,omitempty"`
	// Name of the account. The name of the account must be unique
	Name string 		`json:"name,omitempty"`
	// Created timestamp
	Created int64 		`json:"created,omitempty"`
	// BillingInfo with the billing information of the account
	BillingInfo *AccountBillingInfo `json:"billing_info,omitempty"`
	// State with the state of the account
	State AccountState 	`json:"state,omitempty"`
	// StateInfo in case the account is in a non active state,
	// it contains the information about the reason for this state
	StateInfo string  	`json:"state_info,omitempty"`
}

func NewAccountFromGRPC (account *grpc_account_go.Account) *Account{
	if account == nil {
		return nil
	}

	id := GenerateUUID()
	return &Account{
		AccountId:		id,
		Name: 			account.Name,
		Created: 		time.Now().Unix(),
		BillingInfo: 	NewAccountBillingInfoFromGRPC(id, account.BillingInfo),
		State: 			AccountStateFromGRPC[account.State],
		StateInfo: 		account.StateInfo,
	}
}