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
	"github.com/gocql/gocql"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/system-model/internal/pkg/entities"
	"github.com/rs/zerolog/log"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	"sync"
)

const organizationTable = "Organizations"
const organizationTablePK = "id"
const organizationClusterTable = "Organization_Clusters"
const organizationNodeTable = "Organization_Nodes"
const organizationDescriptorTable = "organization_appdescriptors"
const organizationInstanceTable = "Organization_appinstances"
const organizationUserTable = "Organization_Users"
const organizationRoleTable = "Organization_Roles"
const organizationTableIndex = "name"

const rowNotFound = "not found"

type ScyllaOrganizationProvider struct {
	Address  string
	Port     int
	Keyspace string
	Session  *gocql.Session
	sync.Mutex
}

func NewScyllaOrganizationProvider(address string, port int, keyspace string) *ScyllaOrganizationProvider {
	org := ScyllaOrganizationProvider{Address: address, Port: port, Keyspace: keyspace, Session: nil}
	org.connect()
	return &org
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

	if sp.Session != nil {
		sp.Session.Close()
		sp.Session = nil
	}
}

// check that the session is created
func (sp *ScyllaOrganizationProvider) checkConnection() derrors.Error {
	if sp.Session == nil {
		return derrors.NewGenericError("Session not created")
	}
	return nil
}

func (sp *ScyllaOrganizationProvider) checkAndConnect() derrors.Error {

	err := sp.checkConnection()
	if err != nil {
		// try to reconnect
		log.Info().Msg("session no created, trying to reconnect...")
		err = sp.connect()
		if err != nil {
			return err
		}
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

func (sp *ScyllaOrganizationProvider) unsafeExists(organizationID string) (bool, derrors.Error) {

	var returnedId string

	stmt, names := qb.Select(organizationTable).Columns(organizationTablePK).Where(qb.Eq(organizationTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		organizationTablePK: organizationID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if organization exists")
		}
	}

	return true, nil
}

func (sp *ScyllaOrganizationProvider) unsafeExistsByName(name string) (bool, derrors.Error) {

	var returnedId string

	stmt, names := qb.Select(organizationTable).Columns(organizationTableIndex).Where(qb.Eq(organizationTableIndex)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		organizationTableIndex: name})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if organization exists")
		}
	}

	return true, nil
}

func (sp *ScyllaOrganizationProvider) unsafeClusterExists(organizationID string, clusterID string) (bool, derrors.Error) {

	var returnedId string

	stmt, names := qb.Select(organizationClusterTable).Columns("cluster_id").Where(qb.Eq("organization_id")).Where(qb.Eq("cluster_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"cluster_id":      clusterID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if cluster exists")
		}
	}

	return true, nil
}

func (sp *ScyllaOrganizationProvider) unsafeNodeExists(organizationID string, nodeID string) (bool, derrors.Error) {

	var returnedId string

	stmt, names := qb.Select(organizationNodeTable).Columns("node_id").Where(qb.Eq("organization_id")).Where(qb.Eq("node_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"node_id":         nodeID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if node exists")
		}
	}

	return true, nil

}

func (sp *ScyllaOrganizationProvider) unsafeDescriptorExists(organizationID string, appDescriptorID string) (bool, derrors.Error) {

	var returnedId string

	stmt, names := qb.Select(organizationDescriptorTable).Columns("app_descriptor_id").Where(qb.Eq("organization_id")).Where(qb.Eq("app_descriptor_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id":   organizationID,
		"app_descriptor_id": appDescriptorID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if descriptor exists")
		}
	}

	return true, nil

}

func (sp *ScyllaOrganizationProvider) unsafeInstanceExists(organizationID string, appInstanceID string) (bool, derrors.Error) {

	var returnedId string

	stmt, names := qb.Select(organizationInstanceTable).Columns("app_instance_id").Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"app_instance_id": appInstanceID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate id instance exists")
		}
	}

	return true, nil

}

func (sp *ScyllaOrganizationProvider) unsafeUserExists(organizationID string, email string) (bool, derrors.Error) {

	var returnedId string

	stmt, names := qb.Select(organizationUserTable).Columns("email").Where(qb.Eq("organization_id")).Where(qb.Eq("email")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"email":           email})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if user exists")
		}
	}

	return true, nil

}

func (sp *ScyllaOrganizationProvider) unsafeRoleExists(organizationID string, roleID string) (bool, derrors.Error) {

	var returnedId string

	stmt, names := qb.Select(organizationRoleTable).Columns("role_id").Where(qb.Eq("organization_id")).Where(qb.Eq("role_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"role_id":         roleID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if role exists")
		}
	}

	return true, nil

}

// --------------------------------------------------------------------------------------------------------------------

// Add a new organization to the system.
func (sp *ScyllaOrganizationProvider) Add(org entities.Organization) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	// check if the organization exists
	exists, err := sp.unsafeExists(org.ID)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError(org.ID)
	}
	exists, err = sp.unsafeExistsByName(org.Name)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError(org.Name)
	}

	// insert the organization instance
	stmt, names := qb.Insert(organizationTable).Columns("id", "name", "created").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindStruct(org)
	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add organization")
	}

	return nil
}

// Check if an organization exists on the system.
func (sp *ScyllaOrganizationProvider) Exists(organizationID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	var returnedId string

	if err := sp.checkAndConnect(); err != nil {
		return false, err
	}
	stmt, names := qb.Select(organizationTable).Columns(organizationTablePK).Where(qb.Eq(organizationTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		organizationTablePK: organizationID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if organization exists")
		}
	}

	return true, nil
}

func (sp *ScyllaOrganizationProvider) ExistsByName(name string) (bool, derrors.Error) {

	var returnedName string

	if err := sp.checkAndConnect(); err != nil {
		return false, err
	}

	stmt, names := qb.Select(organizationTable).Columns(organizationTableIndex).Where(qb.Eq(organizationTableIndex)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		organizationTableIndex: name})

	err := q.GetRelease(&returnedName)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, conversions.ToDerror(err)
		}
	}

	return true, nil

}

// Get an organization.
func (sp *ScyllaOrganizationProvider) Get(organizationID string) (*entities.Organization, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	var org entities.Organization
	stmt, names := qb.Select(organizationTable).Where(qb.Eq(organizationTablePK)).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		organizationTablePK: organizationID,
	})

	cqlErr := q.GetRelease(&org)
	if cqlErr != nil {
		if cqlErr.Error() == rowNotFound {
			return nil, derrors.NewNotFoundError(organizationID)
		} else {
			return nil, derrors.AsError(cqlErr, "cannot get organization")
		}
	}

	return &org, nil
}

func (sp *ScyllaOrganizationProvider) List() ([]entities.Organization, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	stmt, names := qb.Select(organizationTable).ToCql()

	q := gocqlx.Query(sp.Session.Query(stmt), names)

	organizations := make([]entities.Organization, 0)
	cqlErr := q.SelectRelease(&organizations)

	if cqlErr != nil {
		return nil, derrors.AsError(cqlErr, "cannot list organization")
	}

	return organizations, nil
}

// --------------------------------------------------------------------------------------------------------------------

// AddCluster adds a new cluster ID to the organization.
func (sp *ScyllaOrganizationProvider) AddCluster(organizationID string, clusterID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	exists, err := sp.unsafeExists(organizationID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	// check if the organization exists
	exists, err = sp.unsafeClusterExists(organizationID, clusterID)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError("cluster").WithParams(organizationID, clusterID)
	}

	// insert the organization instance
	stmt, names := qb.Insert(organizationClusterTable).Columns("organization_id", "cluster_id").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"cluster_id":      clusterID})

	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add cluster")
	}

	return nil
}

// ClusterExists checks if a cluster is linked to an organization.
func (sp *ScyllaOrganizationProvider) ClusterExists(organizationID string, clusterID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil {
		return false, err
	}

	var returnedId string

	stmt, names := qb.Select(organizationClusterTable).Columns("cluster_id").Where(qb.Eq("organization_id")).Where(qb.Eq("cluster_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"cluster_id":      clusterID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if cluster exists")
		}
	}

	return true, nil
}

// ListClusters returns a list of clusters in an organization.
func (sp *ScyllaOrganizationProvider) ListClusters(organizationID string) ([]string, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	// 1.-Check if the organization exists
	exists, err := sp.unsafeExists(organizationID)
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

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}

	// check if the cluster exists in the organization
	exists, err := sp.unsafeClusterExists(organizationID, clusterID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("cluster").WithParams(organizationID, clusterID)
	}

	// delete a cluster of an organization
	stmt, _ := qb.Delete(organizationClusterTable).Where(qb.Eq("organization_id")).Where(qb.Eq("cluster_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, clusterID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete cluster")
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// AddNode adds a new node ID to the organization.
func (sp *ScyllaOrganizationProvider) AddNode(organizationID string, nodeID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	exists, err := sp.unsafeExists(organizationID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	// check if the organization exists
	exists, err = sp.unsafeNodeExists(organizationID, nodeID)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError("node").WithParams(organizationID, nodeID)
	}

	// add a node in the organization instance
	stmt, names := qb.Insert(organizationNodeTable).Columns("organization_id", "node_id").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"node_id":         nodeID})

	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add node")
	}

	return nil
}

// NodeExists checks if a node is linked to an organization.
func (sp *ScyllaOrganizationProvider) NodeExists(organizationID string, nodeID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil {
		return false, err
	}

	var returnedId string

	stmt, names := qb.Select(organizationNodeTable).Columns("node_id").Where(qb.Eq("organization_id")).Where(qb.Eq("node_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"node_id":         nodeID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if node exists")
		}
	}

	return true, nil

}

// ListNodes returns a list of nodes in an organization.
func (sp *ScyllaOrganizationProvider) ListNodes(organizationID string) ([]string, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	// 1.-Check if the organization exists
	exists, err := sp.unsafeExists(organizationID)
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

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}

	// check if the node exists in the organization
	exists, err := sp.unsafeNodeExists(organizationID, nodeID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("node").WithParams(organizationID, nodeID)
	}

	// delete the node of an organization
	stmt, _ := qb.Delete(organizationNodeTable).Where(qb.Eq("organization_id")).Where(qb.Eq("node_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, nodeID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete node")
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// AddDescriptor adds a new descriptor ID to a given organization.
func (sp *ScyllaOrganizationProvider) AddDescriptor(organizationID string, appDescriptorID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	exists, err := sp.unsafeExists(organizationID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	// check if the descriptor exists
	exists, err = sp.unsafeDescriptorExists(organizationID, appDescriptorID)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError("appDescriptor").WithParams(organizationID, appDescriptorID)
	}

	// add an app descriptor in the organization instance
	stmt, names := qb.Insert(organizationDescriptorTable).Columns("organization_id", "app_descriptor_id").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id":   organizationID,
		"app_descriptor_id": appDescriptorID})

	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add descriptor")
	}

	return nil
}

// DescriptorExists checks if an application descriptor exists on the system.
func (sp *ScyllaOrganizationProvider) DescriptorExists(organizationID string, appDescriptorID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil {
		return false, err
	}

	var returnedId string

	stmt, names := qb.Select(organizationDescriptorTable).Columns("app_descriptor_id").Where(qb.Eq("organization_id")).Where(qb.Eq("app_descriptor_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id":   organizationID,
		"app_descriptor_id": appDescriptorID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if descriptor exists")
		}
	}

	return true, nil

}

// ListDescriptors returns the identifiers of the application descriptors associated with an organization.
func (sp *ScyllaOrganizationProvider) ListDescriptors(organizationID string) ([]string, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	// 1.-Check if the organization exists
	exists, err := sp.unsafeExists(organizationID)
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

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}

	// check if the descriptor exists in the organization
	exists, err := sp.unsafeDescriptorExists(organizationID, appDescriptorID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("app descriptor").WithParams(organizationID, appDescriptorID)
	}

	// delete the descriptor of an organization
	stmt, _ := qb.Delete(organizationDescriptorTable).Where(qb.Eq("organization_id")).Where(qb.Eq("app_descriptor_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, appDescriptorID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete descriptor")
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// AddInstance adds a new application instance ID to a given organization.
func (sp *ScyllaOrganizationProvider) AddInstance(organizationID string, appInstanceID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	exists, err := sp.unsafeExists(organizationID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	// check if the instance exists
	exists, err = sp.unsafeInstanceExists(organizationID, appInstanceID)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError("app_instance").WithParams(organizationID, appInstanceID)
	}

	// add an app instance in the organization
	stmt, names := qb.Insert(organizationInstanceTable).Columns("organization_id", "app_instance_id").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"app_instance_id": appInstanceID})

	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add instance")
	}

	return nil
}

// InstanceExists checks if an application instance exists on the system.
func (sp *ScyllaOrganizationProvider) InstanceExists(organizationID string, appInstanceID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil {
		return false, err
	}

	var returnedId string

	stmt, names := qb.Select(organizationInstanceTable).Columns("app_instance_id").Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"app_instance_id": appInstanceID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate id instance exists")
		}
	}

	return true, nil

}

// ListInstances returns a the identifiers associate with a given organization.
func (sp *ScyllaOrganizationProvider) ListInstances(organizationID string) ([]string, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	// 1.-Check if the organization exists
	exists, err := sp.unsafeExists(organizationID)
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

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}

	// check if the instance exists in the organization
	exists, err := sp.unsafeInstanceExists(organizationID, appInstanceID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("app descriptor").WithParams(organizationID, appInstanceID)
	}

	// delete the instance of an organization
	stmt, _ := qb.Delete(organizationInstanceTable).Where(qb.Eq("organization_id")).Where(qb.Eq("app_instance_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, appInstanceID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete instance")
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// AddUser adds a new user to the organization.
func (sp *ScyllaOrganizationProvider) AddUser(organizationID string, email string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	exists, err := sp.unsafeExists(organizationID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	// check if the user exists
	exists, err = sp.unsafeUserExists(organizationID, email)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError("user").WithParams(organizationID, email)
	}

	// add an user in the organization
	stmt, names := qb.Insert(organizationUserTable).Columns("organization_id", "email").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"email":           email})

	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add user")
	}

	return nil
}

// UserExists checks if a user is linked to an organization.
func (sp *ScyllaOrganizationProvider) UserExists(organizationID string, email string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil {
		return false, err
	}

	var returnedId string

	stmt, names := qb.Select(organizationUserTable).Columns("email").Where(qb.Eq("organization_id")).Where(qb.Eq("email")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"email":           email})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if user exists")
		}
	}

	return true, nil

}

// ListUser returns a list of users in an organization.
func (sp *ScyllaOrganizationProvider) ListUsers(organizationID string) ([]string, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	// 1.-Check if the organization exists
	exists, err := sp.unsafeExists(organizationID)
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

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}

	// check if the user exists in the organization
	exists, err := sp.unsafeUserExists(organizationID, email)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("user").WithParams(organizationID, email)
	}

	// delete the user of an organization
	stmt, _ := qb.Delete(organizationUserTable).Where(qb.Eq("organization_id")).Where(qb.Eq("email")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, email).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete user")
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// AddRole adds a new role ID to the organization.
func (sp *ScyllaOrganizationProvider) AddRole(organizationID string, roleID string) derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	exists, err := sp.unsafeExists(organizationID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("organization").WithParams(organizationID)
	}

	// check if the role exists
	exists, err = sp.unsafeRoleExists(organizationID, roleID)
	if err != nil {
		return err
	}
	if exists {
		return derrors.NewAlreadyExistsError("role").WithParams(organizationID, roleID)
	}

	// add an user in the organization
	stmt, names := qb.Insert(organizationRoleTable).Columns("organization_id", "role_id").ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"role_id":         roleID})

	cqlErr := q.ExecRelease()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot add role")
	}

	return nil
}

// RoleExists checks if a role is linked to an organization.
func (sp *ScyllaOrganizationProvider) RoleExists(organizationID string, roleID string) (bool, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil {
		return false, err
	}

	var returnedId string

	stmt, names := qb.Select(organizationRoleTable).Columns("role_id").Where(qb.Eq("organization_id")).Where(qb.Eq("role_id")).ToCql()
	q := gocqlx.Query(sp.Session.Query(stmt), names).BindMap(qb.M{
		"organization_id": organizationID,
		"role_id":         roleID})

	err := q.GetRelease(&returnedId)
	if err != nil {
		if err.Error() == rowNotFound {
			return false, nil
		} else {
			return false, derrors.AsError(err, "cannot determinate if role exists")
		}
	}

	return true, nil

}

// ListNodes returns a list of roles in an organization.
func (sp *ScyllaOrganizationProvider) ListRoles(organizationID string) ([]string, derrors.Error) {

	sp.Lock()
	defer sp.Unlock()

	if err := sp.checkAndConnect(); err != nil {
		return nil, err
	}

	// 1.-Check if the organization exists
	exists, err := sp.unsafeExists(organizationID)
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

	// check connection
	err := sp.checkAndConnect()
	if err != nil {
		return err
	}

	// check if the role exists in the organization
	exists, err := sp.unsafeRoleExists(organizationID, roleID)
	if err != nil {
		return err
	}
	if !exists {
		return derrors.NewNotFoundError("role").WithParams(organizationID, roleID)
	}

	// delete the role of the organization
	stmt, _ := qb.Delete(organizationRoleTable).Where(qb.Eq("organization_id")).Where(qb.Eq("role_id")).ToCql()
	cqlErr := sp.Session.Query(stmt, organizationID, roleID).Exec()

	if cqlErr != nil {
		return derrors.AsError(cqlErr, "cannot delete role")
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

func (sp *ScyllaOrganizationProvider) Clear() derrors.Error {

	sp.Lock()
	defer sp.Unlock()

	// check connection
	if err := sp.checkAndConnect(); err != nil {
		return err
	}

	// delete organizations table
	err := sp.Session.Query("TRUNCATE TABLE organizations").Exec()
	if err != nil {
		return derrors.AsError(err, "cannot truncate organizations table")
	}

	// delete organization-cluster table
	err = sp.Session.Query("TRUNCATE TABLE organization_clusters").Exec()
	if err != nil {
		return derrors.AsError(err, "cannot truncate organization_cluster table")
	}

	// delete organization-nodes table
	err = sp.Session.Query("TRUNCATE TABLE organization_nodes").Exec()
	if err != nil {
		return derrors.AsError(err, "cannot truncate organization_nodes cluster")
	}

	// delete organization-descriptors table
	err = sp.Session.Query("TRUNCATE TABLE organization_appdescriptors").Exec()
	if err != nil {
		return derrors.AsError(err, "cannot truncate organization_appdescriptors table")
	}

	// delete organization-instances table
	err = sp.Session.Query("TRUNCATE TABLE organization_appinstances").Exec()
	if err != nil {
		return derrors.AsError(err, "cannot truncate organization_appinstances table")
	}

	// delete organization-users table
	err = sp.Session.Query("TRUNCATE TABLE organization_users").Exec()
	if err != nil {
		return derrors.AsError(err, "cannot truncate organization_users table")
	}

	// delete organization-roles table
	err = sp.Session.Query("TRUNCATE TABLE organization_roles").Exec()
	if err != nil {
		return derrors.AsError(err, "cannot truncate organization_roles table")
	}

	return nil
}
