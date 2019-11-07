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

package user

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for application
type Provider interface {
	// Add a new user to the system.
	Add(user entities.User) derrors.Error
	// Update an existing user in the system
	Update(user entities.User) derrors.Error
	// Exists checks if a user exists on the system.
	Exists(email string) (bool, derrors.Error)
	// Get a user.
	Get(email string) (*entities.User, derrors.Error)
	// Remove a user.
	Remove(email string) derrors.Error
	// Clear
	Clear() derrors.Error
}
