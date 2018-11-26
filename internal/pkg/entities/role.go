/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-role-go"
	"time"
)

type Role struct {
	OrganizationId string `json:"organization_id,omitempty"`
	RoleId         string `json:"role_id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	Internal       bool   `json:"internal"`
	Created        int64  `json:"created,omitempty"`
}

func NewRoleFromGRPC(addRoleRequest *grpc_role_go.AddRoleRequest) *Role {
	uuid := GenerateUUID()
	return &Role{
		OrganizationId: addRoleRequest.OrganizationId,
		RoleId:         uuid,
		Name:           addRoleRequest.Name,
		Description:    addRoleRequest.Description,
		Internal:       addRoleRequest.Internal,
		Created:        time.Now().Unix(),
	}
}

func (r *Role) ToGRPC() *grpc_role_go.Role {
	return &grpc_role_go.Role{
		OrganizationId: r.OrganizationId,
		RoleId:         r.RoleId,
		Name:           r.Name,
		Description:    r.Description,
		Internal:       r.Internal,
		Created:        r.Created,
	}
}

func ValidAddRoleRequest(addRoleRequest *grpc_role_go.AddRoleRequest) derrors.Error {
	if addRoleRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if addRoleRequest.Name == "" {
		return derrors.NewInvalidArgumentError(emptyName)
	}
	return nil
}

func ValidRemoveRoleRequest(removeRoleRequest *grpc_role_go.RemoveRoleRequest) derrors.Error {
	if removeRoleRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if removeRoleRequest.RoleId == "" {
		return derrors.NewInvalidArgumentError(emptyRoleId)
	}
	return nil
}

func ValidRoleID(roleID *grpc_role_go.RoleId) derrors.Error {
	if roleID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError(emptyOrganizationId)
	}
	if roleID.RoleId == "" {
		return derrors.NewInvalidArgumentError(emptyRoleId)
	}
	return nil
}
