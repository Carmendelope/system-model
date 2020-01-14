/*
 * Copyright 2020 Nalej
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

package organization_setting

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

type Provider interface {
	// Add a new setting for an organization.
	Add(setting entities.OrganizationSetting) derrors.Error
	// Check if a setting is defined for an organization
	Exists(organizationID string, key string) (bool, derrors.Error)
	// Get a setting organization.
	Get(organizationID string, key string) (*entities.OrganizationSetting, derrors.Error)
	// List all the settings of an organization.
	List(organizationID string) ([]entities.OrganizationSetting, derrors.Error)
	// Update a setting of an organization
	Update(setting entities.OrganizationSetting) derrors.Error
	// Remove deletes a given setting.
	Remove(organizationID string, key string) derrors.Error

	Clear() derrors.Error
}