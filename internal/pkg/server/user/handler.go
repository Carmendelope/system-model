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
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

// AddUser adds a new user to a given organization.
func (h *Handler) AddUser(ctx context.Context, addUserRequest *grpc_user_go.AddUserRequest) (*grpc_user_go.User, error) {
	log.Debug().Str("organizationID", addUserRequest.OrganizationId).
		Str("email", addUserRequest.Email).Msg("add user")
	vErr := entities.ValidAddUserRequest(addUserRequest)
	if vErr != nil {
		log.Error().Str("trace", vErr.DebugReport()).Msg("invalid add user request")
		return nil, conversions.ToGRPCError(vErr)
	}
	added, err := h.Manager.AddUser(addUserRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add user")
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("organizationID", addUserRequest.OrganizationId).
		Str("email", addUserRequest.Email).Msg("user has been added")
	return added.ToGRPC(), nil
}

func (h *Handler) Update(ctx context.Context, request *grpc_user_go.UpdateUserRequest) (*grpc_common_go.Success, error) {
	vErr := entities.ValidUpdateUserRequest(request)
	if vErr != nil {
		log.Error().Str("trace", vErr.DebugReport()).Msg("invalid update user request")
		return nil, conversions.ToGRPCError(vErr)
	}
	err := h.Manager.UpdateUser(request)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot update user")
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// GetUser returns an existing user.
func (h *Handler) GetUser(ctx context.Context, userID *grpc_user_go.UserId) (*grpc_user_go.User, error) {
	vErr := entities.ValidUserID(userID)
	if vErr != nil {
		log.Error().Str("trace", vErr.DebugReport()).Msg("invalid user identifier")
		return nil, conversions.ToGRPCError(vErr)
	}
	user, err := h.Manager.GetUser(userID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get user")
		return nil, conversions.ToGRPCError(err)
	}
	return user.ToGRPC(), nil
}

// GetUsers retrieves the list of users of a given organization.
func (h *Handler) GetUsers(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_user_go.UserList, error) {
	vErr := entities.ValidOrganizationID(organizationID)
	if vErr != nil {
		log.Error().Str("trace", vErr.DebugReport()).Msg("invalid organization identifier")
		return nil, conversions.ToGRPCError(vErr)
	}
	users, err := h.Manager.GetUsers(organizationID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get users")
		return nil, conversions.ToGRPCError(err)
	}
	userList := make([]*grpc_user_go.User, 0)
	for _, u := range users {
		userList = append(userList, u.ToGRPC())
	}
	result := &grpc_user_go.UserList{
		Users: userList,
	}
	return result, nil
}

// RemoveUser removes a given user from an organization.
func (h *Handler) RemoveUser(ctx context.Context, removeRequest *grpc_user_go.RemoveUserRequest) (*grpc_common_go.Success, error) {
	log.Debug().Str("organizationID", removeRequest.OrganizationId).
		Str("email", removeRequest.Email).Msg("remove user")
	vErr := entities.ValidRemoveUserRequest(removeRequest)
	if vErr != nil {
		log.Error().Str("trace", vErr.DebugReport()).Msg("invalid remove user request")
		return nil, conversions.ToGRPCError(vErr)
	}
	err := h.Manager.RemoveUser(removeRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot remove user")
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("organizationID", removeRequest.OrganizationId).
		Str("email", removeRequest.Email).Msg("user has been removed")
	return &grpc_common_go.Success{}, nil
}
