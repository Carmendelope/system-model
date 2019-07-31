/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */
package entities

import "github.com/nalej/grpc-project-go"

type ProjectState int32

const (
	ProjectState_Active ProjectState = iota + 1
	ProjectState_Deactivated
)

var ProjectStateToGRPC = map[ProjectState]grpc_project_go.ProjectState{
	ProjectState_Active: grpc_project_go.ProjectState_ACTIVE,
	ProjectState_Deactivated : grpc_project_go.ProjectState_DEACTIVATED,
}
var ProjectStateFromGRPC = map[grpc_project_go.ProjectState]ProjectState {
	grpc_project_go.ProjectState_ACTIVE: ProjectState_Active,
	grpc_project_go.ProjectState_DEACTIVATED: ProjectState_Deactivated,
}

type Project struct {
	// ProjectId with the project identifier
	ProjectId string `json:"project_id,omitempty"`
	// OwnerAccountId with the account identifier of the owner of the project
	OwnerAccountId string `json:"owner_account_id,omitempty"`
	// Name with the name of the project
	Name string `json:"name,omitempty"`
	// Created timestamp
	Created int64 `json:"created,omitempty"`
	// State with the state of the project
	State ProjectState `json:"state,omitempty"`
	// StateInfo in case the project is in a non active state,
	// it contains the information about the reason for this state
	StateInfo            string   `json:"state_info,omitempty"`
}

func NewProjectToGRPC(project *grpc_project_go.Project) *Project {
	if project == nil {
		return nil
	}
	return &Project{
		ProjectId:		GenerateUUID(),
		OwnerAccountId: project.OwnerAccountId,
		Name: 			project.Name,
		Created: 		project.Created,
		State: 			ProjectStateFromGRPC[project.State],
		StateInfo: 		project.StateInfo,
	}
}

func (p *Project) ToGRPC() *grpc_project_go.Project {
	if p == nil {
		return nil
	}
	return &grpc_project_go.Project{
		ProjectId: p.ProjectId,
		OwnerAccountId: p.OwnerAccountId,
		Name: p.Name,
		Created: p.Created,
		State: ProjectStateToGRPC[p.State],
		StateInfo: p.StateInfo,
	}
}