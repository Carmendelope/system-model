/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package organization

import (
	"context"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
)

type Handler struct{}

// NewHandler creates a new Handler.
func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) AddOrganization(context.Context, *grpc_organization_go.AddOrganizationRequest) (*grpc_organization_go.Organization, error) {
	return nil, nil
}
// GetOrganization retrieves the profile information of a given organization.
func (h *Handler) GetOrganization(context.Context, *grpc_organization_go.OrganizationId) (*grpc_organization_go.Organization, error) {
	return nil, nil
}
// UpdateOrganization updates the public information of an organization.
func (h *Handler) UpdateOrganization(context.Context, *grpc_organization_go.UpdateOrganizationRequest) (*grpc_common_go.Success, error) {
	return nil, nil
}

