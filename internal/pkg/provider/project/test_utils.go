/*
 * Copyright (C)  2019 Nalej - All Rights Reserved
 */

package project

import (
	"github.com/nalej/system-model/internal/pkg/entities"
	"time"
)

func CreateProject() *entities.Project {
	return &entities.Project{
		ProjectId:      entities.GenerateUUID(),
		OwnerAccountId: entities.GenerateUUID(),
		Name:           "Test project",
		Created:        time.Now().Unix(),
		State:          entities.ProjectState_Active,
		StateInfo:      "active info",
	}
}
