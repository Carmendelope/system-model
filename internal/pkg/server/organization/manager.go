/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package organization

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
)

type Manager struct {
	Provider organization.Provider
}

func NewManager(provider organization.Provider) Manager{
	return Manager{provider}
}

func (m *Manager) AddOrganization(toAdd grpc_organization_go.AddOrganizationRequest) (* entities.Organization, derrors.Error) {
	newOrg := entities.NewOrganization(toAdd.Name)
	err := m.Provider.Add(*newOrg)
	if err != nil {
		return nil, err
	}
	return newOrg, nil
}

func (m *Manager) GetOrganization(orgID grpc_organization_go.OrganizationId) (* entities.Organization, derrors.Error) {
	return m.Provider.Get(orgID.OrganizationId)
}