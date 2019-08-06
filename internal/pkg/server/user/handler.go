/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package user

import (
	"context"
	"github.com/nalej/grpc-account-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-project-go"
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
		return nil, conversions.ToGRPCError(vErr)
	}
	added, err := h.Manager.AddUser(addUserRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("organizationID", addUserRequest.OrganizationId).
		Str("email", addUserRequest.Email).Msg("user has been added")
	return added.ToGRPC(), nil
}

func (h *Handler) Update(ctx context.Context, request *grpc_user_go.UpdateUserRequest) (*grpc_common_go.Success, error) {
	vErr := entities.ValidUpdateUserRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	err := h.Manager.UpdateUser(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// GetUser returns an existing user.
func (h *Handler) GetUser(ctx context.Context, userID *grpc_user_go.UserId) (*grpc_user_go.User, error) {
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
func (h *Handler) GetUsers(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_user_go.UserList, error) {
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
		return nil, conversions.ToGRPCError(vErr)
	}
	err := h.Manager.RemoveUser(removeRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("organizationID", removeRequest.OrganizationId).
		Str("email", removeRequest.Email).Msg("user has been removed")
	return &grpc_common_go.Success{}, nil
}

func (h *Handler) UpdateContactInfo(ctx context.Context, request *grpc_user_go.UpdateContactInfoRequest) (*grpc_common_go.Success, error){

	vErr := entities.ValidUpdateContactInfoRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	err := h.Manager.UpdateContactInfo(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (h *Handler) AddAccountUser(ctx context.Context, request *grpc_user_go.AddAccountUserRequest) (*grpc_user_go.AccountUser, error){

	vErr := entities.ValidAddAccountUserRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	accountUser, err := h.Manager.AddAccountUser(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return accountUser.ToGRPC(), nil

}
func (h *Handler) RemoveAccountUser(ctx context.Context, request *grpc_user_go.AccountUserId) (*grpc_common_go.Success, error){
	vErr := entities.ValidAccountUserId(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	err := h.Manager.RemoveAccountUser(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}
func (h *Handler) UpdateAccountUser(ctx context.Context, request *grpc_user_go.AccountUserUpdateRequest) (*grpc_user_go.AccountUser, error){
	vErr := entities.ValidAccountUserUpdateRequest(request)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}

	accountUser, err := h.Manager.UpdateAccountUser(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return accountUser.ToGRPC(), nil

}
func (h *Handler) ListAccountsUser(ctx context.Context, in *grpc_account_go.AccountId) (*grpc_user_go.AccountUserList, error){return nil, nil}
// ---------------------------------------------------------------------------------------------------------------------
func (h *Handler) AddAccountUserInvite(ctx context.Context, in *grpc_user_go.AddAccountInviteRequest) (*grpc_user_go.AccountUserInvite, error){return nil, nil}
func (h *Handler) GetAccountUserInvite(ctx context.Context, in *grpc_user_go.AccountUserInviteId) (*grpc_user_go.AccountUserInvite, error){return nil, nil}
func (h *Handler) RemoveAccountUserInvite(ctx context.Context, in *grpc_user_go.AccountUserInviteId) (*grpc_common_go.Success, error){return nil, nil}
func (h *Handler) ListAccountUserInvites(ctx context.Context, in *grpc_user_go.UserId) (*grpc_user_go.AccountInviteList, error){return nil, nil}
// ---------------------------------------------------------------------------------------------------------------------
func (h *Handler) AddProjectUser(ctx context.Context, in *grpc_user_go.AddProjectUserRequest) (*grpc_user_go.ProjectUser, error){return nil, nil}
func (h *Handler) RemoveProjectUser(ctx context.Context, in *grpc_user_go.ProjectUserId) (*grpc_common_go.Success, error){return nil, nil}
func (h *Handler) UpdateProjectUser(ctx context.Context, in *grpc_user_go.ProjectUserUpdateRequest) (*grpc_user_go.ProjectUser, error){return nil, nil}
func (h *Handler) ListProjectsUser(ctx context.Context, in *grpc_project_go.ProjectId) (*grpc_user_go.ProjectUserList, error){return nil, nil}
// ---------------------------------------------------------------------------------------------------------------------
func (h *Handler) AddProjectUserInvite(ctx context.Context, in *grpc_user_go.AddProjectInviteRequest) (*grpc_user_go.ProjectUserInvite, error){return nil, nil}
func (h *Handler) GetProjectUserInvite(ctx context.Context, in *grpc_user_go.ProjectUserInviteId) (*grpc_user_go.ProjectUserInvite, error){return nil, nil}
func (h *Handler) RemoveProjectUserInvite(ctx context.Context, in *grpc_user_go.ProjectUserInviteId) (*grpc_common_go.Success, error){return nil, nil}
func (h *Handler) ListProjectUserInvites(ctx context.Context, in *grpc_user_go.UserId) (*grpc_user_go.ProjectInviteList, error){return nil, nil}
