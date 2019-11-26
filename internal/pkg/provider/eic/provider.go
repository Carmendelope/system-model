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

package eic

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

// Provider for application
type Provider interface {
	// Add a new edge controller to the system.
	Add(eic entities.EdgeController) derrors.Error
	// Update the information of an edge controller.
	Update(eic entities.EdgeController) derrors.Error
	// Exists checks if an EIC exists on the system.
	Exists(edgeControllerID string) (bool, derrors.Error)
	// Get an EIC.
	Get(edgeControllerID string) (*entities.EdgeController, derrors.Error)
	// List the EIC in a given organization
	List(organizationID string) ([]entities.EdgeController, derrors.Error)
	// Remove an EIC
	Remove(edgeControllerID string) derrors.Error
	// Clear all assets
	Clear() derrors.Error
}
