/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package account

import (
	"github.com/nalej/grpc-account-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

// Handler structure for the account requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}


// AddAccount adds a new account in the system
// Once the account is added, it will be active to be able to operate in it
func (h *Handler)AddAccount(ctx context.Context, request *grpc_account_go.AddAccountRequest) (*grpc_account_go.Account, error){
	log.Debug().Str("Name", request.Name).Msg("add account")
	err := entities.ValidateAddAccountRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}

	added, err := h.Manager.AddAccount(request)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("Name", request.Name).Str("account_id", added.AccountId).Msg("account added")
	return added.ToGRPC(), nil
}

// GetAccount retrieves a given account
func (h *Handler)GetAccount(ctx context.Context, request *grpc_account_go.AccountId) (*grpc_account_go.Account, error){

	err := entities.ValidateAccountId(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	account, err := h.Manager.GetAccount(request)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	return account.ToGRPC(), nil
	}

// ListAccounts retrieves a list of all the accounts in the system. This method is only intended to be used by
// management API as the users will not be able to list other accounts
func (h *Handler)ListAccounts(ctx context.Context, request *grpc_common_go.Empty) (*grpc_account_go.AccountList, error){

	accounts, err := h.Manager.ListAccounts()
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	list :=  make ([]*grpc_account_go.Account, 0)

	for _, account := range accounts {
		list = append(list, account.ToGRPC())
	}

	return &grpc_account_go.AccountList{Accounts:list}, nil
}

// UpdateAccount updates the information of an account
func (h *Handler)UpdateAccount(ctx context.Context, request *grpc_account_go.UpdateAccountRequest) (*grpc_common_go.Success, error){

	err := entities.ValidateUpdateAccountRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.UpdateAccount(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// UpdateAccountBillingInfo updates the billing info of an account
func (h *Handler)UpdateAccountBillingInfo(ctx context.Context, request *grpc_account_go.UpdateAccountBillingInfoRequest) (*grpc_common_go.Success, error){

	err := entities.ValidateUpdateAccountBillingInfoRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.UpdateAccountBillingInfo(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}