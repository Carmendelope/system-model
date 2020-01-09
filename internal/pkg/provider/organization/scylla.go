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
	"github.com/gocql/gocql"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/scylladb-utils/pkg/scylladb"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"sync"
)

const organizationTable = "Organizations"
const organizationClusterTable = "Organization_Clusters"
const organizationNodeTable = "Organization_Nodes"
const organizationDescriptorTable = "organization_appdescriptors"
const organizationInstanceTable = "Organization_appinstances"
const organizationUserTable = "Organization_Users"
const organizationRoleTable = "Organization_Roles"

// PKs
const organizationTablePK = "id"
const organizationTableIndex = "name"

// columns
var organizationTableColumns = []string{"id", "name", "full_address", "city", "state", "country", "zip_code", "photo_base64", "created"}
var organizationTableColumnsNoPK = []string{"name", "full_address", "city", "state", "country", "zip_code", "photo_base64", "created"}
var organizationClusterTableColumns = []string{"organization_id", "cluster_id"}
var organizationNodeTableColumns = []string{"organization_id", "node_id"}
var organizationDescriptorTableColumns = []string{"organization_id", "app_descriptor_id"}
var organizationInstanceTableColumns = []string{"organization_id", "app_instance_id"}
var organizationUserTableColumns = []string{"organization_id", "email"}
var organizationRoleTableColumns = []string{"organization_id", "role_id"}


type ScyllaOrganizationProvider struct {
	scylladb.ScyllaDB
	sync.Mutex
}

func NewScyllaOrganizationProvider(address string, port int, keyspace string) *ScyllaOrganizationProvider {
	provider := ScyllaOrganizationProvider{
		ScyllaDB: scylladb.ScyllaDB{
			Address:  address,
			Port:     port,
			Keyspace: keyspace,
		},
	}
	provider.Connect()
	return &provider
}

// connect to the database
func (sp *ScyllaOrganizationProvider) connect() derrors.Error {

	// connect to the cluster
	conf := gocql.NewCluster(sp.Address)
	conf.Keyspace = sp.Keyspace
	conf.Port = sp.Port

	session, err := conf.CreateSession()
	if err != nil {
		log.Error().Str("provider", "ScyllaOrganizationProvider").Str("trace", conversions.ToDerror(err).DebugReport()).Msg("unable to connect")
		return derrors.AsError(err, "cannot connect")
	}

	sp.Session = session

	return nil
}

// disconnect from the database
func (sp *ScyllaOrganizationProvider) Disconnect() {

	sp.Lock()
	defer sp.Unlock()

}

func (sp *ScyllaOrganizationProvider) createOrganizationClusterPKMap(OrganizationID string, clusterID string) map[string]interface{} {

	res := map[string]interface{}{
		"organization_id": OrganizationID,
		"cluster_id":      clusterID,
	}

	return res
}

func (sp *ScyllaOrganizationProvider) createOrganizationNodePKMap(OrganizationID string, nodeID string) map[string]interface{} {

	res := map[string]interface{}{
		"organization_id": OrganizationID,
		"node_id":         nodeID,
	}

	return res
}

func (sp *ScyllaOrganizationProvider) createOrganizationDescriptorPKMap(OrganizationID string, appDescriptorId string) map[string]interface{} {

	res := map[string]interface{}{
		"organization_id": OrganizationID,
		"app_descriptor_id":         appDescriptorId,
	}

	return res
}

func (sp *ScyllaOrganizationProvider) createOrganizationInstanceKMap(OrganizationID string, appInstanceId string) map[string]interface{} {

	res := map[string]interface{}{
		"organization_id": OrganizationID,
		"app_instance_id":         appInstanceId,
	}

	return res
}

func (sp *ScyllaOrganizationProvider) createOrganizationUserKMap(OrganizationID string, email string) map[string]interface{} {

	res := map[string]interface{}{
		"organization_id": OrganizationID,
		"email":         email,
	}

	return res
}

func (sp *ScyllaOrganizationProvider) createOrganizationRoleKMap(OrganizationID string, roleId string) map[string]interface{} {

	res := map[string]interface{}{
		"organization_id": OrganizationID,
		"role_id":         roleId,
	}

	return res
}

// --------------------------------------------------------------------------------------------------------------------

// Add a new organization to the system.
func (sp *ScyllaOrganizationProvider) Add(org entities.Organization) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	//check if exists an organization with the same name
	exists, err := sp.UnsafeGenericExist(organizationTable, organizationTableIndex, org.Name)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError(org.Name)
	}
	return sp.UnsafeAdd(organizationTable, organizationTablePK, org.ID, organizationTableColumns, org)

}

// Check if an organization exists on the system.
func (sp *ScyllaOrganizationProvider) Exists(organizationID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeGenericExist(organizationTable, organizationTablePK, organizationID)

}

func (sp *ScyllaOrganizationProvider) unsafeExistsByName(name string) (bool, derrors.Error) {
	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return false, err
	}

	stmt, names := qb.Select(organizationTable).Columns(organizationTableIndex).Where(qb.Eq(organizationTableIndex)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).Bind(name)

	returned := make([]string, 0)
	cqlErr := q.SelectRelease(&returned)

	if cqlErr != nil {
		return false, derrors.AsError(cqlErr, "cannot list organization")
	}

	return len(returned) > 0, nil
}

func (sp *ScyllaOrganizationProvider) ExistsByName(name string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	return sp.unsafeExistsByName(name)

}

// Get an organization.
func (sp *ScyllaOrganizationProvider) Get(organizationID string) (*entities.Organization, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	var organization interface{} = &entities.Organization{}

	err := sp.UnsafeGet(organizationTable, organizationTablePK, organizationID, organizationTableColumns, &organization)
	if err != nil {
		return nil, err
	}
	return organization.(*entities.Organization), nil

}

func (sp *ScyllaOrganizationProvider) List() ([]entities.Organization, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select(organizationTable).Columns(organizationTableColumns...).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names)

	organizations := make([]entities.Organization, 0)
	cqlErr := q.SelectRelease(&organizations)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list organization")
	}

	return organizations, nil
}

func (sp *ScyllaOrganizationProvider) Update(org entities.Organization) derrors.Error {
	sp.Lock()
	defer sp.Unlock()

	// 1.- Check if organization exists
	exists, err := sp.UnsafeGenericExist(organizationTable, organizationTablePK, org.ID)
	if err != nil {
		return err
	}
	if ! exists {
		return derrors.NewNotFoundError(org.ID)
	}

	// 2.- get the organization to check if the name is being updated
	var retrieved interface{} = &entities.Organization{}
	err = sp.UnsafeGet(organizationTable, organizationTablePK, org.ID, organizationTableColumns, &retrieved)
	if err != nil {
		return err
	}

	// 3.- Check the name
	if retrieved.(*entities.Organization).Name != org.Name {
		log.Debug().Msg("The name is being updated")

		exists, err := sp.unsafeExistsByName(org.Name)
		if err != nil {
			return err
		}
		if exists {
			return derrors.NewAlreadyExistsError("unable to update the organization").WithParams(org.Name)
		}
	}
	// 4.- Update
	return sp.UnsafeUpdate(organizationTable, organizationTablePK, org.ID, organizationTableColumnsNoPK, org)

}


// --------------------------------------------------------------------------------------------------------------------

// AddCluster adds a new cluster ID to the organization.
func (sp *ScyllaOrganizationProvider) AddCluster(organizationID string, clusterID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	exists, err := sp.UnsafeGenericExist(organizationTable, organizationTablePK, organizationID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	// insert the organization instance
	pkColumn := sp.createOrganizationClusterPKMap(organizationID, clusterID)

	record := entities.NewOrganizationCluster(organizationID, clusterID)
	return sp.UnsafeCompositeAdd(organizationClusterTable, pkColumn, organizationClusterTableColumns, record)

}

// ClusterExists checks if a cluster is linked to an organization.
func (sp *ScyllaOrganizationProvider) ClusterExists(organizationID string, clusterID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createOrganizationClusterPKMap(organizationID, clusterID)
	return sp.UnsafeGenericCompositeExist(organizationClusterTable, pkColumn)

}

// ListClusters returns a list of clusters in an organization.
func (sp *ScyllaOrganizationProvider) ListClusters(organizationID string) ([]string, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	// 1.-Check if the organization exists
	exists, err := sp.UnsafeGenericExist(organizationTable, organizationTablePK, organizationID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	stmt, names := qb.Select(organizationClusterTable).Columns("cluster_id").Where(qb.Eq("organization_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
	})

	clusters := make([]string, 0)
	cqlErr := q.SelectRelease(&clusters)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list clusters")
	}

	return clusters, nil
}

// DeleteCluster removes a cluster from an organization.
func (sp *ScyllaOrganizationProvider) DeleteCluster(organizationID string, clusterID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createOrganizationClusterPKMap(organizationID, clusterID)
	return sp.UnsafeCompositeRemove(organizationClusterTable, pkColumn)

}

// --------------------------------------------------------------------------------------------------------------------

// AddNode adds a new node ID to the organization.
func (sp *ScyllaOrganizationProvider) AddNode(organizationID string, nodeID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	exists, err := sp.UnsafeGenericExist(organizationTable, organizationTablePK, organizationID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	// insert the organization instance
	pkColumn := sp.createOrganizationNodePKMap(organizationID, nodeID)

	record := entities.NewOrganizationNode(organizationID, nodeID)
	return sp.UnsafeCompositeAdd(organizationNodeTable, pkColumn, organizationNodeTableColumns, record)
}

// NodeExists checks if a node is linked to an organization.
func (sp *ScyllaOrganizationProvider) NodeExists(organizationID string, nodeID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createOrganizationNodePKMap(organizationID, nodeID)
	return sp.UnsafeGenericCompositeExist(organizationNodeTable, pkColumn)

}

// ListNodes returns a list of nodes in an organization.
func (sp *ScyllaOrganizationProvider) ListNodes(organizationID string) ([]string, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.CheckAndConnect(); err != nil {
		return nil, err
	}

	// 1.-Check if the organization exists
	exists, err := sp.UnsafeGenericExist(organizationTable, organizationTablePK, organizationID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	stmt, names := qb.Select(organizationNodeTable).Columns("node_id").Where(qb.Eq("organization_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
	})

	nodes := make([]string, 0)
	cqlErr := q.SelectRelease(&nodes)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list nodes")
	}

	return nodes, nil
}

// DeleteNode removes a node from an organization.
func (sp *ScyllaOrganizationProvider) DeleteNode(organizationID string, nodeID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createOrganizationNodePKMap(organizationID, nodeID)
	return sp.UnsafeCompositeRemove(organizationNodeTable, pkColumn)
}

// --------------------------------------------------------------------------------------------------------------------

// AddDescriptor adds a new descriptor ID to a given organization.
func (sp *ScyllaOrganizationProvider) AddDescriptor(organizationID string, appDescriptorID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	exists, err := sp.UnsafeGenericExist(organizationTable, organizationTablePK, organizationID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	// insert the organization instance
	pkColumn := sp.createOrganizationDescriptorPKMap(organizationID, appDescriptorID)

	record := entities.NewOrganizationDescriptor(organizationID, appDescriptorID)
	return sp.UnsafeCompositeAdd(organizationDescriptorTable, pkColumn, organizationDescriptorTableColumns, record)
}

// DescriptorExists checks if an application descriptor exists on the system.
func (sp *ScyllaOrganizationProvider) DescriptorExists(organizationID string, appDescriptorID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createOrganizationDescriptorPKMap(organizationID, appDescriptorID)
	return sp.UnsafeGenericCompositeExist(organizationDescriptorTable, pkColumn)

}

// ListDescriptors returns the identifiers of the application descriptors associated with an organization.
func (sp *ScyllaOrganizationProvider) ListDescriptors(organizationID string) ([]string, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	// 1.-Check if the organization exists
	exists, err := sp.UnsafeGenericExist(organizationTable, organizationTablePK, organizationID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	stmt, names := qb.Select(organizationDescriptorTable).Columns("app_descriptor_id").Where(qb.Eq("organization_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
	})

	descriptors := make([]string, 0)
	cqlErr := q.SelectRelease(&descriptors)
	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list descriptors")
	}

	return descriptors, nil
}

// DeleteDescriptor removes a descriptor from an organization
func (sp *ScyllaOrganizationProvider) DeleteDescriptor(organizationID string, appDescriptorID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createOrganizationDescriptorPKMap(organizationID, appDescriptorID)
	return sp.UnsafeCompositeRemove(organizationDescriptorTable, pkColumn)
}

// --------------------------------------------------------------------------------------------------------------------

// AddInstance adds a new application instance ID to a given organization.
func (sp *ScyllaOrganizationProvider) AddInstance(organizationID string, appInstanceID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	exists, err := sp.UnsafeGenericExist(organizationTable, organizationTablePK, organizationID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	// add an app instance in the organization
	pkColumn := sp.createOrganizationInstanceKMap(organizationID, appInstanceID)

	record := entities.NewOrganizationInstance(organizationID, appInstanceID)
	return sp.UnsafeCompositeAdd(organizationInstanceTable, pkColumn, organizationInstanceTableColumns, record)
}

// InstanceExists checks if an application instance exists on the system.
func (sp *ScyllaOrganizationProvider) InstanceExists(organizationID string, appInstanceID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createOrganizationInstanceKMap(organizationID, appInstanceID)
	return sp.UnsafeGenericCompositeExist(organizationInstanceTable, pkColumn)

}

// ListInstances returns a the identifiers associate with a given organization.
func (sp *ScyllaOrganizationProvider) ListInstances(organizationID string) ([]string, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	// 1.-Check if the organization exists
	exists, err := sp.UnsafeGenericExist(organizationTable, organizationTablePK, organizationID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	stmt, names := qb.Select(organizationInstanceTable).Columns("app_instance_id").Where(qb.Eq("organization_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
	})

	instances := make([]string, 0)
	cqlErr := q.SelectRelease(&instances)
	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list instances")
	}

	return instances, nil
}

// DeleteInstance removes an instance from an organization
func (sp *ScyllaOrganizationProvider) DeleteInstance(organizationID string, appInstanceID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createOrganizationInstanceKMap(organizationID, appInstanceID)
	return sp.UnsafeCompositeRemove(organizationInstanceTable, pkColumn)
}

// --------------------------------------------------------------------------------------------------------------------

// AddUser adds a new user to the organization.
func (sp *ScyllaOrganizationProvider) AddUser(organizationID string, email string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	exists, err := sp.UnsafeGenericExist(organizationTable, organizationTablePK, organizationID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	// add an app instance in the organization
	pkColumn := sp.createOrganizationUserKMap(organizationID, email)

	record := entities.NewOrganizationUser(organizationID, email)
	return sp.UnsafeCompositeAdd(organizationUserTable, pkColumn, organizationUserTableColumns, record)
}

// UserExists checks if a user is linked to an organization.
func (sp *ScyllaOrganizationProvider) UserExists(organizationID string, email string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createOrganizationUserKMap(organizationID, email)
	return sp.UnsafeGenericCompositeExist(organizationUserTable, pkColumn)

}

// ListUser returns a list of users in an organization.
func (sp *ScyllaOrganizationProvider) ListUsers(organizationID string) ([]string, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	// 1.-Check if the organization exists
	exists, err := sp.UnsafeGenericExist(organizationTable, organizationTablePK, organizationID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	stmt, names := qb.Select(organizationUserTable).Columns("email").Where(qb.Eq("organization_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
	})

	users := make([]string, 0)
	cqlErr := q.SelectRelease(&users)
	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list users")
	}

	return users, nil
}

// DeleteUser removes a user from an organization.
func (sp *ScyllaOrganizationProvider) DeleteUser(organizationID string, email string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createOrganizationUserKMap(organizationID, email)
	return sp.UnsafeCompositeRemove(organizationUserTable, pkColumn)
}

// --------------------------------------------------------------------------------------------------------------------

// AddRole adds a new role ID to the organization.
func (sp *ScyllaOrganizationProvider) AddRole(organizationID string, roleID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	exists, err := sp.UnsafeGenericExist(organizationTable, organizationTablePK, organizationID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	// add an app instance in the organization
	pkColumn := sp.createOrganizationRoleKMap(organizationID, roleID)

	record := entities.NewOrganizationRole(organizationID, roleID)
	return sp.UnsafeCompositeAdd(organizationRoleTable, pkColumn, organizationRoleTableColumns, record)
}

// RoleExists checks if a role is linked to an organization.
func (sp *ScyllaOrganizationProvider) RoleExists(organizationID string, roleID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createOrganizationRoleKMap(organizationID, roleID)
	return sp.UnsafeGenericCompositeExist(organizationRoleTable, pkColumn)

}

// ListNodes returns a list of roles in an organization.
func (sp *ScyllaOrganizationProvider) ListRoles(organizationID string) ([]string, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	// 1.-Check if the organization exists
	exists, err := sp.UnsafeGenericExist(organizationTable, organizationTablePK, organizationID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	stmt, names := qb.Select(organizationRoleTable).Columns("role_id").Where(qb.Eq("organization_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
	})

	roles := make([]string, 0)
	cqlErr := q.SelectRelease(&roles)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list roles")
	}

	return roles, nil
}

// DeleteRole removes a role from an organization.
func (sp *ScyllaOrganizationProvider) DeleteRole(organizationID string, roleID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	pkColumn := sp.createOrganizationRoleKMap(organizationID, roleID)
	return sp.UnsafeCompositeRemove(organizationRoleTable, pkColumn)
}

// --------------------------------------------------------------------------------------------------------------------

func (sp *ScyllaOrganizationProvider) Clear() derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	return sp.UnsafeClear([]string{organizationTable, organizationNodeTable, organizationRoleTable, organizationUserTable,
		organizationClusterTable, organizationDescriptorTable, organizationInstanceTable})

}
