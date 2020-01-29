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
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-organization-go"
	"time"
)

type Organization struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	FullAddress string `json:"full_address"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	ZipCode     string `json:"zip_code"`
	PhotoBase64 string `json:"photo_base64"`
	Created     int64  `json:"created"`
}

func NewOrganization(name string, email string, fullAddress string, city string, state string, country string, zipCode string, photo string) *Organization {
	uuid := GenerateUUID()
	return &Organization{uuid, name, email, fullAddress, city, state, country,
		zipCode, photo, time.Now().Unix()}
}

func (o *Organization) String() string {
	return fmt.Sprintf("%#v", o)
}

func (o *Organization) ToGRPC() *grpc_organization_go.Organization {
	return &grpc_organization_go.Organization{
		OrganizationId: o.ID,
		Name:           o.Name,
		Email:          o.Email,
		FullAddress:    o.FullAddress,
		City:           o.City,
		Country:        o.Country,
		State:          o.State,
		ZipCode:        o.ZipCode,
		PhotoBase64:    o.PhotoBase64,
		Created:        o.Created,
	}
}
func (o *Organization) ApplyUpdate(toUpdate *grpc_organization_go.UpdateOrganizationRequest) {

	if toUpdate.UpdateName {
		o.Name = toUpdate.Name
	}
	if toUpdate.UpdateEmail {
		o.Email = toUpdate.Email
	}
	if toUpdate.UpdateFullAddress {
		o.FullAddress = toUpdate.FullAddress
	}
	if toUpdate.UpdateCity {
		o.City = toUpdate.City
	}
	if toUpdate.UpdateCountry {
		o.Country = toUpdate.Country
	}
	if toUpdate.UpdateState {
		o.State = toUpdate.State
	}
	if toUpdate.UpdatePhoto {
		o.PhotoBase64 = toUpdate.PhotoBase64
	}
	if toUpdate.UpdateZipCode {
		o.ZipCode = toUpdate.ZipCode
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
	if toUpdate.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if toUpdate.UpdateName && toUpdate.Name == "" {
		return derrors.NewInvalidArgumentError(emptyName)

	}
	return nil
}

type OrganizationCluster struct {
	OrganizationId string `json:"organization_id"`
	ClusterId      string `json:"cluster_id"`
}

func NewOrganizationCluster(org string, cluster string) *OrganizationCluster {
	return &OrganizationCluster{org, cluster}
}

type OrganizationNode struct {
	OrganizationId string `json:"organization_id"`
	NodeId         string `json:"node_id"`
}

func NewOrganizationNode(org string, node string) *OrganizationNode {
	return &OrganizationNode{org, node}
}

type OrganizationDescriptor struct {
	OrganizationId  string `json:"organization_id"`
	AppDescriptorId string `json:"app_descriptor_id"`
}

func NewOrganizationDescriptor(org string, appDescriptorID string) *OrganizationDescriptor {
	return &OrganizationDescriptor{org, appDescriptorID}
}

type OrganizationInstance struct {
	OrganizationId string `json:"organization_id"`
	AppInstanceId  string `json:"app_instance_id"`
}

func NewOrganizationInstance(org string, appInstanceID string) *OrganizationInstance {
	return &OrganizationInstance{org, appInstanceID}
}

type OrganizationUser struct {
	OrganizationId string `json:"organization_id"`
	Email          string `json:"email"`
}

func NewOrganizationUser(org string, email string) *OrganizationUser {
	return &OrganizationUser{org, email}
}

type OrganizationRole struct {
	OrganizationId string `json:"organization_id"`
	RoleId         string `json:"role_id"`
}

func NewOrganizationRole(org string, roleId string) *OrganizationRole {
	return &OrganizationRole{org, roleId}
}
