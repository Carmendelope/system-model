/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package role

import (
	"context"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-role-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
)

// Handler structure for the role requests.
type Handler struct{
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

// AddRole adds a new role to a given organization.
func (h * Handler) AddRole(ctx context.Context, addRoleRequest *grpc_role_go.AddRoleRequest) (*grpc_role_go.Role, error) {
	log.Debug().Str("organizationID", addRoleRequest.OrganizationId).
		Str("name", addRoleRequest.Name).Msg("add role")
	vErr := entities.ValidAddRoleRequest(addRoleRequest)
	if vErr != nil{
		return nil, conversions.ToGRPCError(vErr)
	}
	added, err := h.Manager.AddRole(addRoleRequest)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("roleID", added.RoleId).Msg("role has been added")
	return added.ToGRPC(), nil
}

// GetRole returns an existing role.
func (h * Handler) GetRole(ctx context.Context, roleID *grpc_role_go.RoleId) (*grpc_role_go.Role, error) {
	vErr := entities.ValidRoleID(roleID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	role, err := h.Manager.GetRole(roleID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return role.ToGRPC(), nil
}

// ListRoles retrieves the list of roles of a given organization.
func (h * Handler) ListRoles(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_role_go.RoleList, error) {
	log.Debug().Str("organizationID", organizationID.OrganizationId).Msg("list roles")
	vErr := entities.ValidOrganizationID(organizationID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	roles, err := h.Manager.ListRoles(organizationID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	roleList := make([]*grpc_role_go.Role, 0)
	for _, r := range roles {
		roleList = append(roleList, r.ToGRPC())
	}
	result := &grpc_role_go.RoleList{
		Roles:                roleList,
	}
	return result, nil
}

// RemoveRole removes a given role from an organization.
func (h * Handler) RemoveRole(ctx context.Context, removeRoleRequest *grpc_role_go.RemoveRoleRequest) (*grpc_common_go.Success, error) {
	log.Debug().Str("organizationID", removeRoleRequest.OrganizationId).
		Str("roleID", removeRoleRequest.RoleId).Msg("remove role")
	vErr := entities.ValidRemoveRoleRequest(removeRoleRequest)
	if vErr != nil{
		return nil, conversions.ToGRPCError(vErr)
	}
	err := h.Manager.RemoveRole(removeRoleRequest)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}


