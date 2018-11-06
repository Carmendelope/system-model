/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package user

import (
	"context"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-user-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
)

// Handler structure for the user requests.
type Handler struct{
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

// AddUser adds a new user to a given organization.
func (h*Handler) AddUser(ctx context.Context, addUserRequest *grpc_user_go.AddUserRequest) (*grpc_user_go.User, error){
	log.Debug().Str("organizationID", addUserRequest.OrganizationId).Str("roleID", addUserRequest.RoleId).
		Str("email", addUserRequest.Email).Msg("add user")
	vErr := entities.ValidAddUserRequest(addUserRequest)
	if vErr != nil{
		return nil, conversions.ToGRPCError(vErr)
	}
	added, err := h.Manager.AddUser(addUserRequest)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("organizationID", addUserRequest.OrganizationId).
		Str("email", addUserRequest.Email).Msg("user has been added")
	return added.ToGRPC(), nil
}

// GetUser returns an existing user.
func (h*Handler) GetUser(ctx context.Context, userID *grpc_user_go.UserId) (*grpc_user_go.User, error){
	vErr := entities.ValidUserID(userID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	user, err := h.Manager.GetUser(userID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return user.ToGRPC(), nil
}

// GetUsers retrieves the list of users of a given organization.
func (h*Handler) GetUsers(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_user_go.UserList, error){
	log.Debug().Str("organizationID", organizationID.OrganizationId).Msg("list users")
	vErr := entities.ValidOrganizationID(organizationID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	users, err := h.Manager.GetUsers(organizationID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	userList := make([]*grpc_user_go.User, 0)
	for _, u := range users {
		userList = append(userList, u.ToGRPC())
	}
	result := &grpc_user_go.UserList{
		Users:                userList,
	}
	return result, nil
}

// RemoveUser removes a given user from an organization.
func (h*Handler) RemoveUser(ctx context.Context, removeRequest *grpc_user_go.RemoveUserRequest) (*grpc_common_go.Success, error){
	log.Debug().Str("organizationID", removeRequest.OrganizationId).
		Str("email", removeRequest.Email).Msg("remove user")
	vErr := entities.ValidRemoveUserRequest(removeRequest)
	if vErr != nil{
		return nil, conversions.ToGRPCError(vErr)
	}
	err := h.Manager.RemoveUser(removeRequest)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("organizationID", removeRequest.OrganizationId).
		Str("email", removeRequest.Email).Msg("user has been removed")
	return &grpc_common_go.Success{}, nil
}
