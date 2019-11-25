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

package organization

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
)

// Manager structure with the required providers for organization operations.
type Manager struct {
	Provider organization.Provider
}

// NewManager creates a Manager using a set of providers.
func NewManager(provider organization.Provider) Manager {
	return Manager{provider}
}

// AddOrganization adds a new organization to the system.
func (m *Manager) AddOrganization(toAdd grpc_organization_go.AddOrganizationRequest) (*entities.Organization, derrors.Error) {
	newOrg := entities.NewOrganization(toAdd.Name)

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
