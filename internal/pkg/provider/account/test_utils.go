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
	"github.com/nalej/system-model/internal/pkg/entities"
	"time"
)

func CreateBillingInfo(accountID string) *entities.AccountBillingInfo {
	return &entities.AccountBillingInfo{
		AccountId:      accountID,
		FullName:       "user1",
		CompanyName:    "company test",
		Address:        "Address 10",
		AdditionalInfo: "info",
	}
}

func CreateAccount() *entities.Account {
	id := entities.GenerateUUID()
	return &entities.Account{
		AccountId:   id,
		Name:        "account test",
		Created:     time.Now().Unix(),
		BillingInfo: CreateBillingInfo(id),
		State:       entities.AccountState_Active,
		StateInfo:   "active info",
	}
}
