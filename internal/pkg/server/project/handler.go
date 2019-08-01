/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package project

import (
	"context"
	"github.com/nalej/grpc-account-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-project-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}

// AddProject adds a new project to a given account
func (h *Handler)AddProject(ctx context.Context, request *grpc_project_go.AddProjectRequest) (*grpc_project_go.Project, error){
	log.Debug().Str("account_id", request.AccountId).Str("Name", request.Name).Msg("add project")
	err := entities.ValidateAddProjectRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	added, err := h.Manager.AddProject(request)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("Name", request.Name).Str("account_id", added.OwnerAccountId).
		Str("project_id", added.ProjectId).Msg("account added")
	return added.ToGRPC(), nil
}

// GetProject retrieves a given project
func (h *Handler)GetProject(ctx context.Context, request *grpc_project_go.ProjectId) (*grpc_project_go.Project, error){

	err := entities.ValidateProjectId(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	account, err := h.Manager.GetProject(request)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	return account.ToGRPC(), nil

}
// RemoveProject removes a given project
func (h *Handler)RemoveProject(ctx context.Context, request *grpc_project_go.ProjectId) (*grpc_common_go.Success, error){
	err := entities.ValidateProjectId(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.RemoveProject(request)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}
// ListAccountProjects list the projects of a given account
func (h *Handler)ListAccountProjects(ctx context.Context, request *grpc_account_go.AccountId) (*grpc_project_go.ProjectList, error){

	err := entities.ValidateAccountId(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	projects, err := h.Manager.ListAccountProjects(request)
	if err != nil{
		return nil, conversions.ToGRPCError(err)
	}

	list :=  make ([]*grpc_project_go.Project, 0)
	for _, project := range projects {
		list = append(list, project.ToGRPC())
	}

	return &grpc_project_go.ProjectList{Projects:list}, nil
}
// UpdateProject updates the project information
func (h *Handler)UpdateProject(ctx context.Context, request *grpc_project_go.UpdateProjectRequest) (*grpc_common_go.Success, error){

	err := entities.ValidateUpdateProjectRequest(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.UpdateProject(request)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}