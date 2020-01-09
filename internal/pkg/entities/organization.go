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
	FullAddress string `json:"full_address"`
	City        string `json:"city"`
	State       string `json:"state"`
	Country     string `json:"country"`
	ZipCode     string `json:"zip_code"`
	PhotoBase64 string `json:"photo_base64"`
	Created     int64  `json:"created"`
}

func NewOrganization(name string, fullAddress string, city string, state string, country string, zipCode string, photo string) *Organization {
	uuid := GenerateUUID()
	return &Organization{uuid, name, fullAddress, city, state, country,
		zipCode, photo, time.Now().Unix()}
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
