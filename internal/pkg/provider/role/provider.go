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
 *
 */

package role

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for application
type Provider interface {
	// Add a new role to the system.
	Add(role entities.Role) derrors.Error
	// Update an existing role in the system
	Update(role entities.Role) derrors.Error
	// Exists checks if a role exists on the system.
	Exists(roleID string) (bool, derrors.Error)
	// Get a role.
	Get(roleID string) (*entities.Role, derrors.Error)
	// Remove a role
	Remove(roleID string) derrors.Error
	//clear roles
	Clear() derrors.Error
}
