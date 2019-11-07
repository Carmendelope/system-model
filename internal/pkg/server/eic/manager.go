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

package eic

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-inventory-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/eic"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
)

// Manager structure with the required providers for application operations.
type Manager struct {
	ControllerProvider eic.Provider
	OrgProvider        organization.Provider
}

// NewManager creates a Manager using a set of providers.
func NewManager(controllerProvider eic.Provider, orgProvider organization.Provider) Manager {
	return Manager{controllerProvider, orgProvider}
}

func (m *Manager) Add(request *grpc_inventory_go.AddEdgeControllerRequest) (*entities.EdgeController, derrors.Error) {
	exists, err := m.OrgProvider.Exists(request.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(request.OrganizationId)
	}

	toAdd := entities.NewEdgeControllerFromGRPC(request)
	err = m.ControllerProvider.Add(*toAdd)
	if err != nil {
		return nil, err
	}
	return toAdd, nil
}

func (m *Manager) List(organizationID *grpc_organization_go.OrganizationId) ([]entities.EdgeController, derrors.Error) {
	exists, err := m.OrgProvider.Exists(organizationID.OrganizationId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organizationID").WithParams(organizationID.OrganizationId)
	}
	controllers, err := m.ControllerProvider.List(organizationID.OrganizationId)
	if err != nil {
		return nil, err
	}
	return controllers, nil
}

func (m *Manager) Remove(edgeControllerID *grpc_inventory_go.EdgeControllerId) derrors.Error {
	exists, err := m.OrgProvider.Exists(edgeControllerID.OrganizationId)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organizationID").WithParams(edgeControllerID.OrganizationId)
	}
	retrieved, err := m.ControllerProvider.Get(edgeControllerID.EdgeControllerId)
	if err != nil {
		return err
	}
	if retrieved.OrganizationId != edgeControllerID.OrganizationId {
		return derrors.NewNotFoundError("organization_id & asset_id").WithParams(edgeControllerID.OrganizationId, edgeControllerID.EdgeControllerId)
	}
	return m.ControllerProvider.Remove(edgeControllerID.EdgeControllerId)
}

func (m *Manager) Update(request *grpc_inventory_go.UpdateEdgeControllerRequest) (*entities.EdgeController, derrors.Error) {
	retrieved, err := m.ControllerProvider.Get(request.EdgeControllerId)
	if err != nil {
		return nil, err
	}
	if retrieved.OrganizationId != request.OrganizationId {
		return nil, derrors.NewNotFoundError("organization_id & asset_id").WithParams(request.OrganizationId, request.EdgeControllerId)
	}
	retrieved.ApplyUpdate(request)
	err = m.ControllerProvider.Update(*retrieved)
	if err != nil {
		return nil, err
	}
	return retrieved, nil
}

func (m *Manager) Get(edgeControllerID *grpc_inventory_go.EdgeControllerId) (*entities.EdgeController, derrors.Error) {
	retrieved, err := m.ControllerProvider.Get(edgeControllerID.EdgeControllerId)
	if err != nil {
		return nil, err
	}
	if retrieved.OrganizationId != edgeControllerID.OrganizationId {
		return nil, derrors.NewNotFoundError("organization_id & edge_controller_id").WithParams(edgeControllerID.OrganizationId, edgeControllerID.EdgeControllerId)
	}
	return retrieved, nil
}
