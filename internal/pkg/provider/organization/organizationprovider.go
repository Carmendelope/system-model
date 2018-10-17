/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package organization

import (
	"github.com/nalej/derrors"
	"github.com/nalej/system-model/internal/pkg/entities"
)

type Provider interface {
	// Add a new organization to the system.
	Add(org entities.Organization) derrors.Error
	// Check if an organization exists on the system.
	Exists(organizationID string) bool
	// Get an organization.
	Get(organizationID string) (* entities.Organization, derrors.Error)

	// AddCluster adds a new cluster ID to the organization.
	AddCluster(organizationID string, clusterID string) derrors.Error
	// ClusterExists checks if a cluster is linked to an organization.
	ClusterExists(organizationID string, clusterID string) bool
	// ListClusters returns a list of clusters in an organization.
	ListClusters(organizationID string) ([]string, derrors.Error)
	// DeleteCluster removes a cluster from an organization.
	DeleteCluster(organizationID string, clusterID string) derrors.Error

	// AddNode adds a new node ID to the organization.
	AddNode(organizationID string, nodeID string) derrors.Error
	// NodeExists checks if a node is linked to an organization.
	NodeExists(organizationID string, nodeID string) bool
	// ListNodes returns a list of nodes in an organization.
	ListNodes(organizationID string) ([]string, derrors.Error)
	// DeleteNode removes a node from an organization.
	DeleteNode(organizationID string, nodeID string) derrors.Error

	// AddDescriptor adds a new descriptor ID to a given organization.
	AddDescriptor(organizationID string, appDescriptorID string) derrors.Error
	// DescriptorExists checks if an application descriptor exists on the system.
	DescriptorExists(organizationID string, appDescriptorID string) bool
	// ListDescriptors returns the identifiers of the application descriptors associated with an organization.
	ListDescriptors(organizationID string) ([]string, derrors.Error)
	// DeleteDescriptor removes a descriptor from an organization
	DeleteDescriptor(organizationID string, appDescriptorID string) derrors.Error

	// AddInstance adds a new application instance ID to a given organization.
	AddInstance(organizationID string, appInstanceID string) derrors.Error
	// InstanceExists checks if an application instance exists on the system.
	InstanceExists(organizationID string, appInstanceID string) bool
	// ListInstances returns a the identifiers associate with a given organization.
	ListInstances(organizationID string) ([]string, derrors.Error)
	// DeleteInstance removes an instance from an organization
	DeleteInstance(organizationID string, appInstanceID string) derrors.Error
}
