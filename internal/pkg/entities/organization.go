/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
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
	result := make([] *grpc_organization_go.Organization, 0, len(list))
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
