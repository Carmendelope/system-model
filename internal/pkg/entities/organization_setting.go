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
 *
 */

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-organization-go"
)

type OrganizationSetting struct {
	OrganizationId string `json:"organization_id"`
	Key            string `json:"key"`
	Value          string `json:"value"`
	Description    string `json:"description"`
}

func NewOrganizationSetting(organizationId string, key string, value string, description string) *OrganizationSetting {
	return &OrganizationSetting{
		OrganizationId: organizationId,
		Key:            key,
		Value:          value,
		Description:    description,
	}
}

func NewOrganizationSettingFromGRPC(addRequest *grpc_organization_go.AddSettingRequest) *OrganizationSetting{
	return NewOrganizationSetting(addRequest.OrganizationId, addRequest.Key, addRequest.Value, addRequest.Description)
}

func (o *OrganizationSetting) ToGRPC() *grpc_organization_go.OrganizationSetting {
	return &grpc_organization_go.OrganizationSetting{
		OrganizationId: o.OrganizationId,
		Key:            o.Key,
		Value:          o.Value,
		Description:    o.Description,
	}
}

func (o *OrganizationSetting) ApplyUpdate(toUpdate *grpc_organization_go.UpdateSettingRequest) {

	if toUpdate.UpdateDescription {
		o.Description = toUpdate.Description
	}
	if toUpdate.UpdateValue {
		o.Value = toUpdate.Value
	}
}

func ValidateAddSettingRequest(addRequest *grpc_organization_go.AddSettingRequest) derrors.Error{

	if addRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if addRequest.Key == "" {
		return derrors.NewInvalidArgumentError(emptyKey)
	}
	return nil
}

func ValidateUpdateSettingRequest(updateRequest *grpc_organization_go.UpdateSettingRequest) derrors.Error{

	if updateRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if updateRequest.Key == "" {
		return derrors.NewInvalidArgumentError(emptyKey)
	}
	return nil
}

func ValidateSettingKey(in *grpc_organization_go.SettingKey) derrors.Error{

	if in.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if in.Key == "" {
		return derrors.NewInvalidArgumentError(emptyKey)
	}
	return nil
}