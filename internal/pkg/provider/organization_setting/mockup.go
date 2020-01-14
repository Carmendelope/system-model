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
	"sync"
)

type MockupOrganizationSettingProvider struct {
	sync.Mutex
	settings map[string][]entities.OrganizationSetting

}

func NewMockupOrganizationSettingProvider() *MockupOrganizationSettingProvider{
	return &MockupOrganizationSettingProvider{
		settings: make(map[string][]entities.OrganizationSetting, 0),
	}
}

func (m *MockupOrganizationSettingProvider) unsafeExists(organizationId string, key string) bool {
	settings, exists := m.settings[organizationId]
	if ! exists {
		return false
	}
	for _, setting := range settings {
		if setting.Key == key {
			return true
		}
	}
	return false
}

// Add a new setting for an organization.
func (m *MockupOrganizationSettingProvider) Add(setting entities.OrganizationSetting) derrors.Error{

	m.Lock()
	defer m.Unlock()

	// check if already exists
	exists := m.unsafeExists(setting.OrganizationId, setting.Key)
	if exists {
		return derrors.NewAlreadyExistsError("setting").WithParams(setting.OrganizationId, setting.Key)
	}
	// check if the organization has settings
	settings, exists := m.settings[setting.OrganizationId]
	if !exists {
		m.settings[setting.OrganizationId] = []entities.OrganizationSetting{setting}
	}else{
		m.settings[setting.OrganizationId] = append(settings, setting)
	}

	return nil
}
// Check if a setting is defined for an organization
func (m *MockupOrganizationSettingProvider) Exists(organizationID string, key string) (bool, derrors.Error){
	m.Lock()
	defer m.Unlock()
	return m.unsafeExists(organizationID, key), nil
}
// Get a setting organization.
func (m *MockupOrganizationSettingProvider) Get(organizationID string, key string) (*entities.OrganizationSetting, derrors.Error) {
	m.Lock()
	defer m.Unlock()

	settings, exists := m.settings[organizationID]
	if ! exists {
		return nil, derrors.NewNotFoundError("setting").WithParams(organizationID, key)
	}
	for _, setting := range settings {
		if setting.Key == key {
			return &setting, nil
		}
	}
	return nil, derrors.NewNotFoundError("setting").WithParams(organizationID, key)

}
// List all the settings of an organization.
func (m *MockupOrganizationSettingProvider) List(organizationID string) ([]entities.OrganizationSetting, derrors.Error){
	settings, exists := m.settings[organizationID]
	if ! exists {
		emptyList := make ([]entities.OrganizationSetting, 0)
		return emptyList, nil
	}
	return settings, nil
}
// Update a setting of an organization
func (m *MockupOrganizationSettingProvider) Update(setting entities.OrganizationSetting) derrors.Error {
	m.Lock()
	defer m.Unlock()

	settings, exists := m.settings[setting.OrganizationId]
	if ! exists {
		return  derrors.NewNotFoundError("setting").WithParams(setting.OrganizationId, setting.Key)
	}
	for i:=0; i< len(settings); i++ {
		if settings[i].Key == setting.Key {
			settings[i] = setting
			return nil
		}
	}
	return derrors.NewNotFoundError("setting").WithParams(setting.OrganizationId, setting.Key)
}

func (m *MockupOrganizationSettingProvider)	Remove(organizationID string, key string) derrors.Error{
	m.Lock()
	defer m.Unlock()

	settings, exists := m.settings[organizationID]
	if ! exists {
		return  derrors.NewNotFoundError("setting").WithParams(organizationID, key)
	}

	newSettingsList := make ([]entities.OrganizationSetting, 0)
	found := false
	for _, set := range settings {
		if set.Key != key {
			newSettingsList = append(newSettingsList, set)
		}else{
			found = true
		}

	}
	if !found {
		return derrors.NewNotFoundError("setting").WithParams(organizationID, key)
	}else{
		if len(newSettingsList) == 0 {
			delete(m.settings, organizationID)
		}else{
			m.settings[organizationID] = newSettingsList
		}
	}
	return nil
}

func (m *MockupOrganizationSettingProvider)	Clear() derrors.Error{
	m.settings =  make(map[string][]entities.OrganizationSetting, 0)
	return nil
}
