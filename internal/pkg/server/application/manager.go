/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package application

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/nalej/system-model/internal/pkg/provider/application"
	"github.com/nalej/system-model/internal/pkg/provider/organization"
)

type Manager struct {
	OrgProvider organization.Provider
	AppProvider application.Provider
}

func NewManager(orgProvider organization.Provider, appProvider application.Provider) Manager {
	return Manager{orgProvider, appProvider}
}

func (m * Manager) AddAppDescriptor(addRequest grpc_application_go.AddAppDescriptorRequest) (* entities.AppDescriptor, derrors.Error) {
	exists := m.OrgProvider.Exists(addRequest.OrganizationId)
	if !exists{
		return nil, derrors.NewNotFoundError("organizationID").WithParams(addRequest.OrganizationId)
	}
	descriptor := NewApp
	added, err := m.AppProvider.AddDescriptor(descriptor)
}

