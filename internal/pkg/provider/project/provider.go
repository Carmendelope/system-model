/*
 * Copyright 2019 Nalej
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

package project

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for project
type Provider interface {
	// Add a new project to the system.
	Add(project entities.Project) derrors.Error
	// Update the information of a project.
	Update(project entities.Project) derrors.Error
	// Exists checks if a project exists on the system.
	Exists(accountID string, projectID string) (bool, derrors.Error)
	// check if there is a project in the account with the received name
	ExistsByName(accountID string, name string) (bool, derrors.Error)
	// Get a project.
	Get(accountID string, projectID string) (*entities.Project, derrors.Error)
	// Remove a project
	Remove(accountID string, projectID string) derrors.Error
	// List all the projects of an account
	ListAccountProjects(accountID string) ([]entities.Project, derrors.Error)
	// Clear all projects
	Clear() derrors.Error
}
