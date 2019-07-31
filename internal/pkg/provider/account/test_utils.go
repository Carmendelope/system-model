/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package account

import (
	"github.com/nalej/system-model/internal/pkg/entities"
	"time"
)

func CreateBillingInfo(accountID string) *entities.AccountBillingInfo{
	return &entities.AccountBillingInfo{
		AccountId: accountID,
		FullName: "user1",
		CompanyName: "company test",
		Address: "Address 10",
		AdditionalInfo: "info",
	}
}

func CreateAccount() *entities.Account{
	id := entities.GenerateUUID()
	return &entities.Account{
		AccountId: id,
		Name: "account test",
		Created: time.Now().Unix(),
		BillingInfo: CreateBillingInfo(id),
		State: entities.Active,
		StateInfo: "active info",
	}
}