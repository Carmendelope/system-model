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

package organization

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
)

// Handler structure for the organization requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

// AddOrganization adds a new organization to the system.
func (h *Handler) AddOrganization(ctx context.Context, addOrganizationRequest *grpc_organization_go.AddOrganizationRequest) (*grpc_organization_go.Organization, error) {
	log.Debug().Str("name", addOrganizationRequest.Name).Msg("add organization")
	err := entities.ValidAddOrganizationRequest(addOrganizationRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("invalid add organization request")
		return nil, conversions.ToGRPCError(err)
	}
	org, err := h.Manager.AddOrganization(*addOrganizationRequest)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot add organization")
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("organizationID", org.ID).Msg("organization has been added")
	return org.ToGRPC(), nil
}

// GetOrganization retrieves the profile information of a given organization.
func (h *Handler) GetOrganization(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_organization_go.Organization, error) {
	retrieved, err := h.Manager.GetOrganization(*organizationID)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot get organization")
		return nil, conversions.ToGRPCError(err)
	}
	return retrieved.ToGRPC(), nil
}

// ListOrganizations returns the list of organizations in the system.
func (h *Handler) ListOrganizations(ctx context.Context, _ *grpc_common_go.Empty) (*grpc_organization_go.OrganizationList, error) {
	retrieved, err := h.Manager.ListOrganization()
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot list organizations")
		return nil, conversions.ToGRPCError(err)
	}
	return entities.OrganizationListToGRPC(retrieved), nil

}

// UpdateOrganization updates the public information of an organization.
func (h *Handler) UpdateOrganization(ctx context.Context, updateOrganizationRequest *grpc_organization_go.UpdateOrganizationRequest) (*grpc_common_go.Success, error) {
	notImplemented := derrors.NewUnimplementedError("update organization")
	log.Error().Str("trace", notImplemented.DebugReport()).Msg("cannot update organization")
	return nil, conversions.ToGRPCError(notImplemented)
}
