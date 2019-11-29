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

package entities

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-organization-go"
	"time"
)

type Organization struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Created int64  `json:"created"`
}

func NewOrganization(name string) *Organization {
	uuid := GenerateUUID()
	return &Organization{uuid, name, time.Now().Unix()}
}

func (o *Organization) String() string {
	return fmt.Sprintf("%#v", o)
}

func (o *Organization) ToGRPC() *grpc_organization_go.Organization {
	return &grpc_organization_go.Organization{
		OrganizationId: o.ID,
		Name:           o.Name,
		Created:        o.Created,
	}
}

func OrganizationListToGRPC(list []Organization) *grpc_organization_go.OrganizationList {
	result := make([]*grpc_organization_go.Organization, 0, len(list))
	for _, el := range list {
		result = append(result, el.ToGRPC())
	}
	return &grpc_organization_go.OrganizationList{Organizations: result}
}

func ValidOrganizationID(organizationID *grpc_organization_go.OrganizationId) derrors.Error {
	if organizationID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	return nil
}

func ValidAddOrganizationRequest(toAdd *grpc_organization_go.AddOrganizationRequest) derrors.Error {
	if toAdd.Name != "" {
		return nil
	}
	return derrors.NewInvalidArgumentError("organization required fields missing")
}

func ValidUpdateOrganization(toUpdate *grpc_organization_go.UpdateOrganizationRequest) derrors.Error {
	return nil
}
