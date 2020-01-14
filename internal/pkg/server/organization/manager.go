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
 */

package organization

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
	"github.com/nalej/system-model/internal/pkg/provider/organization_setting"
)

// Manager structure with the required providers for organization operations.
type Manager struct {
	Provider organization.Provider
	SettingProvider organization_setting.Provider
}

// NewManager creates a Manager using a set of providers.
func NewManager(provider organization.Provider, settingProvider organization_setting.Provider) Manager {
	return Manager{Provider:provider, SettingProvider:settingProvider}
}

// AddOrganization adds a new organization to the system.
func (m *Manager) AddOrganization(toAdd grpc_organization_go.AddOrganizationRequest) (*entities.Organization, derrors.Error) {
	newOrg := entities.NewOrganization(toAdd.Name, toAdd.FullAddress, toAdd.City,
		toAdd.State , toAdd.Country, toAdd.ZipCode, toAdd.PhotoBase64)

	exists, err := m.Provider.ExistsByName(newOrg.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, derrors.NewAlreadyExistsError(newOrg.Name)
	}

	err = m.Provider.Add(*newOrg)
	if err != nil {
		return nil, err
	}
	return newOrg, nil
}

// GetOrganization retrieves the profile information of a given organization.
func (m *Manager) GetOrganization(orgID grpc_organization_go.OrganizationId) (*entities.Organization, derrors.Error) {
	return m.Provider.Get(orgID.OrganizationId)
}

// ListOrganization retrieves the profile information of a given organization.
func (m *Manager) ListOrganization() ([]entities.Organization, derrors.Error) {
	return m.Provider.List()
}

func (m *Manager) UpdateOrganization (newOrg *grpc_organization_go.UpdateOrganizationRequest) derrors.Error {

	org, err := m.Provider.Get(newOrg.OrganizationId)
	if err != nil {
		return err
	}
	// if the name is going to be updated to a different name..
	if newOrg.UpdateName && newOrg.Name != org.Name {
		// check if there is an organization with the new name
		exists, err := m.Provider.ExistsByName(newOrg.Name)
		if err != nil {
			return err
		}
		if exists {
			return derrors.NewAlreadyExistsError("name").WithParams(newOrg.Name)
		}
	}


	org.ApplyUpdate(newOrg)
	return m.Provider.Update(*org)

}

// AddSetting adds a new setting for the organization
func (m *Manager) AddSetting(addRequest *grpc_organization_go.AddSettingRequest) (*entities.OrganizationSetting, derrors.Error){

	// check if the organization exists
	exists, err := m.Provider.Exists(addRequest.OrganizationId)
	if err != nil {
		return nil, err
	}
	if ! exists {
		return nil, derrors.NewNotFoundError("organization").WithParams(addRequest.OrganizationId)
	}
	setting := entities.NewOrganizationSettingFromGRPC(addRequest)
	err = m.SettingProvider.Add(*setting)
	if err != nil {
		return nil, err
	}
	return setting, nil
}

// GetSetting returns an OrganizationSetting
func (m *Manager) GetSetting(in *grpc_organization_go.SettingKey) (*entities.OrganizationSetting,  derrors.Error){
	// check if the organization exists
	exists, err := m.Provider.Exists(in.OrganizationId)
	if err != nil {
		return nil, err
	}
	if ! exists {
		return nil, derrors.NewNotFoundError("organization").WithParams(in.OrganizationId)
	}

	return m.SettingProvider.Get(in.OrganizationId, in.Key)

}
// ListSettings returns a list of settings of an organization
func (m *Manager) ListSettings(in *grpc_organization_go.OrganizationId) ([]entities.OrganizationSetting,  derrors.Error){
	// check if the organization exists
	exists, err := m.Provider.Exists(in.OrganizationId)
	if err != nil {
		return nil, err
	}
	if ! exists {
		return nil, derrors.NewNotFoundError("organization").WithParams(in.OrganizationId)
	}

	return m.SettingProvider.List(in.OrganizationId)
}
// UpdateSetting update the value and/or the description of a setting
func (m *Manager) UpdateSetting(updateRequest *grpc_organization_go.UpdateSettingRequest) derrors.Error{

	setting, err := m.SettingProvider.Get(updateRequest.OrganizationId, updateRequest.Key)
	if err != nil {
		return err
	}
	setting.ApplyUpdate(updateRequest)
	return m.SettingProvider.Update(*setting)
}
// RemoveSetting removes a given setting of an organization
func (m *Manager) RemoveSetting(key *grpc_organization_go.SettingKey) derrors.Error{
	// check if the organization exists
	exists, err := m.Provider.Exists(key.OrganizationId)
	if err != nil {
		return  err
	}
	if ! exists {
		return derrors.NewNotFoundError("organization").WithParams(key.OrganizationId)
	}

	return m.SettingProvider.Remove(key.OrganizationId, key.Key)
}