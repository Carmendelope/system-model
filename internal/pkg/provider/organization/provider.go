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
	"github.com/nalej/system-model/internal/pkg/entities"
)

type Provider interface {
	// Add a new organization to the system.
	Add(org entities.Organization) derrors.Error
	// Check if an organization exists on the system.
	Exists(organizationID string) (bool, derrors.Error)
	// Check if an organization with this name exists on the system
	ExistsByName(name string) (bool, derrors.Error)
	// Get an organization.
	Get(organizationID string) (*entities.Organization, derrors.Error)
	// List the set of organizations.
	List() ([]entities.Organization, derrors.Error)

	// AddCluster adds a new cluster ID to the organization.
	AddCluster(organizationID string, clusterID string) derrors.Error
	// ClusterExists checks if a cluster is linked to an organization.
	ClusterExists(organizationID string, clusterID string) (bool, derrors.Error)
	// ListClusters returns a list of clusters in an organization.
	ListClusters(organizationID string) ([]string, derrors.Error)
	// DeleteCluster removes a cluster from an organization.
	DeleteCluster(organizationID string, clusterID string) derrors.Error

	// AddNode adds a new node ID to the organization.
	AddNode(organizationID string, nodeID string) derrors.Error
	// NodeExists checks if a node is linked to an organization.
	NodeExists(organizationID string, nodeID string) (bool, derrors.Error)
	// ListNodes returns a list of nodes in an organization.
	ListNodes(organizationID string) ([]string, derrors.Error)
	// DeleteNode removes a node from an organization.
	DeleteNode(organizationID string, nodeID string) derrors.Error

	// AddDescriptor adds a new descriptor ID to a given organization.
	AddDescriptor(organizationID string, appDescriptorID string) derrors.Error
	// DescriptorExists checks if an application descriptor exists on the system.
	DescriptorExists(organizationID string, appDescriptorID string) (bool, derrors.Error)
	// ListDescriptors returns the identifiers of the application descriptors associated with an organization.
	ListDescriptors(organizationID string) ([]string, derrors.Error)
	// DeleteDescriptor removes a descriptor from an organization
	DeleteDescriptor(organizationID string, appDescriptorID string) derrors.Error

	// AddInstance adds a new application instance ID to a given organization.
	AddInstance(organizationID string, appInstanceID string) derrors.Error
	// InstanceExists checks if an application instance exists on the system.
	InstanceExists(organizationID string, appInstanceID string) (bool, derrors.Error)
	// ListInstances returns a the identifiers associate with a given organization.
	ListInstances(organizationID string) ([]string, derrors.Error)
	// DeleteInstance removes an instance from an organization
	DeleteInstance(organizationID string, appInstanceID string) derrors.Error

	// AddUser adds a new user to the organization.
	AddUser(organizationID string, email string) derrors.Error
	// UserExists checks if a user is linked to an organization.
	UserExists(organizationID string, email string) (bool, derrors.Error)
	// ListUser returns a list of users in an organization.
	ListUsers(organizationID string) ([]string, derrors.Error)
	// DeleteUser removes a user from an organization.
	DeleteUser(organizationID string, email string) derrors.Error

	// AddRole adds a new role ID to the organization.
	AddRole(organizationID string, roleID string) derrors.Error
	// RoleExists checks if a role is linked to an organization.
	RoleExists(organizationID string, roleID string) (bool, derrors.Error)
	// ListNodes returns a list of roles in an organization.
	ListRoles(organizationID string) ([]string, derrors.Error)
	// DeleteRole removes a role from an organization.
	DeleteRole(organizationID string, roleID string) derrors.Error

	Clear() derrors.Error
}
